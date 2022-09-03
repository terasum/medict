import Vue from 'vue';
import VueRouter from 'vue-router';
import Vuex from 'vuex';
import Buefy from 'buefy'

import routes from './routes';
import { __RANDOM_KEY__ } from './utils/random_key';


import 'normalize.css/normalize.css';
// import '@fortawesome/fontawesome-free/js/brands.min.js';
import '@fortawesome/fontawesome-free/js/solid.min.js';
import '@fortawesome/fontawesome-free/js/fontawesome.min.js';
import 'buefy/dist/buefy.min.css';

// customer css
import './renderer.scss';

Vue.config.productionTip = false
Vue.config.devtools = true;


// cleanup ipc listener, make sure this invoke before import store
// cleanUpListeneres();

// make sure this import after than use vuex
import store from './store';

Vue.use(Buefy);

// Create the router instance and pass the `routes` option
// You can pass in additional options here, but let's
// keep it simple for now.
const router = new VueRouter({
  routes, // short for `routes: routes`
});
// default view as mainWindow
router.push({ path: store.state.defaultWindow });
Vue.use(VueRouter);

import './renderer.init';
import App from './App.vue';


// Create and mount the root instance.
// Make sure to inject the router with the router option to make the
// whole app router-aware.
const appMain = new Vue({ router, store , render: h => h(App)});

// remove skeleton
let skeleton = document.querySelector('#skeleton-wrapper')
if (skeleton) {
  skeleton.innerHTML = ''
}

appMain.$mount('#app');
