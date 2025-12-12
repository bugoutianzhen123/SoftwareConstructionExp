<script setup>
import { ref, onMounted } from "vue";
import { api } from "../api";
import FilterBar from "../components/FilterBar.vue";
const userId = Number(sessionStorage.getItem("user_id") || 0);
const apps = ref([]);
const appStatus = ref("submitted");
const error = ref("");
const projects = ref([]);
const projectId = ref("");
const busyId = ref(0);
const msgs = ref({});
const page = ref(1);
const pageSize = ref(10);
const loading = ref(false);
const openDocs = ref({});
const docsMap = ref({});
async function load() {
  try {
    loading.value = true;
    const [ps, av] = await Promise.all([
      api.listProjects(String(userId)),
      api.listApplications(
        `status=${appStatus.value || "submitted"}&fast=1&page=${
          page.value
        }&page_size=${pageSize.value}`
      ),
    ]);
    projects.value = ps;
    apps.value = av;
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}
function filtered() {
  const owned = apps.value.filter(
    (v) => v.project && Number(v.project.teacher_id) === userId
  );
  if (!projectId.value) return owned;
  const pid = Number(projectId.value);
  return owned.filter((v) => v.project && Number(v.project.id) === pid);
}
async function updateStatus(appId, status) {
  try {
    if (
      !confirm(status === "approved" ? "确认通过该申请？" : "确认拒绝该申请？")
    )
      return;
    busyId.value = appId;
    const res = await api.updateApplicationStatus(appId, status);
    msgs.value[appId] = "已更新";
    window.dispatchEvent(
      new CustomEvent("toast", { detail: { msg: "操作成功", type: "success" } })
    );
    await load();
  } catch (e) {
    msgs.value[appId] = e.message || "更新失败";
    error.value = e.message;
  } finally {
    busyId.value = 0;
  }
}
onMounted(load);

function nextPage() {
  if (apps.value.length === pageSize.value) {
    page.value++;
    load();
  }
}
function prevPage() {
  if (page.value > 1) {
    page.value--;
    load();
  }
}
function exportCSV() {
  const rows = filtered().map((v) => ({
    id: v.application.id,
    student: v.student.name,
    project: v.project.title,
    status: v.application.status,
    score: (v.score * 100).toFixed(0) + "%",
  }));
  const header = ["申请ID", "学生", "项目", "状态", "匹配度"];
  const csv = [
    header.join(","),
    ...rows.map((r) =>
      [r.id, r.student, r.project, r.status, r.score]
        .map((x) => `"${String(x).replace(/"/g, '""')}"`)
        .join(",")
    ),
  ].join("\n");
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "applications.csv";
  a.click();
  URL.revokeObjectURL(url);
}

async function toggleDocs(appId) {
  if (!openDocs.value[appId]) {
    openDocs.value[appId] = true;
    try {
      docsMap.value[appId] = await api.listDocuments(appId);
    } catch (e) {
      error.value = e.message;
      docsMap.value[appId] = [];
    }
  } else {
    openDocs.value[appId] = !openDocs.value[appId];
  }
}
async function downloadDoc(d) {
  try {
    const res = await fetch(`http://localhost:8080/api/documents/download?id=${d.id}`, { headers: { Authorization: `Bearer ${sessionStorage.getItem('token')}` } });
    if (!res.ok) { const data = await res.json().catch(()=>({})); throw new Error(data.error||'下载失败') }
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = d.name || 'document';
    a.click();
    URL.revokeObjectURL(url);
  } catch(e) { error.value = e.message }
}
</script>

<template>
  <section class="card">
    <h3>学生申请</h3>
    <div class="toolbar">
      <span style="margin-right: 8px">筛选</span>
      <FilterBar
        :options="[
          { value: 'submitted', label: 'submitted' },
          { value: 'approved', label: 'approved' },
          { value: 'rejected', label: 'rejected' },
        ]"
        @update:option="
          (v) => {
            appStatus = v;
            load();
          }
        "
      />
      <select v-model="projectId">
        <option value="">全部项目</option>
        <option v-for="p in projects" :key="p.id" :value="p.id">
          {{ p.title }}
        </option>
      </select>
      <select v-model.number="pageSize">
        <option :value="10">10条/页</option>
        <option :value="20">20条/页</option>
        <option :value="50">50条/页</option>
      </select>
      <button class="btn secondary" @click="load">刷新</button>
      <button class="btn secondary" @click="exportCSV">导出CSV</button>
    </div>
    <div v-if="loading">加载中…</div>
    <ul class="list">
      <li v-for="v in filtered()" :key="v.application.id">
        申请ID {{ v.application.id }} | 学生 {{ v.student.name }} | 项目
        {{ v.project.title }} | 匹配度 {{ (v.score * 100).toFixed(0) }}%
        <div class="toolbar" style="margin-top: 6px">
          <button
            class="btn"
            :disabled="busyId === v.application.id"
            @click="updateStatus(v.application.id, 'approved')"
          >
            通过
          </button>
          <button
            class="btn secondary"
            :disabled="busyId === v.application.id"
            @click="updateStatus(v.application.id, 'rejected')"
          >
            拒绝
          </button>
          <button class="btn secondary" @click="toggleDocs(v.application.id)">文档</button>
          <span
            style="margin-left: 8px; color: #4a4"
            v-if="msgs[v.application.id]"
            >{{ msgs[v.application.id] }}</span
          >
        </div>
        <div v-if="openDocs[v.application.id]" class="card" style="margin-top:6px">
          <ul class="list">
            <li v-if="!(docsMap[v.application.id]||[]).length">暂无文档</li>
            <li v-for="d in docsMap[v.application.id]||[]" :key="d.id">
              {{ d.name }}
              <button class="btn secondary" @click="downloadDoc(d)">下载</button>
            </li>
          </ul>
        </div>
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
.error {
  color: #d33;
  margin-top: 8px;
}
</style>
