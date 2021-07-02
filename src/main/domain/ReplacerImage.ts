import { replaceAll } from '../../utils/stringutils';
import { Replacer, ResourceFn, LookupFn } from './Replacer';

const IMAGE_REG = /src=\"((\S+)\.(png|jpg|gif|jpeg|svg))\"/gi;
const IMAGE_REG_IDX = 1;

export class ImageReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): string {
    console.log('***************  REPLACE IMAGES [START] ****************');
    if (!html || !html.matchAll) {
      return html;
    }

    let matches = html.matchAll(IMAGE_REG);
    const keySet = new Set<string>();
    for (const match of matches) {
      let rawKey = match[IMAGE_REG_IDX];
      keySet.add(rawKey);
      const keyStart = match.index;
      const keyEnd = (match.index || 0) + match[IMAGE_REG_IDX].length;
      console.log(
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
      console.log(`[REP IMG]: query resource key ${resourceKey}`);
      const queryResult = resourceFn(resourceKey);
      console.log(`[REP IMG]: query resource return result ${queryResult}`);

      if (queryResult && queryResult.definition) {
        console.log(
          `[REP IMG]: replace html rawkey ${rawkey} => ${queryResult.definition}`
        );
        html = replaceAll(html, rawkey, queryResult.definition);
      }
    }
    console.log('***************  REPLACE IMAGES [END] ****************');
    return html;
  }
}
