<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import FilterBar from '../components/FilterBar.vue'
const router = useRouter()
const userId = Number(sessionStorage.getItem('user_id')||0)
const projects = ref([])
const error = ref('')
const q = ref('')
const archivedFilter = ref('active')
async function loadProjects() { try { projects.value = await api.listProjects(String(userId)) } catch(e){ error.value=e.message } }
function goCreate(){ router.push('/teacher/projects/create') }
function filtered(){ if(!q.value) return projects.value; const s=q.value.toLowerCase(); return projects.value.filter(p => (p.title||'').toLowerCase().includes(s) || (p.description||'').toLowerCase().includes(s) || (p.tags||[]).some(t => (t||'').toLowerCase().includes(s))) }
onMounted(loadProjects)
async function removeProject(id){
  if (!confirm('确认归档该项目？')) return
  try { await api.archiveProject(id, true); window.dispatchEvent(new CustomEvent('toast', { detail: { msg: '项目已归档', type: 'success' } })); await loadProjects() } catch(e){ error.value=e.message }
}
function editProject(id){ router.push(`/teacher/projects/${id}/edit`) }
</script>

<template>
  <section class="card">
    <h3>项目管理</h3>
    <div class="toolbar">
      <button class="btn" @click="goCreate">发布项目</button>
      <FilterBar placeholder="搜索标题/描述/标签" @update:query="v=>q.value=v" />
      <select v-model="archivedFilter">
        <option value="active">在架</option>
        <option value="archived">归档</option>
        <option value="all">全部</option>
      </select>
      <button class="btn secondary" @click="loadProjects">刷新列表</button>
    </div>
    <h4 style="margin:12px 0 6px">项目列表</h4>
    <ul class="list">
      <li v-for="p in filtered().filter(x => archivedFilter==='all' ? true : archivedFilter==='archived' ? x.archived : !x.archived)" :key="p.id">
        <a href="#" @click.prevent="router.push(`/teacher/projects/${p.id}`)"><b>{{ p.title }}</b></a>
        <span v-if="p.archived" style="margin-left:8px;color:#888">[已归档]</span>
        <div class="toolbar" style="margin-top:6px">
          <button class="btn secondary" @click="editProject(p.id)" :disabled="p.archived">编辑</button>
          <button class="btn secondary" v-if="!p.archived" @click="removeProject(p.id)">归档</button>
          <button class="btn secondary" v-else @click="async ()=>{ try { await api.archiveProject(p.id, false); window.dispatchEvent(new CustomEvent('toast', { detail: { msg: '项目已恢复', type: 'success' } })); await loadProjects() } catch(e){ error.value=e.message } }">恢复</button>
        </div>
      </li>
    </ul>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
