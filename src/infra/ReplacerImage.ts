import { replaceAll } from '../utils/stringutils';
import { Replacer, ResourceFn, LookupFn } from './Replacer';
import { logger } from '../utils/logger';

const IMAGE_REG = /src=\"((\S+)\.(png|jpg|gif|jpeg|svg))\"/gi;
const IMAGE_REG_IDX = 1;

export class ImageReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): {keyText: string, definition:string} {
    logger.info('[REP IMG]:REPLACE IMAGES [START]');
    if (!html || !html.matchAll) {
      return  {keyText, definition:html};
    }

    let matches = html.matchAll(IMAGE_REG);
    const keySet = new Set<string>();
    for (const match of matches) {
      let rawKey = match[IMAGE_REG_IDX];
      keySet.add(rawKey);
      const keyStart = match.index;
      const keyEnd = (match.index || 0) + match[IMAGE_REG_IDX].length;
      logger.info(
        `[REP IMG]: matched  raw key ${rawKey} start=${keyStart} end=${keyEnd}.`
      );
    }

    for (let rskey of keySet) {
      const rawkey = rskey;
      let resourceKey = rskey;
      if (!resourceKey.startsWith('\\') && !resourceKey.startsWith('/')) {
        resourceKey = '\\' + resourceKey;
      }

      resourceKey = replaceAll(resourceKey, '/', '\\');
      logger.info(`[REP IMG]: query resource key ${resourceKey}`);
      const queryResult = resourceFn(resourceKey);
      logger.info(queryResult, `[REP IMG]: query resource return result:`);

      if (queryResult && queryResult.definition) {
        logger.info(
          `[REP IMG]: replace html rawkey ${rawkey} => ${queryResult.definition}`
        );
        html = replaceAll(html, rawkey, queryResult.definition);
      }
    }
    logger.info('[REP IMG]: REPLACE IMAGES [END]');
    return  {keyText, definition:html};
  }
}
