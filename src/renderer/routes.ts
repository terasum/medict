import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/PreferenceWindow.vue';
import DictionaryView from './view/Preference/DictionarySettingsView.vue';
import DebugView from './view/Preference/DebugView.vue';
import SettingsView from './view/Preference/AppSettings.vue';
import PluginsWindow from './view/PluginsWindow.vue';
import TranslateWindow from './view/TranslateWindow.vue';
import AboutView from './view/Preference/AboutView.vue';

export default [
  { path: '/', component: MainWindow },
  { path: '/translate', component: TranslateWindow },
  { path: '/plugins', component: PluginsWindow },
  {
    path: '/preference',
    component: PreferenceWindow,
    children: [
      { path: '', component: DictionaryView },
      {
        path: 'dictSettings',
        component: DictionaryView,
      },
      {
        path: 'settings',
        component: SettingsView,
      },
      {
        path: 'debug',
        component: DebugView,
      },
      {
        path: 'about',
        component: AboutView,
      },
    ],
  },
];
