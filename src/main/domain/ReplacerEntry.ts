import { Replacer, ResourceFn, LookupFn } from './Replacer';
import cheerio from 'cheerio';
import { __RANDOM_KEY__ } from '../../utils/random_key';
import { logger } from '../../utils/logger';
const ENTRY_REG = /href=\"entry:\/\/([\S\s]+?)\"/gi;
const ENTRY_REG_IDX = 1;

export class EntryReplacer implements Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): {keyText:string, definition:string} {
    /// <a href="entry://buzzword">buzzword</a>
    logger.info('[REP @@@ENTRY]: REPLACE @@@ENTRY [START]');
    if (!html || !html.matchAll) {
      return {keyText, definition:html}
    }
    const $ = cheerio.load(html);
    const alist = $('a');

    const isSupportURL = (url: string) => {
      if (!url) return false;
      if (url.startsWith('entry://')) return true;
      return false;
    };

    for (let i = 0; i < alist.length; i++) {
      const href = alist[i].attribs.href;
      if (isSupportURL(href)) {
        let newWord = href.slice('entry://'.length, href.length);

        if (newWord && newWord.indexOf('#') > 0) {
          newWord = newWord.substr(0, newWord.indexOf('#'));
        }
        logger.info(`entry url ${alist[i].attribs.href}, #${newWord}#`);

        const el = $(alist[i]);
        el.attr('href', '#');
        // for security purpose, add some random key
        // 点击链接将会把 entryLinkWord 事件发送到 main-process, 由 main-process 处理完成后
        // 再将结果返回 webview 页面,加入随机字符串是为了保证页面安全
        $(alist[i]).attr(
          'onclick',
          `
        function entry_click__${i}__${__RANDOM_KEY__}() {
          console.log({ dictid: "${dictid}", word: "${newWord}" });
          window.postMessage({
            channel: "entryLinkWord",
            payload: {
              dictid: "${dictid}",
              keyText: "${newWord}",
            }
          });
        }
        entry_click__${i}__${__RANDOM_KEY__}();
        return false;
        `
        );
      }
    }

    logger.info('[REP @@@ENTRY]: REPLACE @@@ENTRY [END]');

    return {keyText, definition:$.html()};
  }
}
