<template>
  <div class="container">
    <section class="dict-settings">
      <div class="row">
        <div class="column">
          <label>词典id<i v-if="!readOnly">(生成)</i></label>
          <input type="text" class="form-control" disabled v-model="id" />
        </div>
        <div class="column" v-if="!readOnly">
          <label>词典id <i>(自定义,4-6字符)</i></label>
          <input
            type="text"
            placeholder="用户自定义"
            :disabled="readOnly"
            v-model="newid"
          />
        </div>
      </div>

      <div class="row" v-if="!readOnly">
        <p class="explain">
          词典id为唯一区分词典的标识，可由系统自动生成，如果用户想自定义id，请自行填写，用户自定义id优先级更高。
        </p>
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
        <div class="column">
          <label>词典别名</label>
          <input
            type="text"
            placeholder="MYDICT"
            v-model="alias"
            :disabled="readOnly"
          />
        </div>
      </div>

      <div class="row" v-if="!readOnly">
        <p class="explain">
          词典名称为用户可读词典名称，词典别名为词典短简称，用于展示在搜索框左侧，长度不可超过12字符。
        </p>
      </div>

      <div class="row">
        <div class="select-file-header">
          <span>mdx词典文件<i v-if="!readOnly">(必选一个)</i></span>
          <button
            class="btn btn-default"
            type="button"
            :disabled="readOnly"
            v-if="!readOnly"
            v-on:click="openMdxFile"
          >
            <i class="fas fa-folder"></i> 打开mdx文件
          </button>
        </div>
      </div>

      <div class="row">
        <div class="select-file-body">
          <p>{{ mdxpath }}</p>
        </div>
      </div>

      <div class="row">
        <div class="select-file-header">
          <span>mdd词典文件<i v-if="!readOnly">(零至多个)</i></span>

          <button
            class="btn btn-default"
            type="button"
            :disabled="readOnly"
            v-if="!readOnly"
            v-on:click="openMddFile"
          >
            <i class="fas fa-folder"></i> 打开mdd文件
          </button>
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
          <label>词典描述<i v-if="!readOnly">(可选)</i></label>
          <textarea
            class="form-control"
            rows="2"
            v-model="description"
            :disabled="readOnly"
          ></textarea>
        </div>
      </div>

      <div class="form-group" v-if="readOnly">
        <label>资源搜索</label>
        <div class="input-group">
          <input
            type="text"
            class="form-control"
            placeholder="resource key"
            aria-label="resource key"
            v-model="resourceKey"
          />
          <button
            class="btn btn-default"
            type="button"
            :disabled="!readOnly"
            v-on:click="searchResource(id)"
          >
            搜索
          </button>
        </div>
      </div>

      <div class="row">
        <Alert 
        :message="alertMessage" 
        :status="alertStatus" 
        :alertActive="alertShow"
        v-on:close-alert="closeAlert"></Alert>
      </div>

      <div class="btn-group" v-if="!readOnly">
        <button class="btn btn-default" v-on:click="showAlert('info', '测试信息 warn 类型', testCallback)">
          显示
        </button>
        <button class="btn btn-default" v-on:click="confirmAdd">添加</button>
      </div>

      <div v-if="showResourceBtn" class="btn-group">
        <button
          class="btn-control btn btn-sm btn-default"
          v-on:click="checkResource(id)"
        >
          <span class="btn-innertext">
          <i class="fas fa-archive"></i>
             查看缓存资源</span>
        </button>
        <button
          class="btn-control btn btn-sm btn-danger"
          v-on:click="deleteDict(id)"
        >
          <span class="icon icon-trash"></span>
          <span class="btn-innertext">
          <i class="fas fa-trash-alt"></i>
            删除词典</span>
        </button>
      </div>
    </section>
  </div>
</template>

<script lang="ts">
// import { SyncMainAPI, AsyncMainAPI } from '../../service.renderer.manifest';
import { listeners } from '../../service.renderer.listener';
import { random_key } from '../../../utils/random_key';
import { StorabeDictionary } from '../../../model/StorableDictionary';
import Alert from '../default/Alert.vue'

import Folder from '../icons/folder.icon.vue';
import Vue from 'vue';

const validate = function (dictData: StorabeDictionary) {
  if (!dictData) {
    return {
      ok: false,
      msg: `词典数据非法, dictData is null)`,
    };
  }
  if (!dictData.id || dictData.id === '' || dictData.id.length > 6) {
    return {
      ok: false,
      msg: `词典id非法: ${dictData.id} (长度<6/${!dictData.id ? 0 : dictData.id.length})`,
    };
  }

  if (!dictData.alias || dictData.alias === '' || dictData.alias.length > 12) {
    return {
      ok: false,
      msg: `词典别名(alias)非法: ${dictData.alias} (长度<12/${!dictData.alias ? 0 : dictData.alias.length}) `,
    };
  }

  if (!dictData.name || dictData.name === '' || dictData.name.length > 64) {
    return {
      ok: false,
      msg: `词典名称非法: ${dictData.name} (长度<64/${!dictData.name ? 0 :dictData.name.length})`,
    };
  }

  if (dictData.description && dictData.description.length > 512) {
    return {
      ok: false,
      msg: `词典描述长度非法: ${dictData.description} (长度<512/${!dictData.name ? 0 :dictData.description.length})`,
    };
  }

  if (!dictData.mdxpath || dictData.mdxpath === '') {
    return { ok: false, msg: `词典 mdxpath 非法: ${dictData.mdxpath}` };
  }

  // if (!dictData.mddpath || dictData.mddpath === '') {
  //   return { ok: false, msg: `词典 mddpath 非法: ${dictData.mddpath}` };
  // }
  if (!dictData.mddpath || dictData.mddpath === '') {
    dictData.mddpath = undefined;
  }

  return { ok: true, msg: '' };
};

export default Vue.extend({
  components: { Folder, Alert },
  props: {
    showResourceBtn: {
      type: Boolean,
      default: false,
    },
    readOnly: {
      type: Boolean,
      default: true,
    },
    dictData: {
      type: Object,
      default: () => {
        return {} as StorabeDictionary;
      },
    },
  },
  data() {
    return {
      alertMessage: '测试信息',
      alertStatus: 'info',
      alertShow: false,
      id: '',
      newid: '',
      alias: '',
      dictName: '',
      description: '',
      mdxpath: '',
      mddpath: '' as string | string[],
      resourceBaseDir: '',
      resourceKey: '',
    };
  },
  computed: {},
  watch: {},
  methods: {
    testCallback() {
      console.log('test callback!!')
    },
    closeAlert() {
      this.alertShow = false;
      this.$emit('alert-closed');
    },
    showAlert(status: string, msg: string, callback?: ()=>any) {
      this.alertMessage = msg;
      this.alertStatus = status;
      this.alertShow = true;
      if(callback) {
        let _this = this;
        let _handler = function() {
          callback();
          _this.$off('alert-closed', _handler);
        }
        this.$on('alert-closed', _handler);
      }
    },
    confirmAdd() {
      const newDict = new StorabeDictionary(
        this.newid || this.id,
        this.alias,
        this.dictName,
        this.mdxpath,
        this.mddpath,
        this.description
      );
      // validate dictionary data
      const ok = validate(newDict);
      if (!ok.ok) {
        this.showAlert('error', ok.msg);
        return;
      }

      console.log(newDict);
      let _this = this;

      this.$store
        .dispatch('asyncAddNewDict', newDict)
        .then((result) => {
          if (result) {
            this.showAlert('info', '添加成功', function() {
                  _this.$emit('modal-should-close');
                });
            this.resetForm();
          } else {
            this.showAlert('error', '添加失败, ID 重复');
          }
        })
        .catch((err) => {
          this.showAlert('error', err);
        });
    },
    openMddFile() {
      // TODO FIX
      // const resultPath = SyncMainAPI.syncShowOpenDialog({
      //   fileExtensions: ['mdd'],
      //   multiFile: true,
      // });
      // console.log(resultPath);
      // if (resultPath && resultPath.length > 0) {
      //   console.log(resultPath);
      //   this.mddpath = resultPath;
      // }
    },
    openMdxFile() {
      // TODO FIX
      // const resultPath = SyncMainAPI.syncShowOpenDialog({
      //   fileExtensions: ['mdx'],
      //   multiFile: false,
      // });
      // console.log(resultPath);
      // if (resultPath && resultPath.length > 0) {
      //   console.log(resultPath[0]);
      //   this.mdxpath = resultPath[0];
      // }
    },
    resetForm() {
      (this.id = random_key(6)), (this.newid = '');
      this.alias = '';
      this.dictName = '';
      this.description = '';
      this.mdxpath = '';
      this.mddpath = '';
      this.resourceBaseDir = '';
    },
    deleteDict(dictid: string) {
      // TODO FIX
      // let _this = this;
      // const response = SyncMainAPI.syncShowComfirmMessageBox({
      //   message: '确认删除?',
      //   type: 'warning',
      //   buttons: ['取消', '确认'],
      //   defaultId: 0,
      //   cancelId: 0,
      // });

      // if (response === 1) {
      //   this.$store
      //     .dispatch('asyncDelNewDict', dictid)
      //     .then((result) => {
      //       if (result) {
      //         console.log(result);
      //         if (result) {
      //           _this.showAlert('info', '删除成功', function() {
      //             _this.$emit('modal-should-close');
      //           });
      //         }
      //       } else {
      //         _this.showAlert('error', '删除失败');
      //       }
      //     })
      //     .catch((err) => {
      //       _this.showAlert('error', '删除失败,' + err);
      //     });
      // } else {
      //   console.log('canceled');
      // }
    },
    checkResource(dictid: string) {
      // TODO FIX
      // const response = AsyncMainAPI.openDictResourceDir(dictid);
      // console.log(response);
    },
    searchResource(dictid: string) {
      // TODO FIX
      // if (!this.resourceKey || this.resourceKey == '') {
      //   this.showAlert('info', 'null');
      //   return;
      // }
      // AsyncMainAPI.loadDictResource({
      //   dictid: dictid,
      //   resourceKey: this.resourceKey,
      // });
    },
  },
  mounted() {
    this.$nextTick(function () {
      this.id = this.dictData.id;
      this.alias = this.dictData.alias;
      this.dictName = this.dictData.name;
      this.description = this.dictData.description;
      this.mdxpath = this.dictData.mdxpath;
      this.mddpath = this.dictData.mddpath;
      this.resourceBaseDir = this.dictData.resourceBaseDir;
      // listener
      listeners.onLoadDictResource((event, arg) => {
        console.log(arg);
        this.showAlert(
          'info',
          `${arg.keyText} (${arg.contentSize}) [${arg.definition}]`
        );
      });
    });
  },
});
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