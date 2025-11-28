<script setup>
import { ref } from 'vue'
import { api, uploadDocument } from '../api'
import FilterBar from '../components/FilterBar.vue'
const uploadAppId = ref('')
const uploadFile = ref(null)
const docs = ref([])
const error = ref('')
const q = ref('')
async function loadDocs() { try { if (!uploadAppId.value) return; docs.value = await api.listDocuments(Number(uploadAppId.value)) } catch(e){ error.value=e.message } }
function filtered(){ if(!q.value) return docs.value; const s=q.value.toLowerCase(); return docs.value.filter(d => (d.name||'').toLowerCase().includes(s)) }
async function onUpload() {
  try {
    if (!uploadAppId.value || !uploadFile.value?.files?.[0]) return
    await uploadDocument(Number(uploadAppId.value), uploadFile.value.files[0])
    await loadDocs()
  } catch(e) { error.value = e.message }
}
</script>

<template>
  <section class="card">
    <h3>成果文档上传</h3>
    <div class="toolbar">
      <label>申请ID<input v-model="uploadAppId" /></label>
      <input type="file" ref="uploadFile" />
      <FilterBar placeholder="搜索文档名" @update:query="v=>q.value=v" />
      <button class="btn" @click="onUpload">上传</button>
      <button class="btn secondary" @click="loadDocs">刷新列表</button>
    </div>
    <ul class="list">
      <li v-for="d in filtered()" :key="d.id">{{ d.name }} ({{ d.path }})</li>
    </ul>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
