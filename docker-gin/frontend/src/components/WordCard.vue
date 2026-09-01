<template>
  <el-card class="word-card" shadow="never">
    <template #header>
      <div class="word-card-header">
        <div>
          <h2>{{ wordData.word || '查询结果' }}</h2>
          <span class="source-tag">{{ sourceText }}</span>
        </div>
        <el-tag type="success">{{ providerText }}</el-tag>
      </div>
    </template>

    <div class="definition">
      <span class="label">释义</span>
      <p>{{ wordData.definition || '暂无释义' }}</p>
    </div>

    <div class="examples">
      <span class="label">例句</span>
      <div v-for="(item, index) in parsedExamples" :key="index" class="example-item">
        <strong>{{ index + 1 }}. {{ item.sentence }}</strong>
        <p>{{ item.translation }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  wordData: {
    type: Object,
    required: true
  }
})

const providerText = computed(() => {
  const map = { deepseek: 'DeepSeek', qianwen: '通义千问' }
  return map[props.wordData.ai_provider] || props.wordData.ai_provider || 'AI'
})

const sourceText = computed(() => (props.wordData.source === 'database' ? '来自我的单词本' : '来自 AI 查询'))

const parsedExamples = computed(() => {
  const raw = props.wordData.examples
  if (Array.isArray(raw)) return raw
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
})
</script>
