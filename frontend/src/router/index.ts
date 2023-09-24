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
import SettingWindow from '@/view/setting/index.vue';
import DocWindow from '@/view/docs/index.vue';

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
import SettingUpdater from "@/view/setting/SettingUpdate.vue";

export default [
  { path: '/', component: MainWindow },
  { path: '/dict', component: DictWindow },
  { path: '/setting', component: SettingWindow, children:[
    
      { path: '', component: SettingDict },
      { path: 'dict', component: SettingDict },
      { path: 'software', component: SettingSoftware },
      { path: 'theme', component: SettingTheme },
      { path: 'plugin', component: SettingPlugin },
      { path: 'terms', component: terms_and_service },
      { path: 'license', component: license_md },
      { path: 'about', component: about_md },
      { path: 'update', component: SettingUpdater },
    
  ]},
  {
    path: '/docs',
    component: DocWindow,
    children: [
      { path: '', component: index_md },
      { path: 'index', component: index_md },
      { path: 'select_and_use', component: select_and_use_md },
      { path: 'faq', component: faq_md },
    ],
  },
];
