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
  <div class="alert" :class="alertActive ? 'alert-show' : 'alert-hide'">
    <div class="alert-content" :style="'margin-top: ' + top + 'px;'">
      <div class="alert-content-header">
        <span class="alert-icon" :class="statusColor">
          <i class="fas fa-exclamation-circle"></i>
        </span>
        <button type="button" class="btn-alert-close" v-on:click="close">
          <i class="fas fa-times"></i>
        </button>
      </div>

      <div class="alert-content-body">
        <span>{{ message }}</span>
      </div>

      <div class="alert-content-footer">
        <button type="button" class="btn-alert-ok" v-on:click="close">
          OK
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";

export default defineComponent({

  components: {},
  props: {
    top: {
      type: Number,
      default: 100,
    },
    message: {
      type: String,
      default: 'default',
    },
    status: {
      type: String,
    },
    alertActive: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {};
  },
  computed: {
    statusColor() {
      switch (this.status || 'info') {
        case 'info': {
          return 'alert-color-info';
        }
        case 'warn': {
          return 'alert-color-warn';
        }
        case 'error': {
          return 'alert-color-error';
        }
        default: {
          return 'alert-color-info';
        }
      }
    },
  },
  watch: {},
  methods: {
    close() {
      this.$emit('close-alert');
    },
  },
});
</script>

<style lang="scss" scoped>
.alert {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #5a6da6b9;
}

.alert-icon {
  padding-left: 10px;
  width: 30px;
  height: 30px;
  font-size: 20px;
  line-height: 30px;
  color: #919191;
}

.alert-content {
  width: 420px;
  background: #fff;
  margin: 0 auto;
  border-radius: 3px;
  box-shadow: 1px 2px 3px #3b5e91de;
  color: #aeaeae;
  font-size: 14px;
}

.alert-content-header {
  height: 30px;
  border-radius: 3px 3px 0 0;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-bottom: 1px solid #c1c1c1;
}

.btn-alert-close {
  display: none;
  font-size: 20px;
  font-weight: lighter;
  border: none;
  outline: none;
  color: rgb(165, 165, 165);
  background: transparent;
}

.alert-content-body {
  min-height: 120px;
  padding: 10px 30px;
}

.alert-content-footer {
  height: 40px;
  border-radius: 0 0 3px 3px;
  display: flex;
  flex-direction: row-reverse;
  background: #f9f9ff;
  padding-top: 8px;
}

.alert-color-info {
  color: #1e7efc;
}
.alert-color-warn {
  color: #fcae1e;
}
.alert-color-error {
  color: #fc381e;
}

.btn-alert-ok {
  width: 60px;
  height: 24px;
  padding: 4px 18px;
  margin-right: 30px;

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

.alert-hide {
  display: none;
}

.alert-show {
  display: block;
}
</style>
