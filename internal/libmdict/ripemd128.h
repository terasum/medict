/********************************************************************\
 *
 *      FILE:     rupemd128.h
 *
 *      CONTENTS: Header file for a sample C-implementation of the
 *                RIPEMD-128 hash-function. This function is a
 *                plug-in substitute for RIPEMD. A 128-bit hash
 *                result is obtained using RIPEMD-128.
 *      TARGET:   any computer with an ANSI C compiler
 *
 *      AUTHOR:   Chen Quan
 *      DATE:     5 Feb 2018
 *      VERSION:  1.0
 *
 *      Ref: Antoon Bosselaers, ESAT-COSIC version rmd128.h
 *
 *
\********************************************************************/

#ifndef MDICT_RIPEMD128_H_ /* make sure this file is read only once */
#define MDICT_RIPEMD128_H_ /* make sure this file is read only once */

#ifndef RMDsize
#define RMDsize 128
#endif

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

/********************************************************************/

/* typedef 8, 16 and 32 bit types, resp.  */
/* adapt these, if necessary,
   for your operating system and compiler */
typedef unsigned char byte;   /* unsigned 8-bit type */
typedef unsigned short word;  /* unsigned 16-bit type */
typedef unsigned int dword32; /* unsigned 32-bit type */

/********************************************************************/

/* macro definitions */

/* collect four bytes into one word: */
#define BYTES_TO_DWORD(strptr)                                               \
  (((dword32) * ((strptr) + 3) << 24) | ((dword32) * ((strptr) + 2) << 16) | \
   ((dword32) * ((strptr) + 1) << 8) | ((dword32) * (strptr)))

/* ROL(x, n) cyclically rotates x over n bits to the left */
/* x must be of an unsigned 32 bits type and 0 <= n < 32. */
#define ROL(x, n) (((x) << (n)) | ((x) >> (32 - (n))))

#define F(x, y, z) (x ^ y ^ z)
#define G(x, y, z) (z ^ (x & (y ^ z)))
#define H(x, y, z) (z ^ (x | ~y))
#define I(x, y, z) (y ^ (z & (x ^ y)))
#define J(x, y, z) (x ^ (y | ~z))

#define K0 0
#define K1 0x5a827999UL
#define K2 0x6ed9eba1UL
#define K3 0x8f1bbcdcUL
#define K4 0xa953fd4eUL
#define K5 0x50a28be6UL
#define K6 0x5c4dd124UL
#define K7 0x6d703ef3UL
#define K8 0x7a6d76e9UL
#define K9 0

/* the eight basic operations FF() through III() */
#define FF(a, b, c, d, x, s)            \
  {                                     \
    (a) += F((b), (c), (d)) + (x) + K0; \
    (a) = ROL((a), (s));                \
  }
#define GG(a, b, c, d, x, s)            \
  {                                     \
    (a) += G((b), (c), (d)) + (x) + K1; \
    (a) = ROL((a), (s));                \
  }
#define HH(a, b, c, d, x, s)            \
  {                                     \
    (a) += H((b), (c), (d)) + (x) + K2; \
    (a) = ROL((a), (s));                \
  }
#define II(a, b, c, d, x, s)            \
  {                                     \
    (a) += I((b), (c), (d)) + (x) + K3; \
    (a) = ROL((a), (s));                \
  }
#define FFF(a, b, c, d, x, s)           \
  {                                     \
    (a) += F((b), (c), (d)) + (x) + K9; \
    (a) = ROL((a), (s));                \
  }
#define GGG(a, b, c, d, x, s)           \
  {                                     \
    (a) += G((b), (c), (d)) + (x) + K7; \
    (a) = ROL((a), (s));                \
  }
#define HHH(a, b, c, d, x, s)           \
  {                                     \
    (a) += H((b), (c), (d)) + (x) + K6; \
    (a) = ROL((a), (s));                \
  }
#define III(a, b, c, d, x, s)           \
  {                                     \
    (a) += I((b), (c), (d)) + (x) + K5; \
    (a) = ROL((a), (s));                \
  }

/********************************************************************/

/* function prototypes */

/*
 *  initializes MDbuffer to "magic constants"
 */
void ripemd128Init(dword32 *digest);

/*
 *  the compression function.
 *  transforms MDbuf using message bytes X[0] through X[15]
 */
void ripemd128compress(dword32 *digest, dword32 *X);

/*
 * ISO7816 message padding
 * return the length after padding, the data will be modified
 */
int ripemd128PaddingISO7816(uint8_t **data, int data_len);

/**
 * simple ripemd128 for string method
 * @param message
 * @param length
 */

byte *ripemd128bytes(uint8_t *message, int length);

#endif /* RMD128H */

/*********************** end of file rmd128.h ***********************/

#ifdef __cplusplus
}
#endif
