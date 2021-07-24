import { Dictionary } from '../../domain/Dictionary';
import { SuggestItem } from '../../../model/SuggestItem';
import { Definition, NullDef } from '../../../model/Definition';
import { StorabeDictionary } from '../../../model/StorableDictionary';
import { StorageService } from './StorageServcice';
import { getConfigJsonPath,getResourceRootPath } from '../../../config/config';
import path from 'path';
import fs from 'fs';
import rimraf from "rimraf";


const storageService = new StorageService(getConfigJsonPath());

const dicts = new Map<string, Dictionary>();

function loadDicts() {
  const dictLists = storageService.getDataByKey('dicts') as any[];

  if (dictLists) {
    dictLists.forEach(dict => {
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
        fs.mkdirSync(fpath)
      }
    });
  }
}
// init load
loadDicts();

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
    dicts.forEach(val => {
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

  deleteOne(dictid: string) {
    dicts.delete(dictid);
    saveToFile(dicts);
    loadDicts();
    const fpath = path.resolve(getResourceRootPath(), dictid);
    if (fs.existsSync(fpath)) {
      rimraf(fpath, function () { console.log("delete directory done", fpath); });
    }
    return true;
  }

  findWordPrecisly(dictid: string, keyText: string, rofset: number) {
    return dicts.get(dictid)?.findWordDefinition(keyText, rofset);
  }

  loadDictResource(dictid: string, keyText: string) {
    const wordDef = dicts.get(dictid)?.findWordResource(keyText);
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

    const tempMap = new Map<string, SuggestItem>();
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
      tempMap.set(word.keyText, {
        id: counter,
        dictid: word.dictid,
        keyText: word.keyText,
        rofset: word.rofset,
      });
      counter++;
    }
    // reassembe
    for (const item of tempMap.values()) {
      result.push(item);
    }
    return result;
  }
}
