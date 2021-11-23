import { LinkReplacer } from '../../infra/ReplacerLink';
import { ImageReplacer } from '../../infra/ReplacerImage';
import { SoundReplacer } from '../../infra/ReplacerSound';
import { CSSReplacer } from '../../infra/ReplacerCSS';
import { EntryReplacer } from '../../infra/ReplacerEntry';
import { JSReplacer } from '../../infra/ReplacerJs';
import { INIReplacer } from '../../infra/ReplacerIni';
import { ResourceFn, LookupFn } from '../../infra/Replacer';

import { logger } from '../../utils/logger';


const replacerChain = [
  new LinkReplacer(),
  new ImageReplacer(),
  new INIReplacer(),
  new SoundReplacer(),
  new CSSReplacer(),
  new JSReplacer(),
  new EntryReplacer(),
];

export class DictContentService {
  /**
   * 
   * 替换词典中的 HTML 内容，用于内容增强
   * @param dictid 词典 ID 用于搜索资源 
   * @param originKeyText 原始键值对
   * @param originHtml  原始HTML 内容
   * @param lookupFn 词典搜索函数，依赖注入
   * @param resourceFn  资源搜索函数，依赖注入
   * @return {sourceKeyText: originKeyText, keyText, definition};
   */
  definitionReplace(dictid: string, originKeyText: string, originHtml: string, lookupFn: LookupFn, resourceFn: ResourceFn) {
    let keyText = originKeyText;
    let definition = originHtml;

    replacerChain.forEach(replacer => {
      let result = replacer.replace(dictid, keyText, definition, lookupFn, resourceFn);
      keyText = result.keyText;
      definition  = result.definition;
    });

    logger.debug(`replace end, sourcekey: ${originKeyText} newkey: ${keyText}`)
    return {sourceKeyText: originKeyText, keyText, definition};
  }
}
