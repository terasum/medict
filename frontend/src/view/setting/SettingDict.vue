<template>
  <div class="setting-main-container" ref="containerRef">
    <div class="setting-main-container-header">
      <h3>词典设置</h3>
    </div>

    <div class="setting-main-container-content">
      <n-card class="setting-section">
          <SettingItem title="词典目录">
             <template #desc>
              词典目录,用于存放所有的词典
            </template>
            <template #action>
            <button class="btn btn-default" @click="OpenDictDir">打开</button>
            <!-- <button class="btn btn-default">修改</button> -->
            </template>
            {{ dictDir }}
          </SettingItem>
      </n-card>

      <n-card class="setting-section">
         <SettingItem title="单词典是否允许外链词典" >
            <template #desc>
            是否允许引入软链接词典
            </template>
            <template #action>
            <button class="btn btn-default">修改</button>
            </template>
            false
          </SettingItem>
          <SettingItem title="单词典组最大词典数" value="3">
            <template #desc>
            为了避免索引建立时间过长，单个词典组需要限制词典数量
            </template>
            <template #action>
            <button class="btn btn-default">修改</button>
            </template>
            10
          </SettingItem>
      </n-card>

      <n-card class="setting-section">
         <SettingItem title="词典内容预置css/js" >
            <template #desc>
              渲染词典内容时在head中预置的css/js
            </template>
            <template #action>
            <button class="btn btn-default">修改</button>
            </template>
            <div class="setting-content-inner">
            <pre>
              <code>
                  {{  presetContent }}
              </code>
            </pre>
            </div>
          </SettingItem>
         
      </n-card>


    </div>
  </div>
</template>
<script lang="ts" setup>
import { NCard, NAffix, NTag } from 'naive-ui';
import { ref,onMounted } from 'vue';
import SettingItem from "@/components/setting/SettingItem.vue";
import {OpenDirOrFile, BaseDictDirectory} from "@/apis/apis";

const containerRef = ref<any>();
const dictDir = ref("");

const presetContent = ref(`
<link href="#{DictName}.css?dict_id=#{DictID}" rel="stylesheet"/>
<script async="" src="#{DictName}.js?dict_id=#{DictID}"/>
`)

function OpenDictDir() {
    if (!dictDir.value || dictDir.value == "") {
        return;
    }
    OpenDirOrFile(dictDir.value).then(()=>{
        console.log("open success")
    }).catch(err =>{
        console.log(err)
    })
}

onMounted(()=>{
    BaseDictDirectory().then(dir =>{
        dictDir.value = dir
    }).catch(err=>{
        console.log(err)
        dictDir.value = ""
    })

})



</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.setting-main-container {
  height: 100%;
  width: 100%;
  overflow-y: auto;
  margin: 0;
  padding: 0;
  .setting-main-container-header {
    display: flex;
    z-index: 99;
    background: #fff;
    padding-left: 10px;
  }
  .setting-main-container-content {
    padding: 10px;
    .setting-section{
        margin: 6px auto;
        .setting-content-inner{
            padding:0;
          pre{
            display: flex;
            background-color: #f1f1f1;
            padding: 10px;
            width: 100%;
            border-radius: 3px;
            code{
            display: block;
              font-size: 12px;
              margin:0;
              padding:0;

            }
          }

        }
    }
  }

}
</style>
