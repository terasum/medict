import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/PreferenceWindow.vue';
import DictionaryView from './view/Preference/DictionarySettingsView.vue';
import DebugView from './view/Preference/DebugView.vue';
import SettingsView from './view/Preference/AppSettings.vue';
import PluginsWindow from './view/PluginsWindow.vue';
import TranslateWindow from './view/TranslateWindow.vue';
import AboutView from './view/Preference/AboutView.vue';
import TranslateSettingView from './view/Preference/TranslateSettingView.vue';

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
        path: 'translateSettings',
        component: TranslateSettingView,
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
