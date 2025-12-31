<script setup>
import { ref, onMounted } from 'vue'

const cards = ref([])
const flipped = ref([])
const matched = ref([])
const moves = ref(0)
const won = ref(false)

const initGame = () => {
  const emojis = ['🍎', '🍌', '🍇', '🍊', '🍓', '🍉', '🥝', '🍒']
  const deck = [...emojis, ...emojis].sort(() => Math.random() - 0.5)
  cards.value = deck.map((e, i) => ({ id: i, emoji: e }))
  flipped.value = []
  matched.value = []
  moves.value = 0
  won.value = false
}

const flip = (index) => {
  if (flipped.value.length === 2 || matched.value.includes(index) || flipped.value.includes(index)) return
  
  flipped.value.push(index)
  if (flipped.value.length === 2) {
    moves.value++
    const [a, b] = flipped.value
    if (cards.value[a].emoji === cards.value[b].emoji) {
      matched.value.push(a, b)
      flipped.value = []
      if (matched.value.length === cards.value.length) won.value = true
    } else {
      setTimeout(() => { flipped.value = [] }, 800)
    }
  }
}

const isFlipped = (i) => flipped.value.includes(i) || matched.value.includes(i)

const resetGame = () => initGame()

onMounted(initGame)
defineExpose({ resetGame })
</script>

<template>
  <div class="glass-panel p-4 rounded-xl w-full max-w-sm mx-auto text-center">
    <div class="flex justify-between items-center mb-4">
       <h2 class="text-xl font-bold text-white">🃏 Memory</h2>
       <span class="text-slate-400 text-xs">Steps: {{ moves }}</span>
    </div>

    <div class="grid grid-cols-4 gap-2 mb-4">
       <button 
         v-for="(card, i) in cards" 
         :key="card.id"
         @click="flip(i)"
         class="aspect-square rounded-lg text-2xl flex items-center justify-center transition-all transform"
         :class="isFlipped(i) ? 'bg-blue-600 rotate-0' : 'bg-slate-700 rotate-180'"
       >
         <span v-if="isFlipped(i)">{{ card.emoji }}</span>
       </button>
    </div>

    <div v-if="won">
       <p class="text-green-400 font-bold mb-2">Menang!</p>
       <button @click="resetGame" class="bg-blue-600 px-4 py-2 rounded text-white text-sm">Ulang</button>
    </div>
  </div>
</template>
