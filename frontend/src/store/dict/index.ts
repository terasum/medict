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

import { InitDicts, GetAllDicts, SearchWord } from '@/apis/dicts-api';
import { StaticDictServerURL } from '@/apis/apis';

function constructQueryURL(entry) {
  let {
    baseurl,
    dict_id,
    keyword,
    record_start_offset,
    record_end_offset,
    entry_id,
    record_block_data_start_offset,
    record_block_data_compress_size,
    record_block_data_decompress_size,
    keyword_data_start_offset,
    keyword_data_end_offset
} = entry;
  return `${baseurl}/__tcidem_query?dict_id=${dict_id}`+
    `&keyword=${keyword}&record_start_offset=${record_start_offset}`+
    `&entry_id=${entry_id}`+
    `&record_end_offset=${record_end_offset}`+
    `&record_block_data_start_offset=${record_block_data_start_offset}`+
    `&record_block_data_compress_size=${record_block_data_compress_size}`+
    `&record_block_data_decompress_size=${record_block_data_decompress_size}`+
    `&keyword_data_start_offset=${keyword_data_start_offset}`+
    `&keyword_data_end_offset=${keyword_data_end_offset}`;
}

const DefaultContentTemplpate = `
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
    mainContent: btoa(DefaultContentTemplpate),
    mainContentURL: '',
    selectDict: { id: '', name: '', path: '' },
    inputSearchWord: '',

    historyStack: new HistoryStack(),
  }),
  actions: {
    initDicts() {
      return InitDicts();
    },
    // 取得当前词典列表
    queryDictList() {
      return GetAllDicts();
    },
    // 更新当前输入的单词（input) 展示的单词
    updateInputSearchWord(word: string) {
      console.log(`[app-event](store-action), updateInputSearchWord: ${word}`);
      if (!word || word.trim() == '') {
        if (this.selectDict && this.selectDict.id !== '') {
          console.log("[app-event] updateInputSearchWord, selectDict is not empty, update main content")
          this.updateMainContent(this.selectDict.description.description);
        }
        return;
      }
      if (word == this.inputSearchWord) {
        return;
      }
      this.inputSearchWord = word;
    },
    // 搜索单词
    searchWord(word: string) {
      if (this.selectDict.id === '') {
        return;
      }
      if (!word || word.trim() == '') {
        return;
      }

      SearchWord(this.selectDict.id, word).then((res) => {
        console.info('[store-action]{searchWord} success', word, res);
        
        this.updatePendingList(res);
      }).catch((err) => {
        console.info('[store-action]{searchWord} failed', err);
        this.updateSetCurrentDictAsContent();
      });
    },
    // 更新 pending list
    updatePendingList(wordList) {
      console.log(`[app-event](store-action), updatePendingList`, wordList);

      this.queryPendingList = wordList;
      
      if (this.queryPendingList && this.queryPendingList.length > 0) {
        this.locateWord(0);
      } else {
        this.updateSetCurrentDictAsContent()
      }
    },
    // 更新main iframe内容
    updateMainContent(content) {
      // 防止循环嵌入 frame
      if (this.dictApiBaseURL === "") {
        return;
      }
      if (content === '') {
        this.mainContent = btoa(DefaultContentTemplpate);
      } else {
        this.mainContent = content;
      }
    },
    // 更新 main iframe url
    updateMainContentURL(url) {
      // 防止循环嵌入 frame
      if (this.dictApiBaseURL === "") {
        return;
      }
      this.mainContentURL = url;
      if (url === '') {
        this.mainContent = btoa(DefaultContentTemplpate);
      }
    },
    // 更新选中的词典
    updateSelectDict(dictItem) {
      this.selectDict = dictItem;
      if (this.inputSearchWord && this.inputSearchWord.trim() != '') {
        this.searchWord(this.inputSearchWord);
      } else {
        this.updateSetCurrentDictAsContent();
      }
    },
    updateSetCurrentDictAsContent() {
        if (! this.selectDict || this.selectDict.id === '') {
          this.mainContent = btoa(DefaultContentTemplpate);
          return;
        }
        this.updateMainContent(this.selectDict.description.description);
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
              console.log(
                `[app init] static server url is empty, retrying times: ${count}`
              );
              return;
            }
            // browser
            if (url === 'http://localhost:1/') {
              return;
            }
            if (url.startsWith('http://localhost:0/')) {
              console.log(
                `[app init] static server url setting failed, retrying times: ${count}`
              );
              return;
            }
            console.log(
              `[app init] static server url has setting successful, retrying times: ${count}`
            );
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
    // 定位单词并返回释义
    locateWord(entry_idx, skipPushHistory: boolean = false) {
      if (this.dictApiBaseURL === '' || this.selectDict.id === '') {
        console.log(
          "app or dictionary has not ready, skipped"
        );
      }
      if (entry_idx < 0 || entry_idx >= this.queryPendingList.length) {
        return;
      }

      let entry = this.queryPendingList[entry_idx];

      this.updateInputSearchWord(entry.keyword);




      const locateQuerier = {
          baseurl: this.dictApiBaseURL,
          dict_id: this.selectDict.id,
          dict: this.selectDict,
          keyword: entry.keyword,
          record_start_offset: entry.record_start_offset,
          record_end_offset: entry.record_end_offset,
          key_block_idx: entry.key_block_idx,
          entry_id: entry_idx,
          record_block_data_start_offset:entry.record_block_data_start_offset,
          record_block_data_compress_size:entry.record_block_data_compress_size,
          record_block_data_decompress_size:entry.record_block_data_decompress_size,
          keyword_data_start_offset:entry.keyword_data_start_offset,
          keyword_data_end_offset:entry.keyword_data_end_offset,
      }

      if (!skipPushHistory) {
        this.pushHistory(locateQuerier);
      }
      this._locateWord(locateQuerier);

    },
    _locateWord(locateQuerier) {
      console.log("frontend _locateWord", locateQuerier)
      let definitionURL = constructQueryURL(locateQuerier);
      this.updateMainContentURL(definitionURL);
    },
    resetMainContent() {
      this.updateMainContent(btoa(DefaultContentTemplpate));
    },
    pushHistory(qurier: any) {
      if (qurier.key_word === '') {
        return;
      }
      if (qurier.baseurl === '') {
        return;
      }
      if (!this.historyStack.isEmpty() && this.historyStack.peek().key_word === qurier.key_word) {
        return
      }
      
      this.historyStack.push(qurier);
    },
    pushHistoryByEntryIDx(entry_idx){
      if (entry_idx < 0 || entry_idx >= this.queryPendingList.length) {
        return;
      }
      const entry = this.queryPendingList[entry_idx];

      const locateQuerier = {
        baseurl: this.dictApiBaseURL,
        dict_id: this.selectDict.id,
        dict: this.selectDict,
        key_word: entry.key_word,
        record_start_offset: entry.record_start_offset,
        record_end_offset: entry.record_end_offset,
        key_block_idx: entry.key_block_idx,
        entry_id: entry_idx,
    }

      this.pushHistory(locateQuerier);
    },
    backHistory() {
      let locateQuerier = this.historyStack.back();
      if (this.inputSearchWord == locateQuerier.keyword) {
        return;
      }
      this.updateInputSearchWord(locateQuerier.key_word)

      if (this.selectDict.id != locateQuerier.dict_id) {
        this.selectDict = locateQuerier.dict;
      }

      SearchWord(locateQuerier.dict_id, locateQuerier.key_word ).then((res) => {
        console.info('[store-action]{forwardHistory} success', locateQuerier.key_word, res);
        this.queryPendingList = res;
      }).catch((err) => {
        console.info('[store-action]{forwardHistory} failed', err);
      });
      this._locateWord(locateQuerier);
    },

    forwardHistory() {
      let locateQuerier = this.historyStack.forward();
      if (this.inputSearchWord == locateQuerier.key_word) {
        return;
      }

      this.updateInputSearchWord(locateQuerier.key_word)

      if (this.selectDict.id != locateQuerier.dict_id) {
        this.selectDict = locateQuerier.dict;
      }

      SearchWord(locateQuerier.dict_id, locateQuerier.key_word ).then((res) => {
        console.info('[store-action]{forwardHistory} success', locateQuerier.key_word, res);
        this.queryPendingList = res;
      }).catch((err) => {
        console.info('[store-action]{forwardHistory} failed', err);
      });

      this._locateWord(locateQuerier);
    },
  },
});

class HistoryStack {
  items: any[] = [];
  pointer: number = -1;

  push(element: any) {
    console.log('push', this.pointer, this.items);
    if (
      this.items.length > 0 &&
      this.items[this.items.length - 1] === element
    ) {
      return;
    }

    this.items.push(element);
    this.pointer = this.items.length - 1;
  }

  back() {
    console.log('back', this.pointer, this.items);
    if (this.pointer >= 1) {
      this.pointer -= 1;
      return this.items[this.pointer];
    } else if (this.pointer == 0) {
      return this.items[0];
    }
    return '';
  }

  forward() {
    console.log('forward', this.pointer, this.items);
    if (this.pointer < this.items.length - 1) {
      this.pointer += 1;
      return this.items[this.pointer];
    } else if (this.pointer == this.items.length - 1) {
      return this.items[this.pointer];
    }

    return '';
  }

  isEmpty() {
    return this.items.length == 0;
  }

  size() {
    return this.items.length;
  }
  peek(){
    return this.items[this.items.length - 1];
  }
}
