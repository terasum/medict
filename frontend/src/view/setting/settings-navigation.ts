export interface SettingsNavigationItem {
  label: string;
  description: string;
  icon: string;
  to: string;
}

export interface SettingsNavigationGroup {
  title: string;
  items: SettingsNavigationItem[];
}

export const settingsNavigation: SettingsNavigationGroup[] = [
  {
    title: '应用设置',
    items: [
      { label: '词典', description: '目录、加载与内容', icon: 'icon-book', to: '/setting/dict' },
      { label: '通用', description: '搜索与应用行为', icon: 'icon-cog', to: '/setting/software' },
      { label: '外观', description: '主题与显示', icon: 'icon-palette', to: '/setting/theme' },
    ],
  },
  {
    title: '扩展与开发',
    items: [
      { label: '插件', description: '目录与访问权限', icon: 'icon-rocket', to: '/setting/plugin' },
      { label: '调试工具', description: '资源与词典诊断', icon: 'icon-bug', to: '/debug' },
    ],
  },
  {
    title: '帮助与关于',
    items: [
      { label: '使用说明', description: '界面、配置与常见问题', icon: 'icon-help-circled', to: '/docs' },
      { label: '隐私声明', description: '数据与隐私说明', icon: 'icon-feather', to: '/setting/terms' },
      { label: '开源许可', description: '第三方协议', icon: 'icon-cc', to: '/setting/license' },
      { label: '版本更新', description: '检查新版本', icon: 'icon-arrows-ccw', to: '/setting/update' },
      { label: '关于 Medict', description: '版本与开发信息', icon: 'icon-info-circled', to: '/setting/about' },
    ],
  },
];
