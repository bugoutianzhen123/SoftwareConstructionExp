<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'
import FilterBar from '../components/FilterBar.vue'
const myApps = ref([])
const error = ref('')
const q = ref('')
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
async function load() { try { myApps.value = await api.listMyApplications('') } catch(e){ error.value=e.message } }
function filtered(){ if(!q.value) return myApps.value; const s=q.value.toLowerCase(); return myApps.value.filter(v => (v.project.title||'').toLowerCase().includes(s) || (v.application.status||'').toLowerCase().includes(s)) }
onMounted(load)

async function reload(){ try { loading.value=true; myApps.value = await api.listMyApplications(`page=${page.value}&page_size=${pageSize.value}`) } catch(e){ error.value=e.message } finally { loading.value=false } }
function nextPage(){ if (myApps.value.length === pageSize.value) { page.value++; reload() } }
function prevPage(){ if (page.value > 1) { page.value--; reload() } }
</script>

<template>
  <section class="card">
    <h3>我的申请</h3>
    <div class="toolbar">
      <FilterBar placeholder="搜索项目或状态" @update:query="v=>q.value=v" />
      <select v-model.number="pageSize">
        <option :value="10">10条/页</option>
        <option :value="20">20条/页</option>
        <option :value="50">50条/页</option>
      </select>
      <button class="btn secondary" @click="reload">刷新</button>
    </div>
    <div v-if="loading">加载中…</div>
    <ul class="list">
      <li v-for="v in filtered()" :key="v.application.id">
        <b>{{ v.project.title }}</b> | 状态: {{ v.application.status }} | 匹配度: {{ (v.score*100).toFixed(0) }}%
      </li>
    </ul>
    <div class="toolbar">
      <button class="btn secondary" @click="prevPage">上一页</button>
      <span>第 {{ page }} 页</span>
      <button class="btn secondary" @click="nextPage">下一页</button>
    </div>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
