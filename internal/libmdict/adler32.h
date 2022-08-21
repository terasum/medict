#ifndef MDICT_ADLER32_H_
#define MDICT_ADLER32_H_

#include <assert.h>
#include <stddef.h>
#include <stdint.h>

#include <cstdio>
#include <cstdlib>
#include <exception>

typedef unsigned char byte;
typedef unsigned short word16;

uint32_t adler32checksum(const unsigned char *data, uint32_t len);

#endif
