import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/Preference.vue';
import PluginsWindow from './view/Plugins.vue';
import TranslateWindow from './view/Translate.vue';
import DictSettings from './view/DictSettings.vue';

export default [
    { path: "/", component: MainWindow },
    { path: "/preference", component: PreferenceWindow},
    { path: "/plugins", component: PluginsWindow },
    { path: "/translate", component: TranslateWindow },
    { path: "/dictSettings", component: DictSettings },
]