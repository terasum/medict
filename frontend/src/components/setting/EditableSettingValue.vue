<template>
  <div class="editable-setting">
    <div v-if="!editing" class="editable-setting-display">
      <span class="editable-setting-value">{{ displayValue }}</span>
      <button class="btn btn-default" data-test="edit-setting" type="button" @click="startEditing">
        编辑
      </button>
    </div>
    <div v-else class="editable-setting-editor">
      <textarea
        v-if="multiline"
        v-model="draft"
        class="form-control editable-setting-textarea"
        :aria-label="ariaLabel"
        :rows="rows"
      />
      <select v-else-if="options.length" v-model="draft" class="form-control" :aria-label="ariaLabel">
        <option v-for="option in options" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>
      <input
        v-else
        v-model="draft"
        class="form-control"
        :aria-label="ariaLabel"
        :type="inputType"
        :min="min"
      />
      <div class="editable-setting-actions">
        <button class="btn btn-default" data-test="cancel-setting" type="button" :disabled="saving" @click="cancel">
          取消
        </button>
        <button class="btn btn-primary" data-test="save-setting" type="button" :disabled="saving" @click="save">
          {{ saving ? '保存中…' : '保存' }}
        </button>
      </div>
      <span v-if="error" class="editable-setting-error" role="alert">{{ error }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';

interface Option {
  label: string;
  value: string;
}

const props = withDefaults(defineProps<{
  modelValue: string | number | boolean;
  onSave: (value: string) => Promise<void> | void;
  multiline?: boolean;
  rows?: number;
  inputType?: string;
  min?: number;
  options?: Option[];
  ariaLabel?: string;
}>(), {
  multiline: false,
  rows: 5,
  inputType: 'text',
  min: undefined,
  options: () => [],
  ariaLabel: '设置值',
});

const emit = defineEmits<{ 'update:modelValue': [value: string] }>();
const editing = ref(false);
const saving = ref(false);
const error = ref('');
const draft = ref(String(props.modelValue));

const displayValue = computed(() => {
  const match = props.options.find((option) => option.value === String(props.modelValue));
  return match?.label ?? String(props.modelValue);
});

watch(() => props.modelValue, (value) => {
  if (!editing.value) draft.value = String(value);
});

function startEditing() {
  draft.value = String(props.modelValue);
  error.value = '';
  editing.value = true;
}

function cancel() {
  draft.value = String(props.modelValue);
  error.value = '';
  editing.value = false;
}

async function save() {
  saving.value = true;
  error.value = '';
  try {
    await props.onSave(draft.value);
    emit('update:modelValue', draft.value);
    editing.value = false;
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '保存失败，请重试';
  } finally {
    saving.value = false;
  }
}
</script>

<style lang="scss" scoped>
.editable-setting,
.editable-setting-editor {
  width: 100%;
}

.editable-setting-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.editable-setting-value {
  color: var(--c-gray-800);
  overflow-wrap: anywhere;
  white-space: pre-wrap;
}

.editable-setting-editor {
  display: grid;
  gap: 8px;
}

.editable-setting-textarea {
  min-height: 96px;
  resize: vertical;
  font-family: "Fira Code", ui-monospace, monospace;
  font-size: 12px;
}

.editable-setting-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.editable-setting-error {
  color: var(--c-danger, #bd3134);
  font-size: 12px;
}
</style>
