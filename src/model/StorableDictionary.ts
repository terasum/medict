export class StorabeDictionary {

  id: string;
  alias: string;
  name: string;
  mdxpath: string;
  mddpath?: string | string[];
  resourceBaseDir: string;
  description?: string;
  byScanning: boolean;

  constructor(
    id: string,
    alias: string,
    name: string,
    mdxpath: string,
    mddpath?: string | string[],
    description?: string,
    byScanning?: boolean,
  ) {
    this.id = id;
    this.alias = alias;
    this.name = name;
    this.mdxpath = mdxpath;
    this.mddpath = mddpath;
    this.description = description;
    this.resourceBaseDir = '';
    this.byScanning = false;
    if(byScanning) {
      this.byScanning = true;
    }
  }

  static clone(dict: StorabeDictionary) {
    const newDict = new StorabeDictionary(
      dict.id,
      dict.alias,
      dict.name,
      dict.mdxpath,
      dict.mddpath,
      dict.description
    );

    newDict.resourceBaseDir = dict.resourceBaseDir;
    return newDict;
  }
}
