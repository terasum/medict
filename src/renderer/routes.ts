import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/PreferenceWindow.vue';
import PluginsWindow from './view/PluginsWindow.vue';
import TranslateWindow from './view/TranslateWindow.vue';
import DebugWindow from './view/DebugWindow.vue';

import DictSettings from './view/DictSettings.vue';

export default [
  { path: '/', component: MainWindow },
  { path: '/translate', component: TranslateWindow },
  { path: '/plugins', component: PluginsWindow },
  { path: '/preference', component: PreferenceWindow },
  { path: '/debug', component: DebugWindow },

  { path: '/dictSettings', component: DictSettings },
];
