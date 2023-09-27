<style lang="scss">
@import '@/style/variables.scss';
    .header-search-box {
      display: flex;
      height: 54px;
      padding: 0;
      margin: 0;

      .header-navigate-btns {
        height: 54px;
        max-width: 120px;
        padding: 0;
        margin-left: 16px;
        margin-right: 14px;
        display: flex;

        .btn-nav {
          height: 26px;
          width: 26px;
          margin-top: 14px;
          padding: 0;
          text-align: center;
          font-size: 12px;
          color: #333;
          outline: none;
          border: 1px solid #fefefe;
          background-color: #fff;
          box-shadow: none;

          &:active {
            box-shadow: none;
            border: rgba(63, 80, 236, 0.452) 1px solid;
            // background-color: #d80034;
            background-color: #fff;
          }
        }
        .btn-nav-left {
          border-radius: 10px 0px 0px 10px;
          margin-left: 9px;
        }
        .btn-nav-right {
          border-radius: 0px 10px 10px 0px;
          margin-left: 0px;
        }
      }
      .header-search-input {
        height: 54px;
        display: flex;
        flex-direction: column;
        justify-content: center;

        .n-input {
          height: 26px;
          padding: 0 8px;
          margin: 0;
          box-shadow: none;
          font-size: 15px;
          // border: 1px solid #f1f1f1;
          border: none;

          background-color: #fff;
          padding-left: 5px;
          &:active {
            outline: none;
          }
          &:focus {
            outline: none;
          }
        }
      }
    }

</style>
<template>
  <AppFunctions>
   <div class="header-search-box">
        <div class="header-navigate-btns">
          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-left"
          >
            <n-icon><AngleLeft /></n-icon>
          </button>

          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-right"
          >
            <n-icon><AngleRight /></n-icon>
          </button>
        </div>
        <div class="header-search-input">
          <n-input
            type="text"
            size="small"
            placeholder="搜索"
            @keydown.enter="handleChange"
            v-model:value="inputWord"
          >
            <template #suffix>
              <n-icon :component="Search" />
            </template>
          </n-input>
        </div>
      </div>
  </AppFunctions>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { NIcon } from 'naive-ui';
import { Search, AngleLeft, AngleRight } from '@vicons/fa';
import AppFunctions from '@/components/layout/AppFunctions.vue';

import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';
import { useRouter } from "vue-router";

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();
const router = useRouter();

let inputWord = ref('');
let inputActive = ref(false);
let componentEventListenerUnscribe = null;


function listenInputWordUpdate() {
  componentEventListenerUnscribe = dictQueryStore.$onAction(
    ({
      name, // action 名称
      store, // store 实例，类似 `someStore`
      args, // 传递给 action 的参数数组
      after, // 在 action 返回或解决后的钩子
      onError, // action 抛出或拒绝的钩子
    }) => {
      // 这将在 action 成功并完全运行后触发。
      // 它等待着任何返回的 promise
      after((result) => {
        switch (name) {
          case 'updateInputSearchWordRaw': {
            inputWord.value = dictQueryStore.inputSearchWord;
            console.log("update inputWord by event", inputWord.value);
            break;
          }
        }
      });

      // 如果 action 抛出或返回一个拒绝的 promise，这将触发
      onError((error) => {
        
      });
    }
  );

  // 手动删除监听器
  //   unsubscribe()
}

onMounted(() => {
  if (componentEventListenerUnscribe) {
    componentEventListenerUnscribe();
    componentEventListenerUnscribe = null;
  }
  listenInputWordUpdate();
})

///----------------------------
// event listener function
///----------------------------

function handleChange(v) {
  console.info('[app-event](keydown.enter), args:' + inputWord.value);
  if(!uiStore.isSearchInputActive()){
    console.log("[app-event](keydown.enter), input disabled, skipped")
    return;
  }
  dictQueryStore.updateInputSearchWord(inputWord.value.trim());
}
</script>
