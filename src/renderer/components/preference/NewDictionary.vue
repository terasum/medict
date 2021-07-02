<template>
  <form>
    <div class="form-group row">
      <div class="col-6">
        <label>词典id(生成)</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="showDictData.dictIdGen"
        />
      </div>
      <div class="col-6">
        <label>词典id(4英文字符)</label>
        <input
          type="text"
          class="form-control"
          placeholder="用户自定义"
          v-model="showDictData.dictIdCustom"
          :disabled="readOnly"
        />
      </div>
    </div>

    <div class="form-group">
      <label>词典名称</label>
      <input
        type="text"
        class="form-control"
        placeholder="英汉大辞典"
        v-model="showDictData.dictName"
        :disabled="readOnly"
      />
    </div>
    <div class="form-group">
      <label>mdx词典文件</label>
      <input
        v-if="readOnly"
        type="text"
        class="form-control"
        :value="showDictData.dictMdxFilePath"
        :disabled="readOnly"
      />
      <input v-else type="file" class="form-control" />
    </div>
    <div class="form-group">
      <label>mdd词典文件</label>
      <input
        v-if="readOnly"
        type="text"
        class="form-control"
        :value="showDictData.dictMddFilePath"
        :disabled="readOnly"
      />
      <input v-else type="file" class="form-control" />
    </div>
    <div class="form-group">
      <label>词典描述</label>
      <textarea
        class="form-control"
        rows="2"
        v-model="showDictData.dictDescription"
        :disabled="readOnly"
      ></textarea>
    </div>
    <div v-if="showResourceBtn" class="form-group">
      <div class="row">
        <div class="col-4">
          <button class="btn-control btn btn-sm btn-default">
            <span class="icon icon-archive"></span>
            <span class="btn-innertext"> 检视资源</span>
          </button>
        </div>
        <div class="col-4">
          <button class="btn-control btn btn-sm btn-default">
            <span class="icon icon-target"></span>
            <span class="btn-innertext">文件自检</span>
          </button>
        </div>
        <div class="col-4">
          <button class="btn-control btn btn-sm btn-danger">
            <span class="icon icon-trash"></span>
            <span class="btn-innertext">删除词典</span>
          </button>
        </div>
      </div>
    </div>
  </form>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  components: {},
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
      default: {
        dictIdGen: '',
        dictIdCustom: '',
        dictName: '',
        dictMdxFilePath: '',
        dictMddFilePath: '',
        dictDescription: '',
      },
    },
  },
  data() {
    return {
      showDictData: {
        dictIdGen: '',
        dictIdCustom: '',
        dictName: '',
        dictMdxFilePath: '',
        dictMddFilePath: '',
        dictDescription: '',
      },
    };
  },
  computed: {},
  watch: {},
  methods: {},
  mounted() {
    this.$nextTick(function () {
      this.showDictData = this.dictData;
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

.btn-default {
  background-color: #ccc;
  border: 1px solid #999;
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
</style>