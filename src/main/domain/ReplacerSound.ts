import { Replacer, ResourceFn, LookupFn } from './Replacer';
import cheerio from 'cheerio';
import { replaceAll } from '../../utils/stringutils';
import { logger } from '../../utils/logger';

const SOUND_REG = /sound\:\/\/((\S+)\.(mp3|spx|ogg|wav))/gi;
const SOUND_REG_IDX = 0;

export class SoundReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): {keyText: string, definition: string} {
    logger.info('[REP MP3]: REPLACE AUDIO [START]');
    if (!html || !html.matchAll) {
     return {keyText, definition: html}
    }

    const keySet = new Set<string>();
    let matches = html.matchAll(SOUND_REG);
    for (const match of matches) {
      const resourceKey = match[SOUND_REG_IDX];
      keySet.add(resourceKey);
      const keyStart = match.index;
      const keyEnd = (match.index || 0) + match[SOUND_REG_IDX].length;
      logger.info(
        `[REP MP3]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`
      );
    }

    for (let rskey of keySet) {
      const rawkey = rskey;
      let resourceKey = rskey;
      resourceKey = resourceKey.slice('sound:/'.length, resourceKey.length);
      resourceKey = replaceAll(resourceKey, '/', '\\');
      logger.info(`[REP MP3]: query resource key ${resourceKey}`);
      const queryResult = resourceFn(resourceKey);
      logger.info(
        queryResult, `[REP MP3]: query resource result`
      );

      if (queryResult && queryResult.definition) {
        logger.info(
          `[REP MP3]: replace html mp3 rawkey ${rawkey} => ${queryResult.definition}`
        );
        html = replaceAll(html, rawkey, queryResult.definition);
      }
    }

    const isSupportExt = (name: string) => {
      if (!name) return false;
      const ext = name.split('.').pop();
      switch (ext) {
        case 'mp3':
          return true;
        case 'mp4':
          return true;
        case 'ogg':
          return true;
        case 'spx':
          return true;
        case 'wav':
          return true;
      }
      return false;
    };
    const aduioType = (name: string) => {
      if (!name) return 'audio/ogg';
      const ext = name.split('.').pop();
      switch (ext) {
        case 'mp3':
          return 'audio/mpeg';
        case 'mp4':
          return 'audio/mp4';
        case 'ogg':
          return 'audio/ogg';
        case 'spx':
          return 'audio/ogg';
        case 'wav':
          return 'audio/wav';
        case 'acc':
          return 'audio/acc';
      }
      return 'audio/ogg';
    };

    const $ = cheerio.load(html);
    // $('head').append(
    //   `<meta http-equiv="Content-Security-Policy" content="script-src 'self' 'unsafe-inline' https://localhost:4000" >`
    // );
    const alist = $('a');
    for (let i = 0; i < alist.length; i++) {
      if (isSupportExt(alist[i].attribs.href)) {
        const el = $(alist[i]);
        el.append(
          `<audio id="__audio_${i}"><source src="${
            alist[i].attribs.href
          }" type="${aduioType(alist[i].attribs.href)}"><audio>`
        );
        el.attr('href', '#');
        $(alist[i]).attr(
          'onclick',
          `function click__${i}(){let au = document.getElementById("__audio_${i}"); au.play();} click__${i}(); return false;`
        );
      }
    }

    logger.info('[REP MP3]: EPLACE AUDIO [END]');
     return {keyText, definition:$.html()}
  }
}
