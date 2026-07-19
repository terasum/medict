export interface SettingsNavigationItem {
  label: string;
  icon: string;
  to: string;
}

export const settingsNavigation: SettingsNavigationItem[] = [
  { label: '词典设置', icon: 'icon-book', to: '/setting/dict' },
  { label: '通用设置', icon: 'icon-cog', to: '/setting/software' },
  { label: '外观设置', icon: 'icon-palette', to: '/setting/theme' },
  { label: '插件管理', icon: 'icon-rocket', to: '/setting/plugin' },
  { label: '调试工具', icon: 'icon-cog', to: '/setting/debug' },
  { label: '使用说明', icon: 'icon-help-circled', to: '/setting/docs' },
  { label: '关于软件', icon: 'icon-info-circled', to: '/setting/about' },
];
