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

import google_translate from '@vitalets/google-translate-api';
export async function translate(
    appid: string,
    appkey: string,
    from: string,
    to: string,
    query: string

) {
    if (from == "zh") {
        from = 'zh-CN'
    }
    if (from == "jp") {
        from = 'ja'
    }

    if (to == "zh") {
        to = 'zh-CN'
    }
    if (to == "jp") {
        to = 'ja'
    }

    // @ts-ignore
    let request = google_translate(query, { from: from, to: to, tld: 'cn' })
    return request
}
