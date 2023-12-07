<template>
  <div class="container">
    <div class="toolbar">
      <n-grid x-gap="12" :cols="3">
        <n-gi>
          <n-select v-model:value="selectedDict" :options="optionDicts" />
        </n-gi>
        <n-gi :span="2">
          <n-input-group>
            <n-input
              :style="{ width: '70%' }"
              v-model:value="resourceInputValue"
            />
            <n-button type="primary" ghost @click="searchResource">
              搜索
            </n-button>
          </n-input-group>
        </n-gi>
      </n-grid>
    </div>

    <div class="content">
      <ul>
        <li>URL: {{ reqURL }}</li>
        <li>HTTP STATUS: {{ result.status_code }}</li>
        <li>CONTENT LENGTH: {{result.content_length}}</li>
        <li>CONTENT TYPE: {{result.content_type}} </li>
        <li>MSG: {{result.resp_msg}} </li>
      </ul>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  padding: 10px;
  .toolbar {
    margin-bottom: 10px;
  }
  .content {
    display: flex;
    width: 100%;
    min-height: 120px;
  }
}
</style>

<script lang="ts" setup>
import { defineComponent, ref, reactive, onMounted, computed } from 'vue';
import { StaticDictServerURL } from '@/apis/apis';
import { useDictQueryStore } from '@/store/dict';
import axios from 'axios';

const dictQueryStore = useDictQueryStore();
let selectedDict = ref(null);
let optionDicts = reactive([]);
let staticServerUrl = ref('');
let resourceInputValue = ref('');
let result = reactive({
  status_code: 0,
  content_type:"",
  content_length: 0,
  resp_msg:"",
});

let reqURL = computed(() => {
  return composeReqURL();
});

function composeReqURL() {
  var req_url = staticServerUrl.value;
  if (req_url === '') {
    return '';
  }

  var input = resourceInputValue.value;
  if (input.startsWith("http")) {
    return input
  }

  if (input && input !== "" && input.startsWith('/')) {
    let end = req_url.indexOf('/__mdict');
    req_url = req_url.substr(0, end);
    console.log(req_url)
    input = input.substr(1);
  }

  req_url = req_url + '/' + input + '?dict_id=' + selectedDict.value + "&d=0";
  return req_url;
}

function searchResource() {
  console.log('search value', resourceInputValue);
  if (resourceInputValue.value === '') {
    return;
  }

  var req_url = composeReqURL();
  console.log(req_url);

  axios.get(req_url).then((res) => {
    console.log(res);
    result.content_length = res.headers['content-length'];
    result.content_type = res.headers['content-type'];
    result.resp_msg = "success";
  }).catch((err) => {
    result.status_code = 404;
    result.resp_msg = err;
    console.log(err);
  })
  
}

/**
 * init loadings
 */
onMounted(() => {
  if (staticServerUrl.value === '') {
    StaticDictServerURL().then((url) => {
      staticServerUrl.value = url;
    });
  }

  dictQueryStore.queryDictList().then((dicts) => {
    for (let index = 0; index < dicts.length; index++) {
      const dict = dicts[index];
      optionDicts.push({
        label: dict.name,
        value: dict.id,
      });
    }
    if (optionDicts.length > 0 && selectedDict.value === null) {
      selectedDict.value = optionDicts[0].value;
    }
  });
});
</script>
