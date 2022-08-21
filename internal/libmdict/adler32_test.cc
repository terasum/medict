// adler32.cpp - originally written and placed in the public domain by Wei Dai

#include "adler32.h"
int main() {
  //  Adler32 adler32hasher;
  char *str = const_cast<char *>("helloworld");
  uint32_t hash = adler32checksum((unsigned char *)str, 10);
  if (hash != 0x1736043D) {
    printf("adler32 test failed, PANIC\n");
    return -1;
  }
  return 0;
}
