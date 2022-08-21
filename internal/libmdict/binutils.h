#ifndef MDICT_BINUTILS_H_
#define MDICT_BINUTILS_H_

#include <string>

/*
 *
 * convert BigEndian 4 bytes unsigned char array to uint32 integer
 *
 **/

uint8_t be_bin_to_u8(const unsigned char* bin /* 8 bytes char array*/);

uint16_t be_bin_to_u16(const unsigned char* bin /* 4 bytes char array  */);

uint32_t be_bin_to_u32(const unsigned char* bin /* 4 bytes char array  */);

uint64_t be_bin_to_u64(const unsigned char* bin /* 8 bytes char array*/);

int bin_slice(const char* srcByte, int srcByteLen, int offset, int len,
              char* distByte);

void putbytes(const char* bytes, int len, bool hex,
              unsigned long start_ofset = 0);

// ================================================
// code convert part
// ===============================================
std::string le_bin_utf16_to_utf8(const char* bytes, int offset, int len);

std::string be_bin_to_utf8(const char* bytes, unsigned long offset,
                           unsigned long len);

std::string be_bin_to_utf16(const char* bytes, unsigned long offset,
                            unsigned long len);
int bintohex(const char* bin, unsigned long len, char* target);

#endif
