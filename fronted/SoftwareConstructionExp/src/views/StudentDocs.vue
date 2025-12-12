<script setup>
import { ref, onMounted } from 'vue'
import { api, uploadDocument } from '../api'
import FilterBar from '../components/FilterBar.vue'
const approvedApps = ref([])
const selectedAppId = ref('')
const uploadFile = ref(null)
const docs = ref([])
const error = ref('')
const q = ref('')
async function loadApproved() {
  try { approvedApps.value = await api.listMyApplications('status=approved') }
  catch(e){ error.value = e.message }
}
async function loadDocs() {
  try { if (!selectedAppId.value) return; docs.value = await api.listDocuments(Number(selectedAppId.value)) }
  catch(e){ error.value=e.message }
}
function filtered(){ if(!q.value) return docs.value; const s=q.value.toLowerCase(); return docs.value.filter(d => (d.name||'').toLowerCase().includes(s)) }
async function onUpload() {
  try {
    if (!selectedAppId.value || !uploadFile.value?.files?.[0]) return
    await uploadDocument(Number(selectedAppId.value), uploadFile.value.files[0])
    await loadDocs()
  } catch(e) { error.value = e.message }
}
async function downloadDoc(d) {
  try {
    const res = await fetch(`http://localhost:8080/api/documents/download?id=${d.id}`, { headers: { Authorization: `Bearer ${sessionStorage.getItem('token')}` } })
    if (!res.ok) { const data = await res.json().catch(()=>({})); throw new Error(data.error||'下载失败') }
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = d.name || 'document'
    a.click()
    URL.revokeObjectURL(url)
  } catch(e) { error.value = e.message }
}
onMounted(loadApproved)
</script>

<template>
  <section class="card">
    <h3>成果文档上传</h3>
    <div class="toolbar">
      <label>已通过申请
        <select v-model="selectedAppId" @change="loadDocs">
          <option value="">请选择</option>
          <option v-for="v in approvedApps" :key="v.application.id" :value="v.application.id">
            {{ v.project?.title || ('申请 '+v.application.id) }}
          </option>
        </select>
      </label>
      <input type="file" ref="uploadFile" />
      <FilterBar placeholder="搜索文档名" @update:query="v=>q.value=v" />
      <button class="btn" :disabled="!selectedAppId" @click="onUpload">上传</button>
      <button class="btn secondary" :disabled="!selectedAppId" @click="loadDocs">刷新列表</button>
    </div>
    <ul class="list">
      <li v-for="d in filtered()" :key="d.id">{{ d.name }} <button class="btn secondary" @click="downloadDoc(d)">下载</button></li>
    </ul>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.error { color:#d33; margin-top:8px }
</style>
