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

import { IDict } from './types';
import { model } from './model';

export const GetAllDicts = async function (): Promise<Array<IDict>> {
    return new Promise((resolve, reject) =>{
       resolve([])
    })

    // const dicts = await Dicts()
    // return dicts as unknown as Array<IDict>
}


export const LookupWord = async function (dictid:string, word: string) : Promise<model.Resp> {
    return new Promise((resolve, reject) =>{
        resolve(new model.Resp({}))
    })
    // return await Lookup(dictid, word)
}

export const SearchWord = async function(dictid:string, word: string) :Promise<model.Resp> {
    return new Promise((resolve, reject) =>{
        resolve(new model.Resp({}))
    })
    // return await Search(dictid, word)
}

export const LocateWord = async function(dictid:string, keyBlockEntry: model.KeyBlockEntry): Promise<model.Resp> {
    return new Promise((resolve, reject) =>{
        resolve(new model.Resp({}))
    })
    // return await Locate(dictid, keyBlockEntry)
}
