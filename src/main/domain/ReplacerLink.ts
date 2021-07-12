import { Replacer, ResourceFn, LookupFn } from './Replacer';
import { logger } from '../../utils/logger';

import { __RANDOM_KEY__ } from '../../utils/random_key';
const LINK_REG = /^@@@LINK=([\s\S]+)$/i;
const LINK_REG_IDX = 1;

export class LinkReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): string {
    /// @@@LINK=wordy
    logger.info('[REP @@@LINK]: REPLACE @@@LINK [START]')
    if (!html || !html.matchAll) {
      return html;
    }
    if (LINK_REG.test(html)) {
      let matches = html.match(LINK_REG);
      logger.info({ matches });
      if (matches == null || matches.length < 2) {
        return html;
      }
      const oldWord = matches[LINK_REG_IDX];
      let newWord = oldWord;
      if (oldWord.endsWith('\x00')) {
        newWord = newWord.slice(0, newWord.length - '\x00'.length);
      }
      newWord = newWord.trimEnd();
      logger.info(`REPLACE @@@LINK #${newWord}#`);
      logger.info({ newWord });

      const result = lookupFn(newWord);
      if (!result) {
        return 'null';
      }
      logger.info('[REP @@@LINK]: REPLACE @@@LINK [END1]');
      return result?.definition;
    }
    logger.info('[REP @@@LINK]: REPLACE @@@LINK [END2]');
    return html;
  }
}
