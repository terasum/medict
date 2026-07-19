<template>
  <div class="resource-inspector">
    <section class="request-section" aria-labelledby="resource-request-title">
      <div class="section-heading">
        <div>
          <h2 id="resource-request-title">资源查询</h2>
          <p>检查词典服务器中的样式、图片、字体及其他静态资源。</p>
        </div>
        <span class="server-state" :class="`is-${serverConnection}`">
          <span class="server-state-dot" aria-hidden="true" />
          {{ serverStateLabel }}
        </span>
      </div>

      <form class="request-form" @submit.prevent="searchResource">
        <label class="field-group">
          <span class="field-label">目标词典</span>
          <NSelect
            v-model:value="selectedDict"
            size="small"
            :options="optionDicts"
            :disabled="isLoading || optionDicts.length === 0"
            placeholder="暂无可用词典"
          />
        </label>

        <label class="field-group resource-path-field">
          <span class="field-label">资源路径或 URL</span>
          <NInput
            v-model:value="resourceInputValue"
            size="small"
            clearable
            placeholder="例如 images/logo.png 或 /__mdict/style.css"
            :disabled="isLoading"
          />
          <span class="field-hint">支持相对路径、以 / 开头的服务路径或完整 HTTP URL。</span>
        </label>

        <NButton
          class="request-button"
          type="primary"
          size="small"
          attr-type="submit"
          :loading="isLoading"
          :disabled="!canSubmit"
        >
          {{ isLoading ? '请求中…' : '发送请求' }}
        </NButton>
      </form>

      <p v-if="validationMessage" class="validation-message" role="alert">{{ validationMessage }}</p>
    </section>

    <section class="request-preview" aria-labelledby="request-preview-title">
      <div class="section-label-row">
        <h3 id="request-preview-title">请求地址</h3>
        <NButton class="text-button" text type="primary" size="tiny" :disabled="!reqURL" @click="copyRequestURL">
          {{ copied ? '已复制' : '复制 URL' }}
        </NButton>
      </div>
      <code>{{ reqURL || '填写资源路径后将在这里生成最终请求地址' }}</code>
    </section>

    <section class="response-section" aria-labelledby="response-title">
      <div class="section-label-row response-heading">
        <h3 id="response-title">响应结果</h3>
        <span class="response-status" :class="`is-${requestState}`">
          {{ statusLabel }}
        </span>
      </div>

      <div class="response-metrics">
        <div class="response-metric">
          <span>HTTP 状态</span>
          <strong>{{ result.statusCode || '—' }}</strong>
        </div>
        <div class="response-metric">
          <span>内容类型</span>
          <strong>{{ result.contentType || '—' }}</strong>
        </div>
        <div class="response-metric">
          <span>内容长度</span>
          <strong>{{ formattedContentLength }}</strong>
        </div>
        <div class="response-metric">
          <span>请求耗时</span>
          <strong>{{ result.elapsedMs === null ? '—' : `${result.elapsedMs} ms` }}</strong>
        </div>
      </div>

      <div v-if="result.message" class="response-message" :class="`is-${requestState}`">
        {{ result.message }}
      </div>
      <div v-else class="response-empty">发送请求后，响应详情会显示在这里。</div>
    </section>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue';
import { NButton, NInput, NSelect } from 'naive-ui';
import { StaticDictServerURL } from '@/apis/apis';
import { useDictQueryStore } from '@/store/dict';
import { composeResourceRequestURL } from './resource-query';

type RequestState = 'idle' | 'loading' | 'success' | 'error';

interface DictionaryOption {
  label: string;
  value: string;
}

const dictQueryStore = useDictQueryStore();
const selectedDict = ref('');
const optionDicts = ref<DictionaryOption[]>([]);
const staticServerUrl = ref('');
const resourceInputValue = ref('');
const validationMessage = ref('');
const requestState = ref<RequestState>('idle');
const serverConnection = ref<'connecting' | 'ready' | 'error'>('connecting');
const copied = ref(false);
const result = reactive({
  statusCode: 0,
  contentType: '',
  contentLength: 0,
  elapsedMs: null as number | null,
  message: '',
});

const isLoading = computed(() => requestState.value === 'loading');
const canSubmit = computed(() => Boolean(staticServerUrl.value && selectedDict.value && resourceInputValue.value && !isLoading.value));
const reqURL = computed(() => composeResourceRequestURL(staticServerUrl.value, resourceInputValue.value, selectedDict.value));
const formattedContentLength = computed(() => result.contentLength > 0 ? `${result.contentLength.toLocaleString('zh-CN')} B` : '—');
const statusLabel = computed(() => ({
  idle: '尚未请求',
  loading: '请求中',
  success: '请求成功',
  error: '请求失败',
})[requestState.value]);
const serverStateLabel = computed(() => ({
  connecting: '正在连接服务',
  ready: '服务已连接',
  error: '服务连接失败',
})[serverConnection.value]);

function resetResult() {
  result.statusCode = 0;
  result.contentType = '';
  result.contentLength = 0;
  result.elapsedMs = null;
  result.message = '';
}

async function searchResource() {
  validationMessage.value = '';
  if (!resourceInputValue.value) {
    validationMessage.value = '请输入需要查询的资源路径或 URL。';
    return;
  }
  if (!selectedDict.value) {
    validationMessage.value = '请先选择一个词典。';
    return;
  }

  resetResult();
  requestState.value = 'loading';
  const startedAt = performance.now();

  try {
    const response = await fetch(reqURL.value);
    result.statusCode = response.status;
    result.contentLength = Number(response.headers.get('content-length') || 0);
    result.contentType = response.headers.get('content-type') || '';
    result.elapsedMs = Math.round(performance.now() - startedAt);
    result.message = response.ok ? '资源服务器已成功返回响应。' : `资源服务器返回 ${response.status} ${response.statusText || '错误响应'}。`;
    requestState.value = response.ok ? 'success' : 'error';
  } catch (cause) {
    result.elapsedMs = Math.round(performance.now() - startedAt);
    result.message = cause instanceof Error ? cause.message : '请求未能完成，请检查资源地址和词典服务状态。';
    requestState.value = 'error';
  }
}

async function copyRequestURL() {
  if (!reqURL.value) return;
  try {
    await navigator.clipboard.writeText(reqURL.value);
    copied.value = true;
    window.setTimeout(() => { copied.value = false; }, 1400);
  } catch {
    validationMessage.value = '无法复制请求地址，请手动选择并复制。';
  }
}

onMounted(async () => {
  try {
    staticServerUrl.value = await StaticDictServerURL();
    serverConnection.value = 'ready';
  } catch {
    serverConnection.value = 'error';
    validationMessage.value = '无法连接词典资源服务，请重新启动应用后再试。';
  }

  try {
    const dicts = await dictQueryStore.queryDictList();
    optionDicts.value = dicts.map((dict) => ({ label: dict.name, value: dict.id }));
    selectedDict.value = optionDicts.value[0]?.value || '';
  } catch {
    validationMessage.value ||= '无法加载词典列表。';
  }
});
</script>

<style lang="scss" scoped>
.resource-inspector {
  width: 100%;
  color: var(--c-gray-800);
  background: #fff;
  font-size: 13px;
}

.request-section,
.request-preview,
.response-section {
  padding: 16px;
  border-bottom: 1px solid var(--c-gray-300);
}

.response-section { border-bottom: 0; }

.section-heading,
.section-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.section-heading h2,
.section-label-row h3 {
  margin: 0;
  color: var(--c-gray-900);
  font-size: 13px;
  font-weight: 650;
}

.section-heading p {
  margin: 4px 0 0;
  color: var(--c-gray-600);
  font-size: 12px;
}

.server-state,
.response-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  color: var(--c-gray-600);
  font-size: 12px;
}

.server-state-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--c-gray-400);
}

.server-state.is-ready .server-state-dot { background: #3a9b58; }
.server-state.is-error { color: #a12b24; }
.server-state.is-error .server-state-dot { background: #c94b43; }

.request-form {
  display: grid;
  grid-template-columns: minmax(180px, 0.7fr) minmax(300px, 1.5fr) auto;
  align-items: end;
  gap: 12px;
  margin-top: 16px;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.field-label {
  color: var(--c-gray-800);
  font-size: 12px;
  font-weight: 600;
}

.field-hint {
  color: var(--c-gray-500);
  font-size: 11px;
  line-height: 1.35;
}

.request-button {
  min-width: 92px;
  margin-bottom: 23px;
  white-space: nowrap;
}

.request-form :deep(.n-button),
.request-form :deep(.n-input),
.request-form :deep(.n-base-selection) {
  background-image: none;
  box-shadow: none;
}

.request-form :deep(.n-button) {
  border-radius: 5px;
  font-size: 12px;
}

.request-form :deep(.n-input),
.request-form :deep(.n-base-selection) {
  --n-border-radius: 5px !important;
  font-size: 12px;
}

.validation-message {
  margin: 10px 0 0;
  color: #b42318;
  font-size: 12px;
}

.request-preview code {
  display: block;
  margin-top: 10px;
  padding: 9px 10px;
  overflow-x: auto;
  color: var(--c-gray-800);
  background: var(--c-gray-100);
  border: 1px solid var(--c-gray-300);
  border-radius: 4px;
  font-family: "Fira Code", ui-monospace, monospace;
  font-size: 11px;
  line-height: 1.5;
  white-space: nowrap;
}

.text-button {
  font-size: 12px;
}

.text-button :deep(.n-button__content) { font-size: 12px; }

.response-heading { margin-bottom: 12px; }

.response-status {
  padding: 2px 7px;
  background: var(--c-gray-100);
  border-radius: 999px;
}

.response-status.is-success { color: #257942; background: #edf8f0; }
.response-status.is-error { color: #a12b24; background: #fff0ef; }
.response-status.is-loading { color: var(--c-primary); }

.response-metrics {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  border: 1px solid var(--c-gray-300);
}

.response-metric {
  min-width: 0;
  padding: 10px 12px;
  border-right: 1px solid var(--c-gray-300);
}

.response-metric:last-child { border-right: 0; }
.response-metric span {
  display: block;
  margin-bottom: 5px;
  color: var(--c-gray-500);
  font-size: 11px;
}
.response-metric strong {
  display: block;
  overflow: hidden;
  color: var(--c-gray-900);
  font-family: "Fira Code", ui-monospace, monospace;
  font-size: 12px;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.response-message,
.response-empty {
  margin-top: 12px;
  padding: 10px 12px;
  color: var(--c-gray-600);
  background: var(--c-gray-50);
  border-left: 3px solid var(--c-gray-300);
  line-height: 1.5;
}

.response-message.is-success { border-left-color: #3a9b58; }
.response-message.is-error { color: #8f2923; border-left-color: #c94b43; }

@media (max-width: 840px) {
  .request-form { grid-template-columns: 1fr; align-items: stretch; }
  .request-button { width: fit-content; margin-bottom: 0; }
  .response-metrics { grid-template-columns: repeat(2, 1fr); }
  .response-metric:nth-child(2) { border-right: 0; }
  .response-metric:nth-child(-n + 2) { border-bottom: 1px solid var(--c-gray-300); }
}
</style>
