export class StorabeDictionary {
  id: string;
  alias: string;
  name: string;
  mdxpath: string;
  mddpath: string | undefined;
  resourceBaseDir: string;
  constructor(
    id: string,
    alias: string,
    name: string,
    mdxpath: string,
    mddpath?: string
  ) {
    this.id = id;
    this.alias = alias;
    this.name = name;
    this.mdxpath = mdxpath;
    this.mddpath = mddpath;
    this.resourceBaseDir = '';
  }
}
