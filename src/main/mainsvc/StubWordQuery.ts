import { DictService } from './svc/DictionaryService';
import { DictContentService } from './svc/DictionaryContentService';
import { logger } from '../../utils/logger';

const dictService = new DictService();
const dictContentService = new DictContentService();

/**
 * WordService 查词主要服务，依赖词典服务
 */
export class StubWordQuery {
  /**
   * entryLinkWord 异步查询 Entry:// 指向的词
   * @param event 事件源
   * @param arg 请求参数
   * @returns
   */
  entryLinkWord(event: any, arg: { keyText: string; dictid: string }) {
    logger.info('[main-process] WordService.entryLinkWord');
    logger.info(arg);
    const result = dictService.associate(arg.dictid, arg.keyText);
    logger.info('[main-process] WordService.entryLinkWord // result');
    logger.info(result);
    // 先将建议词列表返回
    event.sender.send('onSuggestWord', result);
    if (result.length < 1) {
      return;
    }
    // 建议词列表返回之后，再精确检索keyText

    logger.info(
      "[main-process] suggestWord event.sender.send('onFindWordPrecisly')"
    );

    const wordResult =  dictService.lookup(arg.dictid, arg.keyText);

    // const wordResult = dictService.findWordPrecisly(
    //   result[0].dictid,
    //   result[0].keyText,
    //   result[0].rofset
    // );

    // 如果查询不到精确结果，就返回空值
    // 注意，事件必须是 onFindWordPrecisly
    if (!wordResult) {
      event.sender.send('onFindWordPrecisly', {
        keyText: arg.keyText,
        definition: 'null',
      });
      return;
    }

    // 资源查询函数
    const resFn = (resKey: string, withPayload=false) => {
      return dictService.loadDictResource(arg.dictid, resKey, withPayload);
    };

    // 词查询函数
    const lookupFn = (word: string) => {
      return dictService.lookup(arg.dictid, word);
    };
    // 传入函数的目的是不让 DictionaryContectService 反向依赖 DictionaryService
    // 注意，事件必须是 onFindWordPrecisly
    const replaceResult = dictContentService.definitionReplace(
        arg.dictid,
        wordResult.keyText,
        wordResult.definition,
        lookupFn,
        resFn
      )
    event.sender.send('onFindWordPrecisly', {
      keyText: replaceResult.keyText,
      sourceKeyText: replaceResult.sourceKeyText,
      definition: replaceResult.definition,
    });
  }

  /**
   * suggestWord 查询建议词列表
   * @param event 事件源
   * @param arg 目标词
   * @return 返回关联词列表
   */
  suggestWord(event: any, arg: { dictid: string; word: string }) {
    logger.info(
      "[main-process] suggestWord event.sender.send('onSuggestWord')"
    );
    const result = dictService.associate(arg.dictid, arg.word);

    event.sender.send('onSuggestWord', result);
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
    const result = dictService.findWordPrecisly(
      arg.dictid,
      arg.keyText,
      arg.rofset
    );
    // 如果查询不到精确结果，就返回空值
    if (!result) {
      event.sender.send('onFindWordPrecisly', {
        keyText: arg.keyText,
        definition: 'null',
      });
      return;
    }

    // 资源查询函数
    const resFn = (resKey: string, withPayload=false) => {
      return dictService.loadDictResource(arg.dictid, resKey, withPayload);
    };

    // 词查询函数
    const lookupFn = (word: string) => {
      return dictService.lookup(arg.dictid, word);
    };
    logger.info(
      "[main-process] suggestWord event.sender.send('onFindWordPrecisly')"
    );

    const replaceResult = dictContentService.definitionReplace(
        arg.dictid,
        result.keyText,
        result.definition,
        lookupFn,
        resFn
    )

    event.sender.send('onFindWordPrecisly', {
      keyText: replaceResult.keyText,
      sourceKey: replaceResult.sourceKeyText,
      definition: replaceResult.definition,
    });
  }
  /**
   * loadDictResource 加载词典额外资源
   * @param event 事件源
   * @param arg 资源 keyText
   */
  loadDictResource(event: any, arg: { dictid: string; resourceKey: string }) {
    const result = dictService.loadDictResource(arg.dictid, arg.resourceKey);
    event.sender.send('onLoadDictResource', result);
  }
}
