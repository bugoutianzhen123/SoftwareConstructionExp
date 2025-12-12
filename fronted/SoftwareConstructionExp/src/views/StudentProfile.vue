<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'
const me = ref(null)
const name = ref('')
const email = ref('')
const skillsText = ref('')
const msg = ref('')
const error = ref('')
const nameErr = ref('')
const emailErr = ref('')
async function load() {
  try { const u = await api.me(); me.value = u; name.value = u.name||''; email.value = u.email||''; skillsText.value = (u.skills||[]).join(', ') } catch(e){ error.value=e.message }
}
async function save() {
  error.value=''; msg.value=''
  nameErr.value = ''
  emailErr.value = ''
  if (!name.value) { nameErr.value='姓名不能为空' }
  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!email.value) { emailErr.value='邮箱不能为空' } else if (!emailRe.test(email.value)) { emailErr.value='邮箱格式不正确' }
  if (nameErr.value || emailErr.value) { error.value='请修正表单错误'; return }
  try {
    const skills = skillsText.value.split(',').map(s=>s.trim()).filter(Boolean)
    await api.updateMe({ name: name.value, email: email.value, skills })
    msg.value='已保存'
    window.dispatchEvent(new CustomEvent('toast', { detail: { msg: '保存成功', type: 'success' } }))
    await load()
  } catch(e){ error.value=e.message }
}
onMounted(load)
</script>

<template>
  <section class="card">
    <h3>个人信息</h3>
    <div class="toolbar">
      <button class="btn" @click="save">保存</button>
      <span style="color:#4a4">{{ msg }}</span>
    </div>
    <div class="form-grid">
      <div class="field"><label>姓名</label><input v-model="name" :class="{invalid: !!nameErr}" /><small class="err" v-if="nameErr">{{ nameErr }}</small></div>
      <div class="field"><label>邮箱</label><input v-model="email" type="email" :class="{invalid: !!emailErr}" /><small class="err" v-if="emailErr">{{ emailErr }}</small></div>
      <div class="field span-2"><label>技能(逗号分隔)</label><input v-model="skillsText" placeholder="如: python, nlp" /></div>
    </div>
    <div class="error" v-if="error">{{ error }}</div>
  </section>
</template>

<style scoped>
.form-grid { display:grid; grid-template-columns: 1fr 1fr; gap:12px }
.form-grid .span-2 { grid-column: span 2 }
@media (max-width: 900px) { .form-grid { grid-template-columns: 1fr } .form-grid .span-2 { grid-column: span 1 } }
.field label { display:block; margin-bottom:6px; font-weight:600 }
.error { color:#d33; margin-top:8px }
.err { color:#d33 }
input.invalid { border-color:#d33 }
</style>
