import { Dictionary } from '../../infra/Dictionary';
import { SuggestItem } from '../../model/SuggestItem';
import { Definition, NullDef } from '../../model/Definition';
import { StorabeDictionary } from '../../model/StorableDictionary';
import path from 'path';
import fs from 'fs';
import rimraf from 'rimraf';
import walk from 'walkdir';
import { logger } from '../../utils/logger';
import Configuration from './Configuration.svc';
import { ConfigAccessService } from './ConfigAccess.svc';
import { md5hash } from '../../utils/hash_utils';
import { EventEmitter } from 'keyv';
import _ from 'lodash';

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


export class NewDictService {
  // 已经载入的词典缓存
  dicts: Map<string, StorabeDictionary> = new Map();

  // 当前词典的扫描目录
  dictsBaseDir: string;

  // 当前词典的配置文件中存储的列表
  dictsConfigFile: string

  // 当前词典的资源根目录
  resourceRoot: string

  constructor(resourceRootDir: string, dictsBaseDir: string, dictsConfigFile: string) {
    this.dictsBaseDir = dictsBaseDir;
    this.resourceRoot = resourceRootDir;
    this.dictsConfigFile = dictsConfigFile;
  }

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
   * 分组加载目录中的词典, 并返回可序列化词典
   * @returns 序列化词典
   */
  loadDictsByDir(dictsBaseDir: string) {
    // 扫描根目录
    // walk-through the directory, and copy css/js/font/png files
    walk(dictsBaseDir, { sync: true, no_recurse: true }, (fpath, stat) => {
      // 判断是否合法名称
      // if resource cache dir not exists this file, copy it
      const fileBasename = path.basename(fpath);
      if (!fileBasename && fileBasename == '') {
        return;
      }
      // TODO 如果当前文件夹找到某个mdx，则寻找同名的 mdd 文件

      console.debug("[WORKER] dict scan file path " + fpath);
      if (!stat.isDirectory()) {
        return;
      }

      console.debug("[WORKER] dictionary path " + fpath);

      let mdxflag = false;
      let mdxpath = "";
      let mddpath: string[] = [];
      let currDictFileName = "";

      // 扫描子文件夹
      walk(fpath, { sync: true, no_recurse: true }, (subfile, substat) => {
        console.debug(`[WORKER] scan sub dictionary path, ${subfile}, ${path.extname(subfile)}, ${path.extname(subfile)}`);
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
      console.debug(`[WORDER] dicthash: ${dictId}`, mdxpath, mddpath, fpath);
      // 构造一个可序列化词典
      let newDict = new StorabeDictionary(dictId, dictName, dictName,
        mdxpath, mddpath, fpath, true);
      // 缓存该词典
      this.dicts.set(dictId, newDict);
      // 异步拷贝资源, TODO 不必每次都这么做
      new Promise((resolve => {
        this._copyResources(this.resourceRoot, currDirName, dictId);
      }))

    });
  }


  _copyResources(resourceRootDir: string, dictMdxDir: string, dictid: string) {
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
      // 获取拷贝目标路径
      const resFilePath = path.resolve(resourceRootDir, dictid, fileBasename);
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
        logger.error(`[WORKER]resource-copy, base:[${fileBasename}] successed`);
      });
    });
  }


  // 通过配置文件加载 TODO
  loadDictsByConfig() {
    // logger.info('reload all dictionaries');
    // const dictLists = this.storage.getDataByKey('dicts') as any[];
    // if (!dictLists) {
    //   logger.info('there is no dictionaries in the configuration file');
    //   return;
    // }

    // dictLists.forEach((dict) => {
    //   logger.info(`loading dictionary: ${dict.name}`)
    //   // detect file exists or not
    //   // TODO 不用每次启动都检查
    //   if (!detectValid(dict)) {
    //     this.dicts.delete(dict.id);
    //     // this.saveToFile(this.dicts);
    //     logger.error(`file resource load error, mdxpath: ${dict.mdxpath}, mddpath: ${dict.mddpath}`);
    //     return;
    //   }
    //   this.dicts.set(dict.id, new StorabeDictionary(dict.id, dict.alias, dict.name, dict.mdxpath, dict.mddpath, dict.description));
    // });
  }


}


export class DictStorageService {
  dicts: Map<string, StorabeDictionary>;
  config: Configuration;
  storage: ConfigAccessService;

  constructor() {
    this.dicts = new Map<string, StorabeDictionary>();
    this.config = Configuration.newInstance();
    this.storage = new ConfigAccessService(this.config.getConfigJsonPath());
    this.loadDictsByConfig();
    this.loadDictsByDir();
    console.log('[WORKER] listing dictionaries:')
    console.log(this.dicts);
  }



  loadDictsByDir() {
    let dictBaseDir = this.storage.getDictBaseDir();
    console.log(`[WORKER] dict base directory: ${dictBaseDir}`)
    if (!dictBaseDir || dictBaseDir.length == 0) {
      return;
    }
    // walk-through the directory, and copy css/js/font/png files
    walk(dictBaseDir, { sync: true, no_recurse: true }, (fpath, stat) => {
      // if resource cache dir not exists this file, copy it
      const fileBasename = path.basename(fpath);
      if (!fileBasename && fileBasename == '') {
        return;
      }

      console.debug("[WORKER] fpath " + fpath);
      if (!stat.isDirectory()) {
        return;
      }

      console.debug("[WORKER] dictionary path " + fpath);
      let mdxflag = false;
      let mdxpath = "";
      let mddpath: string[] = [];
      let currDictFileName = "";
      walk(fpath, { sync: true, no_recurse: true }, (subfile, substat) => {
        console.debug("[WORKER] sub dictionary path " + subfile, path.extname(subfile));
        console.debug("[WORKER] sub dictionary ext " + path.extname(subfile));
        if (substat.isDirectory()) {
          return;
        }

        if (path.extname(subfile) === ".mdx") {
          mdxflag = true;
          mdxpath = subfile;
          currDictFileName = path.basename(subfile);
        }

        if (path.extname(subfile) === ".mdd") {
          mddpath.push(subfile);
        }
      });

      if (mdxflag) {
        // valid mdx dict
        let currDirName = path.basename(fpath);
        let dictId = (md5hash(currDirName + '-' + currDictFileName) as string).substring(0, 12);
        let dictName = currDictFileName.substring(0, currDictFileName.length - 4);
        console.debug(`[WORDER] dicthash: ${dictId}`, mdxpath, mddpath, fpath);
        let newDict = new StorabeDictionary(dictId,
          dictName,
          dictName,
          mdxpath,
          mddpath,
          fpath,
          true
        );
        this.dicts.set(dictId, newDict);
      }

    });
  }


  findOne(dictid: string) {
    return this.dicts.get(dictid);
  }

  findAll() {
    const list: StorabeDictionary[] = [];
    this.dicts.forEach((val) => {
      list.push(val);
    });
    return list;
  }

  addOne(dict: Dictionary) {
    if (this.dicts.has(dict.id)) {
      return false;
    }
    this.dicts.set(dict.id, dict);
    this.saveToFile(this.dicts);
    return true;
  }

  deleteOne(dictid: string) {
    if (this.dicts.has(dictid)) {
      let dict = this.dicts.get(dictid);
      if (dict?.byScanning) {
        return;
      }
    }
    this.dicts.delete(dictid);
    this.saveToFile(this.dicts);
    const fpath = path.resolve(this.config.getResourceRootPath(), dictid);
    if (fs.existsSync(fpath)) {
      rimraf(fpath, function () {
        console.log('delete directory done', fpath);
      });
    }
    return true;
  }

  copyResources(mdxDir: string, dictid: string) {
    let that = this;
    // walk-through the directory, and copy css/js/font/png files
    walk(mdxDir, (fpath, stat) => {
      // if resource cache dir not exists this file, copy it
      const fileBasename = path.basename(fpath);
      const resourceFilePath = path.resolve(this.config.getResourceRootPath(), dictid, fileBasename);

      if (!fileBasename && fileBasename == '') {
        return;
      }

      if (
        !fileBasename.endsWith('css') &&
        !fileBasename.endsWith('js') &&
        !fileBasename.endsWith('ttf') &&
        !fileBasename.endsWith('otf') &&
        !fileBasename.endsWith('png') &&
        !fileBasename.endsWith('jpg')
      ) {
        return;
      }

      logger.info(`[RES-DETECT] base:[${fileBasename}] source: [${fpath}]`);
      logger.info(
        `[RES-DETECT] base:[${fileBasename}] dest:   [${resourceFilePath}]`
      );
      if (fs.existsSync(resourceFilePath)) {
        logger.info(`[RES-DETECT] base:[${fileBasename}] exists skipped`);
        return;
      }

      fs.copyFile(fpath, resourceFilePath, (err) => {
        if (err) {
          logger.error(
            `[RES-DETECT] base:[${fileBasename}] copy file failed, ${err}`
          );
          return;
        }
        logger.info(`[RES-DETECT] base:[${fileBasename}] copy file success`);
      });
    });
  }

  // refresh dictionaries
  loadDictsByConfig() {
    logger.info('reload all dictionaries');
    const dictLists = this.storage.getDataByKey('dicts') as any[];
    if (!dictLists) {
      logger.info('there is no dictionaries in the configuration file');
      return;
    }

    dictLists.forEach((dict) => {
      logger.info(`loading dictionary: ${dict.name}`)
      // detect file exists or not
      // TODO 不用每次启动都检查
      if (!detectValid(dict)) {
        this.dicts.delete(dict.id);
        // this.saveToFile(this.dicts);
        logger.error(`file resource load error, mdxpath: ${dict.mdxpath}, mddpath: ${dict.mddpath}`);
        return;
      }
      this.dicts.set(dict.id, new StorabeDictionary(dict.id, dict.alias, dict.name, dict.mdxpath, dict.mddpath, dict.description));
    });
  }

  prehandle() {
    // 预处理
    // TODO 应该在添加词典的时候处理
    // const fpath = path.resolve(this.config.getResourceRootPath(), dict.id);

    // if (!fs.existsSync(fpath)) {
    //   fs.mkdir(fpath, () => {
    //     logger.info(`dictionary resource path created ${fpath}`)
    //   });
    // }

    // copy css/js/fonts files from mdx source directory
    // remender: if the mdx file's directory contains a lot of
    // css/js/fonts, this may cause performance issue
    // const mdxPath = dict.mdxpath;
    // const mdxDir = path.dirname(mdxPath);

    // TODO 不用每次启动都拷贝资源
    // copy resource from mdx.directionary to resource.cache.dictionary
    // this.copyResources(mdxDir, dict.id);

  }

  saveToFile(dicts: Map<string, StorabeDictionary>) {
    this.storage.setDataByKey('dicts', dicts);
  }

}

export class DictService {
  storage: DictStorageService;
  dictCache: Map<string, Dictionary>;
  eventEmitter: EventEmitter;

  constructor() {
    this.storage = new DictStorageService()
    this.dictCache = new Map<string, Dictionary>();
    this.eventEmitter = new EventEmitter();
  }
  // -------- new added methods --------

  async buildIndex() {
    const dictList = this.storage.findAll();
    dictList.forEach((dict) => {
      console.log(dict);
      const newDict = Dictionary.newByStorabe(dict);
      this.eventEmitter.emit('ev-indexing-start', dict.name);
      const wordList = newDict.indexing();
      this.eventEmitter.emit('ev-indexing-end', dict.name);
      console.log(wordList);
    });
  }

  listDicts() {
    return this.storage.findAll();
  }

  // -------- new added methods --------

  /**
   * 加载词典
   * @param dictID 词典Id
   * @returns 返回缓存之后的词典
   */
  loadDict(dictID: string): Dictionary | undefined {
    if (!this.dictCache.has(dictID) && !this.storage.dicts.has(dictID)) {
      return undefined;
    }
    if (this.dictCache.has(dictID)) {
      return this.dictCache.get(dictID);
    }
    let dict = Dictionary.newByStorabe(this.storage.dicts.get(dictID)!);
    this.dictCache.set(dictID, dict);
    return dict;
  }


  /**
   * 精确查词
   * @param dictid 词典ID
   * @param keyText 词典词条
   * @param rofset 偏移量
   * @returns 
   */
  findWordPrecisly(dictid: string, keyText: string, rofset: number) {
    logger.debug(
      `findWordPrecisly, dict: [${dictid}] keyText: ${keyText}, roffset: ${rofset}`
    );
    return this.loadDict(dictid)?.findWordDefinition(keyText, rofset);
  }


  /**
   * 加载词条
   * @param dictid 词典ID
   * @param keyText 词典资源词条
   * @param withPayload 是否返回结果, 或是仅仅返回是否存在
   * @returns 
   */
  loadDictResource(dictid: string, keyText: string, withPayload = false) {
    const wordDef = this.loadDict(dictid)?.findWordResource(keyText, withPayload);
    if (!wordDef || !wordDef.definition) {
      return NullDef(keyText);
    }
    return wordDef;
  }


  /**
   * 查询关联词，返回精确位置
   * @param dictid 词典ID
   * @param keyText 词典词条
   * @returns 
   */
  lookup(dictid: string, keyText: string) {
    // return dicts.get(dictid)?.lookup(keyText) ?? NullDef(keyText);

    const wordDef = this.loadDict(dictid)?.lookup(keyText);
    if (!wordDef || !wordDef.definition) {
      return NullDef(keyText);
    }
    return wordDef as Definition;
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

}



/*******************
 * TOOL METHODS
 * *****************/

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