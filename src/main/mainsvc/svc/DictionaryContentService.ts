import { LinkReplacer } from '../../domain/ReplacerLink';
import { ImageReplacer } from '../../domain/ReplacerImage';
import { SoundReplacer } from '../../domain/ReplacerSound';
import { CSSReplacer } from '../../domain/ReplacerCSS';
import { EntryReplacer } from '../../domain/ReplacerEntry';
import { JSReplacer } from '../../domain/ReplacerJs';
import { logger } from '../../../utils/logger';

import { ResourceFn, LookupFn } from '../../domain/Replacer';

const replacerChain = [
  new LinkReplacer(),
  new ImageReplacer(),
  new SoundReplacer(),
  new CSSReplacer(),
  new JSReplacer(),
  new EntryReplacer(),
];

export class DictContentService {
  definitionReplace(
    dictid: string,
    originKeyText: string,
    originHtml: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ) {
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
