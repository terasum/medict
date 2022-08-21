#include "mdict_extern.h"
#include "mdict.h"

#include <algorithm>

/**
  实现 mdict_extern.h中的方法
 */

#ifdef __cplusplus
extern "C" {
#endif

/**
 init the dictionary
 */
void *mdict_init(const char *dictionary_path) {
    std::string dict_file_path(dictionary_path);
    auto *mydict = new mdict::Mdict(dict_file_path);
    mydict->init();
    return mydict;
}

/**
 lookup a word
 */
void mdict_lookup(void *dict, const char *word, char **result) {
    auto *self = (mdict::Mdict *) dict;
    std::string queryWord(word);
    std::string s = self->lookup(queryWord);

    (*result) = (char *) calloc(sizeof(char), s.size() + 1);
    std::copy(s.begin(), s.end(), (*result));
    (*result)[s.size()] = '\0';
}

void mdict_parse_definition(void *dict, const char *word, unsigned long record_start, char **result) {
    auto *self = (mdict::Mdict *) dict;
    std::string queryWord(word);
    std::string s = self->parse_definition(queryWord, record_start);

    (*result) = (char *) calloc(sizeof(char), s.size() + 1);
    std::copy(s.begin(), s.end(), (*result));
    (*result)[s.size()] = '\0';
}

simple_key_item **mdict_keylist(void *dict, unsigned long *len) {
    auto *self = (mdict::Mdict *) dict;
    auto keylist = self->keyList();

    *len = keylist.size();
    auto *items = new simple_key_item *[keylist.size()];

    for (auto i = 0; i < keylist.size(); i++) {
        items[i] = new simple_key_item;
        auto key_word = (const char *) keylist[i]->key_word.c_str();
        auto key_size = keylist[i]->key_word.size() + 1;
        items[i]->key_word = (char *) malloc(sizeof(char) * key_size);
        strcpy(items[i]->key_word, key_word);
        items[i]->key_word[key_size] = '\0';
        items[i]->record_start = keylist[i]->record_start;
    }

    return items;
}


int free_simple_key_list(simple_key_item **key_items, unsigned long len) {
    if (key_items == nullptr) {
        return 0;
    }

    for (unsigned long i = 0; i < len; i++) {
        if (key_items[i]->key_word != nullptr) {
            free(key_items[i]->key_word);
        }
        delete key_items[i];
    }

    return 0;
}

int mdict_filetype(void *dict) {
    auto *self = (mdict::Mdict *) dict;
    if (self->filetype == "MDX") {
        return 0;
    }
    return 1;
}

/**
suggest  a word
*/
void mdict_suggest(void *dict, char *word, char **suggested_words, int length) {

}

/**
 return a stem
 */
void mdict_stem(void *dict, char *word, char **suggested_words, int length) {}

int mdict_destory(void *dict) {
    auto *self = (mdict::Mdict *) dict;
    delete self;
    return 0;
}

#ifdef __cplusplus
}
#endif
