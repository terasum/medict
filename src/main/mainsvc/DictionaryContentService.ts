import { LinkReplacer } from '../../infra/ReplacerLink';
import { ImageReplacer } from '../../infra/ReplacerImage';
import { SoundReplacer } from '../../infra/ReplacerSound';
import { CSSReplacer } from '../../infra/ReplacerCSS';
import { EntryReplacer } from '../../infra/ReplacerEntry';
import { JSReplacer } from '../../infra/ReplacerJs';
import { INIReplacer } from '../../infra/ReplacerIni';
import { logger } from '../../utils/logger';

import { ResourceFn, LookupFn } from '../../infra/Replacer';

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
