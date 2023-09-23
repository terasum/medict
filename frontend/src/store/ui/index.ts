/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import { defineStore } from 'pinia';

export const useUIStore = defineStore('ui', {
  state: () => ({
    currentTab: 'search',

    progressHint: "",
    progressPercent: 0,
    
  }),
  actions: {
    updateCurrentTab(tabName: string) {
      this.currentTab = tabName;
    },
    isSearchInputActive() {
        return this.currentTab === "search";
    },
    updateProgress(hint:string, progress:number) {
      if (progress > 100)  { 
        progress = 100;
      }
      if (progress < 0) {
        progress = 0
      }

      this.progressHint = hint;
      this.progressPercent = progress;

    }
  },
});
