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

// Typed IPC wrappers (issue #729): each call goes to a typed App method
// (wailsjs binding) instead of the old string-dispatched requestBackend/Dispatch.
import { IDict } from './types';
import { model } from './model';
import * as App from '../../wailsjs/go/main/App';

export const InitDicts = async function (): Promise<model.Resp> {
    const resp = await App.InitDicts();
    return resp.data as model.Resp;
}

export const GetAllDicts = async function (): Promise<Array<IDict>> {
    const resp = await App.GetAllDicts();
    return resp.data as Array<IDict>;
}

// BuildIndex
export const BuildIndex = async function (dictid: string): Promise<model.Resp> {
    const resp = await App.BuildIndexByDictId(dictid);
    return resp.data as model.Resp;
}

export const SearchWord = async function (dictid: string, word: string): Promise<model.Resp> {
    const resp = await App.SearchWord(dictid, word);
    return resp.data as model.Resp;
}
