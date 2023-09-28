<!--

 Copyright (C) 2023 Quan Chen <chenquan_act@163.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

<template>
  <div class="x-space">
    <div class="x-layout">
      <div class="x-layout-main-area">
        <div class="x-layout-sidebar">
          <AppSidebar>
            <div class="dict-groups">
              <div class="dict-group-list">
                <!-- <nav class="nav-group">
                  <h5 class="nav-group-title">Favorites</h5>
                  <a class="nav-group-item active">
                    <span class="icon icon-star"></span>
                    英汉词典组
                  </a>
                  <span class="nav-group-item">
                    <span class="icon icon-star"></span>
                    汉汉词典组
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-star"></span>
                    汉英词典组
                  </span>
                </nav> -->

                <nav class="nav-group">
                  <h5 class="nav-group-title">Groups</h5>
                  <span class="nav-group-item">
                    <span class="icon icon-archive"></span>
                    默认组
                  </span>
                  <!-- <span class="nav-group-item">
                    <span class="icon icon-archive"></span>
                    词典组A
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-archive"></span>
                    词典组2
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-archive"></span>
                    词典组2
                  </span> -->
                </nav>
              </div>
              <div class="dict-group-settings">
                <div class="btn-group">
                  <button class="btn btn-default">
                    <span class="icon icon-plus"></span>
                  </button>

                  <button class="btn btn-default">
                    <span class="icon icon-minus"></span>
                  </button>
                  <button class="btn btn-default">
                    <span class="icon icon-star"></span>
                  </button>
                  <button class="btn btn-default">
                    <span class="icon icon-star-empty"></span>
                  </button>
                </div>
              </div>
            </div>
          </AppSidebar>
        </div>
        <div class="x-layout-content">
          <AppMainContent>
            <div class="dict-main-area">
              <table class="table-striped">
                <thead>
                  <tr>
                    <th>词典名称</th>
                    <th>标题（内置）</th>
                    <th>词典类型</th>
                    <!-- <th>desc</th> -->
                    <th>引擎版本</th>
                    <th>创建时间</th>
                    <th>id</th>
                    <th>所在目录</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in dictsList">
                    <td>{{ item.name }}</td>
                    <td>{{ item.title }}</td>
                    <td>{{ item.dict_type }}</td>
                    <!-- <td>{{ item.desc }}</td> -->
                    <td>{{ item.generateEngineVersion }}</td>
                    <td>{{ item.createDate }}</td>
                    <td>{{ item.id }}</td>
                    <td>{{ item.base_dir }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </AppMainContent>
        </div>
        <div class="x-layout-right-toolbar">
          <AppRightToolbar>
            <div class="dict-toolbar"></div> 
          </AppRightToolbar>
        </div>
      </div>
      <div class="n-layout-footer">
        <AppFooter>
          <span class="dicts-hint">
            当前组拥有词典数: {{ dictsList.length }}
          </span>
        </AppFooter>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from '../../components/layout/AppSidebar.vue';
import AppFooter from '../../components/layout/AppFooter.vue';
import AppFunctions from '../../components/layout/AppFunctions.vue';
import AppMainContent from '../../components/layout/AppMainContent.vue';
import AppRightToolbar from '@/components/layout/AppRightToolbar.vue';

import { NIcon, NButton, NButtonGroup } from 'naive-ui';
import { useUIStore } from '@/store/ui';
import { useDictQueryStore } from '@/store/dict';
import { GetAllDicts } from '@/apis/dicts-api';

import { ref, onMounted } from 'vue';

const dictsList = ref([]);

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();
uiStore.updateCurrentTab('dict');

onMounted(() => {
  GetAllDicts()
    .then((res) => {
      console.log(res);
      dictsList.value = [];
      res.forEach((dict) => {
        console.log(dict);
        dictsList.value.push({
          name: dict.name,
          dict_type: dict.type,
          title: dict.description.title,
          desc: dict.description.description,
          generateEngineVersion: dict.description.generateEngineVersion,
          createDate: dict.description.createDate,
          id: dict.id,
          base_dir: dict.base_dir,
        });
      });
      console.log(dictsList.value);
    })
    .catch((err) => {
      console.log(err);
    });
});
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.x-space {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;

  .x-layout {
    width: 100%;
    height: 100%;
    padding: 0;
    margin: 0;
    .x-layout-main-area {
      display: flex;
      flex-direction: row;
      width: 100%;
      height: calc(100% - $layout-footer-height);
      padding: 0;
      margin: 0;
      margin-right: -60px;

      .x-layout-sidebar {
        width: $layout-left-sidebar-width;
        height: 100%;
        padding: 0;
        margin: 0;
        background-color: #fafafa;
      }

      .x-layout-content {
        width: calc(
          100% - $layout-left-sidebar-width - $layout-right-toolbar-width
        );
        height: 100%;
        padding: 0;
        margin: 0;
      }

      .x-layout-right-toolbar {
        width: $layout-right-toolbar-width;
      }
    }
    .n-layout-footer {
      width: 100%;
      height: $layout-footer-height;
      padding: 0;
      margin: 0;
    }
  }
}

.dict-groups {
  height: 100%;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  .dict-group-list {
    height: calc(100% - 30px);
    overflow: auto;
  }
  .dict-group-settings {
    height: 30px;
    display: flex;
    justify-content: center;
    flex-direction: row;
    background-color: #ededed;
    border-top: 1px solid #ccc;
    .btn-group {
      margin-top: 3px;
    }
  }
}

.dict-main-area {
  width: 100%;
  overflow-x: auto;
  padding-bottom: 20px;
}
.dict-toolbar {
  background-color: #fafafa;
  height: 100%;
  width: 100%;
}
.dicts-hint{
  font-size: 12px;
  height: 20px;
  line-height: 20px;
}
</style>
