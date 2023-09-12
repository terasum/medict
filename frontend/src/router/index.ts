import MainWindow from '@/view/main/MainWindow.vue';
import PreferenceWindow from '@/view/preference/PreferenceWindow.vue';
import DictionaryView from '@/view/preference/DictionarySettingsView.vue';
import DebugView from '@/view/preference/DebugView.vue';
import SettingsView from '@/view/preference/AppSettings.vue';
import DocsWindow from '@/view/docs/DocsWindow.vue';
import AboutView from '@/view/preference/AboutView.vue';
import TranslateSettingView from '@/view/Preference/TranslateSettingView.vue';

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
