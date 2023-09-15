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

import naive from 'naive-ui';
import { createApp } from 'vue';
import { createRouter, createWebHashHistory } from 'vue-router';
import { createPinia } from 'pinia';

import routes from '@/router';
import App from '@/App.vue';

import 'normalize.css/normalize.css';
// 通用字体
import 'vfonts/Lato.css';
// 等宽字体
import 'vfonts/FiraCode.css';

import '@/style/renderer.scss';

import '@/renderer.init';


const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const pinia = createPinia();

const app = createApp(App);

router.push({ path: '/' }); // store.state.defaultWindow });
app.use(router);

app.use(naive);
app.use(pinia);

// remove skeleton
let skeleton = document.querySelector('#skeleton-wrapper');
if (skeleton) {
  skeleton.innerHTML = '';
}

// appMain.$mount('#app');
app.mount('#app');
