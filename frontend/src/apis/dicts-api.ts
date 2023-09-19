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
import {requestBackend} from '@/apis/apis';


export const GetDictCover = async function(dict_id: string, cover_name:string): Promise<model.Resp> {
try{
        let resp = await requestBackend("GetDictCover", {dict_id, cover_name})
        console.log("[dicts-api] GetDictCover: ", resp)
        return resp.data as unknown as model.Resp
    } catch (error) {
        console.error("[dicts-api] GetDictCover: " ,error)
        return Promise.reject(error)
    }
}


export const GetAllDicts = async function (): Promise<Array<IDict>> {
    try{
        let resp = await requestBackend("GetAllDicts", {})
        console.log("[dicts-api] GetAllDicts: ", resp)
        return resp.data as unknown as Array<IDict>
    } catch (error) {
        console.error("[dicts-api] GetAllDicts: " ,error)
        return Promise.reject(error)
    }
}

// BuildIndex
export const BuildIndex = async function (): Promise<model.Resp> {
    try{
        let resp = await requestBackend("BuildIndex", {})
        console.log("[dicts-api] BuildIndex: ", resp)
        return resp.data as unknown as model.Resp
    } catch (error) {
        console.error("[dicts-api] BuildIndex: " ,error)
        return Promise.reject(error)
    }
}


export const LookupWord = async function (dictid:string, word: string) : Promise<model.Resp> {
    try{
        let resp = await requestBackend("LookupWord", {dict_id:dictid, word:word})
        console.log("[dicts-api]", resp)
        return resp.data as unknown as model.Resp
    } catch (error) {
        console.error("[dicts-api]" ,error)
        return Promise.reject(error)
    }
}

export const SearchWord = async function(dictid:string, word: string) :Promise<model.Resp> {
    try{
        let resp = await requestBackend("SearchWord", {dict_id: dictid, word: word})
        console.log("[dicts-api]", resp)
        return resp.data as unknown as model.Resp
    } catch (error) {
        console.error("[dicts-api]" ,error)
        return Promise.reject(error)
    }
}

export const LocateWord = async function(dictid:string, keyBlockEntry: model.KeyBlockEntry): Promise<model.Resp> {
    try{
        let resp = await requestBackend("LocateWord", {dict_id: dictid, key_block_entry: keyBlockEntry})
        console.log("[dicts-api]", resp)
        return resp.data as unknown as model.Resp
    } catch (error) {
        console.error("[dicts-api]" ,error)
        return Promise.reject(error)
    }
}
