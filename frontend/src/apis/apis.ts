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

import { Dispatch, ResourceServerAddr } from '../../wailsjs/go/main/App';
import { model } from '../../wailsjs/go/models';

function objectToPathParams(obj) {
  const params = [];
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      params.push(`${encodeURIComponent(key)}=${encodeURIComponent(obj[key])}`);
    }
  }
  return params.join('&');
}

export const StaticDictServerURL = function (): Promise<string> {
  if (window['go']) {
    return ResourceServerAddr();
  } else {
    Promise.resolve("http://localhost:1")
  }
};


export async function requestBackend(apiName, data): Promise<model.Resp> {
  if (window['go']) {
    return Dispatch(apiName, data);
  } else {
    return Promise.resolve({
      data: [{id:"1", "name":"test"}],
      err: "",
      code: 200,
    })

  }
}
