<script setup>
import { ref, onMounted } from 'vue'

const emit = defineEmits(['gameOver'])

const board = ref(Array(9).fill(null))
const isXNext = ref(true)
const winner = ref(null)

const calculateWinner = (squares) => {
  const lines = [
    [0, 1, 2], [3, 4, 5], [6, 7, 8],
    [0, 3, 6], [1, 4, 7], [2, 5, 8],
    [0, 4, 8], [2, 4, 6]
  ]
  for (let i = 0; i < lines.length; i++) {
    const [a, b, c] = lines[i]
    if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
      return squares[a]
    }
  }
  return null
}

const makeMove = (index) => {
  if (board.value[index] || winner.value) return
  board.value[index] = 'X' // Player is X
  
  if (checkWinOrDraw()) return

  // AI Move
  setTimeout(() => {
    const available = board.value.map((v, i) => v === null ? i : null).filter(v => v !== null)
    if (available.length > 0) {
      const move = available[Math.floor(Math.random() * available.length)]
      board.value[move] = 'O'
      checkWinOrDraw()
    }
  }, 500)
}

const checkWinOrDraw = () => {
  const w = calculateWinner(board.value)
  if (w) {
    winner.value = w
    emit('gameOver', w)
    return true
  }
  if (!board.value.includes(null)) {
    winner.value = 'Draw'
    emit('gameOver', 'Draw')
    return true
  }
  return false
}

const resetGame = () => {
  board.value = Array(9).fill(null)
  winner.value = null
  isXNext.value = true
}

defineExpose({ resetGame })
</script>

<template>
  <div class="glass-panel p-4 rounded-xl w-full max-w-sm mx-auto text-center">
    <h2 class="text-xl font-bold text-white mb-4">⭕ Tic Tac Toe</h2>
    
    <div class="grid grid-cols-3 gap-2 mb-4">
      <button 
        v-for="(cell, i) in board" 
        :key="i"
        @click="makeMove(i)"
        :disabled="cell || winner"
        class="aspect-square bg-slate-800 rounded-lg text-4xl font-bold flex items-center justify-center transition-colors"
        :class="{ 'text-blue-400': cell === 'X', 'text-red-400': cell === 'O', 'hover:bg-slate-700': !cell && !winner }"
      >
        {{ cell }}
      </button>
    </div>

    <div v-if="winner" class="mb-2">
      <p class="text-lg font-bold" :class="winner === 'X' ? 'text-blue-400' : 'text-red-400'">
        {{ winner === 'Draw' ? 'Seri!' : winner + ' Menang!' }}
      </p>
      <button @click="resetGame" class="mt-2 bg-blue-600 px-4 py-2 rounded text-white text-sm">
        🔄 Main Lagi
      </button>
    </div>
    <p v-else class="text-slate-400 text-sm">Giliran Kamu (X)</p>
  </div>
</template>
