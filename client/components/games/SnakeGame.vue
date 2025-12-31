<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  isActive: { type: Boolean, default: false }
})

const emit = defineEmits(['gameOver', 'scoreUpdate'])

const gridSize = 20
const snake = ref([{ x: 10, y: 10 }])
const food = ref({ x: 15, y: 15 })
const direction = ref({ x: 0, y: 0 })
const snakeScore = ref(0)
const snakeGameOver = ref(false)
const snakeGameInterval = ref(null)

const generateFood = () => {
  let newFood
  do {
    newFood = {
      x: Math.floor(Math.random() * gridSize),
      y: Math.floor(Math.random() * gridSize)
    }
  } while (snake.value.some(segment => segment.x === newFood.x && segment.y === newFood.y))
  return newFood
}

const gameOver = () => {
  snakeGameOver.value = true
  if (snakeGameInterval.value) {
    clearInterval(snakeGameInterval.value)
    snakeGameInterval.value = null
  }
  emit('gameOver', snakeScore.value)
}

const moveSnake = () => {
  if (snakeGameOver.value) return

  const head = { ...snake.value[0] }
  head.x += direction.value.x
  head.y += direction.value.y

  // Check wall collision
  if (head.x < 0 || head.x >= gridSize || head.y < 0 || head.y >= gridSize) {
    gameOver()
    return
  }

  // Check self collision
  if (snake.value.some(segment => segment.x === head.x && segment.y === head.y)) {
    gameOver()
    return
  }

  snake.value.unshift(head)

  // Check food collision
  if (head.x === food.value.x && head.y === food.value.y) {
    snakeScore.value += 10
    emit('scoreUpdate', snakeScore.value)
    food.value = generateFood()
  } else {
    snake.value.pop()
  }
}

const startSnakeGame = () => {
  if (snakeGameInterval.value) return
  direction.value = { x: 1, y: 0 } // Start moving right
  snakeGameInterval.value = setInterval(moveSnake, 150)
}

const resetSnakeGame = () => {
  snake.value = [{ x: 10, y: 10 }]
  food.value = generateFood()
  direction.value = { x: 0, y: 0 }
  snakeScore.value = 0
  snakeGameOver.value = false
  emit('scoreUpdate', 0)
  if (snakeGameInterval.value) {
    clearInterval(snakeGameInterval.value)
    snakeGameInterval.value = null
  }
}

const handleKeydown = (e) => {
  if (!props.isActive || snakeGameOver.value) return
  
  // Start on first arrow press
  if (!snakeGameInterval.value && ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(e.key)) {
    startSnakeGame()
  }

  const key = e.key
  if (key === 'ArrowUp' && direction.value.y === 0) direction.value = { x: 0, y: -1 }
  else if (key === 'ArrowDown' && direction.value.y === 0) direction.value = { x: 0, y: 1 }
  else if (key === 'ArrowLeft' && direction.value.x === 0) direction.value = { x: -1, y: 0 }
  else if (key === 'ArrowRight' && direction.value.x === 0) direction.value = { x: 1, y: 0 }
}

const isSnakeCell = (x, y) => snake.value.some(segment => segment.x === x && segment.y === y)
const isFoodCell = (x, y) => food.value.x === x && food.value.y === y
const isSnakeHead = (x, y) => snake.value[0].x === x && snake.value[0].y === y

onMounted(() => {
  if (typeof window !== 'undefined') window.addEventListener('keydown', handleKeydown)
  resetSnakeGame()
})

onUnmounted(() => {
  if (typeof window !== 'undefined') window.removeEventListener('keydown', handleKeydown)
  if (snakeGameInterval.value) clearInterval(snakeGameInterval.value)
})

// Expose reset for parent
defineExpose({ resetSnakeGame })
</script>

<template>
  <div class="glass-panel p-4 rounded-xl w-full max-w-md mx-auto">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-bold text-white flex items-center gap-2">🐍 Snake</h2>
      <div class="glass-panel px-3 py-1 rounded text-sm text-white font-bold">
        Skor: {{ snakeScore }}
      </div>
    </div>

    <!-- Board -->
    <div class="aspect-square bg-slate-900/50 rounded-lg overflow-hidden grid" :style="{ gridTemplateColumns: `repeat(${gridSize}, 1fr)` }">
      <template v-for="y in gridSize" :key="y">
        <div
          v-for="x in gridSize"
          :key="`${x}-${y}`"
          :class="{
            'bg-green-500': isSnakeHead(x - 1, y - 1),
            'bg-green-600': isSnakeCell(x - 1, y - 1) && !isSnakeHead(x - 1, y - 1),
            'bg-red-500': isFoodCell(x - 1, y - 1),
            'bg-slate-800/30': !isSnakeCell(x - 1, y - 1) && !isFoodCell(x - 1, y - 1)
          }"
          class="aspect-square transition-colors duration-100"
        ></div>
      </template>
    </div>

    <!-- Controls / Status -->
    <div class="mt-4 text-center">
      <div v-if="!snakeGameInterval && !snakeGameOver">
        <p class="text-white text-sm">Tekan panah untuk mulai!</p>
      </div>

      <div v-if="snakeGameOver" class="space-y-2">
         <p class="text-red-400 font-bold">Game Over!</p>
         <button @click="resetSnakeGame" class="bg-blue-600 hover:bg-blue-500 text-white text-sm px-4 py-2 rounded-lg">
           🔄 Main Lagi
         </button>
      </div>
    </div>
  </div>
</template>
