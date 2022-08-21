#include "mdict_extern.h"
#include "mdict.h"

#include <sys/time.h>

#include <cstdlib>
#include <iostream>
#include <string>
#include <vector>
#include <gtest/gtest.h>


typedef long long int64;

class Timetool {
public:
    static int64 getSystemTime() {
        timeval tv;
        gettimeofday(&tv, NULL);
        int64 t = tv.tv_sec;
        t *= 1000;
        t += tv.tv_usec / 1000;
        return t;
    }
};


int cmp_word_list() {
    auto *mydict = new mdict::Mdict("../testdict/testdict.mdx");
    mydict->init();

    std::ifstream myfile;
    myfile.open("../testdict/wordlist.txt");

    if (!myfile.is_open()) {
        perror("Error open");
        exit(EXIT_FAILURE);
    }

    std::vector<std::string> words;
    std::string line;
    while (getline(myfile, line)) {
        words.push_back(line);
//        std::cout << line << std::endl;
    }
    std::cout << "total words : " << words.size() << std::endl;

    int foundCount = 0;
    for (auto &word : words) {
        std::string result = mydict->lookup(word);
        if (!result.empty()) {
            foundCount++;
//            std::cout << "found " << word << ", count : " << foundCount << std::endl;
        } else {
            std::cout << word << " not found!" << std::endl;
        }
    }
    std::cout << "foundCount: " << foundCount << ", total: " << words.size() << std::endl;

    assert(foundCount == words.size());
    return 0;
}


TEST(mdict, cmp_word_list) {
    EXPECT_EQ(0,cmp_word_list());
}

int main(int argc, char **argv) {

    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}