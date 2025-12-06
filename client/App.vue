<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";

const url = ref("");
const isLoading = ref(false);
const jobTicket = ref(null);
const showGame = ref(false);
const downloadStatus = ref(""); // "processing", "completed", "failed"
const downloadResult = ref(null);
const selectedGame = ref(""); // "snakes-ladders" or"snake"

// Snakes and Ladders Game State
const playerPosition = ref(1);
const diceValue = ref(null);
const isRolling = ref(false);
const gameMessage = ref("Klik dadu untuk mulai!");
const canRoll = ref(true);

// Snake Game State
const gridSize = 20;
const snake = ref([{ x: 10, y: 10 }]);
const food = ref({ x: 15, y: 15 });
const direction = ref({ x: 0, y: 0 });
const snakeScore = ref(0);
const snakeGameOver = ref(false);
const snakeGameInterval = ref(null);

// Guess Number Game State
const targetNumber = ref(0);
const guessInput = ref("");
const guessHistory = ref([]);
const guessMessage = ref("");
const guessWon = ref(false);
const maxGuesses = 7;

// Memory Card Game State
const memoryCards = ref([]);
const flippedCards = ref([]);
const matchedCards = ref([]);
const memoryMoves = ref(0);
const memoryWon = ref(false);

// Tic Tac Toe Game State
const ticBoard = ref(Array(9).fill(null));
const ticIsXNext = ref(true);
const ticWinner = ref(null);

// 2048 Game State
const game2048Board = ref([]);
const game2048Score = ref(0);
const game2048Over = ref(false);

// Game switching state
const switchingGame = ref(false);
const switchCountdown = ref(3);

// Mouse tracking for 3D effects
const mouseX = ref(0);
const mouseY = ref(0);
const isMouseMoving = ref(false);

// Particle system
const particles = ref([]);



// Ular Tangga Board Configuration (100 cells)
const snakes = {
  98: 28,
  95: 24,
  92: 51,
  83: 19,
  73: 1,
  69: 33,
  64: 36,
  59: 17,
  52: 11,
  48: 9,
  37: 3
};

const ladders = {
  4: 56,
  12: 50,
  14: 55,
  22: 58,
  41: 79,
  54: 88,
  63: 80,
  70: 90,
  80: 99
};

const handleDownload = async () => {
  if (!url.value) return;
  isLoading.value = true;
  showGame.value = true;
  downloadStatus.value = "processing";
  downloadResult.value = null;

  // Random game selection
  const games = ['snakes-ladders', 'snake', 'guess-number', 'memory-card', 'tic-tac-toe', '2048'];
  selectedGame.value = games[Math.floor(Math.random() * games.length)];

  // Reset all games
  playerPosition.value = 1;
  diceValue.value = null;
  gameMessage.value = "Klik dadu untuk mulai!";
  canRoll.value = true;
  resetSnakeGame();
  initGuessNumber();
  initMemoryGame();
  initTicTacToe();
  init2048();

  try {
    // 1. Request Download (Endpoint: /api/queue)
    const response = await fetch("http://localhost:3000/api/queue", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url: url.value }),
    });

    if (!response.ok) throw new Error("Gagal memulai download");

    const data = await response.json();
    jobTicket.value = data.job_id;

    // 2. Start Polling
    pollStatus(data.job_id);
  } catch (error) {
    console.error(error);
    downloadStatus.value = "failed";
    downloadResult.value = { error: error.message };
    isLoading.value = false;
  }
};

const closeGame = () => {
  showGame.value = false;
  downloadStatus.value = "";
  downloadResult.value = null;
  isLoading.value = false;
  selectedGame.value = "";
  if (snakeGameInterval.value) {
    clearInterval(snakeGameInterval.value);
  }
};

const downloadFile = () => {
  if (downloadResult.value && downloadResult.value.filename) {
    window.location.href = `http://localhost:3000/api/download/${downloadResult.value.filename}`;
  }
};

const resetGame = () => {
  if (selectedGame.value === "snakes-ladders") {
    playerPosition.value = 1;
    diceValue.value = null;
    gameMessage.value = "Klik dadu untuk mulai!";
    canRoll.value = true;
  } else if (selectedGame.value === "snake") {
    resetSnakeGame();
  } else if (selectedGame.value === "guess-number") {
    initGuessNumber();
  } else if (selectedGame.value === "memory-card") {
    initMemoryGame();
  } else if (selectedGame.value === "tic-tac-toe") {
    initTicTacToe();
  } else if (selectedGame.value === "2048") {
    init2048();
  }
};



const switchToRandomGame = () => {
  switchingGame.value = true;
  switchCountdown.value = 3;

  const countdownInterval = setInterval(() => {
    switchCountdown.value--;
    if (switchCountdown.value <= 0) {
      clearInterval(countdownInterval);

      const games = ['snakes-ladders', 'snake', 'guess-number', 'memory-card', 'tic-tac-toe', '2048'];
      const availableGames = games.filter(game => game !== selectedGame.value);
      const newGame = availableGames[Math.floor(Math.random() * availableGames.length)];

      // Clean up current game
      if (snakeGameInterval.value) {
        clearInterval(snakeGameInterval.value);
      }

      // Switch to new game
      selectedGame.value = newGame;
      switchingGame.value = false;

      // Initialize the new game
      if (newGame === 'snakes-ladders') {
        playerPosition.value = 1;
        diceValue.value = null;
        gameMessage.value = "Klik dadu untuk mulai!";
        canRoll.value = true;
      } else if (newGame === 'snake') {
        resetSnakeGame();
      } else if (newGame === 'guess-number') {
        initGuessNumber();
      } else if (newGame === 'memory-card') {
        initMemoryGame();
      } else if (newGame === 'tic-tac-toe') {
        initTicTacToe();
      } else if (newGame === '2048') {
        init2048();
      }
    }
  }, 1000);
};

// Snake Game Functions
const resetSnakeGame = () => {
  snake.value = [{ x: 10, y: 10 }];
  food.value = generateFood();
  direction.value = { x: 0, y: 0 };
  snakeScore.value = 0;
  snakeGameOver.value = false;
  if (snakeGameInterval.value) {
    clearInterval(snakeGameInterval.value);
  }
};

const startSnakeGame = () => {
  if (snakeGameInterval.value) return;
  direction.value = { x: 1, y: 0 };
  snakeGameInterval.value = setInterval(moveSnake, 150);
};

const moveSnake = () => {
  if (snakeGameOver.value) return;

  const head = { ...snake.value[0] };
  head.x += direction.value.x;
  head.y += direction.value.y;

  // Check wall collision
  if (head.x < 0 || head.x >= gridSize || head.y < 0 || head.y >= gridSize) {
    gameOver();
    return;
  }

  // Check self collision
  if (snake.value.some(segment => segment.x === head.x && segment.y === head.y)) {
    gameOver();
    return;
  }

  snake.value.unshift(head);

  // Check food collision
  if (head.x === food.value.x && head.y === food.value.y) {
    snakeScore.value += 10;
    food.value = generateFood();
  } else {
    snake.value.pop();
  }
};

const generateFood = () => {
  let newFood;
  do {
    newFood = {
      x: Math.floor(Math.random() * gridSize),
      y: Math.floor(Math.random() * gridSize)
    };
  } while (snake.value.some(segment => segment.x === newFood.x && segment.y === newFood.y));
  return newFood;
};

const gameOver = () => {
  snakeGameOver.value = true;
  if (snakeGameInterval.value) {
    clearInterval(snakeGameInterval.value);
    snakeGameInterval.value = null;
  }
  // Auto switch to next game after 3 seconds
  setTimeout(() => {
    switchToRandomGame();
  }, 3000);
};

const handleKeydown = (e) => {
  if (!showGame.value || selectedGame.value !== "snake" || snakeGameOver.value) return;

  // Start game on first arrow key press
  if (!snakeGameInterval.value && ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(e.key)) {
    startSnakeGame();
  }

  const key = e.key;
  if (key === 'ArrowUp' && direction.value.y === 0) {
    direction.value = { x: 0, y: -1 };
  } else if (key === 'ArrowDown' && direction.value.y === 0) {
    direction.value = { x: 0, y: 1 };
  } else if (key === 'ArrowLeft' && direction.value.x === 0) {
    direction.value = { x: -1, y: 0 };
  } else if (key === 'ArrowRight' && direction.value.x === 0) {
    direction.value = { x: 1, y: 0 };
  }
};

// Add keyboard listener
if (typeof window !== 'undefined') {
  window.addEventListener('keydown', handleKeydown);
}

const isSnakeCell = (x, y) => {
  return snake.value.some(segment => segment.x === x && segment.y === y);
};

const isFoodCell = (x, y) => {
  return food.value.x === x && food.value.y === y;
};

const isSnakeHead = (x, y) => {
  return snake.value[0].x === x && snake.value[0].y === y;
};

// Guess Number Game Functions
const initGuessNumber = () => {
  targetNumber.value = Math.floor(Math.random() * 100) + 1;
  guessInput.value = "";
  guessHistory.value = [];
  guessMessage.value = "Tebak angka antara 1-100!";
  guessWon.value = false;
};

const makeGuess = () => {
  const guess = parseInt(guessInput.value);
  if (isNaN(guess) || guess < 1 || guess > 100) {
    guessMessage.value = "Masukkan angka antara 1-100!";
    return;
  }

  guessHistory.value.push(guess);

  if (guess === targetNumber.value) {
    guessWon.value = true;
    guessMessage.value = `🎉 Benar! Angkanya ${targetNumber.value}! (${guessHistory.value.length} tebakan)`;
    // Auto switch to next game after 3 seconds
    setTimeout(() => {
      switchToRandomGame();
    }, 3000);
  } else if (guessHistory.value.length >= maxGuesses) {
    guessMessage.value = `Game Over! Angkanya ${targetNumber.value}`;
    // Auto switch to next game after 3 seconds
    setTimeout(() => {
      switchToRandomGame();
    }, 3000);
  } else if (guess < targetNumber.value) {
    guessMessage.value = `Terlalu kecil! (${guessHistory.value.length}/${maxGuesses})`;
  } else {
    guessMessage.value = `Terlalu besar! (${guessHistory.value.length}/${maxGuesses})`;
  }

  guessInput.value = "";
};

// Memory Card Game Functions
const initMemoryGame = () => {
  const emojis = ['🍎', '🍌', '🍇', '🍊', '🍓', '🍉', '🥝', '🍒'];
  const cards = [...emojis, ...emojis].sort(() => Math.random() - 0.5);
  memoryCards.value = cards.map((emoji, index) => ({ id: index, emoji, flipped: false }));
  flippedCards.value = [];
  matchedCards.value = [];
  memoryMoves.value = 0;
  memoryWon.value = false;
};

const flipCard = (index) => {
  if (flippedCards.value.length === 2 || matchedCards.value.includes(index) || flippedCards.value.includes(index)) return;

  flippedCards.value.push(index);

  if (flippedCards.value.length === 2) {
    memoryMoves.value++;
    const [first, second] = flippedCards.value;

    if (memoryCards.value[first].emoji === memoryCards.value[second].emoji) {
      matchedCards.value.push(first, second);
      flippedCards.value = [];

      if (matchedCards.value.length === memoryCards.value.length) {
        memoryWon.value = true;
        // Auto switch to next game after 3 seconds
        setTimeout(() => {
          switchToRandomGame();
        }, 3000);
      }
    } else {
      setTimeout(() => {
        flippedCards.value = [];
      }, 800);
    }
  }
};

const isCardFlipped = (index) => {
  return flippedCards.value.includes(index) || matchedCards.value.includes(index);
};

// Tic Tac Toe Game Functions
const initTicTacToe = () => {
  ticBoard.value = Array(9).fill(null);
  ticIsXNext.value = true;
  ticWinner.value = null;
};

const makeMove = (index) => {
  if (ticBoard.value[index] || ticWinner.value) return;

  ticBoard.value[index] = 'X';
  ticWinner.value = calculateWinner(ticBoard.value);

  if (ticWinner.value) {
    // Auto switch to next game after 3 seconds
    setTimeout(() => {
      switchToRandomGame();
    }, 3000);
  } else if (ticBoard.value.some(cell => cell === null)) {
    setTimeout(() => {
      aiMove();
    }, 500);
  }
};

const aiMove = () => {
  const availableMoves = ticBoard.value.map((cell, idx) => cell === null ? idx : null).filter(idx => idx !== null);
  if (availableMoves.length === 0) return;

  const randomMove = availableMoves[Math.floor(Math.random() * availableMoves.length)];
  ticBoard.value[randomMove] = 'O';
  ticWinner.value = calculateWinner(ticBoard.value);

  if (ticWinner.value) {
    // Auto switch to next game after 3 seconds
    setTimeout(() => {
      switchToRandomGame();
    }, 3000);
  }
};

const calculateWinner = (board) => {
  const lines = [
    [0, 1, 2], [3, 4, 5], [6, 7, 8],
    [0, 3, 6], [1, 4, 7], [2, 5, 8],
    [0, 4, 8], [2, 4, 6]
  ];

  for (const [a, b, c] of lines) {
    if (board[a] && board[a] === board[b] && board[a] === board[c]) {
      return board[a];
    }
  }

  if (board.every(cell => cell !== null)) return 'Draw';
  return null;
};

// 2048 Game Functions
const init2048 = () => {
  game2048Board.value = Array(4).fill(null).map(() => Array(4).fill(0));
  game2048Score.value = 0;
  game2048Over.value = false;
  addNewTile();
  addNewTile();
};

const addNewTile = () => {
  const emptyCells = [];
  for (let i = 0; i < 4; i++) {
    for (let j = 0; j < 4; j++) {
      if (game2048Board.value[i][j] === 0) {
        emptyCells.push({ i, j });
      }
    }
  }

  if (emptyCells.length > 0) {
    const { i, j } = emptyCells[Math.floor(Math.random() * emptyCells.length)];
    game2048Board.value[i][j] = Math.random() < 0.9 ? 2 : 4;
  }
};

const move2048 = (direction) => {
  if (game2048Over.value) return;

  let moved = false;
  const newBoard = JSON.parse(JSON.stringify(game2048Board.value));

  if (direction === 'left' || direction === 'right') {
    for (let i = 0; i < 4; i++) {
      const row = direction === 'left' ? newBoard[i] : newBoard[i].reverse();
      const { newRow, score } = mergeRow(row);
      game2048Score.value += score;
      newBoard[i] = direction === 'left' ? newRow : newRow.reverse();
      if (JSON.stringify(newBoard[i]) !== JSON.stringify(game2048Board.value[i])) moved = true;
    }
  } else {
    for (let j = 0; j < 4; j++) {
      const col = newBoard.map(row => row[j]);
      const column = direction === 'up' ? col : col.reverse();
      const { newRow, score } = mergeRow(column);
      game2048Score.value += score;
      const finalCol = direction === 'up' ? newRow : newRow.reverse();
      for (let i = 0; i < 4; i++) {
        newBoard[i][j] = finalCol[i];
      }
      if (JSON.stringify(col) !== JSON.stringify(finalCol)) moved = true;
    }
  }

  if (moved) {
    game2048Board.value = newBoard;
    addNewTile();
    checkGameOver();
  }
};

const mergeRow = (row) => {
  let newRow = row.filter(val => val !== 0);
  let score = 0;

  for (let i = 0; i < newRow.length - 1; i++) {
    if (newRow[i] === newRow[i + 1]) {
      newRow[i] *= 2;
      score += newRow[i];
      newRow.splice(i + 1, 1);
    }
  }

  while (newRow.length < 4) {
    newRow.push(0);
  }

  return { newRow, score };
};

const checkGameOver = () => {
  // Check if player reached 2048
  for (let i = 0; i < 4; i++) {
    for (let j = 0; j < 4; j++) {
      if (game2048Board.value[i][j] === 2048) {
        game2048Over.value = true;
        // Auto switch to next game after 3 seconds
        setTimeout(() => {
          switchToRandomGame();
        }, 3000);
        return;
      }
    }
  }

  // Check if no more moves available
  for (let i = 0; i < 4; i++) {
    for (let j = 0; j < 4; j++) {
      if (game2048Board.value[i][j] === 0) return;
      if (j < 3 && game2048Board.value[i][j] === game2048Board.value[i][j + 1]) return;
      if (i < 3 && game2048Board.value[i][j] === game2048Board.value[i + 1][j]) return;
    }
  }
  game2048Over.value = true;
  // Auto switch to next game after 3 seconds
  setTimeout(() => {
    switchToRandomGame();
  }, 3000);
};

const handle2048Keydown = (e) => {
  if (!showGame.value || selectedGame.value !== '2048') return;

  if (e.key === 'ArrowUp') {
    e.preventDefault();
    move2048('up');
  } else if (e.key === 'ArrowDown') {
    e.preventDefault();
    move2048('down');
  } else if (e.key === 'ArrowLeft') {
    e.preventDefault();
    move2048('left');
  } else if (e.key === 'ArrowRight') {
    e.preventDefault();
    move2048('right');
  }
};

// Add keyboard listener for 2048
if (typeof window !== 'undefined') {
  window.addEventListener('keydown', handle2048Keydown);
}

// Mouse tracking for 3D effects
const handleMouseMove = (e) => {
  mouseX.value = (e.clientX / window.innerWidth - 0.5) * 2;
  mouseY.value = (e.clientY / window.innerHeight - 0.5) * 2;
  isMouseMoving.value = true;

  // Create particle on mouse move (throttled)
  if (Math.random() > 0.95) {
    createParticle(e.clientX, e.clientY);
  }
};

const createParticle = (x, y) => {
  const particle = {
    id: Date.now() + Math.random(),
    x,
    y,
    color: ['#3b82f6', '#8b5cf6', '#ec4899'][Math.floor(Math.random() * 3)]
  };
  particles.value.push(particle);

  setTimeout(() => {
    particles.value = particles.value.filter(p => p.id !== particle.id);
  }, 3000);
};

// Universal tilt effect - same for all elements
const getTiltStyle = (intensity = 1) => {
  if (!isMouseMoving.value) return {};

  const rotateX = mouseY.value * 8 * intensity;
  const rotateY = mouseX.value * -8 * intensity;
  const translateZ = 20 * intensity;

  return {
    transform: `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) translateZ(${translateZ}px)`,
    transition: 'transform 0.1s ease-out',
  };
};

// Parallax effect - same movement for all
const getParallaxStyle = (depth = 1) => {
  if (!isMouseMoving.value) return {};

  const moveX = mouseX.value * 15 * depth;
  const moveY = mouseY.value * 15 * depth;

  return {
    transform: `translate(${moveX}px, ${moveY}px)`,
    transition: 'transform 0.1s ease-out',
  };
};



// Lifecycle hooks
onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('mousemove', handleMouseMove);
  }
});

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('mousemove', handleMouseMove);
  }
});

const rollDice = () => {
  if (!canRoll.value || isRolling.value) return;

  isRolling.value = true;
  canRoll.value = false;
  gameMessage.value = "Melempar dadu...";

  // Animasi dadu
  let count = 0;
  const animationInterval = setInterval(() => {
    diceValue.value = Math.floor(Math.random() * 6) + 1;
    count++;

    if (count > 10) {
      clearInterval(animationInterval);
      const finalDice = Math.floor(Math.random() * 6) + 1;
      diceValue.value = finalDice;
      movePlayer(finalDice);
    }
  }, 100);
};

const movePlayer = (steps) => {
  let newPosition = playerPosition.value + steps;

  // Check if exceeds 100
  if (newPosition > 100) {
    newPosition = 100 - (newPosition - 100);
  }

  // Animate movement
  let currentPos = playerPosition.value;
  const moveInterval = setInterval(() => {
    if (currentPos < newPosition) {
      currentPos++;
      playerPosition.value = currentPos;
    } else {
      clearInterval(moveInterval);
      checkSpecialCell(newPosition);
    }
  }, 200);
};

const checkSpecialCell = (position) => {
  setTimeout(() => {
    // Check for snake
    if (snakes[position]) {
      gameMessage.value = `Kena ular! Turun ke ${snakes[position]}`;
      setTimeout(() => {
        playerPosition.value = snakes[position];
        finishTurn();
      }, 1000);
    }
    // Check for ladder
    else if (ladders[position]) {
      gameMessage.value = `Dapat tangga! Naik ke ${ladders[position]}`;
      setTimeout(() => {
        playerPosition.value = ladders[position];
        finishTurn();
      }, 1000);
    }
    // Normal cell
    else {
      if (position === 100) {
        gameMessage.value = "🎉 Selamat! Kamu menang!";
        // Auto switch to next game after 3 seconds
        setTimeout(() => {
          switchToRandomGame();
        }, 3000);
      } else {
        gameMessage.value = `Posisi: ${position}`;
      }
      finishTurn();
    }
  }, 500);
};

const finishTurn = () => {
  isRolling.value = false;
  if (playerPosition.value < 100) {
    setTimeout(() => {
      canRoll.value = true;
      if (gameMessage.value.includes("Posisi:")) {
        gameMessage.value = "Lempar lagi!";
      }
    }, 1000);
  }
};

const pollStatus = (id) => {
  const interval = setInterval(async () => {
    try {
      // Endpoint: /api/status/:id
      const res = await fetch(`http://localhost:3000/api/status/${id}`);

      if (!res.ok) {
        return;
      }

      const data = await res.json();
      const status = data.status;

      if (status === "completed") {
        clearInterval(interval);
        isLoading.value = false;
        downloadStatus.value = "completed";
        downloadResult.value = {
          filename: data.filename,
          title: data.title || "Video"
        };
      } else if (status === "failed") {
        clearInterval(interval);
        isLoading.value = false;
        downloadStatus.value = "failed";
        downloadResult.value = { error: "Download gagal. Cek URL atau coba lagi." };
      } else {
        console.log("Status:", status);
      }
    } catch (e) {
      clearInterval(interval);
      isLoading.value = false;
      downloadStatus.value = "failed";
      downloadResult.value = { error: "Terjadi kesalahan saat polling status" };
      console.error("Polling error", e);
    }
  }, 2000);
};

// Helper to get cell color
const getCellClass = (cellNumber) => {
  if (cellNumber === playerPosition.value) return 'bg-yellow-500 scale-110 shadow-lg';
  if (snakes[cellNumber]) return 'bg-red-500/30';
  if (ladders[cellNumber]) return 'bg-green-500/30';
  return 'bg-slate-700/50';
};

// Generate board cells (10x10 grid, snake pattern)
const boardCells = computed(() => {
  const cells = [];
  for (let row = 9; row >= 0; row--) {
    const rowCells = [];
    for (let col = 0; col < 10; col++) {
      const cellNumber = row % 2 === 1
        ? row * 10 + (10 - col)
        : row * 10 + col + 1;
      rowCells.push(cellNumber);
    }
    cells.push(rowCells);
  }
  return cells;
});
</script>

``<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4 relative overflow-hidden">
    <!-- Particles -->
    <div class="fixed inset-0 pointer-events-none z-50">
      <div
        v-for="particle in particles"
        :key="particle.id"
        class="particle"
        :style="{
          left: particle.x + 'px',
          top: particle.y + 'px',
          background: particle.color
        }"
      ></div>
    </div>

    <!-- Background Elements with Parallax -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden -z-10">
      <div
        class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-600/20 rounded-full blur-[100px] parallax-layer"
        :style="getParallaxStyle(0.5)"
      ></div>
      <div
        class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-purple-600/20 rounded-full blur-[100px] parallax-layer"
        :style="getParallaxStyle(0.3)"
      ></div>
      <div
        class="absolute top-[50%] left-[50%] w-[30%] h-[30%] bg-pink-600/10 rounded-full blur-[80px] parallax-layer"
        :style="getParallaxStyle(0.7)"
      ></div>
    </div>

    <!-- Game Modal -->
    <div
      v-if="showGame"
      class="fixed inset-0 bg-slate-900/40 backdrop-blur-md z-50 flex items-center justify-center p-4"
    >
      <div class="glass-panel p-4 md:p-6 rounded-2xl max-w-3xl w-full max-h-[95vh] overflow-y-auto relative">
        <!-- Switching Game Notification -->
        <div v-if="switchingGame" class="absolute inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center rounded-2xl">
          <div class="text-center">
            <div class="text-6xl mb-4 animate-bounce">🎮</div>
            <h3 class="text-2xl font-bold text-white mb-2">Game Berikutnya</h3>
            <p class="text-slate-300 text-lg">dalam {{ switchCountdown }} detik...</p>
          </div>
        </div>

        <!-- Close Button -->
        <button
          @click="closeGame"
          class="absolute top-4 right-4 text-slate-400 hover:text-white transition-colors p-2 hover:bg-slate-700/50 rounded-lg z-10"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>



        <!-- Snakes and Ladders Game -->
        <div v-if="selectedGame === 'snakes-ladders'">
          <div class="text-center mb-4">
            <button
              @click="switchToRandomGame"
              class="glass-panel px-4 py-2 rounded-lg hover:bg-slate-600/50 transition-colors mb-3 flex items-center gap-2 mx-auto text-sm text-white"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              Ganti Game
            </button>
            <h2 class="text-2xl font-bold text-white mb-3">🎲 Ular Tangga</h2>
            <div class="flex items-center justify-center gap-3 flex-wrap">
              <div class="glass-panel px-3 py-1.5 rounded-lg">
                <span class="text-slate-400 text-xs">Posisi:</span>
                <span class="text-white font-bold text-lg ml-1">{{ playerPosition }}</span>
              </div>
              <div class="glass-panel px-3 py-1.5 rounded-lg">
                <span class="text-slate-400 text-xs">Target:</span>
                <span class="text-white font-bold text-lg ml-1">100</span>
              </div>
            </div>
          </div>

          <!-- Game Board -->
          <div class="mb-4 bg-slate-800/30 p-2 md:p-3 rounded-xl">
            <div class="grid grid-cols-10 gap-0.5 aspect-square max-w-md mx-auto">
              <template v-for="row in boardCells" :key="row">
                <div
                  v-for="cell in row"
                  :key="cell"
                  :class="getCellClass(cell)"
                  class="flex items-center justify-center text-[10px] font-bold text-white rounded transition-all duration-300 relative"
                >
                  <span class="relative z-10">{{ cell }}</span>
                  <!-- Snake indicator -->
                  <span v-if="snakes[cell]" class="absolute top-0 right-0 text-[6px]">🐍</span>
                  <!-- Ladder indicator -->
                  <span v-if="ladders[cell]" class="absolute top-0 right-0 text-[6px]">🪜</span>
                </div>
              </template>
            </div>
          </div>

          <!-- Dice and Controls -->
          <div class="flex flex-col items-center gap-3">
            <div class="text-center">
              <p class="text-white font-semibold text-base mb-2">{{ gameMessage }}</p>
              <button
                @click="rollDice"
                :disabled="!canRoll || isRolling"
                class="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white font-bold px-6 py-3 rounded-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed text-5xl"
                :class="{ 'animate-bounce': canRoll && !isRolling }"
              >
                {{ diceValue || '🎲' }}
              </button>
            </div>

            <!-- Legend -->
            <div class="flex gap-3 text-xs text-slate-400 flex-wrap justify-center">
              <div class="flex items-center gap-1.5">
                <div class="w-3 h-3 bg-yellow-500 rounded"></div>
                <span>Posisi</span>
              </div>
              <div class="flex items-center gap-1.5">
                <div class="w-3 h-3 bg-red-500/30 rounded"></div>
                <span>🐍 Ular</span>
              </div>
              <div class="flex items-center gap-1.5">
                <div class="w-3 h-3 bg-green-500/30 rounded"></div>
                <span>🪜 Tangga</span>
              </div>
            </div>

            <!-- Download Status -->
            <div v-if="downloadStatus === 'processing'" class="glass-panel px-6 py-3 rounded-lg flex items-center gap-3">
            <svg
              class="animate-spin h-5 w-5 text-blue-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            <span class="text-white">Download sedang diproses...</span>
          </div>

            <!-- Download Completed -->
            <div v-if="downloadStatus === 'completed'" class="space-y-3 w-full max-w-md">
            <div class="glass-panel px-6 py-3 rounded-lg flex items-center gap-3 bg-green-500/20">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="text-green-400"
              >
                <polyline points="20 6 9 17 4 12"></polyline>
 </svg>
              <div class="flex-1 text-left">
                <p class="text-white font-semibold">Download Selesai!</p>
                <p class="text-slate-400 text-sm">{{ downloadResult?.title }}</p>
              </div>
            </div>
            <div class="flex gap-2">
              <button
                @click="downloadFile"
                class="flex-1 bg-green-600 hover:bg-green-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
              >
                📥 Download File
              </button>
              <button
                @click="resetGame"
                class="flex-1 bg-blue-600 hover:bg-blue-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
              >
                🔄 Main Lagi
              </button>
            </div>
          </div>

            <!-- Download Failed -->
            <div v-if="downloadStatus === 'failed'" class="space-y-3 w-full max-w-md">
            <div class="glass-panel px-6 py-3 rounded-lg flex items-center gap-3 bg-red-500/20">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="text-red-400"
              >
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
              <div class="flex-1 text-left">
                <p class="text-white font-semibold">Download Gagal</p>
                <p class="text-slate-400 text-sm">{{ downloadResult?.error }}</p>
              </div>
            </div>
            <button
              @click="resetGame"
              class="w-full bg-blue-600 hover:bg-blue-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
            >
              🔄 Main Lagi
            </button>
          </div>
          </div>
        </div>

        <!-- Snake Game -->
        <div v-if="selectedGame === 'snake'">
          <div class="text-center mb-4">
            <button
              @click="switchToRandomGame"
              class="glass-panel px-4 py-2 rounded-lg hover:bg-slate-600/50 transition-colors mb-3 flex items-center gap-2 mx-auto text-sm text-white"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              Ganti Game
            </button>
            <h2 class="text-2xl font-bold text-white mb-3">🐍 Snake Game</h2>
            <div class="flex items-center justify-center gap-3 flex-wrap">
              <div class="glass-panel px-3 py-1.5 rounded-lg">
                <span class="text-slate-400 text-xs">Skor:</span>
                <span class="text-white font-bold text-lg ml-1">{{ snakeScore }}</span>
              </div>
            </div>
          </div>

          <!-- Snake Game Board -->
          <div class="mb-4 bg-slate-800/30 p-2 md:p-3 rounded-xl">
            <div class="grid gap-0 aspect-square max-w-md mx-auto bg-slate-900/50 rounded-lg overflow-hidden" :style="{ gridTemplateColumns: `repeat(${gridSize}, 1fr)` }">
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
          </div>

          <!-- Snake Game Controls -->
          <div class="flex flex-col items-center gap-3">
            <div v-if="!snakeGameInterval && !snakeGameOver" class="text-center">
              <p class="text-white font-semibold text-base mb-2">Tekan tombol panah untuk mulai!</p>
              <p class="text-slate-400 text-xs">Gunakan ⬆️ ⬇️ ⬅️ ➡️ untuk mengontrol ular</p>
            </div>

            <div v-if="snakeGameOver" class="text-center space-y-3">
              <p class="text-red-400 font-bold text-xl">Game Over!</p>
              <p class="text-white text-base">Skor Akhir: {{ snakeScore }}</p>
              <button
                @click="resetSnakeGame"
                class="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
              >
                🔄 Main Lagi
              </button>
            </div>

            <!-- Download Status -->
            <div v-if="downloadStatus === 'processing'" class="glass-panel px-6 py-3 rounded-lg flex items-center gap-3">
              <svg
                class="animate-spin h-5 w-5 text-blue-400"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              <span class="text-white">Download sedang diproses...</span>
            </div>

            <!-- Download Completed -->
            <div v-if="downloadStatus === 'completed'" class="space-y-3 w-full max-w-md">
              <div class="glass-panel px-6 py-3 rounded-lg flex items-center gap-3 bg-green-500/20">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="text-green-400"
                >
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
                <div class="flex-1 text-left">
                  <p class="text-white font-semibold">Download Selesai!</p>
                  <p class="text-slate-400 text-sm">{{ downloadResult?.title }}</p>
                </div>
              </div>
              <div class="flex gap-2">
                <button
                  @click="downloadFile"
                  class="flex-1 bg-green-600 hover:bg-green-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
                >
                  📥 Download File
                </button>
                <button
                  @click="resetGame"
                  class="flex-1 bg-blue-600 hover:bg-blue-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
                >
                  🔄 Main Lagi
                </button>
              </div>
            </div>

            <!-- Download Failed -->
            <div v-if="downloadStatus === 'failed'" class="space-y-3 w-full max-w-md">
              <div class="glass-panel px-6 py-3 rounded-lg flex items-center gap-3 bg-red-500/20">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="text-red-400"
                >
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="15" y1="9" x2="9" y2="15"></line>
                  <line x1="9" y1="9" x2="15" y2="15"></line>
                </svg>
                <div class="flex-1 text-left">
                  <p class="text-white font-semibold">Download Gagal</p>
                  <p class="text-slate-400 text-sm">{{ downloadResult?.error }}</p>
                </div>
              </div>
              <button
                @click="resetGame"
                class="w-full bg-blue-600 hover:bg-blue-500 text-white font-semibold px-6 py-3 rounded-xl transition-all duration-200"
              >
                🔄 Main Lagi
              </button>
            </div>
          </div>
        </div>

        <!-- Guess Number Game -->
        <div v-if="selectedGame === 'guess-number'">
          <div class="text-center mb-6">
            <button @click="switchToRandomGame" class="glass-panel px-4 py-2 rounded-lg hover:bg-slate-600/50 transition-colors mb-4 flex items-center gap-2 mx-auto text-sm text-white">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              Ganti Game
            </button>
            <h2 class="text-3xl font-bold text-white mb-4">🔢 Tebak Angka</h2>
            <p class="text-lg text-white mb-4">{{ guessMessage }}</p>

            <div v-if="!guessWon && guessHistory.length < maxGuesses" class="max-w-md mx-auto space-y-4">
              <input v-model="guessInput" @keyup.enter="makeGuess" type="number" placeholder="Masukkan tebakan..." class="w-full bg-slate-700/50 text-white px-6 py-4 rounded-xl text-center text-2xl" />
              <button @click="makeGuess" class="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white font-bold px-8 py-4 rounded-xl">Tebak!</button>
            </div>

            <div v-if="guessHistory.length > 0" class="mt-6">
              <p class="text-slate-400 mb-2">Riwayat Tebakan:</p>
              <div class="flex flex-wrap gap-2 justify-center">
                <span v-for="(guess, idx) in guessHistory" :key="idx" class="glass-panel px-4 py-2 rounded-lg text-white">{{ guess }}</span>
              </div>
            </div>

            <button v-if="guessWon || guessHistory.length >= maxGuesses" @click="initGuessNumber" class="mt-6 bg-blue-600 hover:bg-blue-500 text-white font-semibold px-8 py-4 rounded-xl">🔄 Main Lagi</button>
          </div>
        </div>

        <!-- Memory Card Game -->
        <div v-if="selectedGame === 'memory-card'">
          <div class="text-center mb-6">
            <button @click="switchToRandomGame" class="glass-panel px-4 py-2 rounded-lg hover:bg-slate-600/50 transition-colors mb-4 flex items-center gap-2 mx-auto text-sm text-white">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              Ganti Game
            </button>
            <h2 class="text-3xl font-bold text-white mb-4">🃏 Memory Card</h2>
            <p class="text-slate-400">Langkah: {{ memoryMoves }}</p>
          </div>

          <div class="grid grid-cols-4 gap-3 max-w-md mx-auto mb-6">
            <button v-for="(card, index) in memoryCards" :key="card.id" @click="flipCard(index)" :class="{ 'bg-blue-600': isCardFlipped(index), 'bg-slate-700': !isCardFlipped(index) }" class="aspect-square rounded-xl text-4xl flex items-center justify-center transition-all hover:scale-105">
              {{ isCardFlipped(index) ? card.emoji : '?' }}
            </button>
          </div>

          <div v-if="memoryWon" class="text-center">
            <p class="text-green-400 font-bold text-2xl mb-4">🎉 Selamat! Kamu menang dalam {{ memoryMoves }} langkah!</p>
            <button @click="initMemoryGame" class="bg-blue-600 hover:bg-blue-500 text-white font-semibold px-8 py-4 rounded-xl">🔄 Main Lagi</button>
          </div>
        </div>

        <!-- Tic Tac Toe Game -->
        <div v-if="selectedGame === 'tic-tac-toe'">
          <div class="text-center mb-6">
            <button @click="switchToRandomGame" class="glass-panel px-4 py-2 rounded-lg hover:bg-slate-600/50 transition-colors mb-4 flex items-center gap-2 mx-auto text-sm text-white">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              Ganti Game
            </button>
            <h2 class="text-3xl font-bold text-white mb-4">⭕ Tic Tac Toe</h2>
            <p class="text-slate-400 mb-4">{{ ticWinner ? (ticWinner === 'Draw' ? 'Seri!' : `${ticWinner} Menang!`) : 'Giliran Kamu (X)' }}</p>
          </div>

          <div class="grid grid-cols-3 gap-3 max-w-sm mx-auto mb-6">
            <button v-for="(cell, index) in ticBoard" :key="index" @click="makeMove(index)" :disabled="cell || ticWinner" class="aspect-square bg-slate-700 rounded-xl text-5xl font-bold flex items-center justify-center hover:bg-slate-600 disabled:cursor-not-allowed transition-all" :class="{ 'text-blue-400': cell === 'X', 'text-red-400': cell === 'O' }">
              {{ cell }}
            </button>
          </div>

          <div v-if="ticWinner" class="text-center">
            <button @click="initTicTacToe" class="bg-blue-600 hover:bg-blue-500 text-white font-semibold px-8 py-4 rounded-xl">🔄 Main Lagi</button>
          </div>
        </div>

        <!-- 2048 Game -->
        <div v-if="selectedGame === '2048'">
          <div class="text-center mb-6">
            <button @click="switchToRandomGame" class="glass-panel px-4 py-2 rounded-lg hover:bg-slate-600/50 transition-colors mb-4 flex items-center gap-2 mx-auto text-sm text-white">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              Ganti Game
            </button>
            <h2 class="text-3xl font-bold text-white mb-4">🎯 2048</h2>
            <div class="flex gap-4 justify-center mb-4">
              <div class="glass-panel px-4 py-2 rounded-lg">
                <span class="text-slate-400 text-sm">Skor:</span>
                <span class="text-white font-bold text-xl ml-2">{{ game2048Score }}</span>
              </div>
            </div>
            <p class="text-slate-400 text-sm">Gunakan tombol panah untuk bermain</p>
          </div>

          <div class="max-w-md mx-auto mb-6 bg-slate-800/50 p-4 rounded-xl">
            <div class="grid grid-cols-4 gap-3">
              <template v-for="(row, i) in game2048Board" :key="i">
                <div v-for="(cell, j) in row" :key="`${i}-${j}`" class="aspect-square rounded-lg flex items-center justify-center font-bold text-2xl" :class="{ 'bg-slate-700': cell === 0, 'bg-yellow-600': cell === 2, 'bg-yellow-500': cell === 4, 'bg-orange-500': cell === 8, 'bg-orange-600': cell === 16, 'bg-red-500': cell === 32, 'bg-red-600': cell === 64, 'bg-red-700': cell === 128, 'bg-purple-600': cell === 256, 'bg-purple-700': cell === 512, 'bg-blue-600': cell === 1024, 'bg-blue-700': cell === 2048, 'text-white': cell > 0 }">
                  {{ cell || '' }}
                </div>
              </template>
            </div>
          </div>

          <div v-if="game2048Over" class="text-center">
            <p class="text-red-400 font-bold text-2xl mb-4">Game Over!</p>
            <button @click="init2048" class="bg-blue-600 hover:bg-blue-500 text-white font-semibold px-8 py-4 rounded-xl">🔄 Main Lagi</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Container with Universal Mouse Tracking -->
    <main v-if="!showGame" class="w-full max-w-3xl text-center z-10 perspective-container" :style="getTiltStyle(0.5)">
      <!-- Header -->
      <div class="mb-12 space-y-4 animate-slide-up">
        <div class="relative inline-block">
          <h1
            class="text-5xl md:text-7xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 tracking-tight animate-gradient"
            style="transform-style: preserve-3d; transform: translateZ(30px);"
          >
            FastUnduh
          </h1>
          <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg blur opacity-20 animate-pulse-glow" style="transform: translateZ(-10px);"></div>
        </div>
        <p class="text-lg md:text-xl text-slate-300 max-w-xl mx-auto leading-relaxed" style="transform: translateZ(20px);">
          Download video dari berbagai platform dengan kecepatan kilat. Tanpa iklan, tanpa ribet.
        </p>
        <div class="flex items-center justify-center gap-2 text-sm text-slate-400" style="transform: translateZ(15px);">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-green-400">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>100% Gratis & Aman</span>
        </div>
      </div>

      <!-- Input Section -->
      <div class="glass-panel p-2 rounded-2xl input-glow transition-all duration-300 animate-scale-in relative overflow-hidden group glow-on-hover card-3d">
        <div class="absolute inset-0 bg-gradient-to-r from-blue-600/10 to-purple-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
        <div class="flex flex-col md:flex-row gap-2 relative z-10" style="transform: translateZ(10px);">
          <input
            v-model="url"
            type="text"
            placeholder="Tempel link video di sini (YouTube, Instagram, Twitter...)"
            class="flex-1 bg-transparent border-none outline-none text-white placeholder-slate-400 px-6 py-4 text-lg w-full"
            @keyup.enter="handleDownload"
          />
          <button
            @click="handleDownload"
            :disabled="isLoading"
            class="button-3d bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white font-semibold px-8 py-4 rounded-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 min-w-[160px] shadow-lg hover:shadow-blue-500/50"
          >
            <span v-if="!isLoading" class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
              </svg>
              Download
            </span>
            <span v-else class="flex items-center gap-2">
              <svg
                class="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Memproses...
            </span>
          </button>
        </div>
      </div>

      <!-- Features -->
      <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-6 text-slate-400 text-sm">
        <div class="glass-panel p-6 rounded-xl flex flex-col items-center gap-3 transition-all duration-300 cursor-pointer group card-3d glow-on-hover hover:scale-105">
          <div class="p-4 bg-gradient-to-br from-blue-500/20 to-blue-600/20 rounded-full group-hover:from-blue-500/30 group-hover:to-blue-600/30 transition-all duration-300 group-hover:scale-110" style="transform: translateZ(20px);">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="text-blue-400 group-hover:text-blue-300 transition-colors group-hover:rotate-12"
            >
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
            </svg>
          </div>
          <span class="font-semibold text-white text-base" style="transform: translateZ(15px);">Super Cepat</span>
          <span class="text-xs text-center text-slate-400" style="transform: translateZ(10px);">Download dengan kecepatan maksimal</span>
        </div>
        <div class="glass-panel p-6 rounded-xl flex flex-col items-center gap-3 transition-all duration-300 cursor-pointer group card-3d glow-on-hover hover:scale-105">
          <div class="p-4 bg-gradient-to-br from-purple-500/20 to-purple-600/20 rounded-full group-hover:from-purple-500/30 group-hover:to-purple-600/30 transition-all duration-300 group-hover:scale-110" style="transform: translateZ(20px);">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="text-purple-400 group-hover:text-purple-300 transition-colors"
            >
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="7 10 12 15 17 10" />
              <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
          </div>
          <span class="font-semibold text-white text-base">Tanpa Batas</span>
          <span class="text-xs text-center text-slate-400">Download sebanyak yang kamu mau</span>
        </div>
        <div class="glass-panel p-6 rounded-xl flex flex-col items-center gap-3 transition-all duration-300 cursor-pointer group card-3d glow-on-hover hover:scale-105">
          <div class="p-4 bg-gradient-to-br from-green-500/20 to-green-600/20 rounded-full group-hover:from-green-500/30 group-hover:to-green-600/30 transition-all duration-300 group-hover:scale-110" style="transform: translateZ(20px);">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="text-green-400 group-hover:text-green-300 transition-colors"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
            </svg>
          </div>
          <span class="font-semibold text-white text-base">Aman & Privat</span>
          <span class="text-xs text-center text-slate-400">Data kamu tidak disimpan</span>
        </div>
      </div>

      <!-- Footer with About Section -->
      <footer class="mt-16 pb-8">
        <div class="glass-panel p-8 rounded-2xl card-3d glow-on-hover relative overflow-hidden">
          <!-- Animated background gradient -->
          <div class="absolute inset-0 bg-gradient-to-br from-blue-600/5 via-purple-600/5 to-pink-600/5 opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>

          <div class="relative z-10">
            <!-- About Section -->
            <div class="mb-12" style="transform: translateZ(20px);">
              <h3 class="text-3xl font-bold text-center mb-4 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 animate-gradient">
                Tentang FastUnduh
              </h3>
              <p class="text-slate-300 text-center mb-8 max-w-2xl mx-auto leading-relaxed">
                FastUnduh adalah platform download video yang cepat, mudah, dan gratis. Kami mendukung berbagai platform seperti YouTube, Instagram, Twitter, dan banyak lagi. Sambil menunggu download, nikmati 6 game seru yang kami sediakan!
              </p>

              <!-- Stats Cards with 3D -->
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
                <div class="glass-panel p-4 rounded-xl text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(15px);">
                  <div class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-blue-600 mb-2 group-hover:scale-125 transition-transform">6</div>
                  <div class="text-xs text-slate-400 font-medium">Game Seru</div>
                </div>
                <div class="glass-panel p-4 rounded-xl text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(15px);">
                  <div class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-purple-600 mb-2 group-hover:scale-125 transition-transform">∞</div>
                  <div class="text-xs text-slate-400 font-medium">Download Gratis</div>
                </div>
                <div class="glass-panel p-4 rounded-xl text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(15px);">
                  <div class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-green-400 to-green-600 mb-2 group-hover:scale-125 transition-transform">100%</div>
                  <div class="text-xs text-slate-400 font-medium">Aman</div>
                </div>
                <div class="glass-panel p-4 rounded-xl text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(15px);">
                  <div class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 to-yellow-600 mb-2 group-hover:scale-125 transition-transform">0</div>
                  <div class="text-xs text-slate-400 font-medium">Iklan</div>
                </div>
              </div>

              <!-- Platforms -->
              <div class="border-t border-slate-700/50 pt-6 mb-6">
                <h4 class="text-lg font-semibold text-white mb-4 text-center">Platform yang Didukung</h4>
                <div class="flex flex-wrap justify-center gap-3">
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">YouTube</span>
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">Instagram</span>
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">Twitter</span>
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">Facebook</span>
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">TikTok</span>
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">Vimeo</span>
                  <span class="glass-panel px-4 py-2 rounded-lg text-sm text-slate-300 hover:scale-110 hover:text-white transition-all duration-300 cursor-pointer" style="transform: translateZ(10px);">Dan lainnya...</span>
                </div>
              </div>

              <!-- Games Grid -->
              <div class="border-t border-slate-700/50 pt-6 mb-8">
                <h4 class="text-lg font-semibold text-white mb-4 text-center">Game yang Tersedia</h4>
                <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
                  <div class="glass-panel p-4 rounded-lg text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(10px);">
                    <div class="text-3xl mb-2 group-hover:rotate-12 transition-transform">🎲</div>
                    <div class="text-sm text-slate-300 group-hover:text-white transition-colors">Ular Tangga</div>
                  </div>
                  <div class="glass-panel p-4 rounded-lg text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(10px);">
                    <div class="text-3xl mb-2 group-hover:rotate-12 transition-transform">🐍</div>
                    <div class="text-sm text-slate-300 group-hover:text-white transition-colors">Snake Game</div>
                  </div>
                  <div class="glass-panel p-4 rounded-lg text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(10px);">
                    <div class="text-3xl mb-2 group-hover:rotate-12 transition-transform">🔢</div>
                    <div class="text-sm text-slate-300 group-hover:text-white transition-colors">Tebak Angka</div>
                  </div>
                  <div class="glass-panel p-4 rounded-lg text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(10px);">
                    <div class="text-3xl mb-2 group-hover:rotate-12 transition-transform">🃏</div>
                    <div class="text-sm text-slate-300 group-hover:text-white transition-colors">Memory Card</div>
                  </div>
                  <div class="glass-panel p-4 rounded-lg text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(10px);">
                    <div class="text-3xl mb-2 group-hover:rotate-12 transition-transform">⭕</div>
                    <div class="text-sm text-slate-300 group-hover:text-white transition-colors">Tic Tac Toe</div>
                  </div>
          <div class="glass-panel p-4 rounded-lg text-center card-3d hover:scale-110 transition-all duration-300 cursor-pointer group" style="transform: translateZ(10px);">
                    <div class="text-3xl mb-2 group-hover:rotate-12 transition-transform">🎯</div>
                    <div class="text-sm text-slate-300 group-hover:text-white transition-colors">2048</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Footer Info -->
            <div class="border-t border-slate-700/50 pt-8" style="transform: translateZ(15px);">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-6">
                <!-- Brand -->
                <div class="text-center md:text-left group">
                  <h3 class="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-400 mb-3 group-hover:scale-110 transition-transform inline-block">
                    FastUnduh
                  </h3>
                  <p class="text-slate-400 text-sm leading-relaxed">
                    Download video dengan cepat sambil bermain game seru!
                  </p>
                </div>

                <!-- Quick Links -->
                <div class="text-center">
                  <h4 class="text-white font-semibold mb-4 text-base">Quick Links</h4>
                  <div class="flex flex-col gap-3">
                    <a href="#" class="text-slate-400 hover:text-white hover:translate-x-2 transition-all duration-300 text-sm inline-block">→ Beranda</a>
                    <a href="#" class="text-slate-400 hover:text-white hover:translate-x-2 transition-all duration-300 text-sm inline-block">→ Cara Penggunaan</a>
                    <a href="#" class="text-slate-400 hover:text-white hover:translate-x-2 transition-all duration-300 text-sm inline-block">→ FAQ</a>
                  </div>
                </div>

                <!-- Social Media -->
                <div class="text-center md:text-right">
                  <h4 class="text-white font-semibold mb-4 text-base">Ikuti Kami</h4>
                  <div class="flex gap-3 justify-center md:justify-end">
                    <a href="#" class="p-3 bg-slate-800/50 rounded-xl hover:bg-blue-600/30 transition-all duration-300 card-3d hover:scale-125 hover:rotate-12 group">
                      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-slate-400 group-hover:text-blue-400 transition-colors">
                        <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"></path>
                      </svg>
                    </a>
                    <a href="#" class="p-3 bg-slate-800/50 rounded-xl hover:bg-blue-600/30 transition-all duration-300 card-3d hover:scale-125 hover:rotate-12 group">
                      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-slate-400 group-hover:text-blue-400 transition-colors">
                        <path d="M23 3a10.9 10.9 0 0 1-3.14 1.53 4.48 4.48 0 0 0-7.86 3v1A10.66 10.66 0 0 1 3 4s-4 9 5 13a11.64 11.64 0 0 1-7 2c9 5 20 0 20-11.5a4.5 4.5 0 0 0-.08-.83A7.72 7.72 0 0 0 23 3z"></path>
                      </svg>
                    </a>
                    <a href="#" class="p-3 bg-slate-800/50 rounded-xl hover:bg-pink-600/30 transition-all duration-300 card-3d hover:scale-125 hover:rotate-12 group">
                      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-slate-400 group-hover:text-pink-400 transition-colors">
                        <rect x="2" y="2" width="20" height="20" rx="5" ry="5"></rect>
                        <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"></path>
                        <line x1="17.5" y1="6.5" x2="17.51" y2="6.5"></line>
                      </svg>
                    </a>
                  </div>
                </div>
              </div>

              <!-- Copyright -->
              <div class="border-t border-slate-700/50 pt-6 text-center" style="transform: translateZ(10px);">
                <p class="text-slate-400 text-sm mb-2">
                  © 2024 FastUnduh. Made with <span class="text-red-500 animate-pulse inline-block">❤️</span> for video lovers.
                </p>
                <p class="text-slate-500 text-xs leading-relaxed max-w-2xl mx-auto">
                  Disclaimer: FastUnduh tidak menyimpan atau menghosting video apapun. Semua video diunduh langsung dari sumber aslinya.
                </p>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </main>
  </div>
</template>

<style scoped>
/* Additional component-specific styles if needed */
</style>
