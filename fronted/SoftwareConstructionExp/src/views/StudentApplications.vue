<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'
import FilterBar from '../components/FilterBar.vue'
const myApps = ref([])
const error = ref('')
const q = ref('')
async function load() { try { myApps.value = await api.listMyApplications('') } catch(e){ error.value=e.message } }
function filtered(){ if(!q.value) return myApps.value; const s=q.value.toLowerCase(); return myApps.value.filter(v => (v.project.title||'').toLowerCase().includes(s) || (v.application.status||'').toLowerCase().includes(s)) }
onMounted(load)
</script>

<template>
  <section class="card">
    <h3>我的申请</h3>
    <div class="toolbar">
      <FilterBar placeholder="搜索项目或状态" @update:query="v=>q.value=v" />
      <button class="btn secondary" @click="load">刷新</button>
    </div>
    <ul class="list">
      <li v-for="v in filtered()" :key="v.application.id">
        <b>{{ v.project.title }}</b> | 状态: {{ v.application.status }} | 匹配度: {{ (v.score*100).toFixed(0) }}%
      </li>
    </ul>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
