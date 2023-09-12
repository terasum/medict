import { defineStore } from 'pinia';

import { GetAllDicts } from '@/apis/dicts-api';
import { StaticDictServerURL } from '@/apis/static-server-api';

function constructQueryURL(entry) {
  let {
    baseURL, dict_id, key_word, record_start_offset, record_end_offset, key_block_idx, entry_id,
  } = entry;
  return `${baseURL}/__tcidem_query?dict_id=${dict_id}&key_word=${key_word}&record_start_offset=${record_start_offset}&record_end_offset=${record_end_offset}&key_block_idx=${key_block_idx}&entry_id=${entry_id}`;
}

const countDownJs = `
<html>
<head>
<style>
  html, body {
     width: 100%;
    height: 100%;
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
  }),
  actions: {
    queryDictList() {
      return GetAllDicts();
    },
    updateMainContent(content) {
      if (content === "") {
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
      this.selectDict=dictItem;
    },
    setUpDictAPI() {
      StaticDictServerURL().then((url) => {
        console.log('========= static server url ===========');
        console.log(url);
        this.dictApiBaseURL = url;
      });
    },
    locateWord(entry_idx) {
      if (this.dictApiBaseURL === '') {
        console.log('dictionary has not ready, baseurl hasn\'t assigned, skipped');
      }
      if (entry_idx < 0) {
        return;
      }

      if (entry_idx >= this.queryPendingList.length) {
        return;
      }

      if (this.selectDict.id === '') {
        console.log('dictionary has not ready, dictionary id hasn\'t assigned, skipped');
        return;
      }

      let entry = this.queryPendingList[entry_idx];

      this.updateMainContentURL(constructQueryURL(
        {
          baseURL: this.dictApiBaseURL,
          dict_id: this.selectDict.id,
          key_word: entry.key_word,
          record_start_offset: entry.record_start_offset,
          record_end_offset: entry.record_end_offset,
          key_block_idx: entry.key_block_idx,
          entry_id: entry_idx,
        }));

    },
  },
});