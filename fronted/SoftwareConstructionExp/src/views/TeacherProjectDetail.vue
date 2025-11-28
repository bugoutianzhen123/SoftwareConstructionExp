<template>
  <section class="card">
    <div class="toolbar">
      <button class="btn secondary" @click="router.push('/teacher/projects')">返回项目管理</button>
    </div>
    <div v-if="project">
      <h2>{{ project.title }}</h2>
      <p style="white-space:pre-wrap;">{{ project.description }}</p>
      <h3>项目要求</h3>
      <ul class="list">
        <li v-for="(r,i) in project.requirements" :key="i">{{ r }}</li>
      </ul>
      <h3>标签</h3>
      <div class="chips">
        <span v-for="(t,i) in project.tags" :key="i" class="chip">{{ t }}</span>
      </div>
    </div>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
const route = useRoute()
const router = useRouter()
const userId = Number(localStorage.getItem('user_id')||0)
const project = ref(null)
const error = ref('')
async function load(){
  try {
    const list = await api.listProjects(String(userId))
    project.value = list.find(p => String(p.id) === route.params.id)
    if (!project.value) error.value = '未找到该项目'
  } catch(e){ error.value = e.message }
}
onMounted(load)
</script>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
