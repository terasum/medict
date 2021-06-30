import Mdict from 'js-mdict'
import { WordDefinition } from 'js-mdict'
import path from 'path'
import fs from 'fs'
import { getResourceRootPath } from '../../config/config'
import { resourceServerPort } from '../../main/resource.server'

const dicts = new Map<string, Dictionary>()

class Dictionary {
  protected id: string
  protected fname: string
  protected mdxpath: string
  protected mddpath?: string
  protected csspath?: string
  protected jspath?: string
  public dict: Mdict
  public mddDict?: Mdict

  constructor(
    id: string,
    fname: string,
    mdxpath: string,
    mddpath?: string,
    csspath?: string,
    jspath?: string,
  ) {
    this.id = id
    this.fname = fname
    this.mdxpath = mdxpath
    this.mddpath = mddpath
    this.csspath = csspath
    this.jspath = jspath
    this.dict = new Mdict(mdxpath)
    this.mddDict = mddpath ? new Mdict(mddpath) : undefined
  }
}

declare class SuggestItem {
  keyText: string
  rofset: number
  ed?: number | undefined
  id: number
  dictid: string
}

// dicts.set(
//   'macmillan',
//   new Dictionary(
//     'macmillan',
//     'macmillan.mdx',
//     "/Users/chenquan/Workspace/nodejs/js-mdict/mdx/testdict/v2.0/Macmillan English Dictionary/Macmillan English Dictionary.mdx",
//     "/Users/chenquan/Workspace/nodejs/js-mdict/mdx/testdict/v2.0/Macmillan English Dictionary/Macmillan English Dictionary.mdd"
//   ),
// )

// dicts.set(
//   'OxfordStudent',
//   new Dictionary(
//     'OxfordStudent',
//     'OxfordStudent.mdx',
//     "/Users/chenquan/Workspace/nodejs/js-mdict/mdx/testdict/v2.0/Oxford Collocations Dictionary for students of English 2nd/Oxford Collocations Dictionary for students of English 2nd.mdx",
//     "/Users/chenquan/Workspace/nodejs/js-mdict/mdx/testdict/v2.0/Oxford Collocations Dictionary for students of English 2nd/Oxford Collocations Dictionary for students of English 2nd.mdd"
//   ),
// )

dicts.set(
  'oxford',
  new Dictionary(
    'oale8',
    'oale8.mdx',
    '/Users/chenquan/Workspace/nodejs/js-mdict/mdx/testdict/oale8.mdx',
    '/Users/chenquan/Workspace/nodejs/js-mdict/mdx/testdict/oale8.mdd',
  ),
)


function resourceCachePath(dictid: string, resourceKey: string) {
  if (!dictid) {
    dictid = 'temp'
  }
  const resourcePath = resourceKey.split('\\');

  // const fullPath = path.join(__dirname, '../renderer/dicts', dictid, ...resourcePath);
  const fullPath = path.join(getResourceRootPath(), dictid, ...resourcePath)
  const fullDirPath = path.dirname(fullPath);
  if (!fs.existsSync(fullDirPath)) {
    fs.mkdirSync(fullDirPath, { recursive: true });
  }
  return fullPath;
}

function resourceRelativePath(dictid: string, resourceKey: string) {
  if (!dictid) {
    dictid = 'temp'
  }
  const resourcePath = resourceKey.split('\\');
  const fullPath = path.join(dictid, ...resourcePath);
  return 'http://localhost:' + resourceServerPort +'/' + fullPath;
}



export const dictService = {
  findWordPrecisly: (dictid: string, keyText: string, roffset: number) => {
    const result = dicts.get(dictid)?.dict.parse_defination(keyText, roffset)
    if (!result) {
      return { keyText, definition: 'null' }
    }
    return (result as unknown) as { keyText: string; definition: string }
  },
  loadDictResource: (dictid: string, keyText: string) => {
    const nullDef = ({
      keyText,
      definitions: null,
    } as unknown) as WordDefinition
    const result = dicts.get(dictid)?.mddDict?.lookup(keyText) ?? nullDef;
    console.log(' ============= load resource =========== ')
    // console.log(result)

    if (result && result.definition) {
      const filePath = resourceCachePath(dictid, keyText);
      fs.writeFileSync(filePath, Buffer.from(result.definition, 'base64'));
      console.log(`#### main write file ${filePath}`);
    }

    return { keyText, definition: resourceRelativePath(dictid, keyText) }
  },
  lookup: (dictid: string, word: string) => {
    return dicts.get(dictid)?.dict.lookup(word) ?? null
  },
  associate: (word: string) => {
    const result: SuggestItem[] = []
    if (word.trim() == '' || word.length === 0) {
      return result
    }

    const tempMap = new Map<string, SuggestItem>()
    // limits word result upto 50
    let counter = 0
    const limit = 50
    for (const key of dicts.keys()) {
      const words = dicts.get(key)?.dict.associate(word)
      if (!words) {
        continue
      }
      for (let i = 0; i < words?.length ?? 0; i++) {
        if (counter >= limit) {
          break
        }
        const word = words[i]
        // console.log(`set ${key}, ${word.keyText}`)
        tempMap.set(word.keyText, { id: counter, dictid: key, ...word })
        counter++
      }
    }
    // reassembe
    for (const item of tempMap.values()) {
      result.push(item)
    }
    return result
  },
}
