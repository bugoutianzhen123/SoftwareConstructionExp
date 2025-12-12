<script setup>
import { ref, onMounted } from 'vue'
const visible = ref(false)
const msg = ref('')
const type = ref('info')
let timer = null
function show(m, t='info'){
  msg.value = m
  type.value = t
  visible.value = true
  clearTimeout(timer)
  timer = setTimeout(()=>{ visible.value=false }, 2000)
}
onMounted(()=>{
  window.addEventListener('toast', e => {
    const d = e.detail||{}
    show(d.msg||'', d.type||'info')
  })
})
</script>

<template>
  <div v-if="visible" class="toast" :class="type">{{ msg }}</div>
</template>

<style scoped>
.toast { position: fixed; right: 16px; top: 64px; background: #333; color: #fff; padding: 8px 12px; border-radius: 8px; z-index: 1000 }
.toast.success { background: #2e7d32 }
.toast.error { background: #c62828 }
.toast.info { background: #333 }
</style>
