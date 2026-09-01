<template>
  <div class="dashboard-page">
    <header class="topbar glass-card">
      <div>
        <h1>AI 智能单词本</h1>
        <p>输入单词，选择模型，生成释义和例句后手动保存</p>
      </div>
      <div class="user-area">
        <span>你好，{{ userStore.user?.username || '同学' }}</span>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </header>

    <main class="dashboard-grid">
      <section class="query-section glass-card">
        <h2>智能查询</h2>
        <el-form :model="queryForm" label-position="top">
          <el-form-item label="单词">
            <el-input v-model.trim="queryForm.word" size="large" placeholder="例如：innovation" clearable />
          </el-form-item>
          <el-form-item label="AI 提供商">
            <el-select v-model="queryForm.ai_provider" size="large" class="full-width">
              <el-option label="DeepSeek" value="deepseek" />
              <el-option label="通义千问" value="qianwen" />
            </el-select>
          </el-form-item>
          <el-button type="primary" size="large" :loading="queryLoading" @click="handleQuery">
            查询单词
          </el-button>
        </el-form>

        <WordCard v-if="currentWord" :word-data="currentWord" class="result-card" />

        <el-button
          v-if="currentWord && currentWord.source !== 'database'"
          class="save-btn"
          type="success"
          size="large"
          :loading="saveLoading"
          @click="handleSave"
        >
          保存到我的单词本
        </el-button>
      </section>

      <section class="list-section glass-card">
        <div class="section-head">
          <div>
            <h2>我的单词本</h2>
            <p>共 {{ pagination.total }} 条记录</p>
          </div>
          <el-button :loading="listLoading" @click="loadWords">刷新</el-button>
        </div>

        <el-table :data="words" v-loading="listLoading" class="word-table" empty-text="暂无保存的单词">
          <el-table-column prop="word" label="单词" width="130" />
          <el-table-column prop="definition" label="释义" min-width="220" show-overflow-tooltip />
          <el-table-column prop="ai_provider" label="来源" width="110">
            <template #default="{ row }">
              <el-tag>{{ row.ai_provider === 'qianwen' ? '通义千问' : 'DeepSeek' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-popconfirm title="确定删除这个单词吗？" confirm-button-text="删除" cancel-button-text="取消" @confirm="handleDelete(row.id)">
                <template #reference>
                  <el-button type="danger" link>删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          class="pager"
          background
          layout="prev, pager, next, sizes, total"
          :total="pagination.total"
          :page-size="pagination.page_size"
          :current-page="pagination.page"
          :page-sizes="[5, 10, 20, 50]"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </section>
    </main>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import WordCard from '../components/WordCard.vue'
import { deleteWord, getWords, queryWord, saveWord } from '../api'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const queryForm = reactive({
  word: '',
  ai_provider: 'deepseek'
})

const currentWord = ref(null)
const words = ref([])
const queryLoading = ref(false)
const saveLoading = ref(false)
const listLoading = ref(false)

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

async function handleQuery() {
  if (!queryForm.word) {
    ElMessage.warning('请输入要查询的单词')
    return
  }
  queryLoading.value = true
  try {
    const res = await queryWord(queryForm)
    currentWord.value = res
  } finally {
    queryLoading.value = false
  }
}

async function handleSave() {
  if (!currentWord.value) return
  saveLoading.value = true
  try {
    await saveWord({
      word: currentWord.value.word,
      definition: currentWord.value.definition,
      examples: typeof currentWord.value.examples === 'string' ? currentWord.value.examples : JSON.stringify(currentWord.value.examples),
      ai_provider: currentWord.value.ai_provider
    })
    ElMessage.success('保存成功')
    currentWord.value.source = 'database'
    await loadWords()
  } finally {
    saveLoading.value = false
  }
}

async function loadWords() {
  listLoading.value = true
  try {
    const res = await getWords({ page: pagination.page, page_size: pagination.page_size })
    words.value = res.data.list
    pagination.total = res.data.total
    pagination.page = res.data.page
    pagination.page_size = res.data.page_size
  } finally {
    listLoading.value = false
  }
}

async function handleDelete(id) {
  await deleteWord(id)
  ElMessage.success('删除成功')
  if (words.value.length === 1 && pagination.page > 1) {
    pagination.page -= 1
  }
  await loadWords()
}

function handlePageChange(page) {
  pagination.page = page
  loadWords()
}

function handleSizeChange(size) {
  pagination.page_size = size
  pagination.page = 1
  loadWords()
}

function handleLogout() {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  userStore.fetchUser().catch(() => {})
  loadWords()
})
</script>
