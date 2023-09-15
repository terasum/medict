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

const API_URL = 'http://api.fanyi.baidu.com/api/trans/vip/translate';
import axios from 'axios';
import { sign } from './baidu_translate_utils/utils';
import random from '@shuaninfo/random';

var querystring = require('querystring');

export async function translate(
  appid: string,
  appkey: string,
  from: string,
  to: string,
  query: string
) {
  const salt = random({ length: 10, type: 'numeric' });
  const signData = sign(appid, query, salt, appkey);

  var queryData = querystring.stringify({
    q: query,
    from,
    to,
    appid,
    salt,
    sign: signData,
  });

  const newUrl = `${API_URL}?${queryData}`;
  return await axios.get(newUrl);
}
