<template>
  <div class="app-modal-container">
    <b-modal
      custom-class="app-modal"
      v-model="isActive"
      v-on:close="onClose"
      button-size="sm"
      :width="540"
    >
      <div class="app-modal-header">
        <div class="app-modal-header-left">
          <slot name="modal-header"></slot>
        </div>
        <div class="app-modal-header-right">
          <button class="btn btn-close-modal" @click="close"><i class="fas fa-times"></i></button>
        </div>
      </div>

      <slot
      ></slot>

      <div class="app-modal-footer">
        <slot name="modal-footer"> </slot>
      </div>
    </b-modal>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  components: {},
  props: {
    active: {
      type: Boolean,
    },
  },
  // deprecated, to replace with default 'value' in the next breaking change
  model: {
    prop: 'active',
    event: 'update:active',
  },
  data() {
    return {
      isActive: this.active || false,
    };
  },
  watch: {
    active(value) {
      this.isActive = value;
    },
    isActive(value) {
      this.isActive = value;
      this.$emit('update:active', value);
    },
  },
  methods: {
    closeclose() {
      console.log('CLOSECLOSE111111');
    },
    onClose() {
      console.log('modal-close triggered');
      this.$emit('modal-close', this.isActive);
    },
    close() {
      this.isActive = false;
    },
  },
  mounted() {
    let _this = this;
    this.$nextTick(function() {
      console.log('modal-should-close listen');
      _this.$on('modal-should-close', function(){
          console.log('modal-should-close triggerer!!');
        _this.close();
      });
    })
    
  },
});
</script>

<style lang="scss" scoped>
.app-modal {
  background-color: #f3f3f3d3;
  position: fixed;
  left: 0;
  top: 60px;
  width: 100%;
  height: calc(100% - 60px);
  padding: 20px 16px;
  .btn-close-modal{
    border:none;
    outline:none;
    color: #999;
    font-size: 20px;
    line-height: 24px;
    font-weight: lighter;
    background: transparent;
    &:focus{
    color: #333;
    }
    &:active{
    color: #333;
    }
  }

  &::v-deep .modal-background {
    display: none;
  }

  &::v-deep .modal-content {
    background: #fff;
    margin: 0 auto;
    border: 1px solid #cfcfcf;
    margin: 0 auto;
    padding: 10px 15px;
    border-radius: 3px;
    box-shadow: 1px 1px 3px #c1c1c1;
    overflow: hidden;
    max-height: 480px;
  }

  &::v-deep .app-modal-header {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    height: 24px;
    padding: 2px 0 2px 10px;
    height: 40px;
    margin-bottom: 10px;
    border-bottom: 1px solid #879fb6;

    h3 {
      font-size: 20px;
      line-height: 24px;
      color: #3a75c4;
    }

    .app-modal-header-left {
      padding: 0.3rem 0.5rem !important;
      button {
        border: none;
        outline: none;
        width: 22px;
        height: 22px;
        color: #666;
        text-align: center;
        font-size: 18px;
        line-height: 18px;
        padding: 0.1rem;
        border-radius: 5px;
        &:hover {
          color: #999;
          border: 1px solid #ddd;
        }
      }
      .modal-title {
        font-size: 14px;
        line-height: 24px;
        font-weight: normal;
        padding-left: 10px;
      }
    }
    .app-modal-header-right {
      justify-self: right;
    }
  }

  &::v-deep .modal-close {
    display: none;
    &::before {
      content: '确认';
    }
  }
  &::v-deep .app-modal-footer {
    padding: 0.3rem 0.5rem !important;
  }
}
</style>
