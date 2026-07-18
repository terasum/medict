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

// exportCurrentEntry 把当前词条渲染为自包含 HTML(资源内联)并保存到用户选择的
// 路径,便于调试复杂词条(#784)。返回保存路径;用户取消对话框返回 '';
// 后端错误(code != 200)抛错,由调用方 message 提示。
export const exportCurrentEntry = async function (dictid: string, word: string): Promise<string> {
    const resp = await App.ExportCurrentEntry(dictid, word);
    if (resp.code !== 200) {
        throw new Error(resp.err || '导出失败');
    }
    return (resp.data as string) || '';
}

// getDictUserCSS reads the per-dictionary user CSS override (#783). Returns ''
// if no override exists.
export const getDictUserCSS = async (dictid: string): Promise<string> => {
    const resp = await App.GetDictUserCSS(dictid);
    return (resp.data as string) || '';
};

// getDictEditorCSS uses a saved override when present, otherwise the selected
// dictionary's own CSS snapshot prepared by the main process.
export const getDictEditorCSS = async (dictid: string): Promise<string> => {
    const resp = await App.GetDictEditorCSS(dictid);
    if (resp.code !== 200) {
        throw new Error(resp.err || '加载词典 CSS 失败');
    }
    return (resp.data as string) || '';
};

export const openDictCSSWindow = async (dictid: string, dictName: string, word: string): Promise<void> => {
    const resp = await App.OpenDictCSSWindow(dictid, dictName, word);
    if (resp.code !== 200) {
        throw new Error(resp.err || '无法打开 CSS 编辑器');
    }
};

// saveDictUserCSS persists the per-dictionary user CSS override (#783). Empty
// css deletes the sidecar file.
export const saveDictUserCSS = async (dictid: string, css: string): Promise<void> => {
    const resp = await App.SaveDictUserCSS(dictid, css);
    if (resp.code !== 200) {
        throw new Error(resp.err || '保存失败');
    }
};
