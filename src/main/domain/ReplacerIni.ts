import { Replacer, ResourceFn, LookupFn } from './Replacer';
import { replaceAll } from '../../utils/stringutils';
import { logger } from '../../utils/logger';
import cheerio from 'cheerio';
import ini from 'ini';

const INI_REG = /src=\"((\S+)\.ini)\"/gi;
const INI_REG_IDX = 1;

export function extractKeys(html: string) {
  let matches = html.matchAll(INI_REG);
  const keySet = new Set<string>();
  for (const match of matches) {
    let resourceKey = match[INI_REG_IDX];
    keySet.add(resourceKey);
    const keyStart = match.index;
    const keyEnd = (match.index || 0) + match[INI_REG_IDX].length;
    logger.info(
      `[REP INI]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`
    );
  }
  return keySet;
}

export class INIReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): {keyText:string, definition: string} {
    /**
        Found jquery.js start=67 end=76.
     */
    logger.info('[REP INI]: REPLACE  INI [START]');
    if (!html || !html.matchAll) {
      return {keyText, definition:html}
    }

    const keySet = extractKeys(html);

    for (let rskey of keySet) {
      const rawkey = rskey;
      let resourceKey = rskey;
      if (!resourceKey.startsWith('\\') && !resourceKey.startsWith('/')) {
        resourceKey = '\\' + resourceKey;
      }
      resourceKey = replaceAll(resourceKey, '/', '\\');
      logger.info(`[REP INI]: query resource key ${resourceKey}`);
      const queryResult = resourceFn(resourceKey, true);
      logger.info(queryResult, `[REP INI]: query resource result`);
      // IIFE
      let prependJs = `<script type="text/javascript">(function(){`;

      if (queryResult && queryResult.definition) {
        logger.info("[REP INI]: replace INI rawkey", rawkey , queryResult.definition);
        if (queryResult.payload && queryResult.payload.length>0) {
          const iniConfigContent = Buffer.from(queryResult.payload,'base64').toString('utf-8'); 
          console.log(iniConfigContent);
          const iniConfig = ini.parse(iniConfigContent);
          
          for(const iniKey in iniConfig) {
            // skip comment
            if(iniKey && iniKey.trimStart().startsWith('//')){
              continue
            }
            if(iniKey && iniKey.trimStart().startsWith('#')){
              continue
            }
            if(iniKey && iniKey.trimStart().startsWith(';')){
              continue
            }
            prependJs += `if(!window.${iniKey}){ window.${iniKey} = ${iniConfig[iniKey]};}\n`
          }

        }

        // html = replaceAll(html, rawkey, queryResult.definition);
        html = replaceAll(html, rawkey, '#?n=' + rawkey);

      }

      prependJs += `;})(); </script>`
      // prepend 
      const $ = cheerio.load(html);
      const head = $('head');
      head.prepend(prependJs);
      html = $.html();
    }

    logger.info('[REP INI]: REPLACE INI [END]');

    return {keyText, definition:html}
  }
}
