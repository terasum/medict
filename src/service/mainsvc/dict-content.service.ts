import { WordDefinition } from 'js-mdict'
import cheerio from "cheerio";
import { dictService } from './main.service';
import { __RANDOM_KEY__ } from './../../renderer/utils/random_key';


const IMAGE_REG = /src=\"((\S+)\.(png|jpg|gif|jpeg|svg))\"/gi
const IMAGE_REG_IDX = 1

const SOUND_REG = /sound\:\/\/((\S+)\.(mp3|spx|ogg|wav))/gi
const SOUND_REG_IDX = 0

const CSS_REG = /href=\"((\S+)\.css)\"/gi
const CSS_REG_IDX = 1

const JS_REG = /src=\"((\S+)\.js)\"/gi
const JS_REG_IDX = 1

const LINK_REG = /^@@@LINK=([\s\S]+)$/i
const LINK_REG_IDX = 1

const ENTRY_REG = /href=\"entry:\/\/([\S\s]+?)\"/gi
const ENTRY_REG_IDX = 1

type ResourceFn = (
  key: string,
) => WordDefinition | { keyText: string; definition: null }

function replaceAll(str: string, find: string, replace: string) {
  return str.replace(new RegExp(find, 'g'), replace);
}


function replaceImages(dictid: string, keyText:string, html: string, resourceFn: ResourceFn): string {
  console.log('***************  REPLACE IMAGES [START] ****************');
  if (!html || !html.matchAll) {
    return html
  }

  let matches = html.matchAll(IMAGE_REG)
  const keySet = new Set<string>();
  for (const match of matches) {
    let rawKey = match[IMAGE_REG_IDX]
    keySet.add(rawKey);
    const keyStart = match.index
    const keyEnd = (match.index || 0) + match[IMAGE_REG_IDX].length
    console.log(`[REP IMG]: matched  raw key ${rawKey} start=${keyStart} end=${keyEnd}.`);
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
      console.log(`[REP IMG]: replace html rawkey ${rawkey} => ${queryResult.definition}`);
      html = replaceAll(html, rawkey, queryResult.definition)
    }
  }
  console.log('***************  REPLACE IMAGES [END] ****************');
  return html
}

function replaceSound(dictid: string, keyText:string, html: string, resourceFn: ResourceFn): string {
  console.log('***************  REPLACE AUDIO [START] ****************');
  if (!html || !html.matchAll) {
    return html
  }

  const keySet = new Set<string>();
  let matches = html.matchAll(SOUND_REG)
  for (const match of matches) {
    const resourceKey = match[SOUND_REG_IDX];
    keySet.add(resourceKey);
    const keyStart = match.index
    const keyEnd = (match.index || 0) + match[SOUND_REG_IDX].length
    console.log(`[REP MP3]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`);
  }

  for (let rskey of keySet) {
    const rawkey = rskey;
    let resourceKey = rskey;
    resourceKey = resourceKey.slice('sound:/'.length, resourceKey.length)
    resourceKey = replaceAll(resourceKey, '/', '\\');
    console.log(`[REP MP3]: query resource key ${resourceKey}`);
    const queryResult = resourceFn(resourceKey);
    console.log(`[REP MP3]: query resource result ${queryResult}`);

    if (queryResult && queryResult.definition) {
      console.log(`[REP MP3]: replace html mp3 rawkey ${rawkey} => ${queryResult.definition}`);
      html = replaceAll(html, rawkey, queryResult.definition)
    }
  }

  const isSupportExt = (name: string) => {
    if (!name) return false;
    const ext = name.split('.').pop();
    switch (ext) {
      case 'mp3': return true;
      case 'mp4': return true;
      case 'ogg': return true;
      case 'spx': return true;
      case 'wav': return true;
    }
    return false;
  }
    const aduioType = (name: string) => {
    if (!name) return 'audio/ogg';
    const ext = name.split('.').pop();
    switch (ext) {
      case 'mp3': return 'audio/mpeg';
      case 'mp4': return 'audio/mp4';
      case 'ogg': return 'audio/ogg';
      case 'spx': return 'audio/ogg';
      case 'wav': return 'audio/wav';
      case 'acc': return 'audio/acc';
    }
    return 'audio/ogg';
  }

  const $ = cheerio.load(html);
  const alist = $('a');
  for (let i = 0; i < alist.length; i++) {
    if (isSupportExt(alist[i].attribs.href)) {
      const el = $(alist[i]);
      el.append(`<audio id="__audio_${i}"><source src="${alist[i].attribs.href}" type="${aduioType(alist[i].attribs.href)}"><audio>`);
      el.attr('href', '#')
      $(alist[i]).attr('onclick', `function click__${i}(){let au = document.getElementById("__audio_${i}"); au.play();} click__${i}(); return false;`)
    }
  }

  console.log('***************  REPLACE AUDIO [END] ****************');
  return $.html();
}

function replaceCss(dictid: string, keyText:string, html: string, resourceFn: ResourceFn): string {
  /**
        Found oalecd8e.css start=39 end=51.
     */
  console.log('***************  REPLACE CSS [START] ****************');
  if (!html || !html.matchAll) {
    return html
  }
  let matches = html.matchAll(CSS_REG)
  const keySet = new Set<string>();


  for (const match of matches) {
    let resourceKey = match[SOUND_REG_IDX];
    resourceKey = resourceKey.slice('href="'.length, resourceKey.length);
    resourceKey = resourceKey.slice(0, resourceKey.length-1)

    keySet.add(resourceKey);
    const keyStart = match.index
    const keyEnd = (match.index || 0) + match[SOUND_REG_IDX].length
    console.log(`[REP CSS]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`);
  }

  for (let rskey of keySet) {
    const rawkey = rskey;
    let resourceKey = rskey;
    resourceKey = replaceAll(resourceKey, '/', '\\');
    console.log(`[REP CSS]: query resource key ${resourceKey}`);
    const queryResult = resourceFn(resourceKey);
    console.log(`[REP CSS]: query resource result ${queryResult}`);

    if (queryResult && queryResult.definition) {
      console.log(`[REP CSS]: replace html mp3 rawkey ${rawkey} => ${queryResult.definition}`);
      html = replaceAll(html, rawkey, queryResult.definition)
    }
  }
  console.log('***************  REPLACE CSS [END] ****************');
  return html
}

function replaceJs(dictid: string, keyText:string, html: string, resourceFn: ResourceFn): string {
  /**
        Found jquery.js start=67 end=76.
     */
  console.log('***************  REPLACE JS [START] ****************');
  if (!html || !html.matchAll) {
    return html
  }
  let matches = html.matchAll(JS_REG)
  const keySet = new Set<string>();
  for (const match of matches) {
    let resourceKey = match[SOUND_REG_IDX];
    resourceKey = resourceKey.slice('src="'.length, resourceKey.length);
    resourceKey = resourceKey.slice(0, resourceKey.length-1)
    keySet.add(resourceKey);
    const keyStart = match.index
    const keyEnd = (match.index || 0) + match[SOUND_REG_IDX].length
    console.log(`[REP JS]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`);
  }


  for (let rskey of keySet) {
    const rawkey = rskey;
    let resourceKey = rskey;
    resourceKey = replaceAll(resourceKey, '/', '\\');
    console.log(`[REP JS]: query resource key ${resourceKey}`);
    const queryResult = resourceFn(resourceKey);
    console.log(`[REP JS]: query resource result ${queryResult}`);

    if (queryResult && queryResult.definition) {
      console.log(`[REP JS]: replace html mp3 rawkey ${rawkey} => ${queryResult.definition}`);
      html = replaceAll(html, rawkey, queryResult.definition)
    }
  }

  console.log('***************  REPLACE JS [END] ****************');

  return html
}


function replaceLink(dictid: string, keyText:string, html: string, resourceFn: ResourceFn): string {
  /// @@@LINK=wordy
    console.log('***************  REPLACE @@@LINK [START] ****************');
  if (!html || !html.matchAll) {
    return html;
  }
  if (LINK_REG.test(html)) {
    let matches = html.match(LINK_REG);
    console.log(matches);
    if (matches == null||matches.length<2) {
      return html;
    }
    const oldWord = matches[LINK_REG_IDX];
    let newWord = oldWord;
    if (oldWord.endsWith('\x00')) {
      newWord = newWord.slice(0, newWord.length - '\x00'.length);
    }
    newWord = newWord.trimEnd();
    console.log(`REPLACE @@@LINK #${newWord}#`);
    console.log({ newWord });

    const result = dictService.lookup(dictid, newWord);
    if (!result) {
      return 'null';
    }    
    console.log('***************  REPLACE @@@LINK [END1] ****************');
    return result?.definition;
  }
  console.log('***************  REPLACE @@@LINK [END2] ****************');
  return html;
}


function replaceEntry(dictid: string, keyText:string, html: string, resourceFn: ResourceFn): string {
  /// <a href="entry://buzzword">buzzword</a>
    console.log('***************  REPLACE @@@ENTRY [START] ****************');
  if (!html || !html.matchAll) {
    return html;
  }
  // let matches = html.matchAll(ENTRY_REG)
  // const keySet = new Set<string>();
  // for (const match of matches) {
  //   let resourceKey = match[ENTRY_REG_IDX];
  //   // resourceKey = resourceKey.slice('href="'.length, resourceKey.length);
  //   // resourceKey = resourceKey.slice(0, resourceKey.length-1)
  //   keySet.add(resourceKey);
  //   const keyStart = match.index
  //   const keyEnd = (match.index || 0) + match[SOUND_REG_IDX].length
  //   console.log(`[REP ENTRY]: matched raw key ${resourceKey} start=${keyStart} end=${keyEnd}.`);
  // }

  // for (let rskey of keySet) {
  //   const rawkey = rskey;
  //   let resourceKey = rskey;
  //   console.log(`[REP ENTRY]: query resource key ${resourceKey}`);

  //   const queryResult = resourceFn(resourceKey);
  //   console.log(`[REP ENTRY]: query resource result ${queryResult}`);

  //   if (queryResult && queryResult.definition) {
  //     console.log(`[REP ENTRY]: replace html mp3 rawkey ${rawkey} => ${queryResult.definition}`);
  //     html = replaceAll(html, rawkey, queryResult.definition)
  //   }

    const $ = cheerio.load(html);
    const alist = $('a');

    const isSupportURL = (url: string) => {
      if (!url) return false;
      if (url.startsWith('entry://')) return true;
      return false;

  }

  for (let i = 0; i < alist.length; i++) {
    const href = alist[i].attribs.href;
      if (isSupportURL(href)) {
         
        const newWord =href.slice('entry://'.length, href.length)
        console.log(`entry url ${alist[i].attribs.href}, #${newWord}#`)

        const el = $(alist[i]);
        el.attr('href', '#');
        // for security purpose, add some random key
        // 点击链接将会把 entryLinkWord 事件发送到 main-process, 由 main-process 处理完成后
        // 再将结果返回 webview 页面,加入随机字符串是为了保证页面安全
        $(alist[i]).attr('onclick', `
        function entry_click__${i}__${__RANDOM_KEY__}() {
          console.log({ dictid: "${dictid}", word: "${newWord}" });
          window.postMessage({
            channel: "entryLinkWord",
            payload: {
              dictid: "${dictid}",
              word: "${newWord}",
            }
          });
        }
        entry_click__${i}__${__RANDOM_KEY__}();
        return false;
        `)
      }
    }
  
  console.log('***************  REPLACE @@@ENTRY [END] ****************');
  return $.html();
}

export const dictContentService = {
  definitionReplace: (dictid: string, keyText:string, html: string, fn: ResourceFn) => {
    html = replaceLink(dictid, keyText, html, fn)
    html = replaceImages(dictid, keyText, html, fn)
    html = replaceSound(dictid, keyText, html, fn)
    html = replaceCss(dictid, keyText, html, fn)
    html = replaceJs(dictid, keyText, html, fn)
    html = replaceEntry(dictid, keyText, html, fn)
    return html
  },
}
