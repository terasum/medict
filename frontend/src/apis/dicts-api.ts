import { IDict } from './types';
import { Dicts, Lookup, Search, Locate } from '../../wailsjs/go/apis/DictsAPI';
import { model } from '../../wailsjs/go/models'

export const GetAllDicts = async function (): Promise<Array<IDict>> {
    const dicts = await Dicts()
    return dicts as unknown as Array<IDict>
}


export const LookupWord = async function (dictid:string, word: string) : Promise<model.Resp> {
    return await Lookup(dictid, word)
}

export const SearchWord = async function(dictid:string, word: string) :Promise<model.Resp> {
    return await Search(dictid, word)
}

export const LocateWord = async function(dictid:string, keyBlockEntry: model.KeyBlockEntry): Promise<model.Resp> {
    return await Locate(dictid, keyBlockEntry)
}
