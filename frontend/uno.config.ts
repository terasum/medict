/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * GPL-3.0
 */

import { defineConfig, presetWind3 } from 'unocss';
import { palette } from './src/style/tokens';

// UnoCSS config (Tailwind-compatible via presetWind3). Theme tokens come from
// src/style/tokens.ts (single source shared with naive-ui themeOverrides).
// shortcuts factor out the card / section-head / list-row / icon-button patterns
// repeated across components, so migration is consistent and terse.
export default defineConfig({
  presets: [presetWind3({ dark: 'media' })],
  theme: {
    colors: {
      primary: {
        DEFAULT: palette.primary,
        hover: palette.primaryHover,
      },
      danger: palette.danger,
      gray: { ...palette.gray },
    },
    fontFamily: {
      sans: palette.fontSans,
      mono: palette.fontMono,
    },
  },
  shortcuts: {
    // icon button (toolbar tiles)
    'btn-icon':
      'inline-flex items-center justify-center rounded cursor-pointer select-none text-gray-600 hover:bg-gray-100 active:bg-gray-200',
    // surface card (dict entry cards, popups, panels)
    'card': 'bg-white border border-gray-200 rounded-lg shadow-sm',
    // section header bar (词典名标题 / 分组头)
    'section-head':
      'flex items-center gap-2 px-2 text-xs text-gray-500 bg-gray-50 border-b border-gray-200',
    // hoverable list row (生词 / 词典列表条目)
    'list-row':
      'flex items-center gap-3 px-3 py-2 rounded cursor-pointer hover:bg-gray-100',
  },
  preflights: [
    {
      layer: 'base',
      getCSS: () => `
        :root {
          --c-primary: ${palette.primary};
          --c-primary-hover: ${palette.primaryHover};
          --c-danger: ${palette.danger};
        }
        body {
          font-family: ${palette.fontSans.join(',')};
          color: ${palette.gray[800]};
          background: #fff;
        }
      `,
    },
  ],
});
