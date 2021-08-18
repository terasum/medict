<template>
  <div>
    <form>
      <div class="form-group row">
        <div class="col-6">
          <label>词典id<i v-if="!readOnly">(生成)</i></label>
          <input type="text" class="form-control" disabled v-model="id" />
        </div>
        <div class="col-6" v-if="!readOnly">
          <label>词典id <i>(自定义,4-6字符)</i></label>
          <input
            type="text"
            class="form-control"
            placeholder="用户自定义"
            :disabled="readOnly"
            v-model="newid"
          />
        </div>
      </div>

      <div class="form-group row">
        <div class="col-6">
          <label>词典名称</label>
          <input
            type="text"
            class="form-control"
            placeholder="英汉大辞典"
            v-model="dictName"
            :disabled="readOnly"
          />
        </div>
        <div class="col-6">
          <label>词典别名</label>
          <input
            type="text"
            class="form-control"
            placeholder="英汉大辞典"
            v-model="alias"
            :disabled="readOnly"
          />
        </div>
      </div>

      <div class="form-group">
        <div class="select-file-header">
          <span>mdx词典文件<i v-if="!readOnly">(必选一个)</i></span>
          <button
            class="btn btn-default btn-open"
            type="button"
            :disabled="readOnly"
            v-on:click="openMdxFile"
          >
          打开
            <Folder :width="16" :height="16" />
          </button>
        </div>

        <div class="select-file-body">
          <p>{{mdxpath}}</p>
        </div>
      </div>

      <div class="form-group">
        <div class="select-file-header">
        <span>mdd词典文件<i v-if="!readOnly">(零至多个)</i></span>

          <button
            class="btn btn-default"
            type="button"
            :disabled="readOnly"
            v-on:click="openMddFile"
          >
          打开
            <Folder :width="16" :height="16" />
          </button>
        </div>

        <div class="select-file-body">
            <ul class="">
              <li v-for="path_item in mddpath" :key="path_item">
                {{ path_item }}
              </li>
            </ul>
        </div>

      </div>

      <div class="form-group">
        <label>词典描述<i v-if="!readOnly">(可选)</i></label>
        <textarea
          class="form-control"
          rows="2"
          v-model="description"
          :disabled="readOnly"
        ></textarea>
      </div>
    </form>

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


    <div
      id="input-alert"
      class="alert-embbed alert alert-dismissible fade"
      :class="(alertShow ? 'show' : 'hide') + ' alert-' + alertStatus"
      role="alert"
    >
      <span>{{ alertMessage }}</span>
      <button
        type="button"
        class="alert-btn btn-close"
        data-bs-dismiss="alert"
        aria-label="Close"
        v-on:click="hideAlert"
      ></button>
    </div>
    <div class="btn-group" v-if="!readOnly">
      <!-- <button class="btn btn-default" v-on:click="showAlert('info', '测试')">
        显示
      </button> -->
      <button class="btn btn-default" v-on:click="closeModal">取消</button>
      <button class="btn btn-default btn-primary" v-on:click="confirmAdd">
        添加
      </button>
    </div>
    <div v-if="showResourceBtn" class="btn-group">
      <button class="btn-control btn btn-sm btn-default" v-on:click="checkResource(id)">
        <span class="icon icon-archive"></span>
        <span class="btn-innertext"> 检视资源</span>
      </button>
      <button class="btn-control btn btn-sm btn-default">
        <span class="icon icon-target"></span>
        <span class="btn-innertext">文件自检</span>
      </button>
      <button
        class="btn-control btn btn-sm btn-danger"
        v-on:click="deleteDict(id)"
      >
        <span class="icon icon-trash"></span>
        <span class="btn-innertext">删除词典</span>
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { SyncMainAPI, AsyncMainAPI } from '../../service.renderer.manifest';
import {listeners} from '../../service.renderer.listener';
import { random_key } from '../../../utils/random_key';
import { StorabeDictionary } from '../../../model/StorableDictionary';

import Folder from '../icons/folder.icon.vue';
import Vue from 'vue';

const validate = function (dictData: StorabeDictionary) {
  if (!dictData.id || dictData.id === '' || dictData.id.length > 6) {
    return { ok: false, msg: `词典id非法: ${dictData.id} (长度<6/${dictData.id.length})` };
  }

  if (!dictData.alias || dictData.alias === '' || dictData.alias.length > 12) {
    return { ok: false, msg: `词典别名非法: ${dictData.alias} (长度<12/${dictData.alias.length}) ` };
  }

  if (!dictData.name || dictData.name === '' || dictData.name.length > 64) {
    return { ok: false, msg: `词典名称非法: ${dictData.name} (长度<64/${dictData.name.length})` };
  }

  if (dictData.description && dictData.description.length > 512) {
    return {
      ok: false,
      msg: `词典描述长度非法: ${dictData.description} (长度<512/${dictData.description.length})`,
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
  components: { Folder },
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
      mddpath: '' as string|string[],
      resourceBaseDir: '',
      resourceKey:'',
    };
  },
  computed: {},
  watch: {},
  methods: {
    showAlert(status: string, msg: string) {
      this.alertMessage = msg;
      this.alertStatus = status;
      this.alertShow = true;
    },
    hideAlert() {
      this.alertShow = false;
    },
    closeModal() {
      let closeBtn = document.getElementById(
        'add-dictionary-modal___BV_modal_title_'
      )?.nextSibling;

      let closeBtn2 = document.getElementById(
        'dictionary-item-modal___BV_modal_title_'
      )?.nextSibling;

      if (closeBtn) {
        // @ts-ignore
        closeBtn.click();
      }

      if (closeBtn2) {
        // @ts-ignore
        closeBtn2.click();
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
        this.showAlert('danger', ok.msg);
        return;
      }

      console.log(newDict);

      this.$store
        .dispatch('asyncAddNewDict', newDict)
        .then((result) => {
          if (result) {
            this.showAlert('info', '添加成功');
            this.resetForm();
          } else {
            this.showAlert('danger', '添加失败, ID 重复');
          }
        })
        .catch((err) => {
          this.showAlert('danger', err);
        });
    },
    openMddFile() {
      const resultPath = SyncMainAPI.syncShowOpenDialog({fileExtensions: ['mdd'], multiFile: true});
      console.log(resultPath);
      if (resultPath && resultPath.length > 0) {
        console.log(resultPath);
        this.mddpath = resultPath;
      }
    },
    openMdxFile() {
      const resultPath = SyncMainAPI.syncShowOpenDialog({fileExtensions: ['mdx'], multiFile: false});
      console.log(resultPath);
      if (resultPath && resultPath.length > 0) {
        console.log(resultPath[0]);
        this.mdxpath = resultPath[0];
      }
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
      const response = SyncMainAPI.syncShowComfirmMessageBox({
        message: '确认删除?',
        type: 'warning',
        buttons: ['取消', '确认'],
        defaultId: 0,
        cancelId: 0,
      });
      console.log(response);
      // response.then((resp: { response: number; checkboxChecked: boolean }) => {
      if (response === 1) {
        this.$store
          .dispatch('asyncDelNewDict', dictid)
          .then((result) => {
            if (result) {
              console.log(result);
              if (result) {
                this.showAlert('info', '删除成功');
                let that = this;
                setTimeout(() => {
                  that.closeModal();
                }, 300);
              }
            } else {
              this.showAlert('danger', '删除失败');
            }
          })
          .catch((err) => {
            this.showAlert('danger', '删除失败,' + err);
          });
      } else {
        console.log('canceled');
      }
      // });
    },
     checkResource(dictid: string) {
      const response = AsyncMainAPI.openDictResourceDir(dictid);
      console.log(response);
    },
    searchResource(dictid: string) {
      this.hideAlert();
      if(!this.resourceKey || this.resourceKey == '') {
        this.showAlert('info', 'null');
        return;
      }
     AsyncMainAPI.loadDictResource({ dictid: dictid, resourceKey: this.resourceKey } ); 
    }
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
      listeners.onLoadDictResource((event, arg) =>{
         console.log(arg);
         this.showAlert('info', `${arg.keyText} (${arg.contentSize}) [${arg.definition}]`);
      })
    });
  },
});
</script>

<style lang="scss" scoped>
.form-group {
  label {
    color: #737475;
    margin-bottom: 0;
    font-size: 0.875rem;
  }
}
.btn-group {
  float: right;
}

.btn-open {
  height: 26px;

}  

.btn-default {
  background-color: #fff;
  color: #666;
  border: 1px solid #999;
  padding: 0rem 0.475rem;
  min-height: 22px;
  line-height: 22px;
  font-size:0.875rem;

  .icon {
    color: #737475;
  }
  .btn-innertext {
    color: #737475;
  }
  &:hover {
    box-shadow: none;
    cursor: pointer;
  }
  &:focus {
    box-shadow: none;
  }
}

.btn-danger {
  background-color: linear-gradient(to bottom, #bd3134 0, #d84042 100%);
  border: 1px solid #999;
  color: #fff;
  .icon {
    color: #fff;
  }
  .btn-innertext {
    color: #fff;
  }
  &:hover {
    box-shadow: none;
    cursor: pointer;
  }
  &:focus {
    box-shadow: none;
  }
}

.btn-control {
  font-size: 13px;
  line-height: 13px;
  &:hover {
    box-shadow: none;
    cursor: pointer;
  }
  &:focus {
    box-shadow: none;
  }
  .icon {
    float: none;
  }
  .btn-innertext {
    padding: 0 3px 0 5px;
  }
}

.select-file-header {
  height: 26px;
  font-size: 0.875rem;
  color: #666;
  margin-top: 10px;;
  span{
      display: inline-block;
      margin-right: 5px;;
      height: 20px;
      line-height: 20px;
      padding-top:1px;
  }
  button{
    float:right;
    height: 20px;
    line-height: 20px;
  }
}

.select-file-body {
  font-size:0.875rem;
  color:#666;
  border: 1px solid #ccc;
  min-height: 30px;
  border-radius: 2px;
  background: #f1f1f1;
  padding-left: 3px;
  ul{
    list-style: none;
    padding:0;
    margin:0;
    li{
      padding:0;
      margin:0; 
    }
  }
}


.alert-embbed {
  app-region: no-drag;
  -webkit-app-region: no-drag;
  padding: 0.125rem 0.375rem;
  font-size: 14px;
  .alert-btn {
    padding: 0;
    font-size: 0.6rem;
    line-height: 0.6rem;
    height: 26px;
    width: 26px;
  }
}
.alert-hide {
  display: none;
}
.alert-show {
  display: block;
}
</style>