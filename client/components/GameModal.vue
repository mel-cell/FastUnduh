<script setup>
import { ref, watch } from 'vue'
import SnakeGame from './games/SnakeGame.vue'
import TicTacToe from './games/TicTacToe.vue'
import SnakesLadders from './games/SnakesLadders.vue'
import GuessNumber from './games/GuessNumber.vue'
import MemoryCard from './games/MemoryCard.vue'
import Game2048 from './games/Game2048.vue'

const props = defineProps({
  isOpen: Boolean,
  downloadStatus: String, // 'processing', 'completed', 'failed'
  downloadResult: Object
})

const emit = defineEmits(['close', 'download'])

// List of available games
const games = [
  { id: 'snakes-ladders', component: SnakesLadders },
  { id: 'snake', component: SnakeGame },
  { id: 'tic-tac-toe', component: TicTacToe },
  { id: 'guess-number', component: GuessNumber },
  { id: 'memory-card', component: MemoryCard },
  { id: '2048', component: Game2048 },
]

const activeGameId = ref(games[0].id)
const showContinuationPrompt = ref(false)

// Select random game on mount or when reset
const switchGame = () => {
  const others = games.filter(g => g.id !== activeGameId.value)
  const next = others[Math.floor(Math.random() * others.length)]
  activeGameId.value = next.id
  showContinuationPrompt.value = false
}

// Watch download status
watch(() => props.downloadStatus, (newStatus) => {
  if (newStatus === 'completed') {
    showContinuationPrompt.value = true
  }
})

const handleDownloadAndClose = () => {
  emit('download')
  emit('close')
}

const handleDownloadAndContinue = () => {
  emit('download')
  showContinuationPrompt.value = false
}

const close = () => {
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 bg-slate-900/95 backdrop-blur-md z-50 flex flex-col items-center justify-center p-4 overflow-y-auto">
    
    <!-- Game Container -->
    <div class="relative w-full max-w-2xl flex-1 flex flex-col justify-center my-4">
      
      <!-- Top Actions -->
      <div class="flex justify-between items-center mb-6 px-2">
         <button 
           @click="switchGame" 
           class="bg-white/5 hover:bg-white/10 border border-white/10 px-4 py-2 rounded-full text-sm text-white transition-all flex items-center gap-2"
         >
           <span class="text-xl">🎲</span> Ganti Game
         </button>
         
         <button 
           @click="close" 
           class="bg-red-500/10 hover:bg-red-500/20 border border-red-500/20 p-2 rounded-full text-red-300 transition-all"
           title="Tutup"
         >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
         </button>
      </div>

      <!-- Active Game Component -->
      <div class="relative z-10 transition-all duration-300">
        <component 
          :is="games.find(g => g.id === activeGameId)?.component" 
          :isActive="isOpen"
        />
      </div>

    </div>

    <!-- Notification Bar (Bottom Fixed or Floating) -->
    
    <!-- PROCESSING -->
    <div v-if="downloadStatus === 'processing'" class="fixed bottom-8 left-1/2 -translate-x-1/2 bg-slate-800/90 backdrop-blur border border-blue-500/30 px-6 py-4 rounded-2xl flex items-center gap-4 shadow-2xl z-50 min-w-[300px]">
       <div class="bg-blue-500/20 p-2 rounded-full">
         <svg class="animate-spin h-6 w-6 text-blue-400" viewBox="0 0 24 24" fill="none">
           <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
           <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
         </svg>
       </div>
       <div class="flex-1">
         <p class="text-white font-medium">Sedang memproses...</p>
         <p class="text-slate-400 text-xs">Mohon tunggu sebentar</p>
       </div>
    </div>

    <!-- COMPLETED PROMPT -->
    <div v-else-if="showContinuationPrompt && downloadStatus === 'completed'" class="fixed bottom-8 left-1/2 -translate-x-1/2 bg-slate-800/95 backdrop-blur shadow-2xl border border-green-500/50 p-4 rounded-2xl flex flex-col md:flex-row items-center gap-6 animate-slide-up max-w-2xl w-[90%] z-50">
       <div class="flex items-center gap-4 flex-1">
         <div class="bg-green-500 p-3 rounded-xl shadow-lg shadow-green-500/20">
           <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" class="text-white"><polyline points="20 6 9 17 4 12"></polyline></svg>
         </div>
         <div>
           <p class="text-white font-bold text-lg">Download Selesai!</p>
           <p class="text-green-200/80 text-sm truncate max-w-[200px]">{{ downloadResult?.title || 'File siap diambil.' }}</p>
         </div>
       </div>

       <div class="flex gap-3 w-full md:w-auto">
          <button @click="handleDownloadAndClose" class="flex-1 whitespace-nowrap bg-white hover:bg-slate-100 text-slate-900 font-bold px-6 py-3 rounded-xl transition-all shadow-lg transform hover:-translate-y-0.5">
             📥 Ambil & Keluar
          </button>
          <button @click="handleDownloadAndContinue" class="flex-1 whitespace-nowrap bg-slate-700 hover:bg-slate-600 text-white font-semibold px-6 py-3 rounded-xl border border-slate-600 transition-all">
             🎮 Main Terus
          </button>
       </div>
    </div>

    <!-- FAILED -->
    <div v-else-if="downloadStatus === 'failed'" class="fixed bottom-8 left-1/2 -translate-x-1/2 bg-red-900/90 backdrop-blur border border-red-500/50 px-6 py-4 rounded-2xl flex flex-col items-center gap-2 shadow-2xl z-50">
       <p class="text-red-200 font-bold flex items-center gap-2">
         <span class="text-xl">⚠️</span> Gagal Mengunduh
       </p>
       <p class="text-red-300 text-sm">{{ downloadResult?.error }}</p>
       <button @click="close" class="mt-2 text-white bg-red-600/50 hover:bg-red-600 px-4 py-1 rounded-lg text-sm transition-colors">
         Tutup
       </button>
    </div>

  </div>
</template>

<style scoped>
.animate-slide-up {
  animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from { opacity: 0; transform: translate(-50%, 100px); }
  to { opacity: 1; transform: translate(-50%, 0); }
}
</style>
