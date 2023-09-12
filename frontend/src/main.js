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
