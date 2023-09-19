/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import { defineStore } from 'pinia';

import { GetAllDicts, SearchWord } from '@/apis/dicts-api';
import { StaticDictServerURL } from '@/apis/apis';

function constructQueryURL(entry) {
  let {
    baseURL,
    dict_id,
    key_word,
    record_start_offset,
    record_end_offset,
    key_block_idx,
    entry_id,
  } = entry;
  return `${baseURL}/__tcidem_query?dict_id=${dict_id}&key_word=${key_word}&record_start_offset=${record_start_offset}&record_end_offset=${record_end_offset}&key_block_idx=${key_block_idx}&entry_id=${entry_id}`;
}

const countDownJs = `
<html>
<head>
<meta content="width=device-width, initial-scale=1.0" name="viewport" />
<style>
  html, body {
    height: 100%;
    padding:0;
    margin:0;
    user-select: none;
    -moz-user-select: none;
    -webkit-user-select: none;
  }
  #skeleton{
    width: 100%;
    height: 100%;
    color:#696969;
    text-align: center;
    display: flex;
    justify-content: center;
    flex-direction: column;
    
  }
  #skeleton > h1 {
      font-style: italic;
      cursor: default;
    }
</style>
</head>
<body>
<div id="skeleton">
<h1>Medict</h1>
</div>
</body>
</html>

`;

export const useDictQueryStore = defineStore('dictQuery', {
  state: () => ({
    dictApiBaseURL: '',
    queryPendingList: [],
    mainContent: btoa(countDownJs),
    mainContentURL: '',
    selectDict: { id: '', name: '', path: '' },
    inputSearchWord: '',
  }),
  actions: {
    queryDictList() {
      return GetAllDicts();
    },
    updateInputSearchWord(word :string) {
      console.log(`===> updateInputSearchWord: ${word}`)
      if (!word || word.trim() == "") {
        return;
      }
      if (word == this.inputSearchWord) {
        return;
      }
        this.inputSearchWord = word;
        this.searchWord(word)
      
    },
    updateMainContent(content) {
      if (content === '') {
        this.mainContent = btoa(countDownJs);
      } else {
        this.mainContent = content;
      }
    },
    updateMainContentURL(url) {
      this.mainContentURL = url;
    },
    updatePendingList(wordList) {
      console.log('==== updatePendingList ====');
      console.log(wordList);
      this.queryPendingList = wordList;
    },
    updateSelectDict(dictItem) {
      this.selectDict = dictItem;
      if (this.inputSearchWord && this.inputSearchWord.trim() != '') {
        this.searchWord(this.inputSearchWord);
      }
    },
    searchWord(word) {
      if (this.selectDict.id === '') {
        console.log('skipped');
        return;
      }
      if (word === '') {
        console.log('empty word skipped');
      }

      SearchWord(this.selectDict.id, word).then((res) => {
        console.log('=== SearchWord ===');
        console.log(this.selectDict.id);
        console.log(res);
        this.updatePendingList(res);
        if (res.data && res.data.length > 0) {
          this.locateWord(0);
        }
      });
    },
    setUpAPIBaseURL() {
      let count = 0;
      let that = this;
      let inv = setInterval(function () {
        let urlPromise = StaticDictServerURL();

        if (!urlPromise) {
          clearInterval(inv);
          return;
        }

        urlPromise
          .then((url) => {
            if (url === '') {
              console.log(`[init] static server url empty retry ${count}`);
              return;
            }
            // browser
            if (url === 'http://localhost:1/') {
              return;
            }
            if (url.startsWith('http://localhost:0/')) {
              console.log(`[init] static server url empty retry ${count}`);
              return;
            }
            console.log(`[init] static server url success ${count}`);
            that.updateBaseURL(url);
            clearInterval(inv);
          })
          .catch((err) => {
            console.error(err);
            clearInterval(inv);
          });
      }, 1000);
    },
    updateBaseURL(url) {
      console.log(url);
      this.dictApiBaseURL = url;
    },
    locateWord(entry_idx) {
      if (this.dictApiBaseURL === '') {
        console.log(
          "dictionary has not ready, baseurl hasn't assigned, skipped"
        );
      }
      if (entry_idx < 0) {
        return;
      }

      if (entry_idx >= this.queryPendingList.length) {
        return;
      }

      if (this.selectDict.id === '') {
        console.log(
          "dictionary has not ready, dictionary id hasn't assigned, skipped"
        );
        return;
      }

      let entry = this.queryPendingList[entry_idx];

      this.updateMainContentURL(
        constructQueryURL({
          baseURL: this.dictApiBaseURL,
          dict_id: this.selectDict.id,
          key_word: entry.key_word,
          record_start_offset: entry.record_start_offset,
          record_end_offset: entry.record_end_offset,
          key_block_idx: entry.key_block_idx,
          entry_id: entry_idx,
        })
      );
    },
  },
});
