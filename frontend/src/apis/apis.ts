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


import {Dispatch, ResourceServerAddr} from '../../wailsjs/go/main/App';
import {model} from "../../wailsjs/go/models"

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
  return ResourceServerAddr()
}

const APIs = async function() {
  // let baseURL = await (getBaseURL())();
  // return {
  //   "GetAllDicts": {method: "GET", path: baseURL + "/__api/v1/dicts/GetAllDicts"},
  //   "LookupWord":  {method: "GET", path: baseURL + "/__api/v1/dicts/LookupWord"},
  //   "SearchWord":  {method: "GET", path: baseURL + "/__api/v1/dicts/SearchWord"},
  //   "LocateWord":  {method: "GET", path: baseURL + "/__api/v1/dicts/LocateWord"},
  // }
}

export async function requestBackend(apiName, data): Promise<model.Resp> {
  return Dispatch(apiName, data)

  // let apis = await APIs()
  // let api = apis[apiName];
  // if (!api) {
  //   return
  // }
  // let method = api.method;
  // let path = api.path;
  // if (method=== "GET"){
  //   const queryParams = objectToPathParams(data);
  //   const fullUrl = `${path}?${queryParams}`;
  //   return await request.get(fullUrl, {
  //     url: fullUrl,
  //     method: "GET",
  //   })
  // } else {
  //   return await request.get(path, {
  //     url: path,
  //     method: "GET",
  //     data: data,
  //   })
  // }

}