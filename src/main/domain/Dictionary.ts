import Mdict from 'js-mdict';
import path from 'path';
import fs from 'fs';

import { getResourceRootPath } from '../../config/config';
import { resourceServerPort } from '../init.resource.server';
import { SuggestItem } from '../../model/SuggestItem';
import { Definition, NullDef } from '../../model/Definition';
import { StorabeDictionary } from '../../model/StorableDictionary';

const resourceRoot = getResourceRootPath();

/**
 * Dictionary 词典实体类
 * id: 词典唯一ID, 通常是随机生成的UUID
 * alias: 别名，用于展示在下拉框中，不超过6字符
 * name: 词典全名
 * mdxpath: mdx 文件完整路径
 * mddpath: mdd 文件完整路径
 * resourceBaseDir: 资源文件基础路径
 * mdxDict: mdx 词典实例
 * mddDict: mdd 词典实例
 */
export class Dictionary extends StorabeDictionary {
  mdxDict: Mdict;
  mddDict: Mdict | undefined;
  description: string;
  constructor(
    id: string,
    alias: string,
    name: string,
    mdxpath: string,
    mddpath?: string,
    description?: string
  ) {
    super(id, alias, name, mdxpath, mddpath);
    this.mdxDict = new Mdict(mdxpath);
    this.mddDict = mddpath ? new Mdict(mddpath) : undefined;
    this.description = description || 'undefined';
  }

  findWordDefinition(keyText: string, roffset: number) {
    const result = this.mdxDict.parse_defination(keyText, roffset);
    if (!result) {
      return NullDef(keyText);
    }
    return (result as unknown) as Definition;
  }

  findWordResource(keyText: string) {
    const result = this.mddDict?.lookup(keyText) ?? NullDef(keyText);
    console.log(' ============= load resource =========== ');
    if (result && result.definition) {
      const filePath = rscCachePath(this.id, keyText);
      fs.writeFileSync(filePath, Buffer.from(result.definition, 'base64'));
      console.log(`#### main write file ${filePath}`);
    }

    return {
      keyText,
      definition: rscRelativePath(this.id, keyText),
    } as Definition;
  }

  lookup(word: string) {
    return this.mdxDict.lookup(word);
  }

  associate(word: string) {
    const result: SuggestItem[] = [];
    if (word.trim() == '' || word.length === 0) {
      return result;
    }
    // limits word result upto 50
    let counter = 0;
    const limit = 50;
    const words = this.mdxDict.associate(word);
    for (let i = 0; i < words?.length ?? 0; i++) {
      if (counter >= limit) {
        break;
      }
      const word = words[i];
      // console.log(`set ${key}, ${word.keyText}`)
      result.push({ id: counter, dictid: this.id, ...word });
      counter++;
    }
    return result;
  }
}

function rscCachePath(dictid: string, resourceKey: string) {
  if (!dictid || dictid === '') {
    throw new Error(`invalid dictionary id ${dictid}`);
  }
  const resourcePath = resourceKey.split('\\');

  const fullPath = path.join(resourceRoot, dictid, ...resourcePath);
  const fullDirPath = path.dirname(fullPath);
  if (!fs.existsSync(fullDirPath)) {
    fs.mkdirSync(fullDirPath, { recursive: true });
  }
  return fullPath;
}

function rscRelativePath(dictid: string, resourceKey: string) {
  if (!dictid || dictid === '') {
    throw new Error(`invalid dictionary id ${dictid}`);
  }
  const resourcePath = resourceKey.split('\\');
  const fullPath = path.join(dictid, ...resourcePath);
  return 'http://localhost:' + resourceServerPort + '/' + fullPath;
}
