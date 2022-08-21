#include "mdict.h"

#include <algorithm>
#include <cstring>
#include <regex>
#include <utility>

const std::regex re_pattern("(\\s|:|\\.|,|-|_|'|\\(|\\)|#|<|>|!)");

namespace mdict {

// constructor
Mdict::Mdict(std::string fn) noexcept : filename(std::move(fn)) {
  instream = std::ifstream(filename, std::ios::binary);
  if(endsWith(filename, ".mdd")) {
      this->filetype = MDDTYPE;
  } else {
      this->filetype = MDXTYPE;
  }
}

// distructor
Mdict::~Mdict() {
  // close instream
  instream.close();
}

/**
 * transform word into comparable string
 * @param word
 * @return
 */
std::string _s(std::string word) {
  std::string s = std::regex_replace(word, re_pattern, "");
  std::transform(s.begin(), s.end(), s.begin(), ::tolower);
  return s;
}



/***************************************
 *             private part            *
 ***************************************/

/**
 * read header
 */
void Mdict::read_header() {
  // -----------------------------------------
  // 1. [0:4] dictionary header length 4 byte
  // -----------------------------------------

  // header size buffer
  char* head_size_buf = (char*)std::calloc(4, sizeof(char));
  readfile(0, 4, head_size_buf);

  // header byte size convert
  uint32_t header_bytes_size =
      be_bin_to_u32((const unsigned char*)head_size_buf);
  std::free(head_size_buf);
  // assign key block start offset
  this->header_bytes_size = header_bytes_size;
  this->key_block_start_offset = this->header_bytes_size + 8;
  /// passed

  // -----------------------------------------
  // 2. [4: header_bytes_size+4], header buffer
  // -----------------------------------------

  // header buffer
  char* head_buffer = (char*)std::calloc(header_bytes_size, sizeof(char));
  readfile(4, header_bytes_size, head_buffer);
  /// passed

  // -----------------------------------------
  // 3. alder32 checksum
  // -----------------------------------------

  // TODO  version < 2.0 needs to checksum?
  // alder32 checksum buffer
  char* head_checksum_buffer = (char*)std::calloc(4, sizeof(char));
  readfile(header_bytes_size + 4, 4, head_checksum_buffer);
  /// passed

  // TODO skip head checksum for now
  std::free(head_checksum_buffer);

  // -----------------------------------------
  // 4. convert header buffer into utf16 text
  // -----------------------------------------

  // header text utf16
  std::string header_text =
      le_bin_utf16_to_utf8(head_buffer, 0, header_bytes_size - 2);
  if (header_text.empty()) {
    std::cout << "len:" << header_bytes_size << std::endl;
    std::cout << "this mdx file is invalid" << std::endl;
    return;
  }
  /// passed

  // -----------------------------------------
  // 5. parse xml string into map
  // -----------------------------------------

  std::map<std::string, std::string> headinfo = parseXMLHeader(header_text);
  /// passed

  // -----------------------------------------
  // 6. handle header message, set flags
  // -----------------------------------------

  // encrypted flag
  // 0x00 - no encryption
  // 0x01 - encrypt record block
  // 0x02 - encrypt key info block
  if (headinfo.find("Encrypted") == headinfo.end() ||
      headinfo["Encrypted"].empty() || headinfo["Encrypted"] == "No") {
    this->encrypt = ENCRYPT_NO_ENC;
  } else if (headinfo["Encrypted"] == "Yes") {
    this->encrypt = ENCRYPT_RECORD_ENC;
  } else {
    std::string s = headinfo["Encrypted"];
    if (s.at(0) == '2') {
      this->encrypt = ENCRYPT_KEY_INFO_ENC;
    } else if (s.at(0) == '1') {
      this->encrypt = ENCRYPT_RECORD_ENC;
    } else {
      this->encrypt = ENCRYPT_NO_ENC;
    }
  }
  /// passed

  // -------- stylesheet ----------
  // stylesheet attribute if present takes from of:
  // style_number # 1-255
  // style_begin # or ''
  // style_end # or ''
  // TODO: splitstyle info

  // header_info['_stylesheet'] = {}
  // if header_tag.get('StyleSheet'):
  //   lines = header_tag['StyleSheet'].splitlines()
  //   for i in range(0, len(lines), 3):
  //        header_info['_stylesheet'][lines[i]] = (lines[i + 1], lines[i + 2])

  // ---------- version ------------
  // before version 2.0, number is 4 bytes integer
  // version 2.0 and above use 8 bytes
  std::string sver = headinfo["GeneratedByEngineVersion"];
  std::string::size_type sz;  // alias of size_t

  float version = std::stof(sver, &sz);
  this->version = version;

  if (this->version >= 2.0) {
    this->number_width = 8;
    this->number_format = NUMFMT_BE_8BYTESQ;
    this->key_block_info_start_offset = this->key_block_start_offset + 40 + 4;
  } else {
    this->number_format = NUMFMT_BE_4BYTESI;
    this->number_width = 4;
    this->key_block_info_start_offset = this->key_block_start_offset + 16;
  }

  // ---------- encoding ------------
  if (headinfo.find("Encoding") != headinfo.end() ||
      headinfo["Encoding"] == "" || headinfo["Encoding"] == "UTF-8") {
    this->encoding = ENCODING_UTF8;
  } else if (headinfo["Encoding"] == "GBK" ||
             headinfo["Encoding"] == "GB2312") {
    this->encoding = ENCODING_GB18030;
  } else if (headinfo["Encoding"] == "Big5" || headinfo["Encoding"] == "BIG5") {
    this->encoding = ENCODING_BIG5;
  } else if (headinfo["Encoding"] == "utf16" ||
             headinfo["Encoding"] == "utf-16") {
    this->encoding = ENCODING_UTF16;
  } else {
    this->encoding = ENCODING_UTF8;
  }
  // FIX mdd
  if (this->filetype == "MDD") {
      this->encoding = ENCODING_UTF16;
  }
  /// passed
}

/**
 * read key block header, key block header contains a serials number, including
 *
 * key block header info struct:
 * [0:8]/[0:4]   - number of key blocks
 * [8:16]/[4:8]  - number of entries
 * [16:24]/nil - key block info decompressed size (if version >= 2.0,
 * otherwise, this section does not exist)
 * [24:32]/[8:12] - key block info size
 * [32:40][12:16] - key block size
 * note: if version <2.0, the key info buffer size is 4 * 4
 *       otherwise, ths key info buffer size is 5 * 8
 * <2.0  the order of number is same
 */
void Mdict::read_key_block_header() {
  // key block header part
  int key_block_info_bytes_num = 0;
  if (this->version >= 2.0) {
    key_block_info_bytes_num = 8 * 5;
  } else {
    key_block_info_bytes_num = 4 * 4;
  }

  // key block info buffer
  char* key_block_info_buffer = (char*)calloc(
      static_cast<size_t>(key_block_info_bytes_num), sizeof(char));
  // read buffer
  this->readfile(this->key_block_start_offset,
                 static_cast<uint64_t>(key_block_info_bytes_num),
                 key_block_info_buffer);
  //  putbytes(key_block_info_buffer,key_block_info_bytes_num, true);
  /// PASSED

  // TODO key block info encrypted file not support yet
  if (this->encrypt == ENCRYPT_RECORD_ENC) {
    std::cout << "user identification is needed to read encrypted file"
              << std::endl;
    if (key_block_info_buffer) std::free(key_block_info_buffer);
    throw std::invalid_argument("invalid encrypted file");
  }

  // key block header info struct:
  // [0:8]/[0:4]   - number of key blocks
  // [8:16]/[4:8]  - number of entries
  // [16:24]/nil - key block info decompressed size (if version >= 2.0,
  // otherwise, this section does not exist)
  // [24:32]/[8:12] - key block info size
  // [32:40][12:16] - key block size
  // note: if version <2.0, the key info buffer size is 4 * 4
  //       otherwise, ths key info buffer size is 5 * 8
  // <2.0  the order of number is same

  // 1. [0:8]([0:4]) number of key blocks
  char* key_block_nums_bytes =
      (char*)calloc(static_cast<size_t>(this->number_width), sizeof(char));
  int eno = bin_slice(key_block_info_buffer, key_block_info_bytes_num, 0,
                      this->number_width, key_block_nums_bytes);
  if (eno != 0) {
    if (key_block_info_buffer) std::free(key_block_info_buffer);
    if (key_block_nums_bytes) std::free(key_block_nums_bytes);
    std::cout << "eno: " << eno << std::endl;
    throw std::logic_error("get key block bin slice failed");
  }
  /// passed

  uint64_t key_block_num = 0;
  if (this->number_width == 8)
    key_block_num = be_bin_to_u64((const unsigned char*)key_block_nums_bytes);
  else if (this->number_width == 4)
    key_block_num = be_bin_to_u32((const unsigned char*)key_block_nums_bytes);
  if (key_block_nums_bytes) std::free(key_block_nums_bytes);
  /// passed

  // 2. [8:16]  - number of entries
  char* entries_num_bytes =
      (char*)calloc(static_cast<size_t>(this->number_width), sizeof(char));
  eno = bin_slice(key_block_info_buffer, key_block_info_bytes_num,
                  this->number_width, this->number_width, entries_num_bytes);
  if (eno != 0) {
    if (key_block_info_buffer) std::free(key_block_info_buffer);
    if (entries_num_bytes) std::free(entries_num_bytes);
    throw std::logic_error("get key block bin slice failed");
  }
  /// passed

  uint64_t entries_num = 0;
  if (this->number_width == 8)
    entries_num = be_bin_to_u64((const unsigned char*)entries_num_bytes);
  else if (this->number_width == 4)
    key_block_num = be_bin_to_u32((const unsigned char*)entries_num_bytes);
  if (entries_num_bytes) std::free(entries_num_bytes);
  /// passed

  int key_block_info_size_start_offset = 0;

  // 3. [16:24] - key block info decompressed size (if version >= 2.0,
  // otherwise, this section does not exist)
  if (this->version >= 2.0) {
    char* key_block_info_decompress_size_bytes =
        (char*)calloc(static_cast<size_t>(this->number_width), sizeof(char));
    eno = bin_slice(key_block_info_buffer, key_block_info_bytes_num,
                    this->number_width * 2, this->number_width,
                    key_block_info_decompress_size_bytes);
    if (eno != 0) {
      if (key_block_info_buffer) std::free(key_block_info_buffer);
      if (key_block_info_decompress_size_bytes)
        std::free(key_block_info_decompress_size_bytes);
      throw std::logic_error("decode key block decompress size failed");
    }
    /// passed

    uint64_t key_block_info_decompress_size = 0;
    if (this->number_width == 8)
      key_block_info_decompress_size = be_bin_to_u64(
          (const unsigned char*)key_block_info_decompress_size_bytes);
    else if (this->number_width == 4)
      key_block_info_decompress_size = be_bin_to_u32(
          (const unsigned char*)key_block_info_decompress_size_bytes);
    this->key_block_info_decompress_size = key_block_info_decompress_size;
    if (key_block_info_decompress_size_bytes)
      std::free(key_block_info_decompress_size_bytes);
    /// passed

    // key block info size (number) start at 24 ([24:32])
    key_block_info_size_start_offset = this->number_width * 3;
  } else {
    // key block info size (number) start at 24 ([8:12])
    key_block_info_size_start_offset = this->number_width * 2;
  }

  // 4. [24:32] - key block info size
  char* key_block_info_size_buffer =
      (char*)calloc(static_cast<size_t>(this->number_width), sizeof(char));
  eno = bin_slice(key_block_info_buffer, key_block_info_bytes_num,
                  key_block_info_size_start_offset, this->number_width,
                  key_block_info_size_buffer);
  if (eno != 0) {
    if (key_block_info_buffer != nullptr) std::free(key_block_info_buffer);
    if (key_block_info_size_buffer != nullptr)
      std::free(key_block_info_size_buffer);
    throw std::logic_error("decode key block info size failed");
  }

  uint64_t key_block_info_size = 0;
  if (this->number_width == 8)
    key_block_info_size =
        be_bin_to_u64((const unsigned char*)key_block_info_size_buffer);
  else if (this->number_width == 4)
    key_block_info_size =
        be_bin_to_u32((const unsigned char*)key_block_info_size_buffer);
  if (key_block_info_size_buffer != nullptr)
    std::free(key_block_info_size_buffer);
  /// passed

  // 5. [32:40] - key block size
  char* key_block_size_buffer =
      (char*)calloc(static_cast<size_t>(this->number_width), sizeof(char));
  eno = bin_slice(key_block_info_buffer, key_block_info_bytes_num,
                  key_block_info_size_start_offset + this->number_width,
                  this->number_width, key_block_size_buffer);
  if (eno != 0) {
    if (key_block_info_buffer) std::free(key_block_info_buffer);
    if (key_block_size_buffer) std::free(key_block_size_buffer);
    throw std::logic_error("decode key block size failed");
  }
  /// passed

  uint64_t key_block_size = 0;
  if (this->number_width == 8)
    key_block_size = be_bin_to_u64((const unsigned char*)key_block_size_buffer);
  else if (this->number_width == 4)
    key_block_size = be_bin_to_u32((const unsigned char*)key_block_size_buffer);
  if (key_block_size_buffer) std::free(key_block_size_buffer);
  /// passed

  // 6. [40:44] - 4bytes checksum
  // TODO if version > 2.0, skip 4bytes checksum

  // free key block info buffer
  if (key_block_info_buffer != nullptr) std::free(key_block_info_buffer);

  this->key_block_num = key_block_num;
  this->entries_num = entries_num;
  this->key_block_info_size = key_block_info_size;
  this->key_block_size = key_block_size;
  if (this->version >= 2.0) {
    this->key_block_info_start_offset = this->key_block_start_offset + 40 + 4;
  } else {
    this->key_block_info_start_offset = this->key_block_start_offset + 16;
  }
}

/**
 * read key block info
 *
 * it will decode the key block info, and set the key block info list
 * it contains:
 * first key
 * last key
 * comp size
 * decomp size
 * offset
 */
void Mdict::read_key_block_info() {
  // start at this->key_block_info_start_offset
  char* key_block_info_buffer = (char*)calloc(
      static_cast<size_t>(this->key_block_info_size), sizeof(char));
  readfile(this->key_block_info_start_offset,
           static_cast<int>(this->key_block_info_size), key_block_info_buffer);

  // ------------------------------------
  // decode key_block_info
  // ------------------------------------
  decode_key_block_info(key_block_info_buffer, this->key_block_info_size,
                        this->key_block_num, this->entries_num);

  // key block compressed start offset = this->key_block_info_start_offset +
  // key_block_info_size
  this->key_block_compressed_start_offset = static_cast<uint32_t>(
      this->key_block_info_start_offset + this->key_block_info_size);

  /// passed

  char* key_block_compressed_buffer =
      (char*)calloc(static_cast<size_t>(this->key_block_size), sizeof(char));

  readfile(this->key_block_compressed_start_offset,
           static_cast<int>(this->key_block_size), key_block_compressed_buffer);

  // ------------------------------------
  // decode key_block_compressed
  // ------------------------------------
  unsigned long kb_len = this->key_block_size;
  //  putbytes(key_block_compressed_buffer,this->key_block_size, true);

  int err =
      decode_key_block((unsigned char*)key_block_compressed_buffer, kb_len);
  if (err != 0) {
    throw std::runtime_error("decode key block error");
  }

  if (key_block_info_buffer != nullptr) std::free(key_block_info_buffer);
  if (key_block_compressed_buffer != nullptr)
    std::free(key_block_compressed_buffer);
}

/**
 * use ripemd128 as decrypt key, and decrypt the key info data
 * @param data the data which needs to decrypt
 * @param k the decrypt key
 * @param data_len data length
 * @param key_len key length
 */
void fast_decrypt(byte* data, const byte* k, int data_len, int key_len) {
  const byte* key = k;
  //      putbytes((char*)data, 16, true);
  byte* b = data;
  byte previous = 0x36;

  for (int i = 0; i < data_len; ++i) {
    byte t = static_cast<byte>(((b[i] >> 4) | (b[i] << 4)) & 0xff);
    t = t ^ previous ^ ((byte)(i & 0xff)) ^ key[i % key_len];
    previous = b[i];
    b[i] = t;
  }
}

/**
 *
 * decrypt the data, this is a helper function to invoke the fast_decrypt
 * note: don't forget free comp_block !!
 *
 * @param comp_block compressed block buffer
 * @param comp_block_len compressed block buffer size
 * @return the decrypted compressed block
 */
byte* mdx_decrypt(byte* comp_block, const int comp_block_len) {
  byte* key_buffer = (byte*)calloc(8, sizeof(byte));
  memcpy(key_buffer, comp_block + 4 * sizeof(char), 4 * sizeof(char));
  key_buffer[4] = 0x95;  // comp_block[4:8] + [0x95,0x36,0x00,0x00]
  key_buffer[5] = 0x36;

  byte* key = ripemd128bytes(key_buffer, 8);

  fast_decrypt(comp_block + 8 * sizeof(byte), key, comp_block_len - 8,
               16 /* key length*/);

  // finally
  std::free(key_buffer);
  return comp_block;
  /// passed
}

/**
 * split key block into key block list
 *
 * this is for key block (not key block info)
 *
 * @param key_block key block buffer
 * @param key_block_len key block length
 */
std::vector<key_list_item*> Mdict::split_key_block(unsigned char* key_block,
                                                   unsigned long key_block_len,
                                                   unsigned long block_id) {
  // TODO assert checksum
  // uint32_t adlchk = adler32checksum(key_block, key_block_len);
  //  std::cout<<"adler32 chksum: "<<adlchk<<std::endl;
  int key_start_idx = 0;
  int key_end_idx = 0;
  std::vector<key_list_item*> inner_key_list;
  unsigned long entry_acc = 0l;

  while (key_start_idx < key_block_len) {
    // # the corresponding record's offset in record block
    unsigned long record_start = 0;
    int width = 0;
    if (this->version >= 2.0) {
      record_start = be_bin_to_u64(key_block + key_start_idx);
    } else {
      record_start = be_bin_to_u32(key_block + key_start_idx);
    }

    if (this->encoding == 1 /* utf16 */) {
      width = 2;
    } else {
      width = 1;
    }

    // key text ends with '\x00'
    // version >= 2.0 delimiter == '0x0000'
    // else delimiter == '0x00'  (< 2.0)
    int i = key_start_idx + number_width;  // ver > 2.0, move 8, else move 4
    if (i >= key_block_len) {
      throw std::runtime_error("key start idx > key block length");
    }
    while (i < key_block_len) {
      if (encoding == 1 /*ENCODING_UTF16*/) {
        if ((key_block[i] & 0x0f) == 0 &&        /* delimiter = '0000' */
            ((key_block[i] & 0xf0) >> 4) == 0 && /* delimiter = '0000' */
            ((key_block[i + 1] & 0x0f) == 0) &&
            (((key_block[i + 1] & 0xf0) >> 4) == 0)) {
          key_end_idx = i;
          break;
        }
      } else {
        // var a = key_block[i]
        // (a >> 4) & 255                     (01011010 >> 4) & 11111111 ->
        // 00000101 & 11111111 -> 00000101
        //
        if ((key_block[i] & 0xf0) >> 4 == 0 && /* delimiter == '0' */
            (key_block[i] & 0x0f) >> 0 == 0) {
          key_end_idx = i;
          break;
        }
      }

      i += width;
    }
    /// passed

    if (key_end_idx >= key_block_len) {
      key_end_idx = static_cast<int>(key_block_len);
    }

    std::string key_text = "";
    if (this->encoding == 1 /* ENCODING_UTF16 */) {
      // TODO
        key_text = be_bin_to_utf16(
                (const char*)key_block, (key_start_idx + this->number_width),
                static_cast<unsigned long>(key_end_idx - key_start_idx -
                                           this->number_width));
//      throw std::runtime_error("NOT SUPPORT UTF16 YET");
    } else if (this->encoding == 0 /* ENCODING_UTF8 */) {
      key_text = be_bin_to_utf8(
          (const char*)key_block, (key_start_idx + this->number_width),
          static_cast<unsigned long>(key_end_idx - key_start_idx -
                                     this->number_width));
    }
    /// passed
    entry_acc++;
    inner_key_list.push_back(
        new key_list_item(record_start, key_text));

    // TODO add to word list

    // next round
    //    int j = key_end_idx;
    //    while((key_block[j] & 0xff) == 0){
    //      j += width;
    //    }
    key_start_idx = key_end_idx + width;
  }
  return inner_key_list;
}

/**
 * decode key block info by block id use with reduce function
 * @param block_id key_block id
 * @return return key list item
 */
std::vector<key_list_item*> Mdict::decode_key_block_by_block_id(
    unsigned long block_id) {
  // ------------------------------------
  // decode key_block_compressed
  // ------------------------------------

  unsigned long idx = block_id;

  unsigned long comp_size = this->key_block_info_list[idx]->key_block_comp_size;
  unsigned long decomp_size =
      this->key_block_info_list[idx]->key_block_decomp_size;
  unsigned long start_ofset =
      this->key_block_info_list[idx]->key_block_comp_accumulator +
      this->key_block_compressed_start_offset;

  char* key_block_buffer =
      (char*)calloc(static_cast<size_t>(comp_size), sizeof(unsigned char));

  readfile(start_ofset, static_cast<int>(comp_size), key_block_buffer);

  // 4 bytes comp type
  char* key_block_comp_type = (char*)calloc(4, sizeof(char));
  memcpy(key_block_comp_type, key_block_buffer, 4 * sizeof(char));
  // 4 bytes adler checksum of decompressed key block
  uint32_t chksum =
      be_bin_to_u32((unsigned char*)key_block_buffer + 4 * sizeof(char));

  unsigned char* key_block = nullptr;
  std::vector<uint8_t> kb_uncompressed;  // note: ensure kb_uncompressed not
                                         // die when out of uncompress scope

  if ((key_block_comp_type[0] & 255) == 0) {
    // none compressed
    key_block = (unsigned char*)(key_block_buffer + 8 * sizeof(char));
  } else if ((key_block_comp_type[0] & 255) == 1) {
    // 01000000
    // TODO lzo decompress

  } else if ((key_block_comp_type[0] & 255) == 2) {
    // zlib compress
    kb_uncompressed =
        zlib_mem_uncompress(key_block_buffer + 8 * sizeof(char), comp_size);
    if (kb_uncompressed.empty()) {
      throw std::runtime_error("key block decompress failed empty");
    }
    key_block = kb_uncompressed.data();

    uint32_t adler32cs =
        adler32checksum(key_block, static_cast<uint32_t>(decomp_size));
    assert(adler32cs == chksum);
    assert(kb_uncompressed.size() == decomp_size);
  } else {
    throw std::runtime_error("cannot determine the key block compress type");
  }

  // split key
  std::vector<key_list_item*> tlist =
      split_key_block(key_block, decomp_size, idx);
  return tlist;
}

/**
 * decode the key block decode function, will invoke split key block
 *
 * this is for key block (not key block info)
 *
 * @param key_block_buffer
 * @param kb_buff_len
 * @return
 */
int Mdict::decode_key_block(unsigned char* key_block_buffer,
                            unsigned long kb_buff_len) {
  int i = 0;

  for (long idx = 0; idx < this->key_block_info_list.size(); idx++) {
    unsigned long comp_size =
        this->key_block_info_list[idx]->key_block_comp_size;
    unsigned long decomp_size =
        this->key_block_info_list[idx]->key_block_decomp_size;
    unsigned long start_ofset = i;
    // unsigned long end_ofset = i + comp_size;
    // 4 bytes comp type
    char* key_block_comp_type = (char*)calloc(4, sizeof(char));
    memcpy(key_block_comp_type, key_block_buffer, 4 * sizeof(char));
    // 4 bytes adler checksum of decompressed key block
    // TODO  adler32 = unpack('>I', key_block_compressed[start + 4:start +
    // 8])[0]
    uint32_t chksum =
        be_bin_to_u32(key_block_buffer + start_ofset + 4 * sizeof(char));

    unsigned char* key_block = nullptr;

    std::vector<uint8_t> kb_uncompressed;  // note: ensure kb_uncompressed not
                                           // die when out of uncompress scope

    if ((key_block_comp_type[0] & 255) == 0) {
      // none compressed
      key_block = key_block_buffer + 8 * sizeof(char);
    } else if ((key_block_comp_type[0] & 255) == 1) {
      // 01000000
      // TODO lzo decompress

    } else if ((key_block_comp_type[0] & 255) == 2) {
      // zlib compress
      kb_uncompressed =
          zlib_mem_uncompress(key_block_buffer + start_ofset + 8, comp_size);
      if (kb_uncompressed.empty() || kb_uncompressed.size() == 0) {
        throw std::runtime_error("key block decompress failed");
      }
      key_block = kb_uncompressed.data();

      uint32_t adler32cs =
          adler32checksum(key_block, static_cast<uint32_t>(decomp_size));
      assert(adler32cs == chksum);
      assert(kb_uncompressed.size() == decomp_size);
    } else {
      throw std::runtime_error("cannot determine the key block compress type");
    }

    // split key
    std::vector<key_list_item*> tlist =
        split_key_block(key_block, decomp_size, idx);
    key_list.insert(key_list.end(), tlist.begin(), tlist.end());

    // TODO HERE append keys

    // next round
    i += comp_size;
  }
  assert(key_list.size() == this->entries_num);
  /// passed

  this->record_block_info_offset = this->key_block_info_start_offset +
                                   this->key_block_info_size +
                                   this->key_block_size;
  /// passed

  return 0;
}

// note: kb_info_buff_len == key_block_info_compressed_size

/**
 * decode the record block
 * @param record_block_buffer
 * @param rb_len record block buffer length
 * @return
 */
int Mdict::read_record_block_header() {
  /**
   * record block info section
   * decode the record block info section
   * [0:8/4]    - record blcok number
   * [8:16/4:8] - num entries the key-value entries number
   * [16:24/8:12] - record block info size
   * [24:32/12:16] - record block size
   */
  if (this->version >= 2.0) {
    record_block_info_size = 4 * 8;
  } else {
    record_block_info_size = 4 * 4;
  }

  char* record_info_buffer =
      (char*)calloc(record_block_info_size, sizeof(char));
  this->readfile(record_block_info_offset, record_block_info_size,
                 record_info_buffer);
  if (this->version >= 2.0) {
    record_block_number = be_bin_to_u64((unsigned char*)record_info_buffer);
    record_block_entries_number = be_bin_to_u64(
        (unsigned char*)record_info_buffer + number_width * sizeof(char));
    record_block_header_size = be_bin_to_u64(
        (unsigned char*)record_info_buffer + 2 * number_width * sizeof(char));
    record_block_size = be_bin_to_u64((unsigned char*)record_info_buffer +
                                      3 * number_width * sizeof(char));

  } else {
  }
  free(record_info_buffer);
  assert(record_block_entries_number == entries_num);
  /// passed

  /**
   * record_block_header_list:
   * {
   *     compressed size
   *     decompressed size
   * }
   */
  char* record_header_buffer =
      (char*)calloc(record_block_header_size, sizeof(char));
  this->readfile(this->record_block_info_offset + record_block_info_size,
                 record_block_header_size, record_header_buffer);

  unsigned long comp_size = 0l;
  unsigned long uncomp_size = 0l;
  unsigned long size_counter = 0l;

  unsigned long comp_accu = 0l;
  unsigned long decomp_accu = 0l;

  for (unsigned long i = 0; i < record_block_number; ++i) {
    if (this->version >= 2.0) {
      comp_size =
          be_bin_to_u64((unsigned char*)(record_header_buffer + size_counter));
      size_counter += number_width;
      uncomp_size =
          be_bin_to_u64((unsigned char*)(record_header_buffer + size_counter));
      size_counter += number_width;

      this->record_header.push_back(new record_header_item(
          i, comp_size, uncomp_size, comp_accu, decomp_accu));
      // ensure after push
      comp_accu += comp_size;
      decomp_accu += uncomp_size;
    } else {
      // TODO
    }
  }

  free(record_header_buffer);
  assert(this->record_header.size() == this->record_block_number);
  assert(size_counter == this->record_block_header_size);

  record_block_offset = record_block_info_offset + record_block_info_size +
                        record_block_header_size;
  /// passed
  return 0;
}

std::vector<std::pair<std::string, std::string>>
Mdict::decode_record_block_by_rid(unsigned long rid /* record id */) {
  // record block start offset: record_block_offset
  uint64_t record_offset = this->record_block_offset;

  // key list index counter
  unsigned long i = 0l;

  std::vector<uint8_t> record_block_uncompressed_v;
  unsigned char* record_block_uncompressed_b;
  uint64_t checksum = 0l;

  unsigned long idx = rid;

  //  for (int idx = 0; idx < this->record_header.size(); idx++) {
  uint64_t comp_size = record_header[idx]->compressed_size;
  uint64_t uncomp_size = record_header[idx]->decompressed_size;
  uint64_t comp_accu = record_header[idx]->compressed_size_accumulator;
  uint64_t decomp_accu = record_header[idx]->decompressed_size_accumulator;
  uint64_t previous_end = 0;
  uint64_t previous_uncomp_size = 0;
  if (idx > 0) {
      previous_end = record_header[idx-1]->decompressed_size_accumulator;
      previous_uncomp_size = record_header[idx-1]->decompressed_size;
  }
  //  std::cout << "record decomp accu " << decomp_accu << std::endl;

  char* record_block_cmp_buffer = (char*)calloc(comp_size, sizeof(char));

  this->readfile(record_offset + comp_accu, comp_size, record_block_cmp_buffer);
  // 4 bytes, compress type
  char* comp_type_b = (char*)calloc(4, sizeof(char));
  memcpy(comp_type_b, record_block_cmp_buffer, 4 * sizeof(char));
  //    putbytes(comp_type_b, 4, true);
  int comp_type = comp_type_b[0] & 0xff;
  // 4 bytes adler32 checksum
  char* checksum_b = (char*)calloc(4, sizeof(char));
  memcpy(checksum_b, record_block_cmp_buffer + 4, 4 * sizeof(char));
  checksum = be_bin_to_u32((unsigned char*)checksum_b);
  free(checksum_b);

  if (comp_type == 0 /* not compressed TODO*/) {
    throw std::runtime_error("uncompress block not support yet");
  } else {
    char* record_block_decrypted_buff;
    if (this->encrypt == ENCRYPT_RECORD_ENC /* record block encrypted */) {
      // TODO
      throw std::runtime_error("record encrypted not support yet");
    }
    record_block_decrypted_buff = record_block_cmp_buffer + 8 * sizeof(char);
    // decompress
    if (comp_type == 1 /* lzo */) {
      throw std::runtime_error("lzo compress not support yet");
    } else if (comp_type == 2) {
      // zlib compress
      record_block_uncompressed_v =
          zlib_mem_uncompress(record_block_decrypted_buff, comp_size);
      if (record_block_uncompressed_v.empty()) {
        throw std::runtime_error("record block decompress failed size == 0");
      }
      record_block_uncompressed_b = record_block_uncompressed_v.data();
      uint32_t adler32cs = adler32checksum(record_block_uncompressed_b,
                                           static_cast<uint32_t>(uncomp_size));
      assert(record_block_uncompressed_v.size() == uncomp_size);
      assert(adler32cs == checksum);
    } else {
      throw std::runtime_error(
          "cannot determine the record block compress type");
    }
  }

  free(comp_type_b);
  free(record_block_cmp_buffer);
  //    free(record_block_uncompressed_b); /* ensure not free twice*/

  unsigned char* record_block = record_block_uncompressed_b;
  /**
   * 请注意，block 是会有很多个的，而每个block都可能会被压缩
   * 而 key_list中的 record_start,
   * key_text是相对每一个block而言的，end是需要每次解析的时候算出来的
   * 所有的record_start/length/end都是针对解压后的block而言的
   */

  std::vector<std::pair<std::string, std::string>> vec;

  while (i < this->key_list.size()) {
    // TODO OPTIMISE
    unsigned long record_start = key_list[i]->record_start;

    std::string key_text = key_list[i]->key_word;
    // start, skip the keys which not includes in record block
    if (record_start < decomp_accu) {
      i++;
      continue;
    }

    // end important: the condition should be lgt, because, the end bound will
    // be equal to uncompressed size
    // this part ensures the record match to key list bound
    if (record_start - decomp_accu >= uncomp_size) {
      break;
    }

    //    std::cout << "key text: " << key_text << std::endl;
    //    std::cout << "idx: " << idx << std::endl;
    unsigned long upbound = uncomp_size; // - this->key_list[i]->record_start;
    unsigned long expect_end = 0;
    auto expect_start = this->key_list[i]->record_start - decomp_accu;
    if (i < this->key_list.size() - 1) {
      expect_end = this->key_list[i + 1]->record_start - this->key_list[i]->record_start;
      expect_start = this->key_list[i]->record_start - decomp_accu;
    } else {
      // 前一个的 end + size 等于当前这个的开始
      expect_end = this->record_block_size - (previous_end + previous_uncomp_size);
    }
    upbound = expect_end < upbound ? expect_end : upbound;

    std::string def;
      if (this->filetype == "MDD") {
        def = be_bin_to_utf16(
                (char*)record_block, expect_start,
                upbound /* to delete null character*/);
    } else {
        def = be_bin_to_utf8(
                (char*)record_block, expect_start,
                upbound - 1 /* to delete null character*/);
    }
    //    std::cout << "def: " << def << std::endl;
    std::pair<std::string, std::string> vp(key_text, def);
    vec.push_back(vp);
    i++;
  }

  //  assert(size_counter == record_block_size);
  return vec;
}

int Mdict::decode_record_block() {
  // record block start offset: record_block_offset
  uint64_t record_offset = this->record_block_offset;

  uint64_t item_counter = 0;
  uint64_t size_counter = 0l;

  // key list index counter
  unsigned long i = 0l;

  // record offset
  unsigned long offset = 0l;

  std::vector<uint8_t> record_block_uncompressed_v;
  unsigned char* record_block_uncompressed_b;
  uint64_t checksum = 0l;
  for (int idx = 0; idx < this->record_header.size(); idx++) {
    uint64_t comp_size = record_header[idx]->compressed_size;
    uint64_t uncomp_size = record_header[idx]->decompressed_size;
    char* record_block_cmp_buffer = (char*)calloc(comp_size, sizeof(char));
    this->readfile(record_offset, comp_size, record_block_cmp_buffer);
    //    putbytes(record_block_cmp_buffer, 8, true);
    // 4 bytes, compress type
    char* comp_type_b = (char*)calloc(4, sizeof(char));
    memcpy(comp_type_b, record_block_cmp_buffer, 4 * sizeof(char));
    //    putbytes(comp_type_b, 4, true);
    int comp_type = comp_type_b[0] & 0xff;
    // 4 bytes adler32 checksum
    char* checksum_b = (char*)calloc(4, sizeof(char));
    memcpy(checksum_b, record_block_cmp_buffer + 4, 4 * sizeof(char));
    checksum = be_bin_to_u32((unsigned char*)checksum_b);
    free(checksum_b);

    if (comp_type == 0 /* not compressed TODO*/) {
      throw std::runtime_error("uncompress block not support yet");
    } else {
      char* record_block_decrypted_buff;
      if (this->encrypt == ENCRYPT_RECORD_ENC /* record block encrypted */) {
        // TODO
        throw std::runtime_error("record encrypted not support yet");
      }
      record_block_decrypted_buff = record_block_cmp_buffer + 8 * sizeof(char);
      // decompress
      if (comp_type == 1 /* lzo */) {
        throw std::runtime_error("lzo compress not support yet");
      } else if (comp_type == 2) {
        // zlib compress
        record_block_uncompressed_v =
            zlib_mem_uncompress(record_block_decrypted_buff, comp_size);
        if (record_block_uncompressed_v.empty()) {
          throw std::runtime_error("record block decompress failed size == 0");
        }
        record_block_uncompressed_b = record_block_uncompressed_v.data();
        uint32_t adler32cs = adler32checksum(
            record_block_uncompressed_b, static_cast<uint32_t>(uncomp_size));
        assert(adler32cs == checksum);
        assert(record_block_uncompressed_v.size() == uncomp_size);
      } else {
        throw std::runtime_error(
            "cannot determine the record block compress type");
      }
    }

    free(comp_type_b);
    free(record_block_cmp_buffer);
    //    free(record_block_uncompressed_b); /* ensure not free twice*/

    // unsigned char* record_block = record_block_uncompressed_b;
    /**
     * 请注意，block 是会有很多个的，而每个block都可能会被压缩
     * 而 key_list中的 record_start,
     * key_text是相对每一个block而言的，end是需要每次解析的时候算出来的
     * 所有的record_start/length/end都是针对解压后的block而言的
     */
    while (i < this->key_list.size()) {
      unsigned long record_start = key_list[i]->record_start;
      std::string key_text = key_list[i]->key_word;
      if (record_start - offset >= uncomp_size) {
        // overflow
        break;
      }
      unsigned long record_end;
      if (i < this->key_list.size() - 1) {
        record_end = this->key_list[i + 1]->record_start;
      } else {
        record_end = uncomp_size + offset;
      }

      this->key_data.push_back(new record(
          key_text, key_list[i]->record_start, this->encoding, record_offset,
          comp_size, uncomp_size, comp_type, (this->encrypt == 1),
          record_start - offset, record_end - offset));
      i++;
      item_counter++;
    }
    // offset += record_block.length
    offset += uncomp_size;
    size_counter += comp_size;
    record_offset += comp_size;

    //    break;
  }
  assert(size_counter == record_block_size);
  return 0;
}

/**
 * decode the key block info
 * @param key_block_info_buffer the key block info buffer
 * @param kb_info_buff_len the key block buffer length
 * @param key_block_num the key block number
 * @param entries_num the entries number
 * @return
 */
int Mdict::decode_key_block_info(char* key_block_info_buffer,
                                 unsigned long kb_info_buff_len,
                                 int key_block_num, int entries_num) {
  char* kb_info_buff = key_block_info_buffer;

  // key block info offset indicator
  unsigned long data_offset = 0;

  if (this->version >= 2.0) {
    // if version >= 2.0, use zlib compression
    assert(kb_info_buff[0] == 2);
    assert(kb_info_buff[1] == 0);
    assert(kb_info_buff[2] == 0);
    assert(kb_info_buff[3] == 0);
    byte* kb_info_decrypted = (unsigned char*)key_block_info_buffer;
    if (this->encrypt == ENCRYPT_KEY_INFO_ENC) {
      kb_info_decrypted = mdx_decrypt((byte*)kb_info_buff, kb_info_buff_len);
    }

    // finally, we needs to check adler32 checksum
    // key_block_info_compressed[4:8] => adler32 checksum
    //          uint32_t chksum = be_bin_to_u32((unsigned char*) (kb_info_buff +
    //          4));
    //          uint32_t adlercs = adler32checksum(key_block_info_uncomp,
    //          static_cast<uint32_t>(key_block_info_uncomp_len)) & 0xffffffff;
    //
    //          assert(chksum == adlercs);

    /// here passed, key block info is corrected
    // TODO decode key block info compressed into keys list

    // for version 2.0, will compress by zlib, lzo just just for 1.0
    // key_block_info_buff[0:8] => compress_type
    // TODO zlib decompress
    // TODO:
    // if the size of compressed data original data is unknown,
    // we malloc 8 size of source data len, we cannot estimate the original data
    // size
    // but currently, we know the size of key_block_info decompress size, so we
    // use this

    // note: we should uncompress key_block_info_buffer[8:] data, so we need
    // (decrypted + 8, and length -8)
    std::vector<uint8_t> decompress_buff =
        zlib_mem_uncompress(kb_info_decrypted + 8, kb_info_buff_len - 8,
                            this->key_block_info_decompress_size);
    /// uncompress successed
    assert(decompress_buff.size() == this->key_block_info_decompress_size);

    // get key block info list
    //          std::vector<key_block_info*> key_block_info_list;
    /// entries summary, every block has a lot of entries, the sum of entries
    /// should equals entries_number
    unsigned long num_entries_counter = 0;
    // key number counter
    unsigned long counter = 0;

    // current block entries
    unsigned long current_entries = 0;

    unsigned long previous_start_offset = 0;

    int byte_width = 1;
    int text_term = 0;
    if (this->version >= 2.0) {
      byte_width = 2;
      text_term = 1;
    }

    unsigned long comp_acc = 0l;
    unsigned long decomp_acc = 0l;
    while (counter < this->key_block_num) {
      if (this->version >= 2.0) {
        current_entries = be_bin_to_u64(decompress_buff.data() +
                                        data_offset * sizeof(uint8_t));
      } else {
        current_entries = be_bin_to_u32(decompress_buff.data() +
                                        data_offset * sizeof(uint8_t));
      }
      num_entries_counter += current_entries;

      // move offset
      // if version>= 2.0 move forward 8 bytes

      data_offset += this->number_width * sizeof(uint8_t);

      // first key size
      unsigned long first_key_size = 0;

      if (this->version >= 2.0) {
        first_key_size = be_bin_to_u16(decompress_buff.data() +
                                       data_offset * sizeof(uint8_t));
      } else {
        first_key_size = be_bin_to_u8(decompress_buff.data() +
                                      data_offset * sizeof(uint8_t));
      }
      data_offset += byte_width;

      // step_gap means first key start offset to first key end;
      int step_gap = 0;

      if (this->encoding == 1 /* encoding utf16 equals 1*/) {
        step_gap = (first_key_size + text_term) * 2;
      } else {
        step_gap = first_key_size + text_term;
      }

      // DECODE first CODE
      // TODO here minus the terminal character size(1), but we still not sure
      // should minus this or not
      std::string fkey =
          be_bin_to_utf8((char*)(decompress_buff.data() + data_offset), 0,
                         (unsigned long)step_gap - text_term);
      //            std::cout<<"first key: "<<fkey<<std::endl;
      // move forward
      data_offset += step_gap;

      // the last key
      unsigned long last_key_size = 0;

      if (this->version >= 2.0) {
        last_key_size = be_bin_to_u16(decompress_buff.data() +
                                      data_offset * sizeof(uint8_t));
      } else {
        last_key_size = be_bin_to_u8(decompress_buff.data() +
                                     data_offset * sizeof(uint8_t));
      }
      data_offset += byte_width;

      if (this->encoding == 1 /* ENCODING_UTF16 */) {
        step_gap = (last_key_size + text_term) * 2;
      } else {
        step_gap = last_key_size + text_term;
      }

      std::string last_key =
          be_bin_to_utf8((char*)(decompress_buff.data() + data_offset), 0,
                         (unsigned long)step_gap - text_term);

      // move forward
      data_offset += step_gap;

      // ------------
      // key block part
      // ------------

      uint64_t key_block_compress_size = 0;
      if (version >= 2.0) {
        key_block_compress_size =
            be_bin_to_u64(decompress_buff.data() + data_offset);
      } else {
        key_block_compress_size =
            be_bin_to_u32(decompress_buff.data() + data_offset);
      }

      data_offset += this->number_width;

      uint64_t key_block_decompress_size = 0;

      if (version >= 2.0) {
        key_block_decompress_size =
            be_bin_to_u64(decompress_buff.data() + data_offset);
      } else {
        key_block_decompress_size =
            be_bin_to_u32(decompress_buff.data() + data_offset);
      }

      // entries offset move forward
      data_offset += this->number_width;

      key_block_info* kbinfo = new key_block_info(
          fkey, last_key, previous_start_offset, key_block_compress_size,
          key_block_decompress_size, comp_acc, decomp_acc);

      // adjust ofset
      previous_start_offset += key_block_compress_size;
      key_block_info_list.push_back(kbinfo);

      // key block counter
      counter += 1;
      // accumulate
      comp_acc += key_block_compress_size;
      decomp_acc += key_block_decompress_size;

      //          break;
    }
    assert(counter == this->key_block_num);
    assert(num_entries_counter == this->entries_num);

    //    std::vector<key_block_info*>::iterator it;

    //     TODO WORKING HERE
    //    for (auto it = key_block_info_list.begin(); it !=
    //    key_block_info_list.end();
    //         it++) {
    //      std::cout << "fkey : " << (*it)->first_key << std::endl;
    //      std::cout << "lkey : " << (*it)->last_key << std::endl;
    //      std::cout << "comp_size : " << (*it)->key_block_comp_size <<
    //      std::endl;
    //      std::cout << "decomp_size : " << (*it)->key_block_decomp_size
    //                << std::endl;
    //      std::cout << "offset : " << (*it)->key_block_start_offset <<
    //      std::endl;
    //      break;
    //    }

  } else {
    // doesn't compression
    throw std::logic_error("not implements yet");
  }

  //        std::cout<<"data offset: " << data_offset<<std::endl;
  //        assert(data_offset == this->key_block_info_decompress_size);
  this->key_block_body_start =
      this->key_block_info_start_offset + this->key_block_info_size;
  //        std::cout<<"key_block_body offset: " <<
  //        this->key_block_body_start<<std::endl;
  /// here passed
  return 0;
}

/**
 * read in the file from the file stream
 * @param offset the file start offset
 * @param len the byte length needs to read
 * @param buf the target buffer
 */
void Mdict::readfile(uint64_t offset, uint64_t len, char* buf) {
  instream.seekg(offset);
  instream.read(buf, static_cast<std::streamsize>(len));
}

/***************************************
 *             public part             *
 ***************************************/

/**
 * init the dictionary file
 */
void Mdict::init() {
  /* indexing... */
  this->read_header();
  this->printhead();
  this->read_key_block_header();
  this->read_key_block_info();
  this->read_record_block_header();
  //  this->decode_record_block();

  // TODO delete this  this->decode_record_block(); // very slow!!!
}

/**
 * find the key word includes in which block
 * @param phrase
 * @param start
 * @param end
 * @return
 */
long Mdict::reduce0(std::string phrase, unsigned long start,
                    unsigned long end) {  // non-recursive reduce implements
  for (int i = 0; i < end; ++i) {
    std::string first_key = this->key_block_info_list[i]->first_key;
    std::string last_key = this->key_block_info_list[i]->last_key;
    // std::cout << "index : " << i << ", first_key : " << first_key <<
    // ", last_key : " << last_key << std::endl;
    if (phrase.compare(first_key) >= 0 && phrase.compare(last_key) <= 0) {
      //            std::cout << ">>>>>>>>>>>> found index " << i << std::endl;
      return i;
    }
  }
  return -1;
}

long Mdict::reduce1(std::vector<key_list_item*> wordlist,
                    std::string phrase) {  // non-recursive reduce implements
  unsigned long left = 0;
  unsigned long right = wordlist.size() - 1;
  unsigned long mid = 0;
  std::string word = _s(std::move(phrase));

  int comp = 0;
  while (left <= right) {
    mid = left + ((right - left) >> 1);
    // std::cout << "reduce1, mid = " << mid << ", left: " << left << ", right :
    // " <<  right << ", size: " << wordlist.size() << std::endl;
    if (mid >= wordlist.size()) {
      return -1;
    }
    comp = word.compare(_s(wordlist[mid]->key_word));
    if (comp == 0) {
      return mid;
    } else if (comp > 0) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }
  return -1;
}

/**
 *
 * @param wordlist
 * @param phrase
 * @return
 */
long Mdict::reduce2(
    unsigned long record_start) {  // non-recursive reduce implements
  // TODO OPTIMISE
  unsigned long left = 0l;
  unsigned long right = this->record_header.size() - 1;
  unsigned long mid = 0;
  while (left <= right) {
    mid = left + ((right - left) >> 1);
    if (record_start >=
        this->record_header[mid]->decompressed_size_accumulator) {
      left = mid + 1;
    } else if (record_start <
               this->record_header[mid]->decompressed_size_accumulator) {
      right = mid - 1;
    }
  }
  return left - 1;
  // TODO test this
  //  for (unsigned long i = 0; i < record_header.size(); i++) {
  //    if (this->record_header[i]->decompressed_size_accumulator >=
  //    record_start) {
  //      return i == 0 ? 0 : i - 1;
  //    }
  //  }
  return 0;
}

std::string Mdict::reduce3(std::vector<std::pair<std::string, std::string>> vec,
                           std::string phrase) {
  unsigned int left = 0;
  unsigned int right = vec.size() - 1;
  unsigned int mid = 0;
  unsigned int result = 0;
  while (left < right) {
    mid = left + ((right - left) >> 1);
    // std::cout << _s(vec[mid].first) << std::endl;
    if (_s(phrase).compare(_s(vec[mid].first)) > 0) {
      left = mid + 1;
    } else if (_s(phrase).compare(_s(vec[mid].first)) == 0) {
      left = mid;
      break;
    } else {
      right = mid - 1;
    }
  }
  result = left;
  return vec[result].second;
}

/**
 * look the file by word
 * @param word the searching word
 * @return
 */
std::string Mdict::lookup(const std::string word) {
  try {
    // search word in key block info list
    long idx = this->reduce0(_s(word), 0, this->key_block_info_list.size());
    //            std::cout << "==> lookup idx " << idx << std::endl;
    if (idx >= 0) {
      // decode key block by block id
      std::vector<key_list_item*> tlist =
          this->decode_key_block_by_block_id(idx);
      // reduce word id from key list item vector to get the word index of key
      // list
      long word_id = reduce1(tlist, word);
      if (word_id >= 0) {
        // reduce search the record block index by word record start offset
        unsigned long record_block_idx = reduce2(tlist[word_id]->record_start);
        // decode recode by record index
        auto vec = decode_record_block_by_rid(record_block_idx);
        //  for(auto it= vec.begin(); it != vec.end(); ++it){
        //   std::cout<<"word: "<<(*it).first<<" \n def:
        //   "<<(*it).second<<std::endl;
        //  }
        // reduce the definition by word
        std::string def = reduce3(vec, word);
        return def;
      }
    }
  } catch (std::exception& e) {
    std::cout << "==> lookup error" << e.what() << std::endl;
  }
  return std::string();
}

std::string Mdict::parse_definition(const std::string word, unsigned long record_start) {
    // reduce search the record block index by word record start offset
    unsigned long record_block_idx = reduce2(record_start);
    // decode recode by record index
    auto vec = decode_record_block_by_rid(record_block_idx);
    // reduce the definition by word
    std::string def = reduce3(vec, word);
    return def;

}


/**
 * look the file by word
 * @param word the searching word
 * @return
 */
std::vector<key_list_item *>  Mdict::keyList() {
  return this->key_list;
}

bool Mdict::endsWith(std::string const &fullString, std::string const &ending) {
        if (fullString.length() >= ending.length()) {
            return (0 == fullString.compare (fullString.length() - ending.length(), ending.length(), ending));
        } else {
            return false;
        }
}
}  // namespace mdict
