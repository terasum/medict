import { Replacer, ResourceFn, LookupFn } from './Replacer';
import { replaceAll } from '../utils/stringutils';
import { logger } from '../utils/logger';
import axios from 'axios';
import fs from 'fs';
import path from 'path';
import { SyncMainAPI } from '../worker/worksvc/worker.main.svc.manifest';



const CSS_REG = /href=\"((\S+)\.css)\"/gi;
const CSS_REG_IDX = 1;

const URL_REG = /url\(['|"|]?((\S+)\.(png|jpg|gif|jpeg|svg|woff|ttf|otf))['|"|]?\)/gi;
const URL_REG_IDX = 1;

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

function __closureResourceRoot() {
  let _resourceRoot = '';
  return function() {
    if(_resourceRoot === '') {
      _resourceRoot = SyncMainAPI.syncGetResourceRootPath();
    }
    return _resourceRoot
  }
}

export class CSSReplacer implements Replacer {
  _resourceRoot =  __closureResourceRoot();

  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ): { keyText: string, definition: string } {

    /**
        Found oalecd8e.css start=39 end=51.
     */
    logger.info('[REP CSS]: REPLACE CSS [START]');
    if (!html || !html.matchAll) {
      return { keyText, definition: html }
    }
    const keySet = extractKeys(html);

    for (let rskey of keySet) {
      const rawkey = rskey;
      let resourceKey = rskey;
      if (!resourceKey.startsWith('\\') && !resourceKey.startsWith('/')) {
        resourceKey = '\\' + resourceKey;
      }
      resourceKey = replaceAll(resourceKey, '/', '\\');
      const queryResult = resourceFn(resourceKey);
      logger.info(queryResult, `[REP CSS]: query resource result`);

      if (queryResult && queryResult.definition) {
        logger.info(
          `[REP CSS]: replace css rawkey ${rawkey} => ${queryResult.definition}`
        );
        // rawkey: \\dict.css
        // queryResult.definition like http://localhost:40000/19aed952d621/dict.css
        html = replaceAll(html, rawkey, queryResult.definition);
        // async query resource and copy to dest directory
        // 1. 请求获得 css 文件
        axios.get(queryResult.definition).then((result) => {
          if (result.status !== 200) {
            console.log('[REP CSS] css request content failed, try to load from mdd file');
            // TODO load from mdd file
            return;
          }
          // 2. 处理分析 css 文件并把内容进行下载替换
          const cssContent = result.data;
          this.cssContentHandle(dictid, cssContent, resourceFn);
        })
      }
    }
    logger.info('[REP CSS]: REPLACE CSS [END]');
    return { keyText, definition: html }

  }


  async cssContentHandle(dictid: string, cssContent: string, resourceFn: ResourceFn) {
    const queryKeys = this.extractCSSURL(cssContent);
    for (let qkey of queryKeys) {

      // if qkey starts with: //dlweb.sogoucdn.com/
      if (qkey.startsWith('//') || qkey.startsWith('http')) {
        // donot request the online resource
        logger.info('[REP CSS]: EXTRACT CSS CONTENT, skipped online resource:', qkey);
        continue;
      }

      // write content into resource folder
      const targetFilePath = path.resolve(this._resourceRoot(), dictid, qkey);
      if (fs.existsSync(targetFilePath)) {
        logger.info('[REP CSS]: EXTRACT CSS CONTENT, exists, skipped:', qkey, targetFilePath);
      }
      const queryResource = resourceFn(qkey, true, true);
      if (!queryResource.payload) {
        logger.info('[REP CSS]: EXTRACT CSS CONTENT, payload not found:', qkey,);
        continue;
      }

      fs.writeFile(targetFilePath, Buffer.from(queryResource.payload, 'base64'), (err) => {
        if (err) {
          logger.info('[REP CSS]: EXTRACT CSS CONTENT, write failed:', qkey, err);
          return;
        }
      })
    }
  }

  extractCSSURL(html: string) {
    let matches = html.matchAll(URL_REG);
    const rawkeys: string[] = [];
    for (const match of matches) {
      if (match.length > URL_REG_IDX) {
        const cssQueryKey = match[URL_REG_IDX];
        rawkeys.push(cssQueryKey);
      }
    }
    return rawkeys;
  }
}
