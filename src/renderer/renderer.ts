/**
 * This file will automatically be loaded by webpack and run in the "renderer" context.
 * To learn more about the differences between the "main" and the "renderer" context in
 * Electron, visit:
 *
 * https://electronjs.org/docs/tutorial/application-architecture#main-and-renderer-processes
 *
 * By default, Node.js integration in this file is disabled. When enabling Node.js integration
 * in a renderer process, please be aware of potential security implications. You can read
 * more about security risks here:
 *
 * https://electronjs.org/docs/tutorial/security
 *
 * To enable Node.js integration in this file, open up `main.js` and enable the `nodeIntegration`
 * flag:
 *
 * ```
 *  // Create the browser window.
 *  mainWindow = new BrowserWindow({
 *    width: 800,
 *    height: 600,
 *    webPreferences: {
 *      nodeIntegration: true
 *    }
 *  });
 * ```
 */

import Vue from 'vue';
import VueRouter from 'vue-router';
import Vuex from 'vuex';
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';
import routes from './routes';
import { __RANDOM_KEY__ } from '../utils/random_key';

// customer css
import './renderer.scss';

// Import Bootstrap an BootstrapVue CSS files (order is important)
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-vue/dist/bootstrap-vue.min.css';

// use vuex
Vue.use(Vuex);

// make sure this import after than use vuex
import store from './store';

// Make BootstrapVue available throughout your project
Vue.use(BootstrapVue);
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin);

// Create the router instance and pass the `routes` option
// You can pass in additional options here, but let's
// keep it simple for now.
const router = new VueRouter({
  routes, // short for `routes: routes`
});
// default view as mainWindow
router.push({ path: store.state.defaultWindow });
Vue.use(VueRouter);

console.log(
  '👋 This message is being logged by "renderer.ts", included via webpack'
);

import { cleanUpListeneres } from './init.renderersvc.register';

// Create and mount the root instance.
// Make sure to inject the router with the router option to make the
// whole app router-aware.
const app = new Vue({ router, store }).$mount('#app');
// window extended vue
// @ts-ignore
window[`$vue_${__RANDOM_KEY__}`] = app;
window['DISPATCH_REFER_LINK_WORD'] = function (dictid: string, word: string) {
  app['$state'].dispatch('DISPATCH_REFER_LINK_WORD', { dictid, word });
};
// cleanup ipc listener
cleanUpListeneres();
// rpc test TODO delete this
import './rpctest';
