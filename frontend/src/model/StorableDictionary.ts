/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

export declare class DictItem {
  dictid: string;
  dictionary: StorabeDictionary;
}