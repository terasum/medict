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

import MainWindow from '@/view/main/MainWindow.vue';
import PreferenceWindow from '@/view/preference/PreferenceWindow.vue';
import DictionaryView from '@/view/preference/DictionarySettingsView.vue';
import DebugView from '@/view/preference/DebugView.vue';
import SettingsView from '@/view/preference/AppSettings.vue';
import DocsWindow from '@/view/docs/DocsWindow.vue';
import AboutView from '@/view/preference/AboutView.vue';
import TranslateSettingView from '@/view/Preference/TranslateSettingView.vue';

export default [
  { path: '/', component: MainWindow },
  { path: '/docs', component: DocsWindow },
  {
    path: '/preference',
    component: PreferenceWindow,
    children: [
      { path: '', component: DictionaryView },
      {
        path: 'dictSettings',
        component: DictionaryView,
      },
      {
        path: 'translateSettings',
        component: TranslateSettingView,
      },
      {
        path: 'settings',
        component: SettingsView,
      },
      {
        path: 'debug',
        component: DebugView,
      },
      {
        path: 'about',
        component: AboutView,
      },
    ],
  },
];
