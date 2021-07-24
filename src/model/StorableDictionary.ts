export class StorabeDictionary {
  id: string;
  alias: string;
  name: string;
  mdxpath: string;
  mddpath?: string|string[];
  resourceBaseDir: string;
  description?: string;
  constructor(
    id: string,
    alias: string,
    name: string,
    mdxpath: string,
    mddpath?: string|string[],
    description?: string
  ) {
    this.id = id;
    this.alias = alias;
    this.name = name;
    this.mdxpath = mdxpath;
    this.mddpath = mddpath;
    this.description = description;
    this.resourceBaseDir = '';
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
