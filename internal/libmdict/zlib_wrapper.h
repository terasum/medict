#ifndef MDICT_ZIP_WRAPPER_H_
#define MDICT_ZIP_WRAPPER_H_
#include <vector>

#include "deps/miniz/miniz.h"
// #include "zlib.h"

/*
 * 调用zlib解压缩数据
 * uncompress_bound为压缩前的数据长度,如果不知道数据源长度设置为0
 * */
inline std::vector<uint8_t> zlib_mem_uncompress(const void *source,
                                                size_t sourceLen,
                                                size_t uncompress_bound = 0) {
  //   throw_except_if_msg(nullptr==source||0==sourceLen,"invalid source");
  // uncompress_bound为0时将缓冲区设置为sourceLen的8倍长度
  if (!uncompress_bound) uncompress_bound = sourceLen << 3;
  for (;;) {
    std::vector<uint8_t> buffer(uncompress_bound);
    auto destLen = uLongf(buffer.size());
    auto err =
        uncompress(buffer.data(), &destLen,
                   reinterpret_cast<const Bytef *>(source), uLong(sourceLen));
    if (Z_OK == err)
      return std::vector<uint8_t>(buffer.data(), buffer.data() + destLen);
    else if (Z_BUF_ERROR == err) {
      // 缓冲区不足
      uncompress_bound <<= 2;  // 缓冲区放大4倍再尝试
      continue;
    }
    // 其他错误抛出异常
    printf("ZLIBERR %d\n", err);
    return std::vector<uint8_t>();
  }
}

/*
 * 调用zlib解压缩数据
 * */
inline int zlib_mem_uncompress(void *dest, size_t *destLen, const void *source,
                               size_t sourceLen) {
  //  throw_if(nullptr==source||0==sourceLen||nullptr==dest||nullptr==destLen||0==*destLen)
  auto len = uLongf(*destLen);
  auto err =
      uncompress(reinterpret_cast<Bytef *>(dest), &len,
                 reinterpret_cast<const Bytef *>(source), uLong(sourceLen));
  *destLen = size_t(len);
  return err;
}

#endif /* INCLUDE_ZLIB_WRAPPER_H_ */
