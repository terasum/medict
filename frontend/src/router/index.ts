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

import MainWindow from '@/view/main/index.vue';
import DictWindow from '@/view/dict/index.vue';
import PluginsWindow from '@/view/plugins/index.vue';
import SettingWindow from '@/view/setting/index.vue';
import BookmarksWindow from '@/view/bookmarks/index.vue';

import index_md from '@/assets/docs/index.md';
import select_and_use_md from '@/assets/docs/select_and_use_dict.md';
import faq_md from '@/assets/docs/faq.md';
import terms_and_service from '@/assets/docs/terms_and_service.md';
import license_md from '@/assets/docs/license.md';
import about_md from '@/view/about/index.vue';

import SettingDict from "@/view/setting/SettingDict.vue";
import SettingSoftware from "@/view/setting/SettingSoftware.vue";
import SettingTheme from "@/view/setting/SettingTheme.vue";
import SettingPlugin from "@/view/setting/SettingPlugin.vue";
import SettingDebug from "@/view/setting/SettingDebug.vue";
import SettingDocs from "@/view/setting/SettingDocs.vue";

export default [
  { path: '/', component: MainWindow },
  { path: '/dict', component: DictWindow },
  { path: '/setting', component: SettingWindow, children:[
    
      { path: '', redirect: '/setting/dict' },
      { path: 'dict', component: SettingDict },
      { path: 'software', component: SettingSoftware },
      { path: 'theme', component: SettingTheme },
      { path: 'plugin', component: SettingPlugin },
      { path: 'debug', component: SettingDebug },
      {
        path: 'docs',
        component: SettingDocs,
        children: [
          { path: '', redirect: '/setting/docs/index' },
          { path: 'index', component: index_md },
          { path: 'select-and-use', component: select_and_use_md },
          { path: 'faq', component: faq_md },
          { path: 'privacy', component: terms_and_service },
          { path: 'license', component: license_md },
        ],
      },
      { path: 'terms', redirect: '/setting/docs/privacy' },
      { path: 'license', redirect: '/setting/docs/license' },
      { path: 'about', component: about_md },
      { path: 'update', redirect: '/setting/about' },
    
  ]},
  { path: '/bookmarks', component: BookmarksWindow },
  { path: '/plugins', component: PluginsWindow },
  { path: '/debug', redirect: '/setting/debug' },
  { path: '/docs', redirect: '/setting/docs/index' },
  { path: '/docs/index', redirect: '/setting/docs/index' },
  { path: '/docs/select_and_use', redirect: '/setting/docs/select-and-use' },
  { path: '/docs/faq', redirect: '/setting/docs/faq' },
  { path: '/docs/privacy', redirect: '/setting/docs/privacy' },
  { path: '/docs/license', redirect: '/setting/docs/license' },
];
