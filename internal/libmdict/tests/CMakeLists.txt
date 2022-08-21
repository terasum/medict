#设置CMAKE最小版本
CMAKE_MINIMUM_REQUIRED(VERSION 3.10)

enable_testing()
include_directories(${PROJECT_SOURCE_DIR})
include_directories(${PROJECT_SOURCE_DIR}/deps/googletest/googletest/include)
include_directories(${PROJECT_SOURCE_DIR}/deps/googletest/googlemock/include)

add_subdirectory(testdict)

add_executable(test_reduce0 test_lookup_01.cpp)
target_link_libraries(test_reduce0 GTest GTestMain mdict miniz)
add_test(NAME test_reduce0 COMMAND test_reduce0)

add_executable(test_adler32 test_adler32.cpp)
target_link_libraries(test_adler32 GTest GTestMain mdict)
add_test(NAME test_adler32 COMMAND test_adler32)


add_executable(test_wordlist test_wordlist.cpp)
target_link_libraries(test_wordlist GTest GTestMain mdict miniz)
add_test(NAME test_wordlist COMMAND test_wordlist)