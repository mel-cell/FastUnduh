<script setup>
import { ref } from 'vue'

const targetNumber = ref(Math.floor(Math.random() * 100) + 1)
const guessInput = ref("")
const guessHistory = ref([])
const message = ref("Tebak 1-100!")
const won = ref(false)
const maxGuesses = 7

const makeGuess = () => {
  const g = parseInt(guessInput.value)
  if (isNaN(g) || g < 1 || g > 100) return

  guessHistory.value.push(g)
  if (g === targetNumber.value) {
    won.value = true
    message.value = `🎉 Benar! Angkanya ${g}`
  } else if (guessHistory.value.length >= maxGuesses) {
    message.value = `Gagal! Angkanya ${targetNumber.value}`
  } else if (g < targetNumber.value) {
    message.value = "Terlalu KECIL! 🔼"
  } else {
    message.value = "Terlalu BESAR! 🔽"
  }
  guessInput.value = ""
}

const resetGame = () => {
  targetNumber.value = Math.floor(Math.random() * 100) + 1
  guessInput.value = ""
  guessHistory.value = []
  message.value = "Tebak 1-100!"
  won.value = false
}

defineExpose({ resetGame })
</script>

<template>
  <div class="glass-panel p-6 rounded-xl w-full max-w-sm mx-auto text-center">
    <h2 class="text-xl font-bold text-white mb-4">🔢 Tebak Angka</h2>
    <p class="text-white mb-4 font-medium">{{ message }}</p>

    <div v-if="!won && guessHistory.length < maxGuesses" class="flex gap-2 mb-4">
      <input 
        v-model="guessInput" 
        type="number" 
        @keyup.enter="makeGuess"
        class="bg-slate-700 text-white px-3 py-2 rounded-lg w-full text-center" 
        placeholder="?" 
      />
      <button @click="makeGuess" class="bg-blue-600 px-4 py-2 rounded-lg text-white font-bold">Tekan</button>
    </div>

    <div v-if="guessHistory.length > 0" class="flex flex-wrap gap-2 justify-center mt-4">
       <span v-for="g in guessHistory" :key="g" class="bg-slate-800 px-2 py-1 rounded text-slate-300 text-xs">{{ g }}</span>
    </div>

    <button v-if="won || guessHistory.length >= maxGuesses" @click="resetGame" class="mt-4 bg-purple-600 px-6 py-2 rounded-lg text-white text-sm">
      Main Lagi
    </button>
  </div>
</template>
