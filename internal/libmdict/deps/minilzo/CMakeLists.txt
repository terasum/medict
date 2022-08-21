SET(CMAKE_OSX_DEPLOYMENT_TARGET "11.0" CACHE STRING "Minimum OS X deployment version")

PROJECT(minilzo C)
cmake_minimum_required(VERSION 3.10)
if(CMAKE_BUILD_TYPE STREQUAL "")
	# CMake defaults to leaving CMAKE_BUILD_TYPE empty. This screws up
	# differentiation between debug and release builds.
	set(CMAKE_BUILD_TYPE "Release" CACHE STRING
			"Choose the type of build, options are: None (CMAKE_CXX_FLAGS or \
CMAKE_C_FLAGS used) Debug Release RelWithDebInfo MinSizeRel." FORCE)
endif ()

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/target/bin)

set(minilzo_SOURCE minilzo.c testmini.c)

add_library(minilzo ${minilzo_SOURCE})

# add_executable(miniz_tester ${MINIZ_TESTER_SRC_LIST})
# target_link_libraries(miniz_tester miniz)

