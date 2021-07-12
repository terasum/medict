 import Vue from 'vue';
 import VueRouter from 'vue-router';
 import Vuex from 'vuex';
 import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';
 import routes from '../routes';
 
 // Import Bootstrap an BootstrapVue CSS files (order is important)
 import 'bootstrap/dist/css/bootstrap.min.css';
 import 'bootstrap-vue/dist/bootstrap-vue.min.css';
 // use vuex
 Vue.use(Vuex);
 // Make BootstrapVue available throughout your project
Vue.use(BootstrapVue);
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin);

 
 // make sure this import after than use vuex
 import store from '../store';
 
 // Create the router instance and pass the `routes` option
 // You can pass in additional options here, but let's
 // keep it simple for now.
 const router = new VueRouter({
   routes // short for `routes: routes`
 })
 // default view as mainWindow
 router.push({ path: '/dictSettings'});
 Vue.use(VueRouter);
 
 
console.log("this is console from dict-settings.window.ts");

 // Create and mount the root instance.
 // Make sure to inject the router with the router option to make the
 // whole app router-aware.
 const app = new Vue({router, store}).$mount('#app');
 
 
 
 