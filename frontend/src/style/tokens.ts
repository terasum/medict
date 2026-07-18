/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * GPL-3.0
 */

// Design tokens — the SINGLE source of truth for the app's visual language.
// Consumed by BOTH UnoCSS (uno.config.ts theme) and naive-ui (App.vue
// themeOverrides), so the component library and the custom utility CSS share one
// palette. Replaces the fragmented hardcoded colors (#11a8ff / #2a7de1 / #0645ad
// / ~10 neutral grays) audited across the components.

export const palette = {
  // brand / link / active — matches naive-ui primaryColor (was #326cb8 in App.vue)
  primary: '#326cb8',
  primaryHover: '#2a5fa6',
  // destructive / active-red (legacy $theme-function-box-active-color)
  danger: '#bd3134',
  // neutral scale: collapses the ~10 ad-hoc grays (#999/#aaa/#bbb/#ccc/#ddd/
  // #eee/#f1f1f1/#fafafa/#f6f8fa/#dcdcdc/#d0d7de/…) into one ramp.
  gray: {
    50: '#f6f6f7',
    100: '#f1f1f2',
    200: '#e3e3e3',
    300: '#dcdcdc',
    400: '#bbbbbb',
    500: '#999999',
    600: '#777777',
    700: '#666666',
    800: '#444444',
    900: '#333333',
  },
  fontSans: [
    'system-ui',
    '-apple-system',
    '"Segoe UI"',
    '"Microsoft YaHei"',
    '"PingFang SC"',
    'sans-serif',
  ],
  fontMono: ['"Fira Code"', 'ui-monospace', 'SFMono-Regular', 'monospace'],
} as const;

// BRAND is the one constant both UnoCSS and naive-ui read for the primary color.
export const BRAND = palette.primary;
