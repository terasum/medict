import { Dictionary } from '../../main/domain/Dictionary';

export declare class DictItem {
  id: string;
  alias: string;
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
  suggestWords: string[];
  historyStack: any[];
  currentWord: string;
  currentLookupWord: string;
  currentContent: string;
  currentSelectDict: DictItem;
}
