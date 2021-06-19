import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/Preference.vue';
import PluginsWindow from './view/Plugins.vue';

export default [
    { path: "/", component: MainWindow },
    { path: "/preference", component: PreferenceWindow},
    { path: "/plugins", component: PluginsWindow }
]