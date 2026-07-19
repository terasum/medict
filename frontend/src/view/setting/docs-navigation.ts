import type { MarkdownNavigationItem } from '@/components/setting/MarkdownNavigation.vue';

export const docsNavigation: MarkdownNavigationItem[] = [
  { label: '界面介绍', icon: 'icon-window', to: '/setting/docs/index' },
  { label: '词典配置与使用', icon: 'icon-book', to: '/setting/docs/select-and-use' },
  { label: '常见问题', icon: 'icon-help-circled', to: '/setting/docs/faq' },
  { label: '隐私声明', icon: 'icon-feather', to: '/setting/docs/privacy' },
  { label: '开源许可', icon: 'icon-cc', to: '/setting/docs/license' },
];
