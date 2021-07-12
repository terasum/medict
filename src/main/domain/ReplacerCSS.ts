import { Replacer, ResourceFn, LookupFn } from './Replacer';
import { replaceAll } from '../../utils/stringutils';
import { logger } from '../../utils/logger';
const CSS_REG = /href=\"((\S+)\.css)\"/gi;
const CSS_REG_IDX = 1;

export function extractKeys(html: string) {
  let matches = html.matchAll(CSS_REG);
  const keySet = new Set<string>();

  for (const match of matches) {
    let resourceKey = match[CSS_REG_IDX];
    // resourceKey = resourceKey.slice('href="'.length, resourceKey.length);
    // resourceKey = resourceKey.slice(0, resourceKey.length - 1);

    keySet.add(resourceKey);
    const keyStart = match.index;
    const keyEnd = (match.index || 0) + match[CSS_REG_IDX].length;
    logger.info(
      `[REP CSS]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`
    );
  }
  return keySet;
}

export class CSSReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): string {
    /**
        Found oalecd8e.css start=39 end=51.
     */
    logger.info('[REP CSS]: REPLACE CSS [START]');
    if (!html || !html.matchAll) {
      return html;
    }
    const keySet = extractKeys(html);

    for (let rskey of keySet) {
      const rawkey = rskey;
      let resourceKey = rskey;
      resourceKey = replaceAll(resourceKey, '/', '\\');
      logger.info(`[REP CSS]: query resource key ${resourceKey}`);
      const queryResult = resourceFn(resourceKey);
      logger.info(queryResult, `[REP CSS]: query resource result`);

      if (queryResult && queryResult.definition) {
        logger.info(
          `[REP CSS]: replace css rawkey ${rawkey} => ${queryResult.definition}`
        );
        html = replaceAll(html, rawkey, queryResult.definition);
      }
    }
    logger.info('[REP CSS]: REPLACE CSS [END]');
    return html;
  }
}
