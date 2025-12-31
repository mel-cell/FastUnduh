<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({ isActive: Boolean })

const board = ref([])
const score = ref(0)
const gameOver = ref(false)

const initGame = () => {
  board.value = Array(4).fill(null).map(() => Array(4).fill(0))
  score.value = 0
  gameOver.value = false
  addTile(); addTile()
}

const addTile = () => {
  const empty = []
  board.value.forEach((r, i) => r.forEach((c, j) => { if (c === 0) empty.push({x:i, y:j}) }))
  if (empty.length) {
    const {x, y} = empty[Math.floor(Math.random() * empty.length)]
    board.value[x][y] = Math.random() < 0.9 ? 2 : 4
  }
}

const move = (dir) => {
  if (gameOver.value) return
  let moved = false
  const newBoard = JSON.parse(JSON.stringify(board.value))
  
  const merge = (row) => {
    let r = row.filter(v => v)
    let s = 0
    for(let i=0; i<r.length-1; i++) {
        if(r[i] === r[i+1]) { r[i]*=2; s+=r[i]; r.splice(i+1, 1) }
    }
    while(r.length < 4) r.push(0)
    return { r, s }
  }

  if (dir === 'left' || dir === 'right') {
    newBoard.forEach((row, i) => {
      const p = dir === 'left' ? row : [...row].reverse()
      const { r, s } = merge(p)
      score.value += s
      newBoard[i] = dir === 'left' ? r : r.reverse()
    })
  } else {
    for(let j=0; j<4; j++) {
      const col = newBoard.map(r => r[j])
      const p = dir === 'up' ? col : col.reverse()
      const { r, s } = merge(p)
      score.value += s
      const f = dir === 'up' ? r : r.reverse()
      newBoard.forEach((row, i) => { row[j] = f[i] })
    }
  }

  if (JSON.stringify(newBoard) !== JSON.stringify(board.value)) {
    board.value = newBoard; addTile(); checkOver()
  }
}

const checkOver = () => {
   // Simplified check
   if (!board.value.flat().includes(0)) gameOver.value = true // Naive check
}

const handleKey = (e) => {
  if (!props.isActive) return
  if (e.key === 'ArrowUp') { e.preventDefault(); move('up') }
  else if (e.key === 'ArrowDown') { e.preventDefault(); move('down') }
  else if (e.key === 'ArrowLeft') { e.preventDefault(); move('left') }
  else if (e.key === 'ArrowRight') { e.preventDefault(); move('right') }
}

const resetGame = () => initGame()

onMounted(() => {
  if(typeof window !== 'undefined') window.addEventListener('keydown', handleKey)
  initGame()
})

onUnmounted(() => {
  if(typeof window !== 'undefined') window.removeEventListener('keydown', handleKey)
})

defineExpose({ resetGame })
</script>

<template>
  <div class="glass-panel p-4 rounded-xl w-full max-w-sm mx-auto text-center">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-bold text-white">🎯 2048</h2>
      <span class="bg-slate-700 px-2 py-1 rounded text-white text-sm">Target: 2048</span>
    </div>
    
    <div class="bg-slate-800 p-2 rounded-lg grid grid-cols-4 gap-2 aspect-square">
       <div 
         v-for="(cell, i) in board.flat()" 
         :key="i"
         class="rounded flex items-center justify-center font-bold text-lg transition-all"
         :class="cell ? 'bg-orange-500 text-white' : 'bg-slate-700/50'"
       >
         {{ cell || '' }}
       </div>
    </div>

    <div class="mt-4 flex justify-between items-center">
       <span class="text-white font-bold">Skor: {{ score }}</span>
       <button v-if="gameOver" @click="resetGame" class="bg-blue-600 px-3 py-1 rounded text-white text-xs">Ulang</button>
    </div>
  </div>
</template>
