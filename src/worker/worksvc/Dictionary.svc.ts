import { Dictionary } from '../../infra/Dictionary';
import { SuggestItem } from '../../model/SuggestItem';
import { Definition, NullDef } from '../../model/Definition';
import { StorabeDictionary } from '../../model/StorableDictionary';
import { DictContentService } from '../worksvc/DictionaryContent.svc';
import { logger } from '../../utils/logger';
import { md5hash } from '../../utils/hash_utils';
import { EventEmitter } from 'keyv';

import _ from 'lodash';
import path from 'path';
import fs from 'fs';
import walk from 'walkdir';

export const EV_START_INDEXING = 'EV_START_INDEXING';
export const EV_END_INDEXING = 'EV_END_INDEXING';

/**
 * 词典类需要完成如下用例功能：
 * 1. 精确查询
 * 2. 模糊查询
 * 3. 分组查询
 * 4. 按目录扫描
 * 5. 按目录分组扫描
 * 6. 按照词典内容展示词典信息
 * 
 * 故词典类应当拥有三大实体: 1. 目录  2. 可序列化词典  3. 索引后内存词典
 * 目录(1) -> 可序列化词典(1)
 * 可序列化词典(1) ->  索引后内存词典(1)
 * 
 * 构建方式
 * 词典目录路径 -> 提供给 DictService -> 扫描 -> 索引 -> 查询
 */


export class DictService {
  // 已经载入的可持久化词典缓存
  dicts: Map<string, StorabeDictionary> = new Map();

  // 已经载入的已经索引的词典缓存
  dictCache: Map<string, Dictionary> = new Map();

  // 当前词典的扫描目录
  dictsBaseDir: string;

  // 当前词典的配置文件中存储的列表
  dictsConfigFile: string

  // 当前词典的资源根目录
  resourceRoot: string

  // 事件抛出器
  eventEmitter = new EventEmitter();

  // dictContentService 内容处理服务
  dictContentService: DictContentService

  constructor(resourceRootDir: string, dictsBaseDir: string, dictsConfigFile: string) {
    this.dictsBaseDir = dictsBaseDir;
    this.resourceRoot = resourceRootDir;
    this.dictsConfigFile = dictsConfigFile;
    this.eventEmitter = new EventEmitter();
    this.dictContentService = new DictContentService();
  }

  /**
   * ------------ 可持久化词典部分  -------------- 
   */

  // 从目录扫描，因此不需要增删
  /**
   * 查找一个可持久化词典实例
   * @param dictid 词典id
   * @returns 返回一个持久化词典实例
   */
  findOne(dictid: string) {
    return this.dicts.get(dictid);
  }

  /**
   * 返回所有可持久化词典列表
   * @returns 可持久化词典列表
   */
  findAll() {
    const list: StorabeDictionary[] = [];
    this.dicts.forEach((val) => {
      list.push(val);
    });
    return list;
  }



  /**
   * ------------ 可查询内存词典部分  -------------- 
   */
  listDicts() {
    const list: Dictionary[] = [];
    this.dictCache.forEach((dict) => {
      list.push(dict);
    });
    return list;
  }

  async autoIndexing() {
    console.log('[WORKER] auto indexing...')
    this.loadDictsByDir(this.dictsBaseDir);
    this.dicts.forEach((val) => {
      let dict = Dictionary.newByStorabe(val);
      // 构建索引
      new Promise((resolve) => {
        this.eventEmitter.emit(EV_START_INDEXING, dict.name)
        dict.indexing();
        this.dictCache.set(dict.id, dict);
        this.eventEmitter.emit(EV_END_INDEXING, dict.name)
      })
    })
  }

  // -------- new added methods --------

  /**
   * 加载词典
   * @param dictID 词典Id
   * @returns 返回缓存之后的词典
   */
  loadDict(dictID: string): Dictionary | undefined {
    if (!this.dictCache.has(dictID) && this.dicts.has(dictID)) {
      return undefined;
    }

    if (this.dictCache.has(dictID)) {
      return this.dictCache.get(dictID);
    }

    let dict = Dictionary.newByStorabe(this.dicts.get(dictID)!);
    // 构建索引
    new Promise((resolve) => {
      this.eventEmitter.emit(EV_START_INDEXING, dict.name)
      dict.indexing();
      this.eventEmitter.emit(EV_END_INDEXING, dict.name)
    })

    this.dictCache.set(dictID, dict);
    return dict;
  }

  /**
   * 精确查词
   * @param dictid 词典ID
   * @param keyText 词典词条
   * @param roffset 偏移量
   * @returns 
   */
  lookupPrecisly(dictid: string, keyText: string, roffset: number) {
    return this.loadDict(dictid)?.findWordDefinition(keyText, roffset);
  }

  /**
   * 查询关联词，返回精确位置, 如果索引建立完成，则通过索引查询，如果未建立完成，通过传统方法查询(有误差)
   * @param dictid 词典ID
   * @param keyText 词典词条
   * @returns 
   */
  lookup(dictid: string, keyText: string) {
    const wordDef = this.loadDict(dictid)?.lookup(keyText);
    if (!wordDef || !wordDef.definition) {
      return NullDef(keyText);
    }
    return wordDef as Definition;
  }

  /**
   * 加载词条资源
   * @param dictid 词典ID
   * @param keyText 词典资源词条
   * @param withPayload 是否返回结果, 或是仅仅返回是否存在
   * @returns 
   */
  loadResource(dictid: string, keyText: string, isWithPayload = false) {
    const wordDef = this.loadDict(dictid)?.findWordResource(keyText, isWithPayload);
    if (!wordDef || !wordDef.definition) {
      return NullDef(keyText);
    }
    return wordDef;
  }

/**
   * 处理原生HTML返回结果
   * @param dictid 词典ID
   * @param keyText 关键词
   * @param definition 原生返回结果
   * @returns 
   */
  definitionReplace(dictid: string, keyText: string, definition: string) {
    // 资源查询函数
    const resFn = (resKey: string, withPayload = false) => {
      return this.loadResource(dictid, resKey, withPayload);
    };

    // 词查询函数
    const lookupFn = (word: string) => {
      return this.lookup(dictid, word);
    };

    logger.info("[main-process] suggestWord event.sender.send('onFindWordPrecisly')");

    return this.dictContentService.definitionReplace(
      dictid,
      keyText,
      definition,
      lookupFn,
      resFn
    )
  }


  /**
   * 关联查询
   * @param dictid 词典ID
   * @param word 查询词
   * @returns 返回模糊查询结果
   */
  associate(dictid: string, word: string) {
    const result: SuggestItem[] = [];
    if (word.trim() == '' || word.length === 0) {
      return result;
    }

    // limits word result upto 50
    let counter = 0;
    const limit = 50;
    const dict = this.loadDict(dictid);
    if (!dict) {
      return [];
    }
    const words = dict.associate(word);
    for (let i = 0; i < words?.length ?? 0; i++) {
      if (counter >= limit) {
        break;
      }
      const word = words[i];
      // logger.info(`set ${key}, ${word.keyText}`)
      result.push({
        id: counter,
        dictid: word.dictid,
        keyText: word.keyText,
        rofset: word.rofset,
      });
      counter++;
    }
    return result;
  }

  /**
   *  ---------------- 词典加载部分 --------------
   * 
   */

  /**
   * 分组加载目录中的词典, 并返回可序列化词典
   * @returns 序列化词典
   */
  loadDictsByDir(dictsBaseDir: string) {
    // 扫描根目录
    // walk-through the directory, and copy css/js/font/png files
    console.info("[WORKER] dict scan root dict dir " + dictsBaseDir);
    walk(dictsBaseDir, { sync: true, no_recurse: true }, (fpath, stat) => {
      // 判断是否合法名称
      // if resource cache dir not exists this file, copy it
      const fileBasename = path.basename(fpath);
      if (!fileBasename && fileBasename == '') {
        return;
      }
      // TODO 如果当前文件夹找到某个mdx，则寻找同名的 mdd 文件

      console.info("[WORKER] dict scan file path " + fpath);
      if (!stat.isDirectory()) {
        return;
      }

      console.info("[WORKER] dictionary path " + fpath);

      let mdxflag = false;
      let mdxpath = "";
      let mddpath: string[] = [];
      let currDictFileName = "";

      // 扫描子文件夹
      walk(fpath, { sync: true, no_recurse: true }, (subfile, substat) => {
        console.info(`[WORKER] scan sub dictionary path, ${subfile}, ${path.extname(subfile)}, ${path.extname(subfile)}`);
        // 三级目录跳过
        if (substat.isDirectory()) {
          return;
        }
        // 如果是合法文件则将其确认并处理
        if (path.extname(subfile) === ".mdx") {
          mdxflag = true;
          mdxpath = subfile;
          currDictFileName = path.basename(subfile);
        }

        if (path.extname(subfile) === ".mdd") {
          mddpath.push(subfile);
        }

      });
      // 提前返回
      if (!mdxflag) {
        return;
      }
      // 如果是合法词典则缓存
      // 当前子目录名称
      let currDirName = path.basename(fpath);
      // 随机生成一个词典id
      let dictId = (md5hash(currDirName + '-' + currDictFileName) as string).substring(0, 12);
      // 随机生成一个词典名称， TODO 应该从词典描述中读取
      let dictName = currDictFileName.substring(0, currDictFileName.length - 4);
      // 打印当前词典的 hash
      console.info(`[WORDER] dicthash: ${dictId}`, mdxpath, mddpath, fpath);
      // 构造一个可序列化词典
      let newDict = new StorabeDictionary(dictId, dictName, dictName,
        mdxpath, mddpath, fpath, true);
      // 缓存该词典
      this.dicts.set(dictId, newDict);
      // 异步拷贝资源, TODO 不必每次都这么做
      new Promise((resolve => {
        this._copyResources(this.resourceRoot, fpath, dictId);
      }))

    });
  }


  async _copyResources(resourceRootDir: string, dictMdxDir: string, dictid: string) {
    let that = this;
    // walk-through the directory, and copy css/js/font/png files
    walk(dictMdxDir, (fpath, stat) => {
      // if resource cache dir not exists this file, copy it
      const fileBasename = path.basename(fpath);
      if (!fileBasename && fileBasename == '') {
        return;
      }

      if (stat.isDirectory()) {
        return;
      }

      if (fileBasename.startsWith('.')) {
        return;
      }

      // 如果是非法文件则 return
      if (['css', 'js', 'ttf', 'otf', 'png', 'jpg', 'gif'].every((val) => !fileBasename.endsWith(val))) {
        return;
      }
      // 判断目标文件夹是否存在
      const resDirPath = path.resolve(resourceRootDir, dictid)
      if(!fs.existsSync(resDirPath)) {
        logger.info(`[WORKER]resource-copy, mkdir ${resDirPath}`);
        fs.mkdirSync(resDirPath, {recursive: true})
      }

      // 获取拷贝目标路径
      const resFilePath = path.resolve(resDirPath, fileBasename);
      logger.info(`[WORKER]resource-copy, base:[${fileBasename}] source: [${fpath}] => dest ${resFilePath}`);

      if (fs.existsSync(resFilePath)) {
        logger.info(`[WORKER]resource-copy, base:[${fileBasename}] already existed, skipped`);
        return;
      }

      fs.copyFile(fpath, resFilePath, (err) => {
        if (err) {
          logger.error(`[WORKER]resource-copy, base:[${fileBasename}] failed, ${err}`);
          return;
        }
        logger.info(`[WORKER]resource-copy, base:[${fileBasename}] successed`);
      });
      // 针对 css 文件，需要解析器内部内容并检索查询内部




    });
  }


  // 通过配置文件加载 TODO
  loadDictsByConfig() {

  }

}


/**
 * ----------- TOOL METHODS ---------------
 */

/**
* 判断当前词典的mdx路径是否合法
* @param dict 词典结构
* @returns boolean
*/
function isMdxValid(dict: any) {
  // detects mdx path
  if (!dict.mdxpath || !fs.existsSync(dict.mdxpath)) {
    return false;
  }
  return true;
}

function isMddValid(dict: any) {
  if (typeof dict.mddpath === 'string') {
    return fs.existsSync(dict.mddpath)
  }

  // detect dict.mddpath
  if (dict.mddpath instanceof Array) {
    if (dict.mddpath.length == 0) {
      return true;
    }

    if (dict.mddpath.length > 0) {
      let flag = false;
      dict.mddpath.forEach((mddpath: string) => {
        if (flag) {
          return;
        }
        if (!mddpath || !fs.existsSync(mddpath)) {
          flag = true;
        }
      });
      // if flag == true, means there are some mddpath is invalid
      return !flag;
    }
  }

  if (dict.mddpath == undefined) {
    return true;
  }

  return false;

}

// detect mdd/mdx file exists or not
// if exists return true, else return false
function detectValid(dict: any) {
  return isMdxValid(dict) && isMddValid(dict);
}