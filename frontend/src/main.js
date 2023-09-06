import {createApp} from 'vue';
import naive from 'naive-ui';

import routes from './routes';
import { createRouter,createWebHashHistory } from 'vue-router'

import { __RANDOM_KEY__ } from './utils/random_key';


import 'normalize.css/normalize.css';
// 通用字体
import 'vfonts/Lato.css'
// 等宽字体
import 'vfonts/FiraCode.css'

import './renderer.scss';

import './renderer.init';

import App from './App.vue';


// cleanup ipc listener, make sure this invoke before import store
// cleanUpListeneres();

// make sure this import after than use vuex
// import store from './store';


const router = createRouter({
  history: createWebHashHistory(),
  routes, 
})

// default view as mainWindow


const app = createApp(App);

router.push({ path: "/"}); // store.state.defaultWindow });
app.use(router);

app.use(naive)

// remove skeleton
let skeleton = document.querySelector('#skeleton-wrapper')
if (skeleton) {
  skeleton.innerHTML = ''
}

// appMain.$mount('#app');
app.mount('#app')