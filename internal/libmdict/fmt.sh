#!/bin/sh
set -e
cd $(dirname "$0")
clang-format -style=Google -i *.cc *.h
