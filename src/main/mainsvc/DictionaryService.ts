import { Dictionary } from '../../infra/Dictionary';
import { SuggestItem } from '../../model/SuggestItem';
import { Definition, NullDef } from '../../model/Definition';
import { StorabeDictionary } from '../../model/StorableDictionary';
import { StorageService } from './StorageServcice';
import { getConfigJsonPath, getResourceRootPath } from '../../config/config';
import path from 'path';
import fs from 'fs';
import rimraf from 'rimraf';
import walk from 'walkdir';
import { logger } from '../../utils/logger';

const storageService = new StorageService(getConfigJsonPath());

const dicts = new Map<string, Dictionary>();

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
    if (dict.mddpath.length == 0){
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

  if(dict.mddpath == undefined){
    return true;
  }

  return false;

}

// detect mdd/mdx file exists or not
// if exists return true, else return false
function detectValid(dict: any) {
  return isMdxValid(dict) && isMddValid(dict);
}

function copyResources(mdxDir: string, dictid: string) {
  // walk-through the directory, and copy css/js/font/png files
  walk(mdxDir, function (fpath, stat) {
    // if resource cache dir not exists this file, copy it
    const fileBasename = path.basename(fpath);
    const resourceFilePath = path.resolve(
      getResourceRootPath(),
      dictid,
      fileBasename
    );

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
function loadDicts() {
  logger.info('reload all dictionaries');
  const dictLists = storageService.getDataByKey('dicts') as any[];

  if (dictLists) {
    dictLists.forEach((dict) => {
      // detect file exists or not
      if (!detectValid(dict)) {
        dicts.delete(dict.id);
        saveToFile(dicts);
        logger.error(`file resource load error, mdxpath: ${dict.mdxpath}, mddpath: ${dict.mddpath}`);
        return;
      }

      dicts.set(
        dict.id,
        new Dictionary(
          dict.id,
          dict.alias,
          dict.name,
          dict.mdxpath,
          dict.mddpath,
          dict.description
        )
      );

      const fpath = path.resolve(getResourceRootPath(), dict.id);
      if (!fs.existsSync(fpath)) {
        fs.mkdirSync(fpath);
      }

      // copy css/js/fonts files from mdx source directory
      // remender: if the mdx file's directory contains a lot of
      // css/js/fonts, this may cause performance issue
      const mdxPath = dict.mdxpath;
      const mdxDir = path.dirname(mdxPath);
      // copy resource from mdx.directionary to resource.cache.dictionary
      copyResources(mdxDir, dict.id);
    });
  }
}

// init load, this function will load during startup process (main-process)
// so, if load dictionary failed, it will block main-process
// we should try-catch the errors
// synchronized load dictionaries  process

// new Promise((resolve) => {
//   let dict = loadDicts();
//   resolve(dict);
// }).then(() =>{
//   logger.info("dictionary reloaded")
// }).catch(error =>{
//   logger.error(`raise error durring load dictionaries, ${error}`);
// })


function saveToFile(dicts: Map<string, Dictionary>) {
  const storageList = [];
  for (let redict of dicts.values()) {
    storageList.push({
      id: redict.id,
      alias: redict.alias,
      name: redict.name,
      mdxpath: redict.mdxpath,
      mddpath: redict.mddpath,
      description: redict.description,
      resourceBaseDir: redict.resourceBaseDir,
    } as StorabeDictionary);
  }

  storageService.setDataByKey('dicts', storageList);
}

export class DictService {
  findOne(dictid: string) {
    return dicts.get(dictid);
  }

  findAll() {
  
    const list: StorabeDictionary[] = [];
    dicts.forEach((val) => {
      list.push(val);
    });
    return list;
  }

  addOne(dict: Dictionary) {
    if (dicts.has(dict.id)) {
      return false;
    }
    dicts.set(dict.id, dict);
    saveToFile(dicts);
    loadDicts();
    return true;
  }

  deleteOne(dictid: string, reload: boolean = true) {
    dicts.delete(dictid);
    saveToFile(dicts);
    if (reload) {
      loadDicts();
    }
    const fpath = path.resolve(getResourceRootPath(), dictid);
    if (fs.existsSync(fpath)) {
      rimraf(fpath, function () {
        console.log('delete directory done', fpath);
      });
    }
    return true;
  }

  findWordPrecisly(dictid: string, keyText: string, rofset: number) {
    logger.debug(
      `findWordPrecisly, dict: [${dictid}] keyText: ${keyText}, roffset: ${rofset}`
    );
    return dicts.get(dictid)?.findWordDefinition(keyText, rofset);
  }

  loadDictResource(dictid: string, keyText: string, withPayload=false) {
    const wordDef = dicts.get(dictid)?.findWordResource(keyText, withPayload);
    if (!wordDef || !wordDef.definition) {
      return NullDef(keyText);
    }
    return wordDef;
  }

  lookup(dictid: string, keyText: string) {
    // return dicts.get(dictid)?.lookup(keyText) ?? NullDef(keyText);

    const wordDef = dicts.get(dictid)?.lookup(keyText);
    if (!wordDef || !wordDef.definition) {
      return NullDef(keyText);
    }
    return wordDef as Definition;
  }
  associate(dictid: string, word: string) {
    const result: SuggestItem[] = [];
    if (word.trim() == '' || word.length === 0) {
      return result;
    }

    // limits word result upto 50
    let counter = 0;
    const limit = 50;
    const dict = dicts.get(dictid);
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
