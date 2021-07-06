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
        <label>mdx词典文件</label>
        <div class="input-group">
          <input
            type="text"
            class="form-control"
            placeholder="mdx path"
            aria-label="mdx path"
            v-model="mdxpath"
            :disabled="readOnly"
          />
          <button
            class="btn btn-default"
            type="button"
            :disabled="readOnly"
            v-on:click="openMdxFile"
          >
            <Folder :width="16" :height="16" />
          </button>
        </div>
      </div>

      <div class="form-group">
        <label>mdd词典文件</label>
        <div class="input-group">
          <input
            type="text"
            class="form-control"
            placeholder="mdd path"
            aria-label="mdd path"
            v-model="mddpath"
            :disabled="readOnly"
          />
          <button
            class="btn btn-default"
            type="button"
            :disabled="readOnly"
            v-on:click="openMddFile"
          >
            <Folder :width="16" :height="16" />
          </button>
        </div>
      </div>

      <div class="form-group">
        <label>词典描述</label>
        <textarea
          class="form-control"
          rows="2"
          v-model="description"
          :disabled="readOnly"
        ></textarea>
      </div>
    </form>
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
      <button class="btn btn-default" v-on:click="showAlert('info', '测试')">
        显示
      </button>
      <button class="btn btn-default" v-on:click="closeModal">取消</button>
      <button class="btn btn-default btn-primary" v-on:click="confirmAdd">
        添加词典
      </button>
    </div>
    <div v-if="showResourceBtn" class="btn-group">
      <button class="btn-control btn btn-sm btn-default">
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
import { SyncMainAPI } from '../../service.renderer.manifest';
import Folder from '../icons/folder.icon.vue';
import Vue from 'vue';
import { StorabeDictionary } from '../../../model/StorableDictionary';
import { random_key } from '../../../utils/random_key';

const validate = function (dictData: StorabeDictionary) {
  if (!dictData.id || dictData.id === '' || dictData.id.length > 6) {
    return { ok: false, msg: `词典id非法: ${dictData.id}` };
  }

  if (!dictData.alias || dictData.alias === '' || dictData.alias.length > 12) {
    return { ok: false, msg: `词典 alias 非法: ${dictData.alias}` };
  }

  if (!dictData.name || dictData.name === '' || dictData.name.length > 64) {
    return { ok: false, msg: `词典 name 非法: ${dictData.name}` };
  }

  if (dictData.description && dictData.description.length > 512) {
    return {
      ok: false,
      msg: `词典 description 长度 非法: ${dictData.description}`,
    };
  }

  if (!dictData.mdxpath || dictData.mdxpath === '') {
    return { ok: false, msg: `词典 mdxpath 非法: ${dictData.mdxpath}` };
  }

  if (!dictData.mddpath || dictData.mddpath === '') {
    return { ok: false, msg: `词典 mddpath 非法: ${dictData.mddpath}` };
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
      mddpath: '',
      resourceBaseDir: '',
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
      if (closeBtn) {
        // @ts-ignore
        closeBtn.click();
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
      const result = SyncMainAPI.dictAddOne({ dict: newDict });
      if (result) {
        this.showAlert('info', '添加成功');
        this.resetForm();
      } else {
        this.showAlert('danger', '添加失败, ID 重复');
      }
    },
    openMddFile() {
      const resultPath = SyncMainAPI.syncShowOpenDialog(['mdd']);
      console.log(resultPath);
      if (resultPath && resultPath.length > 0) {
        console.log(resultPath[0]);
        this.mddpath = resultPath[0];
      }
    },
    openMdxFile() {
      const resultPath = SyncMainAPI.syncShowOpenDialog(['mdx']);
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
      const result = SyncMainAPI.dictDeleteOne({ dictid: dictid });
      console.log(result);
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
    });
  },
});
</script>

<style lang="scss" scoped>
.form-group {
  label {
    color: #737475;
    margin-bottom: 0;
  }
}
.btn-group {
  float: right;
}

.btn-default {
  background-color: #ccc;
  border: 1px solid #999;
  padding: 0rem 0.475rem;
  min-height: 26px;
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
  background-color: red;
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