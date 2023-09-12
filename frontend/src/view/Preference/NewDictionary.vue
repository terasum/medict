<template>
<!--  <div class="container">
    <section class="dict-settings">
      <div class="row">
        <div class="column">
          <label>词典id<i v-if="!readOnly">(生成)</i></label>
          <input type="text" class="form-control" disabled v-model="id" />
        </div>
      </div>
      <div class="row">
        <div class="column">
          <label>词典名称</label>
          <input
            type="text"
            placeholder="英汉大辞典"
            v-model="dictName"
            :disabled="readOnly"
          />
        </div>
      </div>
      <div class="row">
        <div class="select-file-header">
          <span>mdx词典文件</span>
        </div>
      </div>
      <div class="row">
        <div class="select-file-body">
          <p>{{ mdxpath }}</p>
        </div>
      </div>
      <div class="row">
        <div class="select-file-header">
          <span>mdd词典文件</span>
        </div>
      </div>
      <div class="row">
        <div class="select-file-body">
          <ul class="">
            <li v-for="path_item in mddpath" :key="path_item">
              {{ path_item }}
            </li>
          </ul>
        </div>
      </div>
      <div class="row">
        <div class="column">
          <label>词典描述</label>
          <textarea
            class="form-control"
            rows="3"
            v-model="description"
            :disabled="readOnly"
          ></textarea>
        </div>
      </div>

      <div class="row">
      <div v-if="showResourceBtn" class="btn-group">
        <button
          class="btn-control btn btn-sm btn-default"
          v-on:click="checkResource(id)"
        >
          <span class="btn-innertext">
          <i class="fas fa-archive"></i>
             查看缓存词典资源</span>
        </button>
      </div>
      </div>
    </section>
  </div> -->
</template>

<script lang="ts">
// import { StorabeDictionary } from '@/model/StorableDictionary';
// import Alert from '@/components/default/Alert.vue'

// import Folder from '@/components/icons/folder.icon.vue';
// import {defineComponent} from "vue";
// export default defineComponent({

//   components: { Folder, Alert },
//   props: {
//     showResourceBtn: {
//       type: Boolean,
//       default: false,
//     },
//     readOnly: {
//       type: Boolean,
//       default: true,
//     },
//     dictData: {
//       type: Object,
//       default: () => {
//         return {} as StorabeDictionary;
//       },
//     },
//   },
//   data() {
//     return {
//       alertMessage: '测试信息',
//       alertStatus: 'info',
//       alertShow: false,
//       id: '',
//       newid: '',
//       alias: '',
//       dictName: '',
//       description: '',
//       mdxpath: '',
//       mddpath: '' as string | string[],
//       resourceBaseDir: '',
//       resourceKey: '',
//     };
//   },
//   computed: {},
//   watch: {},
//   methods: {
//     resetForm() {
//       // (this.id = random_key(6)), (this.newid = '');
//       this.alias = '';
//       this.dictName = '';
//       this.description = '';
//       this.mdxpath = '';
//       this.mddpath = '';
//       this.resourceBaseDir = '';
//     },
//     checkResource(dictid: string) {
//       // TODO FIX
//       // windowOpenApi.openDictResourceDir(dictid)
//     },
    
//   },  mounted() {
//     this.$nextTick(function () {
//       this.id = this.dictData.id;
//       this.alias = this.dictData.alias;
//       this.dictName = this.dictData.name;
//       this.description = this.dictData.description;
//       this.mdxpath = this.dictData.mdxpath;
//       this.mddpath = this.dictData.mddpath;
//       this.resourceBaseDir = this.dictData.resourceBaseDir;
//     });
//   },
// });
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  max-height: 420px;
  overflow-y: auto;
}

.row {
  display: flex;
  flex-direction: row;
  width: 100%;
  margin-top: 2px;
  margin-bottom: 2px;
}

.column {
  display: flex;
  flex-direction: column;
  width: 100%;
  padding: 1px;
}

.dict-settings label {
  margin-bottom: 3px;
  color: #777;
  font-size: 12px;
}

.dict-settings input {
  &:focus {
    outline: none;
  }
  border: 1px solid #c1c1c1;
  border-radius: 3px;
  // height: 32px;
  font-size: 14px;
  padding-left: 4px;
}
.dict-settings .explain {
  color: #777;
  font-size: 12px;
  padding: 2px 4px;
}

.btn {
  background-color: rgb(252, 252, 252);
  background-image: linear-gradient(to bottom, #fcfcfc 0, #f1f1f1 100%);
  border-radius: 4px;
  border: 1px solid transparent;
  box-shadow: 0 1px 1px rgb(0 0 0 / 6%);
  border-color: #c2c0c2 #c2c0c2 #a19fa1;
  text-align: center;
  color: #666;
  font-size: 12px;
  display: inline-block;
  padding: 3px 8px;
  margin-bottom: 0;
  white-space: nowrap;
  &:active {
    background-color: #ddd;
    background-image: none;
  }
}

.select-file-header {
  height: 26px;
  font-size: 12px;
  color: #777;
  display: flex;
  width:100%;
  justify-content: space-between;
  span{
    font-size: 12px;
    line-height: 26px;
    display: block;
  }
}

.select-file-body {
  width: 100%;
  font-size: 14px;
  color: #666;
  border: 1px solid #ccc;
  min-height: 30px;
  border-radius: 2px;
  background: #f1f1f1;
  padding-left: 3px;
  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    li {
      padding: 0;
      margin: 0;
    }
  }
}

</style>