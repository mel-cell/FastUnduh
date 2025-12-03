<script setup>
import { ref } from "vue";

const url = ref("");
const isLoading = ref(false);
const jobTicket = ref(null);

const handleDownload = async () => {
  if (!url.value) return;
  isLoading.value = true;

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
    alert("Terjadi kesalahan: " + error.message);
    isLoading.value = false;
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

        // Tampilkan info sukses
        const filename = data.filename;
        const title = data.title || "Video";

        // Trigger download otomatis atau tampilkan tombol
        // Kita redirect aja ke endpoint download
        window.location.href = `http://localhost:3000/api/download/${filename}`;

        alert(`Download Selesai! File: ${title}`);
      } else if (status === "failed") {
        clearInterval(interval);
        isLoading.value = false;
        alert("Download Gagal. Cek URL atau coba lagi.");
      } else {
        console.log("Status:", status);
      }
    } catch (e) {
      clearInterval(interval);
      isLoading.value = false;
      console.error("Polling error", e);
    }
  }, 2000);
};
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4 relative overflow-hidden">
    <!-- Background Elements -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden -z-10">
      <div
        class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-600/20 rounded-full blur-[100px]"
      ></div>
      <div
        class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-purple-600/20 rounded-full blur-[100px]"
      ></div>
    </div>

    <!-- Main Container -->
    <main class="w-full max-w-3xl text-center z-10">
      <!-- Header -->
      <div class="mb-12 space-y-4">
        <h1
          class="text-5xl md:text-7xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-400 tracking-tight"
        >
          FastUnduh
        </h1>
        <p class="text-lg md:text-xl text-slate-400 max-w-xl mx-auto">
          Download video dari berbagai platform dengan kecepatan kilat. Tanpa iklan, tanpa ribet.
        </p>
      </div>

      <!-- Input Section -->
      <div class="glass-panel p-2 rounded-2xl input-glow transition-all duration-300">
        <div class="flex flex-col md:flex-row gap-2">
          <input
            v-model="url"
            type="text"
            placeholder="Tempel link video di sini (YouTube, Instagram, Twitter...)"
            class="flex-1 bg-transparent border-none outline-none text-white placeholder-slate-500 px-6 py-4 text-lg w-full"
            @keyup.enter="handleDownload"
          />
          <button
            @click="handleDownload"
            :disabled="isLoading"
            class="bg-blue-600 hover:bg-blue-500 text-white font-semibold px-8 py-4 rounded-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 min-w-[160px]"
          >
            <span v-if="!isLoading">Download</span>
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

      <!-- Features / Footer -->
      <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-6 text-slate-400 text-sm">
        <div class="flex flex-col items-center gap-2">
          <div class="p-3 bg-slate-800/50 rounded-full">
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
              class="text-blue-400"
            >
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
            </svg>
          </div>
          <span class="font-medium">Super Cepat</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <div class="p-3 bg-slate-800/50 rounded-full">
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
              class="text-purple-400"
            >
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="7 10 12 15 17 10" />
              <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
          </div>
          <span class="font-medium">Tanpa Batas</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <div class="p-3 bg-slate-800/50 rounded-full">
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
              class="text-green-400"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
            </svg>
          </div>
          <span class="font-medium">Aman & Privat</span>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
/* Additional component-specific styles if needed */
</style>
