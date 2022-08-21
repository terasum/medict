#设置工程名称

CMAKE_POLICY(SET CMP0048 NEW)
SET(CMAKE_OSX_DEPLOYMENT_TARGET "11.0" CACHE STRING "Minimum OS X deployment version")
PROJECT(MDictCpp VERSION 1.0.0)

#设置CMAKE最小版本
CMAKE_MINIMUM_REQUIRED(VERSION 3.10)

#设置构建类型，及相应的编译选项
SET(CMAKE_BUILD_TYPE "Debug")

# SET(CMAKE_OSX_SYSROOT "/Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk")
# SET(CMAKE_OSX_DEPLOYMENT_TARGET "11.0" CACHE STRING "Minimum OS X deployment version")
 SET(CMAKE_CXX_FLAGS_DEBUG "$ENV{CXXFLAGS} -O0 -Wall -g -ggdb --std=c++11 -mmacosx-version-min=11.0")
 SET(CMAKE_CXX_FLAGS_RELEASE "$ENV{CXXFLAGS} -O3 -Wall --std=c++11 -mmacosx-version-min=11.0")

#设置执行文件输出目录
SET(EXECUTABLE_OUTPUT_PATH ${CMAKE_CURRENT_BINARY_DIR}/bin)

#设置库输出路径
SET(LIBRARY_OUTPUT_PATH ${CMAKE_CURRENT_BINARY_DIR}/lib)

#头文件搜索目录
INCLUDE_DIRECTORIES(/usr/local/include ${PROJECT_SOURCE_DIR} .)


#链接库搜索路径
LINK_DIRECTORIES(/usr/local/lib ${CMAKE_CURRENT_BINARY_DIR}/lib)

# 添加子目录
ADD_SUBDIRECTORY(deps)

# ADD_SUBDIRECTORY(tests)


ADD_EXECUTABLE(mydict mdict.cc binutils.cc xmlutils.cc ripemd128.c adler32.cc mdict_extern.cc mdict_extern_test.cc)
TARGET_LINK_LIBRARIES(mydict miniz)

# most important library
ADD_LIBRARY(mdict mdict.cc binutils.cc xmlutils.cc ripemd128.c adler32.cc mdict_extern.cc)

# adler32 test
ADD_EXECUTABLE(adler32_test adler32.cc adler32_test.cc binutils.cc)

# ripemd128 test
ADD_EXECUTABLE(ripemd128_test ripemd128_test.c ripemd128.c)

# xmlutils
ADD_EXECUTABLE(xmlutils_test xmlutils.cc xmlutils_test.cc)

ADD_EXECUTABLE(binutils_test binutils.cc binutils_test.cc)
