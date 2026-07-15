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

// 用户偏好(medict.toml)的读写包装。设置回写基础设施:任何「记住用户偏好」类
// 需求(多词典激活集、主题、字号等)都通过这两个函数落盘,无需后端逐项加字段。
//
// 注意:viper 会把所有 key 小写化,所以 getPreferences 返回的 key 是小写的;
// 保存时 key 大小写不敏感,建议统一用小写(如 multidictids)避免混淆。
import * as App from '../../wailsjs/go/main/App';

export type Preferences = Record<string, any>;

// getPreferences 读取全部已持久化设置(medict.toml + 运行时覆盖)。
export const getPreferences = async (): Promise<Preferences> => {
    const resp = await App.GetPreferences();
    return (resp.data as Preferences) || {};
};

// savePreferences 把给定 key/value 合并写入 medict.toml(保留既有 key)。
// 失败(code != 200)抛错,由调用方 try/catch + message 提示。
export const savePreferences = async (prefs: Preferences): Promise<void> => {
    const resp = await App.SavePreferences(prefs);
    if (resp.code !== 200) {
        throw new Error(resp.err || '保存设置失败');
    }
};
