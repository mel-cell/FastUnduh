<script setup>
import { ref, computed } from 'vue'

const props = defineProps({ isActive: Boolean })

const playerPosition = ref(1)
const diceValue = ref(null)
const isRolling = ref(false)
const canRoll = ref(true)
const gameMessage = ref("Klik dadu untuk mulai!")

const snakes = { 98: 28, 95: 24, 92: 51, 83: 19, 73: 1, 69: 33, 64: 36, 59: 17, 52: 11, 48: 9, 37: 3 }
const ladders = { 4: 56, 12: 50, 14: 55, 22: 58, 41: 79, 54: 88, 63: 80, 70: 90, 80: 99 }

// Board Cells
const boardCells = computed(() => {
  const cells = []
  for (let row = 9; row >= 0; row--) {
    const rowCells = []
    for (let col = 0; col < 10; col++) {
      const cellNumber = row % 2 === 1 ? row * 10 + (10 - col) : row * 10 + col + 1
      rowCells.push(cellNumber)
    }
    cells.push(rowCells)
  }
  return cells
})

const getCellClass = (cellNumber) => {
  if (cellNumber === playerPosition.value) return 'bg-yellow-500 scale-110 shadow-lg z-20'
  if (snakes[cellNumber]) return 'bg-red-500/30'
  if (ladders[cellNumber]) return 'bg-green-500/30'
  return 'bg-slate-700/50'
}

const rollDice = () => {
  if (!canRoll.value || isRolling.value) return
  isRolling.value = true
  canRoll.value = false
  gameMessage.value = "Melempar..."

  let count = 0
  const animationInterval = setInterval(() => {
    diceValue.value = Math.floor(Math.random() * 6) + 1
    count++
    if (count > 10) {
      clearInterval(animationInterval)
      const finalDice = Math.floor(Math.random() * 6) + 1
      diceValue.value = finalDice
      movePlayer(finalDice)
    }
  }, 100)
}

const movePlayer = (steps) => {
  let newPosition = playerPosition.value + steps
  if (newPosition > 100) newPosition = 100 - (newPosition - 100)

  let currentPos = playerPosition.value
  const moveInterval = setInterval(() => {
    if (currentPos < newPosition) {
      currentPos++
      playerPosition.value = currentPos
    } else if (currentPos > newPosition) { // Handle overshoot logic visually if needed, but logic above handles calc
        currentPos--
        playerPosition.value = currentPos
    } else {
      clearInterval(moveInterval)
      checkSpecialCell(newPosition)
    }
  }, 200)
}

const checkSpecialCell = (position) => {
  setTimeout(() => {
    if (snakes[position]) {
      gameMessage.value = `🐍 Kena ular! Turun ke ${snakes[position]}`
      setTimeout(() => { playerPosition.value = snakes[position]; finishTurn() }, 1000)
    } else if (ladders[position]) {
      gameMessage.value = `🪜 Dapat tangga! Naik ke ${ladders[position]}`
      setTimeout(() => { playerPosition.value = ladders[position]; finishTurn() }, 1000)
    } else {
      if (position === 100) gameMessage.value = "🎉 Menang!"
      else gameMessage.value = `Posisi: ${position}`
      finishTurn()
    }
  }, 500)
}

const finishTurn = () => {
  isRolling.value = false
  if (playerPosition.value < 100) {
    setTimeout(() => {
      canRoll.value = true
      if (gameMessage.value.includes("Posisi:")) gameMessage.value = "Lempar lagi!"
    }, 1000)
  }
}

const resetGame = () => {
  playerPosition.value = 1
  diceValue.value = null
  gameMessage.value = "Mulai!"
  canRoll.value = true
}

defineExpose({ resetGame })
</script>

<template>
  <div class="glass-panel p-4 rounded-xl w-full max-w-lg mx-auto">
     <h2 class="text-xl font-bold text-white mb-2 text-center">🎲 Ular Tangga</h2>
     
     <div class="grid grid-cols-10 gap-0.5 aspect-square bg-slate-800/50 p-1 rounded-lg">
       <div v-for="row in boardCells" :key="row[0]" class="contents">
          <div 
             v-for="cell in row" 
             :key="cell"
             :class="getCellClass(cell)"
             class="relative flex items-center justify-center text-[8px] sm:text-[10px] font-bold text-white rounded transition-all"
          >
             {{ cell }}
             <span v-if="snakes[cell]" class="absolute top-0 right-0 text-[6px]">🐍</span>
             <span v-if="ladders[cell]" class="absolute top-0 right-0 text-[6px]">🪜</span>
          </div>
       </div>
     </div>

     <div class="flex items-center justify-between mt-4 bg-slate-900/40 p-3 rounded-lg">
        <div class="text-white text-sm font-medium">{{ gameMessage }}</div>
        <button 
           @click="rollDice" 
           :disabled="!canRoll" 
           class="bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white font-bold px-4 py-2 rounded-lg text-xl w-16 h-12 flex items-center justify-center"
           :class="{ 'animate-bounce': canRoll && !isRolling }"
        >
           {{ diceValue || '🎲' }}
        </button>
     </div>
  </div>
</template>
