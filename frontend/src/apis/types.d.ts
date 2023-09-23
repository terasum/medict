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

export interface IDict {
    id: string
    name: string
    discrption?: string
    imageURL?: string
    mdxFileURL?: string
    mddFileURL?: string
}

export interface IWordEntry {
    word: string,
    keyText: string,
    record_start: number,
    record_end: number,
}

export interface IArticleData {
    id: number
    status: string
    title: string
    abstractContent: string
    fullContent: string
    sourceURL: string
    imageURL: string
    timestamp: string | number
    platforms: string[]
    disableComment: boolean
    importance: number
    author: string
    reviewer: string
    type: string
    pageviews: number
  }
  
  export interface IRoleData {
    key: string
    name: string
    description: string
    routes: any
  }
  
  export interface ITransactionData {
    orderId: string
    timestamp: string | number
    username: string
    price: number
    status: string
  }
  
  export interface IUserData {
    id: number
    username: string
    password: string
    name: string
    email: string
    phone: string
    avatar: string
    introduction: string
    roles: string[]
  }