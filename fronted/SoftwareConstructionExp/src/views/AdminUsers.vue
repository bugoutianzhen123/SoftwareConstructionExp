<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'
import FilterBar from '../components/FilterBar.vue'
const users = ref([])
const error = ref('')
const q = ref('')
const role = ref('')
async function load() { try { users.value = await api.listUsers('') } catch(e){ error.value=e.message } }
function filtered(){
  const list = role.value? users.value.filter(u => u.role===role.value) : users.value
  if(!q.value) return list
  const s = q.value.toLowerCase()
  return list.filter(u => (u.name||'').toLowerCase().includes(s) || (u.email||'').toLowerCase().includes(s) || (u.role||'').toLowerCase().includes(s))
}
onMounted(load)
</script>

<template>
  <section class="card">
    <h3>用户列表</h3>
    <div class="toolbar">
      <FilterBar placeholder="搜索姓名/邮箱/角色" @update:query="v=>q.value=v" :options="[{value:'student',label:'student'},{value:'teacher',label:'teacher'},{value:'admin',label:'admin'}]" @update:option="v=>role.value=v" />
      <button class="btn secondary" @click="load">刷新</button>
    </div>
    <ul class="list">
      <li v-for="u in filtered()" :key="u.id">{{ u.id }} - {{ u.name }} ({{ u.email }}) - {{ u.role }}</li>
    </ul>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
