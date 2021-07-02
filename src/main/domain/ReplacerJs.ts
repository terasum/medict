import { Replacer, ResourceFn, LookupFn } from './Replacer';
import { replaceAll } from '../../utils/stringutils';

const JS_REG = /src=\"((\S+)\.js)\"/gi;
const JS_REG_IDX = 1;

export class JSReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): string {
    /**
        Found jquery.js start=67 end=76.
     */
    console.log('***************  REPLACE JS [START] ****************');
    if (!html || !html.matchAll) {
      return html;
    }
    let matches = html.matchAll(JS_REG);
    const keySet = new Set<string>();
    for (const match of matches) {
      let resourceKey = match[JS_REG_IDX];
      resourceKey = resourceKey.slice('src="'.length, resourceKey.length);
      resourceKey = resourceKey.slice(0, resourceKey.length - 1);
      keySet.add(resourceKey);
      const keyStart = match.index;
      const keyEnd = (match.index || 0) + match[JS_REG_IDX].length;
      console.log(
        `[REP JS]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`
      );
    }

    for (let rskey of keySet) {
      const rawkey = rskey;
      let resourceKey = rskey;
      resourceKey = replaceAll(resourceKey, '/', '\\');
      console.log(`[REP JS]: query resource key ${resourceKey}`);
      const queryResult = resourceFn(resourceKey);
      console.log(`[REP JS]: query resource result ${queryResult}`);

      if (queryResult && queryResult.definition) {
        console.log(
          `[REP JS]: replace html mp3 rawkey ${rawkey} => ${queryResult.definition}`
        );
        html = replaceAll(html, rawkey, queryResult.definition);
      }
    }

    console.log('***************  REPLACE JS [END] ****************');

    return html;
  }
}
