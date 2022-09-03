import MainWindow from './view/MainWindow.vue';
import PreferenceWindow from './view/PreferenceWindow.vue';
import DictionaryView from './view/Preference/DictionarySettingsView.vue';
import DebugView from './view/Preference/DebugView.vue';
import SettingsView from './view/Preference/AppSettings.vue';
import DocsWindow from './view/DocsWindow.vue';
import AboutView from './view/Preference/AboutView.vue';
import TranslateSettingView from './view/Preference/TranslateSettingView.vue';

export default [
  { path: '/', component: MainWindow },
  { path: '/docs', component: DocsWindow },
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
