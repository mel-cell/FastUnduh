<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  owner: { type: String, default: 'mel' },
  repo: { type: String, default: 'FastUnduh' } // Changed default to something likely
})

const collaborators = ref([])
const loading = ref(true)

// Fallback data if API fails or rate limited
const fallbackCollaborators = [
  { id: 1, login: 'mel', html_url: '#', avatar_url: 'https://avatars.githubusercontent.com/u/9919?s=200&v=4' }, // Generic avatar or user's
  { id: 2, login: 'contributor', html_url: '#', avatar_url: 'https://avatars.githubusercontent.com/u/583231?v=4' },
]

const fetchCollaborators = async () => {
  try {
    const response = await fetch(`https://api.github.com/repos/${props.owner}/${props.repo}/contributors`)
    if (!response.ok) throw new Error('API Error')
    collaborators.value = await response.json()
  } catch (err) {
    console.log("Using fallback contributors due to API limit")
    collaborators.value = fallbackCollaborators
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCollaborators()
})
</script>

<template>
  <div class="glass-panel p-6 rounded-2xl mt-12 w-full max-w-4xl border border-white/5 mx-auto">
    <div class="flex items-center gap-3 mb-6">
      <div class="p-2 bg-blue-500/20 rounded-lg">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-blue-400"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
      </div>
      <h3 class="text-lg font-bold text-white">Kontributor Proyek</h3>
    </div>
    
    <div v-if="loading" class="flex gap-4 animate-pulse">
       <div class="w-12 h-12 bg-slate-700/50 rounded-full" v-for="i in 3" :key="i"></div>
    </div>

    <div v-else class="flex flex-wrap gap-4">
      <a 
        v-for="user in collaborators" 
        :key="user.id" 
        :href="user.html_url" 
        target="_blank"
        class="group relative"
      >
        <img 
          :src="user.avatar_url" 
          :alt="user.login" 
          class="w-12 h-12 rounded-full border-2 border-slate-700 group-hover:border-blue-400 transition-all transform group-hover:scale-110 object-cover bg-slate-800"
        >
        <!-- Tooltip -->
        <div class="absolute -bottom-8 left-1/2 -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-slate-900 border border-slate-700 text-white text-xs px-2 py-1 rounded-md whitespace-nowrap pointer-events-none z-10 shadow-xl">
          {{ user.login }}
        </div>
      </a>
    </div>
  </div>
</template>

<style scoped>
.glass-panel {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
}
</style>
