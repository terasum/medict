/********************************************************************\
 *
 *      FILE:     rmd128.c
 *
 *      CONTENTS: A sample C-implementation of the RIPEMD-128
 *                hash-function. This function is a plug-in substitute 
 *                for RIPEMD. A 160-bit hash result is obtained using
 *                RIPEMD-160.
 *      TARGET:   any computer with an ANSI C compiler
 *
 *      AUTHOR:   Antoon Bosselaers, ESAT-COSIC
 *      DATE:     1 March 1996
 *      VERSION:  1.0
 *
 *      Copyright (c) Katholieke Universiteit Leuven
 *      1996, All Rights Reserved
 *
\********************************************************************/

/*  header files */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "ripemd128.h"

#ifdef __cplusplus
extern "C" {
#endif

/********************************************************************/

void ripemd128Init(dword32 *digest)
{
  digest[0] = 0x67452301UL;
  digest[1] = 0xefcdab89UL;
  digest[2] = 0x98badcfeUL;
  digest[3] = 0x10325476UL;

  return;
}


/********************************************************************/

void ripemd128compress(dword32 *digest, dword32 *X)
{
  dword32 aa,aaa,bb,bbb,cc,ccc,dd,ddd;
  aa = aaa = digest[0];
  bb = bbb = digest[1];
  cc = ccc = digest[2];
  dd = ddd = digest[3];

  /* round 1 */
  FF(aa, bb, cc, dd, X[ 0], 11);
  FF(dd, aa, bb, cc, X[ 1], 14);
  FF(cc, dd, aa, bb, X[ 2], 15);
  FF(bb, cc, dd, aa, X[ 3], 12);
  FF(aa, bb, cc, dd, X[ 4],  5);
  FF(dd, aa, bb, cc, X[ 5],  8);
  FF(cc, dd, aa, bb, X[ 6],  7);
  FF(bb, cc, dd, aa, X[ 7],  9);
  FF(aa, bb, cc, dd, X[ 8], 11);
  FF(dd, aa, bb, cc, X[ 9], 13);
  FF(cc, dd, aa, bb, X[10], 14);
  FF(bb, cc, dd, aa, X[11], 15);
  FF(aa, bb, cc, dd, X[12],  6);
  FF(dd, aa, bb, cc, X[13],  7);
  FF(cc, dd, aa, bb, X[14],  9);
  FF(bb, cc, dd, aa, X[15],  8);

  /* round 2 */
  GG(aa, bb, cc, dd, X[ 7],  7);
  GG(dd, aa, bb, cc, X[ 4],  6);
  GG(cc, dd, aa, bb, X[13],  8);
  GG(bb, cc, dd, aa, X[ 1], 13);
  GG(aa, bb, cc, dd, X[10], 11);
  GG(dd, aa, bb, cc, X[ 6],  9);
  GG(cc, dd, aa, bb, X[15],  7);
  GG(bb, cc, dd, aa, X[ 3], 15);
  GG(aa, bb, cc, dd, X[12],  7);
  GG(dd, aa, bb, cc, X[ 0], 12);
  GG(cc, dd, aa, bb, X[ 9], 15);
  GG(bb, cc, dd, aa, X[ 5],  9);
  GG(aa, bb, cc, dd, X[ 2], 11);
  GG(dd, aa, bb, cc, X[14],  7);
  GG(cc, dd, aa, bb, X[11], 13);
  GG(bb, cc, dd, aa, X[ 8], 12);

  /* round 3 */
  HH(aa, bb, cc, dd, X[ 3], 11);
  HH(dd, aa, bb, cc, X[10], 13);
  HH(cc, dd, aa, bb, X[14],  6);
  HH(bb, cc, dd, aa, X[ 4],  7);
  HH(aa, bb, cc, dd, X[ 9], 14);
  HH(dd, aa, bb, cc, X[15],  9);
  HH(cc, dd, aa, bb, X[ 8], 13);
  HH(bb, cc, dd, aa, X[ 1], 15);
  HH(aa, bb, cc, dd, X[ 2], 14);
  HH(dd, aa, bb, cc, X[ 7],  8);
  HH(cc, dd, aa, bb, X[ 0], 13);
  HH(bb, cc, dd, aa, X[ 6],  6);
  HH(aa, bb, cc, dd, X[13],  5);
  HH(dd, aa, bb, cc, X[11], 12);
  HH(cc, dd, aa, bb, X[ 5],  7);
  HH(bb, cc, dd, aa, X[12],  5);

  /* round 4 */
  II(aa, bb, cc, dd, X[ 1], 11);
  II(dd, aa, bb, cc, X[ 9], 12);
  II(cc, dd, aa, bb, X[11], 14);
  II(bb, cc, dd, aa, X[10], 15);
  II(aa, bb, cc, dd, X[ 0], 14);
  II(dd, aa, bb, cc, X[ 8], 15);
  II(cc, dd, aa, bb, X[12],  9);
  II(bb, cc, dd, aa, X[ 4],  8);
  II(aa, bb, cc, dd, X[13],  9);
  II(dd, aa, bb, cc, X[ 3], 14);
  II(cc, dd, aa, bb, X[ 7],  5);
  II(bb, cc, dd, aa, X[15],  6);
  II(aa, bb, cc, dd, X[14],  8);
  II(dd, aa, bb, cc, X[ 5],  6);
  II(cc, dd, aa, bb, X[ 6],  5);
  II(bb, cc, dd, aa, X[ 2], 12);

  /* parallel round 1 */
  III(aaa, bbb, ccc, ddd, X[ 5],  8);
  III(ddd, aaa, bbb, ccc, X[14],  9);
  III(ccc, ddd, aaa, bbb, X[ 7],  9);
  III(bbb, ccc, ddd, aaa, X[ 0], 11);
  III(aaa, bbb, ccc, ddd, X[ 9], 13);
  III(ddd, aaa, bbb, ccc, X[ 2], 15);
  III(ccc, ddd, aaa, bbb, X[11], 15);
  III(bbb, ccc, ddd, aaa, X[ 4],  5);
  III(aaa, bbb, ccc, ddd, X[13],  7);
  III(ddd, aaa, bbb, ccc, X[ 6],  7);
  III(ccc, ddd, aaa, bbb, X[15],  8);
  III(bbb, ccc, ddd, aaa, X[ 8], 11);
  III(aaa, bbb, ccc, ddd, X[ 1], 14);
  III(ddd, aaa, bbb, ccc, X[10], 14);
  III(ccc, ddd, aaa, bbb, X[ 3], 12);
  III(bbb, ccc, ddd, aaa, X[12],  6);

  /* parallel round 2 */
  HHH(aaa, bbb, ccc, ddd, X[ 6],  9);
  HHH(ddd, aaa, bbb, ccc, X[11], 13);
  HHH(ccc, ddd, aaa, bbb, X[ 3], 15);
  HHH(bbb, ccc, ddd, aaa, X[ 7],  7);
  HHH(aaa, bbb, ccc, ddd, X[ 0], 12);
  HHH(ddd, aaa, bbb, ccc, X[13],  8);
  HHH(ccc, ddd, aaa, bbb, X[ 5],  9);
  HHH(bbb, ccc, ddd, aaa, X[10], 11);
  HHH(aaa, bbb, ccc, ddd, X[14],  7);
  HHH(ddd, aaa, bbb, ccc, X[15],  7);
  HHH(ccc, ddd, aaa, bbb, X[ 8], 12);
  HHH(bbb, ccc, ddd, aaa, X[12],  7);
  HHH(aaa, bbb, ccc, ddd, X[ 4],  6);
  HHH(ddd, aaa, bbb, ccc, X[ 9], 15);
  HHH(ccc, ddd, aaa, bbb, X[ 1], 13);
  HHH(bbb, ccc, ddd, aaa, X[ 2], 11);

  /* parallel round 3 */
  GGG(aaa, bbb, ccc, ddd, X[15],  9);
  GGG(ddd, aaa, bbb, ccc, X[ 5],  7);
  GGG(ccc, ddd, aaa, bbb, X[ 1], 15);
  GGG(bbb, ccc, ddd, aaa, X[ 3], 11);
  GGG(aaa, bbb, ccc, ddd, X[ 7],  8);
  GGG(ddd, aaa, bbb, ccc, X[14],  6);
  GGG(ccc, ddd, aaa, bbb, X[ 6],  6);
  GGG(bbb, ccc, ddd, aaa, X[ 9], 14);
  GGG(aaa, bbb, ccc, ddd, X[11], 12);
  GGG(ddd, aaa, bbb, ccc, X[ 8], 13);
  GGG(ccc, ddd, aaa, bbb, X[12],  5);
  GGG(bbb, ccc, ddd, aaa, X[ 2], 14);
  GGG(aaa, bbb, ccc, ddd, X[10], 13);
  GGG(ddd, aaa, bbb, ccc, X[ 0], 13);
  GGG(ccc, ddd, aaa, bbb, X[ 4],  7);
  GGG(bbb, ccc, ddd, aaa, X[13],  5);

  /* parallel round 4 */
  FFF(aaa, bbb, ccc, ddd, X[ 8], 15);
  FFF(ddd, aaa, bbb, ccc, X[ 6],  5);
  FFF(ccc, ddd, aaa, bbb, X[ 4],  8);
  FFF(bbb, ccc, ddd, aaa, X[ 1], 11);
  FFF(aaa, bbb, ccc, ddd, X[ 3], 14);
  FFF(ddd, aaa, bbb, ccc, X[11], 14);
  FFF(ccc, ddd, aaa, bbb, X[15],  6);
  FFF(bbb, ccc, ddd, aaa, X[ 0], 14);
  FFF(aaa, bbb, ccc, ddd, X[ 5],  6);
  FFF(ddd, aaa, bbb, ccc, X[12],  9);
  FFF(ccc, ddd, aaa, bbb, X[ 2], 12);
  FFF(bbb, ccc, ddd, aaa, X[13],  9);
  FFF(aaa, bbb, ccc, ddd, X[ 9], 12);
  FFF(ddd, aaa, bbb, ccc, X[ 7],  5);
  FFF(ccc, ddd, aaa, bbb, X[10], 15);
  FFF(bbb, ccc, ddd, aaa, X[14],  8);

  /* combine results */
  ddd += cc + digest[1];               /* final result for MDbuf[0] */
  digest[1] = digest[2] + dd + aaa;
  digest[2] = digest[3] + aa + bbb;
  digest[3] = digest[0] + bb + ccc;
  digest[0] = ddd;
}


/********************************************************************/

int ripemd128PaddingISO7816(uint8_t **data, int data_len) {
  // padding ISO7816
  // message += 0x80 += 0x00 * length
  int padding_length = ((data_len % 64) < 56 ? 56 : 120) - data_len % 64;
  //printf("paddinglength: %d\n", padding_length);

  uint8_t* padding = (uint8_t *) calloc((size_t) (padding_length), sizeof(uint8_t));
  padding[0] = 0x80; // 0x80

//	printf("padding:");
//	for(int i =0;i<padding_length;i++)
//	   printf("%d ", padding[i]);
//	printf("\n");

  // 1 byte for 0x80
  uint8_t * new_data = (uint8_t *) calloc((size_t) (padding_length + data_len), sizeof(uint8_t));
  // concat data and padding
  memcpy(new_data, *data, data_len);
  memcpy(new_data + data_len * sizeof(uint8_t), padding, padding_length);

  // add length bits
  // bytes_len concat
  uint32_t bytes_len = (uint32_t) data_len; //u32 -> 4 bytes
  bytes_len = bytes_len << 3;
  // malloc x_data memory
  // x_data = data || padding_data(0x80,0x00...0x00<length>) || length_bits
  uint8_t * length_bits = (uint8_t *)calloc(2* sizeof(uint32_t), sizeof(uint8_t)); // 8bytes data
  // length_bits: [length:\x00\x00\x00\x00]
  length_bits[0] = (uint8_t) bytes_len;
  length_bits[1] = (uint8_t) (bytes_len >> 8);
  length_bits[2] = (uint8_t) (bytes_len >> 16);
  length_bits[3] = (uint8_t) (bytes_len >> 24);

  int xdata_length = padding_length +  data_len + 2 * sizeof(uint32_t) /* 2 * 4 * 8*/;
  uint8_t* x_data = (uint8_t *) calloc((size_t) xdata_length, sizeof(uint8_t));
  memcpy(x_data, new_data, data_len + padding_length);
  // concat length bits
  memcpy(x_data + (data_len + padding_length) * sizeof(uint8_t), length_bits , 2 * sizeof(uint32_t));
//	char* tmpdata = *data;
//	if (*tmpdata) free(*tmpdata);
  if (new_data) free(new_data);
  free(padding);
  free(length_bits);

//	printf("xdata:");
//	for(int i =0;i<xdata_length;i++)
//	   printf("%02x", x_data[i]);
//	printf("\n");

  // little endian to big endian
//	uint32_t testd = x_data[3];
//	testd = testd << 8 | x_data[2];
//
//  testd = testd << 8 | x_data[1];
//  testd = testd << 8 | x_data[0];
//  printf("testd:%u\n" ,testd);

  (*data) = x_data;

  return data_len + padding_length + 2 * sizeof(uint32_t);
}


/*
 * returns ripemd128(message, length)
 * message should be a string terminated by '\0'
 */
byte *ripemd128(const byte *message, int len)
{
   dword32       digest[RMDsize/32];   /* contains (A, B, C, D(, E))   */
   static byte   hashcode[RMDsize/8];  /* for final hash-value         */
   dword32       X[16];                /* current 16-word chunk        */
   word          i;                    /* counter                      */
   dword32       length;               /* length in bytes of message   */
   dword32       nbytes;               /* # of bytes not yet processed */

   /* initialize */
   ripemd128Init(digest);

   length = (dword32) len;

   /* process message in 16-word chunks */
   for (nbytes=length; nbytes > 63; nbytes-=64) {
      for (i=0; i<16; i++) {
         X[i] = BYTES_TO_DWORD(message);
//         printf("%ld ", X[i]);
         message += 4;
      }
//      printf("\nX[i] ");
//      for(int j = 0; j < 16; j++) {
//         printf("%ld ", X[j]);
//      }
//
//      printf("\n");
//      printf("\nMDBUFFER[i] ");
//      for(int j = 0; j < 4; j++) {
//         printf("%lx ", digest[j]);
//      }
//      printf("\n");

      ripemd128compress(digest, X);

// printf("\nMDBuffer: ");
//      for(int j = 0; j < 4; j++) {
//         printf("%lx ", digest[j]);
//      }

   }  /* length mod 64 bytes left */


   /* finish: */
//   MDfinish(MDbuf, message, length, 0);

//   printf("\nMDBuffer: ");
//      for(int j = 0; j < 4; j++) {
//         printf("%lx ", digest[j]);
//      }
//      printf("\n");

   for (i=0; i<RMDsize/8; i+=4) {
      hashcode[i]   = (byte) (digest[i >> 2] & 0x000000ff);         /* implicit cast to byte  */
      hashcode[i+1] = (byte) ((digest[i >> 2] >> 8 ) & 0x000000ff);  /*  extracts the 8 least  */
      hashcode[i+2] = (byte) ((digest[i >> 2] >> 16) & 0x000000ff);  /*  significant bits.     */
      hashcode[i+3] = (byte) ((digest[i >> 2] >> 24) & 0x000000ff);
   }
//   for (i=0; i<RMDsize/8; i+=4) {
//     printf("%02x(%d)",hashcode[i]  ,hashcode[i]);         /* implicit cast to byte  */
//     printf("%02x(%d)",hashcode[i+1],hashcode[i+1]);  /*  extracts the 8 least  */
//     printf("%02x(%d)",hashcode[i+2],hashcode[i+2]);  /*  significant bits.     */
//     printf("%02x(%d)",hashcode[i+3],hashcode[i+3]);
//   }

   return (byte *)hashcode;
}

/********************************************************************/

/********************************************************************/

byte* ripemd128bytes(uint8_t *message, int length)
{
//   int i;
//   printf("\n* message: %s\n  hashcode: ", message);

   int rlen = ripemd128PaddingISO7816(&message, length);
//   printf("\n* padding  length: %d\n  : ", rlen);
//   printf("\n* message  length: %d\n  : ", (int) strlen(message));

   byte* hashcode = ripemd128(message, rlen);

   for (int i=0; i< RMDsize/8 /*ripemd size / 8 = 128 / 8*/; i++) {
     hashcode[i] = (byte)(hashcode[i] & 255);
//     printf("%d ", hashcode[i]);
   }
//   printf("\n");
//   for (i=0; i< RMDsize/8; i++) {
//      printf("%02x ", hashcode[i] );
//   }
   return hashcode;
}



/************************ end of file rmd128.c **********************/

#ifdef __cplusplus
}
#endif
