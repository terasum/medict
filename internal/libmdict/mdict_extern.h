#ifndef mdict_extern_h
#define mdict_extern_h


#ifdef __cplusplus
extern "C" {
#endif

#include "mdict_simple_key.h"
// ------------------------

/**
 init the dictionary
 */
void *mdict_init(const char *dictionary_path);

/**
 lookup a word
 */
void mdict_lookup(void *dict, const char *word, char **result);

void mdict_parse_definition(void *dict, const char *word, unsigned long record_start, char **result);

simple_key_item **mdict_keylist(void *dict, unsigned long *len);

int free_simple_key_list(simple_key_item **key_items, unsigned long len);

// 0 mdx, 1 mdd
int mdict_filetype(void *dict);

/**
suggest  a word
*/
void mdict_suggest(void *dict, char *word, char **suggested_words, int length);

/**
 return a stem
 */
void mdict_stem(void *dict, char *word, char **suggested_words, int length);

int mdict_destory(void *dict);

//-------------------------

#ifdef __cplusplus
}
#endif

#endif /* mdict_extern_h */
