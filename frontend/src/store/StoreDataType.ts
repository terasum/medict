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
