#include "adler32.h"

/// ADLER-32 checksum calculations
class Adler32 {
 public:
  static const int DIGESTSIZE = 4;

  static uint32_t be_bin_to_u32(
      const unsigned char *bin /* 4 bytes char array  */) {
    uint32_t n = 0;
    for (int i = 0; i < 3; i++) {
      n = n | (unsigned int)bin[i];
      n = n << 8;
    }
    n = n | (unsigned int)bin[3];
    return n;
  }

  Adler32() { Reset(); }

  void Update(const byte *input, size_t length);

  void TruncatedFinal(byte *hash, size_t size);

  unsigned int DigestSize() const { return DIGESTSIZE; }

  /// \brief Computes the hash of the current message
  /// \param digest a pointer to the buffer to receive the hash
  /// \details Final() restarts the hash for a new message.
  /// \pre <tt>COUNTOF(digest) == DigestSize()</tt> or <tt>COUNTOF(digest) ==
  /// HASH::DIGESTSIZE</tt> ensures
  ///   the output byte buffer is large enough for the digest.
  void Final(byte *digest) { TruncatedFinal(digest, DigestSize()); }

 private:
  void Reset() {
    m_s1 = 1;
    m_s2 = 0;
  }

  word16 m_s1, m_s2;
};

void Adler32::Update(const byte *input, size_t length) {
  const unsigned long BASE = 65521;

  unsigned long s1 = m_s1;
  unsigned long s2 = m_s2;

  if (length % 8 != 0) {
    do {
      s1 += *input++;
      s2 += s1;
      length--;
    } while (length % 8 != 0);

    if (s1 >= BASE) s1 -= BASE;
    s2 %= BASE;
  }

  while (length > 0) {
    s1 += input[0];
    s2 += s1;
    s1 += input[1];
    s2 += s1;
    s1 += input[2];
    s2 += s1;
    s1 += input[3];
    s2 += s1;
    s1 += input[4];
    s2 += s1;
    s1 += input[5];
    s2 += s1;
    s1 += input[6];
    s2 += s1;
    s1 += input[7];
    s2 += s1;

    length -= 8;
    input += 8;

    if (s1 >= BASE) s1 -= BASE;
    if (length % 0x8000 == 0) s2 %= BASE;
  }

  assert(s1 < BASE);
  assert(s2 < BASE);

  m_s1 = (word16)s1;
  m_s2 = (word16)s2;
}

void Adler32::TruncatedFinal(byte *hash, size_t size) {
  // ThrowIfInvalidTruncatedSize(size);
  if (size != DIGESTSIZE) {
    // TODO enhance exception logic
    throw std::exception();
  }

  switch (size) {
    default:
      hash[3] = byte(m_s1);
    // fall through
    case 3:
      hash[2] = byte(m_s1 >> 8);
    // fall through
    case 2:
      hash[1] = byte(m_s2);
    // fall through
    case 1:
      hash[0] = byte(m_s2 >> 8);
    // fall through
    case 0:;
      ;
      // fall through
  }

  Reset();
}

uint32_t adler32checksum(const unsigned char *data, uint32_t len) {
  Adler32 adler32Hasher;
  adler32Hasher.Update(data, len);
  char *hash = (char *)calloc(4, sizeof(char));
  adler32Hasher.Final(reinterpret_cast<byte *>(hash));
  uint32_t chksum = Adler32::be_bin_to_u32((unsigned char *)hash);
  if (hash) free(hash);
  return chksum;
}
