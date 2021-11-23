import { logger } from '../../utils/logger';
import { WordQueryService } from '../mainsvc/WordQueryService';

/**
 * WordService 查词主要服务，依赖词典服务
 */
export class WordQueryApi {
  queryService: WordQueryService;

  constructor() {
    this.queryService = new WordQueryService();
  }

  async __asyncSendReturn(event: any, channel: string, data: any, error?: any) {
    if (!error) {
      event.sender.send(channel, { data: undefined, error: error });
    } else {
      event.sender.send(channel, { data: data, error: undefined });
    }
  }


  /**
   * entryLinkWord 异步查询 Entry:// 指向的词
   * @param event 事件源
   * @param arg 请求参数
   * @returns
   */
  entryLinkWord(event: any, arg: { keyText: string; dictid: string }) {
    logger.info('[main-process] WordService.entryLinkWord');
    logger.info(arg);
    try {
      // 先将建议词列表返回
      const result = this.queryService.suggestWord(arg.dictid, arg.keyText);
      if (result.length < 1) {
        throw new Error('unable to query any suggest word' + arg.keyText);
      }
      this.__asyncSendReturn(event, 'onSuggestWorkd', result);

      // 查询精确结果
      const wordResult = this.queryService.lookup(arg.dictid, arg.keyText);
      // 如果查询不到精确结果，说明没有该词，就返回空值
      // 注意，事件必须是 onFindWordPrecisly
      if (!wordResult) {
        throw new Error(`lookup word ${arg.keyText} failed, return null`);
      }
      // 传入函数的目的是不让 DictionaryContectService 反向依赖 DictionaryService
      // 注意，事件必须是 onFindWordPrecisly
      const replaceResult = this.queryService.definitionReplace(arg.dictid, wordResult.keyText, wordResult.definition);

      this.__asyncSendReturn(event, 'onFindWordPrecisly', {
        keyText: replaceResult.keyText,
        sourceKeyText: replaceResult.sourceKeyText,
        definition: replaceResult.definition,
      });
    } catch (error: any) {
      this.__asyncSendReturn(event, 'onSuggestWorkd', undefined, error);
    }

  }

  /**
   * suggestWord 查询建议词列表
   * @param event 事件源
   * @param arg 目标词
   * @return 返回关联词列表
   */
  suggestWord(event: any, arg: { dictid: string; word: string }) {
    try {
      const result = this.queryService.suggestWord(arg.dictid, arg.word);
      this.__asyncSendReturn(event, 'onSuggestWord', result);
    } catch (error) {
      this.__asyncSendReturn(event, 'onSuggestWord', undefined, error);
    }
  }

  /**
   * findWordPrecisly 精确查询某个词
   * @param event 事件源
   * @param arg 词详细信息
   */
  findWordPrecisly(
    event: any,
    arg: { keyText: string; dictid: string; rofset: number }
  ) {
    logger.info('[main-process] WordService.findWordPrecisly');
    try {
      const precislyResult = this.queryService.findWordPrecisly(arg.keyText, arg.dictid, arg.rofset);
      this.__asyncSendReturn(event, 'onFindWordPrecisly', {
        keyText: precislyResult.keyText,
        sourceKey: precislyResult.sourceKeyText,
        definition: precislyResult.definition,
      });
    } catch (error) {
      this.__asyncSendReturn(event, 'onFindWordPrecisly', undefined, error)
    }
  }

  /**
   * loadDictResource 加载词典额外资源
   * @param event 事件源
   * @param arg 资源 keyText
   */
  loadDictResource(event: any, arg: { dictid: string; resourceKey: string }) {
    try {
      const result = this.queryService.loadDictResource(arg.dictid, arg.resourceKey);
      this.__asyncSendReturn(event, 'onLoadDictResource', result);
    } catch (error) {
      this.__asyncSendReturn(event, 'onLoadDictResource', undefined, error);
    }
  }
}
