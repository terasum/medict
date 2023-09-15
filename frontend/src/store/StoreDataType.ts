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

import { Dictionary } from '../model/Dictionary';

export declare class DictItem {
  id: string;
  name: string;
}

export declare class StoreDataType {
  defaultWindow: string;
  // defaultWindow: '/',
  headerData: {
    currentTab: string;
    // currentTab: '词典',
  };
  sideBarData: {
    selectedWordIdx: number;
    candidateWordNum: number;
  };

  dictionaries: Dictionary[];
  dictBaseDir: string;
  suggestWords: any[];
  historyStack: any[];
  currentWord: { dictid: string, word:string };
  currentLookupWord: string;
  currentActualWord: string;
  currentContent: string;
  currentSelectDict: DictItem;
  translateApi: {
    baidu: {
      appid: string;
      appkey: string;
    };
  };
}
