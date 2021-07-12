import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/PreferenceWindow.vue';
import DictionaryView from './view/Preference/DictionarySettingsView.vue';
import DebugView from './view/Preference/DebugView.vue';
import PluginsWindow from './view/PluginsWindow.vue';
import TranslateWindow from './view/TranslateWindow.vue';

import DictSettings from './view/DictSettings.vue';

export default [
  { path: '/', component: MainWindow },
  { path: '/translate', component: TranslateWindow },
  { path: '/plugins', component: PluginsWindow },
  {
    path: '/preference',
    component: PreferenceWindow,
    children: [
      { path: '', component: DebugView },
      {
        path: 'dictSettings',
        component: DictionaryView,
      },
      {
        path: 'debug',
        component: DebugView,
      },
    ],
  },

  { path: '/dictSettings', component: DictSettings },
];
