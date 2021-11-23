import { DictService } from './DictionaryService';
import { DictContentService } from './DictionaryContentService';
import { logger } from '../../utils/logger';

const dictService = new DictService();
const dictContentService = new DictContentService();

/**
 * WordService 查词主要服务，依赖词典服务
 */
export class WordQueryService {
  /**
   * 查询建议词列表
   * @param dictid 词典ID
   * @param word 查询词
   * @returns 返回建议词列表
   */
  suggestWord(dictid: string, word: string) {
    logger.info("[main-process] suggestWord event.sender.send('onSuggestWord')");
    return dictService.associate(dictid, word);
  }

  /**
   * 处理原生HTML返回结果
   * @param dictid 词典ID
   * @param keyText 关键词
   * @param definition 原生返回结果
   * @returns 
   */
  definitionReplace(dictid: string, keyText: string, definition: string) {
    // 资源查询函数
    const resFn = (resKey: string, withPayload = false) => {
      return dictService.loadDictResource(dictid, resKey, withPayload);
    };

    // 词查询函数
    const lookupFn = (word: string) => {
      return dictService.lookup(dictid, word);
    };

    logger.info("[main-process] suggestWord event.sender.send('onFindWordPrecisly')");

    return dictContentService.definitionReplace(
      dictid,
      keyText,
      definition,
      lookupFn,
      resFn
    )
  }

  /**
   * 单个keyText范围查询，返回列表
   * @param dictid dictionary id
   * @param keyText 待查询词
   * @returns 返回单个词的结果
   */
  lookup(dictid: string, keyText: string) {
    logger.info("[main-process] suggestWord event.sender.send('onFindWordPrecisly')");
    const wordResult = dictService.lookup(dictid, keyText);
    return wordResult;
  }

  /**
   * findWordPrecisly 精确查询某个词
   * @param event 事件源
   * @param arg 词详细信息
   */
  findWordPrecisly(keyText: string, dictid: string, rofset: number) {
    logger.info('[main-process] WordService.findWordPrecisly');

    const result = dictService.findWordPrecisly(
      dictid,
      keyText,
      rofset
    );

    // 如果查询不到精确结果，就返回空值
    if (!result) {
      throw new Error('cannot find workd precisly: ' + keyText);
    }
    return this.definitionReplace(dictid, result.keyText, result.definition);
  }

  /**
   * loadDictResource 加载词典额外资源
   * @param event 事件源
   * @param arg 资源 keyText
   */
  loadDictResource(dictid: string, resourceKey: string) {
    return dictService.loadDictResource(dictid, resourceKey);
  }
}
