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

interface ISettings {
    title: string // Overrides the default title
    showSettings: boolean // Controls settings panel display
    showTagsView: boolean // Controls tagsview display
    showSidebarLogo: boolean // Controls siderbar logo display
    fixedHeader: boolean // If true, will fix the header component
    errorLog: string[] // The env to enable the errorlog component, default 'production' only
    sidebarTextTheme: boolean // If true, will change active text color for sidebar based on theme
    devServerPort: number // Port number for webpack-dev-server
    mockServerPort: number // Port number for mock server
  }
  
  // You can customize below settings :)
  const settings: ISettings = {
    title: 'Vue Typescript Admin',
    showSettings: true,
    showTagsView: true,
    fixedHeader: false,
    showSidebarLogo: false,
    errorLog: ['production'],
    sidebarTextTheme: true,
    devServerPort: 9527,
    mockServerPort: 9528
  }
  
  export default settings