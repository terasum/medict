import Mdict from 'js-mdict';

import { SuggestItem } from './SuggestItem';
import { Definition, NullDef } from './Definition';
import { StorabeDictionary } from './StorableDictionary';
import FuzzyTrie from './indexing/FuzzyTrie';
// import Configuration  from './Configuration'

// const config = Configuration.newInstance();

export declare class DictProps {
  id: string
  alias: string
  name: string
  mdxpath: string
  mddpath?: string | string[]
  description?: string
}

/**
 * Dictionary 词典实体类
 * id: 词典唯一ID, 通常是随机生成的UUID
 * alias: 别名，用于展示在下拉框中，不超过6字符
 * name: 词典全名
 * mdxpath: mdx 文件完整路径
 * mddpath: mdd 文件完整路径
 * resourceBaseDir: 资源文件基础路径
 * mdxDict: mdx 词典实例
 * mddDict: mdd 词典实例
 */
export class Dictionary extends StorabeDictionary {
  mdxDict: Mdict;
  mddDicts: Mdict[];
  description: string;
  mdxIndex: FuzzyTrie = new FuzzyTrie();
  mddIndex: FuzzyTrie[] = [];
  indexed: boolean = false;

  public static newByStorabe(dict: StorabeDictionary) {
    return new Dictionary(dict);
  }

  constructor(arg: DictProps) {
    super(arg.id, arg.alias, arg.name, arg.mdxpath, arg.mddpath);
    this.mdxDict = new Mdict(arg.mdxpath);
    this.mddDicts = [];

    if (arg.mddpath && typeof arg.mddpath === 'string' && arg.mddpath != '' && arg.mddpath.length > 0) {
      this.mddDicts.push(new Mdict(arg.mddpath));
    }

    if (arg.mddpath && isArray(arg.mddpath) && arg.mddpath.length > 0) {
      for (let i = 0; i < arg.mddpath.length; i++) {
        this.mddDicts.push(new Mdict(arg.mddpath[i]));
      }
    }

    this.description = arg.description || 'undefined';
  }

  indexing() {
    if (!this.mdxDict) {
      return;
    }
    // index mdx keys
    const keyWords = this.mdxDict.rangeKeyWords();
    console.log('[web worker] words length', keyWords.length);
    // for (let word of keyWords) {
    //   this.mdxIndex.add(word.keyText, word);
    // }

    // indexing mdd
    // this.mddDicts.forEach(mdd => {
    //   const mddTrie = new FuzzyTrie();
    //   const mddKey = mdd.rangeKeyWords();
    //   for (let key of mddKey) {
    //     mddTrie.add(key.keyText, key);
    //   }
    //   this.mddIndex.push(mddTrie)
    // });
    // 将词典标记为已经索引完成
    this.indexed = true;
  }

  lookupIndex(word: string) {
    const wordIdx = this.mdxIndex.has(word);
    if (wordIdx) {
      return {
        keyText: wordIdx!.getData().keyText,
        rofset: wordIdx!.getData().recordStartOffset
      }
    }
    return undefined;
  }

  lookup2(word: string) {
    const wordIndex = this.mdxIndex.has(word);
    if (!wordIndex || wordIndex == null || !wordIndex.getData()) {
      return NullDef(word);
    }

    const result = this.mdxDict.parse_defination(wordIndex.getData().keyText, wordIndex.getData().roffset);
    if (!result) {
      return NullDef(word);
    }
    return (result as unknown) as Definition;

  }

  associate2(word: string) {
    const result: SuggestItem[] = [];
    if (word.trim() == '' || word.length === 0) {
      return result;
    }
    // limits word result upto 50
    let counter = 0;
    const limit = 50;
    const pfxWords = this.mdxIndex.prefix(word);
    if (!pfxWords || pfxWords.length == 0) {
      return [];
    }
    for (let i = 0; i < pfxWords?.length ?? 0; i++) {
      if (counter >= limit) {
        break;
      }
      const word = pfxWords[i];
      result.push({ id: counter, dictid: this.id, ...word.data });
      counter++;
    }
    return result;
  }


  findWordDefinition(keyText: string, roffset: number) {
    const result = this.mdxDict.parse_defination(keyText, roffset);
    if (!result) {
      return NullDef(keyText);
    }
    return (result as unknown) as Definition;
  }

  findWordResource(keyText: string, withPayload = true) {
    let result = NullDef(keyText);
    for (let i = 0; i < this.mddDicts.length; i++) {
      const tempMdd = this.mddDicts[i];
      if (!tempMdd) {
        continue;
      }
      let tempResult = tempMdd.lookup(keyText);
      if (!tempResult || !tempResult.definition) {
        continue;
      }
      const returnResult =  {
        keyText: tempResult.keyText,
        contentSize: tempResult.definition.length,
        payload: withPayload ? tempResult.definition : undefined,
      } as Definition;
      return returnResult;
    }
    return result;
  }

  lookup(word: string) {
    return this.mdxDict.lookup(word);
  }

  associate(word: string) {
    const result: SuggestItem[] = [];
    if (word.trim() == '' || word.length === 0) {
      return result;
    }
    // limits word result upto 50
    let counter = 0;
    const limit = 50;
    const words = this.mdxDict.associate(word);
    for (let i = 0; i < words?.length ?? 0; i++) {
      if (counter >= limit) {
        break;
      }
      const word = words[i];
      // logger.info(`set ${key}, ${word.keyText}`)
      result.push({ id: counter, dictid: this.id, ...word });
      counter++;
    }
    return result;
  }
}

/***
 * -----------------------
 *  辅助方法类
 * -----------------------
 */
function isArray(o: any) {
  return Object.prototype.toString.call(o) == '[object Array]';
}

