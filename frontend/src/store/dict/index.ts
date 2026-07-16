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
import { getPreferences, savePreferences } from '@/apis/config';

/**
 * 历史栈条目类型。统一字段命名为 `keyword`，避免历史上 `keyword` / `key_word`
 * 混用导致的去重失效与后退/前进读不到字段的问题（issue #699）。
 */
export interface HistoryEntry {
  baseurl: string;
  dict_id: string;
  dict: any;
  keyword: string;
  record_start_offset: number;
  record_end_offset: number;
  key_block_idx?: number;
  entry_id: number;
  record_block_data_start_offset: number;
  record_block_data_compress_size: number;
  record_block_data_decompress_size: number;
  keyword_data_start_offset: number;
  keyword_data_end_offset: number;
  // 多词典查询历史条目:multi=true 时用 keyword + multiDicts 重放,offsets 字段忽略。
  multi?: boolean;
  multiDicts?: any[];
}

// MultiResult 是多词典堆叠里的一节:某词典对该词首个匹配的释义 URL(无匹配则 empty)。
export interface MultiResult {
  dictId: string;
  dictName: string;
  url: string;
  empty: boolean;
}

// buildEntryURL 把一个词条匹配(entry,含 keyword 与各 offset)拼成 Gin 的释义查询 URL。
// 抽出来供单词典(locateWord)、多词典(searchWordMulti)与悬停弹窗共用,dictId 显式传入。
export function buildEntryURL(baseurl: string, dictId: string, entry: any, entryId: number): string {
  return `${baseurl}/__tcidem_query?dict_id=${dictId}` +
    `&keyword=${entry.keyword}&record_start_offset=${entry.record_start_offset}` +
    `&entry_id=${entryId}` +
    `&record_end_offset=${entry.record_end_offset}` +
    `&record_block_data_start_offset=${entry.record_block_data_start_offset}` +
    `&record_block_data_compress_size=${entry.record_block_data_compress_size}` +
    `&record_block_data_decompress_size=${entry.record_block_data_decompress_size}` +
    `&keyword_data_start_offset=${entry.keyword_data_start_offset}` +
    `&keyword_data_end_offset=${entry.keyword_data_end_offset}`;
}

function constructQueryURL(entry: HistoryEntry) {
  return buildEntryURL(entry.baseurl, entry.dict_id, entry, entry.entry_id);
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

/**
 * searchWord 的请求序号（last-write-wins）。
 * 快速连续 Enter 或 iframe entry:// 跳转会触发多次 searchWord，
 * 仅最后一次（序号最大者）的响应会真正写入 pending list，避免乱序覆盖。
 */
let searchWordRequestId = 0;

export const useDictQueryStore = defineStore('dictQuery', {
  state: () => ({
    dictApiBaseURL: '',
    queryPendingList: [] as any[],
    mainContent: btoa(DefaultContentTemplpate),
    mainContentURL: '',
    selectDict: { id: '', name: '', path: '' } as any,
    inputSearchWord: '',

    historyStack: new HistoryStack(),

    // 多词典查询(类 GoldenDict):开关、激活词典集、堆叠结果。
    multiMode: false,
    multiSelectedDicts: [] as any[],
    multiResults: [] as MultiResult[],
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
      // 多词典模式:分发到 searchWordMulti(对激活集每个词典各查一次)。
      if (this.multiMode) {
        this.searchWordMulti(word);
        return;
      }
      if (this.selectDict.id === '') {
        return;
      }
      if (!word || word.trim() == '') {
        return;
      }

      // last-write-wins：分配本次请求的序号，响应回来时若已过期则丢弃
      const reqId = ++searchWordRequestId;
      SearchWord(this.selectDict.id, word).then((res) => {
        if (reqId !== searchWordRequestId) {
          console.info('[store-action]{searchWord} stale response, ignoring', word);
          return;
        }
        console.info('[store-action]{searchWord} success', word, res);

        this.updatePendingList(res);
      }).catch((err) => {
        if (reqId !== searchWordRequestId) {
          console.info('[store-action]{searchWord} stale error, ignoring', word);
          return;
        }
        console.info('[store-action]{searchWord} failed', err);
        this.updateSetCurrentDictAsContent();
      });
    },
    // 更新 pending list
    updatePendingList(wordList: any) {
      console.log(`[app-event](store-action), updatePendingList`, wordList);

      this.queryPendingList = wordList;
      
      if (this.queryPendingList && this.queryPendingList.length > 0) {
        this.locateWord(0);
      } else {
        this.updateSetCurrentDictAsContent()
      }
    },
    // 更新main iframe内容
    updateMainContent(content: string) {
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
    updateMainContentURL(url: string) {
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
    updateSelectDict(dictItem: any) {
      this.selectDict = dictItem;
      if (this.inputSearchWord && this.inputSearchWord.trim() != '') {
        this.searchWord(this.inputSearchWord);
      } else {
        this.updateSetCurrentDictAsContent();
      }
    },
    // ===================== 多词典查询(类 GoldenDict)=====================
    // 开/关多词典模式。开启时若激活集为空,用当前 selectDict 播种,并对当前词重跑;
    // 关闭时清空多结果并切回单词典视图。
    setMultiMode(on: boolean) {
      this.multiMode = on;
      if (on) {
        if (this.multiSelectedDicts.length === 0 && this.selectDict && this.selectDict.id !== '') {
          this.multiSelectedDicts = [this.selectDict];
        }
        if (this.inputSearchWord && this.inputSearchWord.trim() !== '') {
          this.searchWordMulti(this.inputSearchWord);
        }
      } else {
        this.multiResults = [];
        if (this.selectDict && this.selectDict.id !== '' && this.inputSearchWord && this.inputSearchWord.trim() !== '') {
          this.searchWord(this.inputSearchWord);
        }
      }
      this._persistMulti();
    },
    // 在激活集里增/减一个词典,并对当前词重跑。
    toggleMultiDict(dictItem: any) {
      const i = this.multiSelectedDicts.findIndex((d) => d.id === dictItem.id);
      if (i >= 0) {
        this.multiSelectedDicts.splice(i, 1);
      } else {
        this.multiSelectedDicts.push(dictItem);
      }
      if (this.inputSearchWord && this.inputSearchWord.trim() !== '' && this.multiSelectedDicts.length > 0) {
        this.searchWordMulti(this.inputSearchWord);
      } else {
        this.multiResults = [];
      }
      this._persistMulti();
    },
    // 把当前多词典开关 + 激活集持久化到 medict.toml(fire-and-forget,失败仅告警)。
    _persistMulti() {
      savePreferences({
        multimode: this.multiMode,
        multidictids: this.multiSelectedDicts.map((d) => d.id),
      }).catch((e: any) => console.warn('[multi-dict] persist failed', e));
    },
    // 启动时按已保存的偏好恢复多词典开关与激活集(用全量词典列表把 id 还原成对象)。
    async restoreMultiSelection(allDicts: any[]) {
      try {
        const prefs = await getPreferences();
        const ids = Array.isArray(prefs.multidictids) ? prefs.multidictids.map(String) : [];
        const idset = new Set(ids);
        this.multiSelectedDicts = allDicts.filter((d) => idset.has(String(d.id)));
        // 仅在曾经开过多词典、且激活集非空时恢复开关,避免首次用户被强制进入多模式
        if (prefs.multimode === true && this.multiSelectedDicts.length > 0) {
          this.multiMode = true;
        }
      } catch (e) {
        console.warn('[multi-dict] restore selection failed', e);
      }
    },
    // 多词典搜索:对激活集每个词典并行 SearchWord,各取首个匹配拼成释义 URL;无匹配则 empty。
    searchWordMulti(word: string) {
      if (this.multiSelectedDicts.length === 0) {
        this.multiResults = [];
        return;
      }
      if (!word || word.trim() === '') {
        return;
      }
      const baseurl = this.dictApiBaseURL;
      const reqId = ++searchWordRequestId;
      const dicts = [...this.multiSelectedDicts];
      Promise.all(
        dicts.map((d) =>
          SearchWord(d.id, word)
            .then((res: any) => ({ dict: d, res: res as any[] }))
            .catch(() => ({ dict: d, res: [] as any[] }))
        )
      ).then((settled) => {
        if (reqId !== searchWordRequestId) {
          console.info('[store-action]{searchWordMulti} stale response, ignoring', word);
          return;
        }
        const results: MultiResult[] = settled.map(({ dict, res }) => {
          const list = Array.isArray(res) ? res : [];
          const name = dict.name || (dict.description && dict.description.title) || dict.id;
          if (list.length > 0) {
            return { dictId: dict.id, dictName: name, url: buildEntryURL(baseurl, dict.id, list[0], 0), empty: false };
          }
          return { dictId: dict.id, dictName: name, url: '', empty: true };
        });
        this.multiResults = results;
        // 主词典(激活集首项)的完整匹配列表供 sidebar 浏览。
        const primary = dicts[0];
        this.queryPendingList = settled.find((s) => s.dict.id === primary.id)?.res || [];
        this.pushHistoryMulti(word, dicts);
      }).catch((err) => {
        if (reqId !== searchWordRequestId) return;
        console.info('[store-action]{searchWordMulti} failed', err);
        this.multiResults = [];
      });
    },
    // 多模式:用主词典匹配列表的第 entryIdx 条更新堆叠中主词典那节的释义 URL。
    locateInMultiPrimary(entryIdx: number) {
      if (entryIdx < 0 || entryIdx >= this.queryPendingList.length) return;
      const primary = this.multiSelectedDicts[0];
      if (!primary) return;
      const entry = this.queryPendingList[entryIdx];
      const i = this.multiResults.findIndex((r) => r.dictId === primary.id);
      if (i >= 0) {
        this.multiResults[i] = {
          ...this.multiResults[i],
          url: buildEntryURL(this.dictApiBaseURL, primary.id, entry, entryIdx),
          empty: false,
        };
      }
    },
    // 压一条多词典历史(keyword + 当时的激活集),供后退/前进重放。
    pushHistoryMulti(keyword: string, dicts: any[]) {
      if (!keyword || dicts.length === 0) return;
      if (!this.historyStack.isEmpty() && this.historyStack.peek().keyword === keyword && this.historyStack.peek().multi) {
        return;
      }
      this.historyStack.push({
        baseurl: this.dictApiBaseURL,
        dict_id: dicts[0].id,
        dict: dicts[0],
        keyword,
        record_start_offset: 0,
        record_end_offset: 0,
        entry_id: 0,
        record_block_data_start_offset: 0,
        record_block_data_compress_size: 0,
        record_block_data_decompress_size: 0,
        keyword_data_start_offset: 0,
        keyword_data_end_offset: 0,
        multi: true,
        multiDicts: dicts,
      });
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
      const MAX_RETRIES = 30;
      let that = this;
      let inv = setInterval(function () {
        count++;
        if (count > MAX_RETRIES) {
          console.error(
            `[app init] static server url polling exceeded max retries (${MAX_RETRIES}), giving up`
          );
          clearInterval(inv);
          return;
        }

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
              console.log(
                `[app init] static server url is placeholder, retrying times: ${count}`
              );
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
    updateBaseURL(url: string) {
      console.log(url);
      this.dictApiBaseURL = url;
    },
    // 定位单词并返回释义
    locateWord(entry_idx: number, skipPushHistory: boolean = false) {
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




      const locateQuerier: HistoryEntry = {
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
    _locateWord(locateQuerier: HistoryEntry) {
      console.log("frontend _locateWord", locateQuerier)
      let definitionURL = constructQueryURL(locateQuerier);
      this.updateMainContentURL(definitionURL);
    },
    resetMainContent() {
      this.updateMainContent(btoa(DefaultContentTemplpate));
    },
    pushHistory(qurier: HistoryEntry) {
      if (!qurier.keyword || qurier.keyword === '') {
        return;
      }
      if (qurier.baseurl === '') {
        return;
      }
      if (!this.historyStack.isEmpty() && this.historyStack.peek().keyword === qurier.keyword) {
        return
      }

      this.historyStack.push(qurier);
    },
    pushHistoryByEntryIDx(entry_idx: number){
      if (entry_idx < 0 || entry_idx >= this.queryPendingList.length) {
        return;
      }
      const entry = this.queryPendingList[entry_idx];

      const locateQuerier: HistoryEntry = {
        baseurl: this.dictApiBaseURL,
        dict_id: this.selectDict.id,
        dict: this.selectDict,
        keyword: entry.keyword,
        record_start_offset: entry.record_start_offset,
        record_end_offset: entry.record_end_offset,
        key_block_idx: entry.key_block_idx,
        entry_id: entry_idx,
        record_block_data_start_offset: entry.record_block_data_start_offset,
        record_block_data_compress_size: entry.record_block_data_compress_size,
        record_block_data_decompress_size: entry.record_block_data_decompress_size,
        keyword_data_start_offset: entry.keyword_data_start_offset,
        keyword_data_end_offset: entry.keyword_data_end_offset,
    }

      this.pushHistory(locateQuerier);
    },
    backHistory() {
      let locateQuerier: HistoryEntry | '' = this.historyStack.back();
      if (!locateQuerier || this.inputSearchWord == locateQuerier.keyword) {
        return;
      }
      this.updateInputSearchWord(locateQuerier.keyword)

      // 多词典历史条目:还原激活集并重跑多词典搜索,不走单词典 locate 分支。
      if (locateQuerier.multi && locateQuerier.multiDicts) {
        this.multiSelectedDicts = locateQuerier.multiDicts;
        this.searchWordMulti(locateQuerier.keyword);
        return;
      }

      if (this.selectDict.id != locateQuerier.dict_id) {
        this.selectDict = locateQuerier.dict;
      }

      SearchWord(locateQuerier.dict_id, locateQuerier.keyword ).then((res) => {
        console.info('[store-action]{backHistory} success', locateQuerier.keyword, res);
        this.queryPendingList = res as any;
      }).catch((err) => {
        console.info('[store-action]{backHistory} failed', err);
      });
      this._locateWord(locateQuerier);
    },

    forwardHistory() {
      let locateQuerier: HistoryEntry | '' = this.historyStack.forward();
      if (!locateQuerier || this.inputSearchWord == locateQuerier.keyword) {
        return;
      }

      this.updateInputSearchWord(locateQuerier.keyword)

      // 多词典历史条目:还原激活集并重跑多词典搜索,不走单词典 locate 分支。
      if (locateQuerier.multi && locateQuerier.multiDicts) {
        this.multiSelectedDicts = locateQuerier.multiDicts;
        this.searchWordMulti(locateQuerier.keyword);
        return;
      }

      if (this.selectDict.id != locateQuerier.dict_id) {
        this.selectDict = locateQuerier.dict;
      }

      SearchWord(locateQuerier.dict_id, locateQuerier.keyword ).then((res) => {
        console.info('[store-action]{forwardHistory} success', locateQuerier.keyword, res);
        this.queryPendingList = res as any;
      }).catch((err) => {
        console.info('[store-action]{forwardHistory} failed', err);
      });

      this._locateWord(locateQuerier);
    },
  },
});

class HistoryStack {
  items: HistoryEntry[] = [];
  pointer: number = -1;

  push(element: HistoryEntry) {
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

  back(): HistoryEntry | '' {
    console.log('back', this.pointer, this.items);
    if (this.pointer >= 1) {
      this.pointer -= 1;
      return this.items[this.pointer];
    } else if (this.pointer == 0) {
      return this.items[0];
    }
    return '';
  }

  forward(): HistoryEntry | '' {
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
  peek(): HistoryEntry {
    return this.items[this.items.length - 1];
  }
}
