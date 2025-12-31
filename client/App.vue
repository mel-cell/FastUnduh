<script setup>
import { ref } from "vue";
import GameModal from "./components/GameModal.vue";
import Collaborators from "./components/Collaborators.vue";

const url = ref("");
const isLoading = ref(false);
const showGame = ref(false);
const downloadStatus = ref("");
const downloadResult = ref(null);
const jobTicket = ref(null);

const handleDownload = async () => {
  if (!url.value) return;
  isLoading.value = true;
  showGame.value = true;
  downloadStatus.value = "processing";
  downloadResult.value = null;

  try {
    const response = await fetch("/api/queue", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url: url.value }),
    });

    if (!response.ok) throw new Error("Gagal memulai download. Cek URL Anda.");

    const data = await response.json();
    jobTicket.value = data.job_id;
    pollStatus(data.job_id);
  } catch (error) {
    console.error(error);
    downloadStatus.value = "failed";
    downloadResult.value = { error: error.message };
    isLoading.value = false;
  }
};

const pollStatus = (id) => {
  const interval = setInterval(async () => {
    try {
      const res = await fetch(`/api/status/${id}`);
      if (!res.ok) return;

      const data = await res.json();
      const status = data.status;

      if (status === "completed") {
        clearInterval(interval);
        isLoading.value = false;
        downloadStatus.value = "completed";
        downloadResult.value = {
          filename: data.filename,
          title: data.title || "Video Berhasil Diunduh"
        };
      } else if (status === "failed") {
        clearInterval(interval);
        isLoading.value = false;
        downloadStatus.value = "failed";
        downloadResult.value = { error: "Proses gagal di server." };
      }
    } catch (e) {
      clearInterval(interval);
      isLoading.value = false;
      downloadStatus.value = "failed";
      downloadResult.value = { error: "Koneksi terputus." };
    }
  }, 2000);
};

const startDownloadFile = () => {
  if (downloadResult.value && downloadResult.value.filename) {
    window.location.href = `/api/download/${downloadResult.value.filename}`;
  }
};

const closeGameModal = () => {
  showGame.value = false;
};
</script>

<template>
  <div class="min-h-screen flex flex-col items-center p-4 relative overflow-x-hidden text-slate-200 selection:bg-blue-500/30">
    
    <!-- Background Gradients -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden -z-10 bg-slate-900">
      <div class="absolute top-[-20%] left-[20%] w-[60%] h-[60%] bg-blue-600/10 rounded-full blur-[120px] animate-pulse-slow"></div>
      <div class="absolute top-[40%] right-[-10%] w-[50%] h-[50%] bg-purple-600/10 rounded-full blur-[120px] animate-pulse-slow delay-1000"></div>
    </div>

    <!-- Game Modal -->
    <GameModal 
      :isOpen="showGame"
      :downloadStatus="downloadStatus"
      :downloadResult="downloadResult"
      @close="closeGameModal"
      @download="startDownloadFile"
    />

    <!-- Main Content -->
    <main class="w-full max-w-5xl flex flex-col items-center text-center z-0 px-4 mt-16 md:mt-24">
      
      <!-- Brand / Header -->
      <div class="mb-14 animate-fade-in-down flex flex-col items-center">
         <div class="mb-6 p-4 bg-gradient-to-br from-slate-800 to-slate-900 rounded-3xl border border-white/10 shadow-2xl shadow-blue-900/20 transform hover:-translate-y-1 transition-transform duration-500">
            <!-- New Icon: Download Cloud -->
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-blue-400 drop-shadow-lg"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
         </div>
         <h1 class="text-5xl md:text-7xl font-extrabold text-white tracking-tight mb-5 drop-shadow-xl">
           Fast<span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-400">Unduh</span>
         </h1>
         <p class="text-lg md:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed">
           Download video dari platform favoritmu seketika.
         </p>
      </div>

      <!-- Input Section with Glow -->
      <div class="w-full max-w-2xl relative group mb-20 z-20">
        <div class="absolute -inset-1 bg-gradient-to-r from-blue-500 to-purple-500 rounded-2xl blur opacity-30 group-hover:opacity-60 transition duration-500 will-change-transform"></div>
        <div class="relative bg-slate-900/90 backdrop-blur-3xl p-2 rounded-2xl border border-white/10 flex flex-col md:flex-row gap-2 shadow-2xl ring-1 ring-white/5">
           <!-- Icon inside input -->
           <div class="absolute left-6 top-1/2 -translate-y-1/2 text-slate-500 pointer-events-none hidden md:block">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path></svg>
           </div>
           
           <input
             v-model="url"
             type="text"
             placeholder="Tempel link video di sini..."
             class="flex-1 bg-transparent border-none outline-none text-white placeholder-slate-500 pl-4 md:pl-12 pr-4 py-4 text-base md:text-lg w-full font-medium"
             @keyup.enter="handleDownload"
           />
           <button
             @click="handleDownload"
             :disabled="isLoading"
             class="bg-blue-600 hover:bg-blue-500 text-white font-bold px-8 py-3.5 rounded-xl transition-all duration-300 transform md:hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 min-w-[140px] shadow-lg shadow-blue-500/25"
           >
             <span v-if="!isLoading" class="flex items-center gap-2">
               Download 
               <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
             </span>
             <span v-else class="flex items-center gap-2">
               <svg class="animate-spin h-5 w-5 text-white/50" viewBox="0 0 24 24" fill="none">
                 <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                 <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
               </svg>
               Proses
             </span>
           </button>
        </div>
      </div>

      <!-- Features Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full max-w-5xl mb-12">
         <!-- Card 1 -->
         <div class="feature-card group">
            <div class="mb-4 p-3 bg-yellow-500/10 rounded-2xl group-hover:bg-yellow-500/20 transition-colors">
               <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-yellow-400"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
            </div>
            <h3 class="font-bold text-white text-lg mb-2">Super Cepat</h3>
            <p class="text-slate-400 text-sm leading-relaxed">Server berkecepatan tinggi menjamin file Anda siap dalam hitungan detik.</p>
         </div>
         
         <!-- Card 2 -->
         <div class="feature-card group">
            <div class="mb-4 p-3 bg-purple-500/10 rounded-2xl group-hover:bg-purple-500/20 transition-colors">
               <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-400"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
            </div>
            <h3 class="font-bold text-white text-lg mb-2">Mini Games</h3>
            <p class="text-slate-400 text-sm leading-relaxed">Bosan menunggu? Mainkan game klasik selagi file Anda diproses.</p>
         </div>

         <!-- Card 3 -->
         <div class="feature-card group">
            <div class="mb-4 p-3 bg-green-500/10 rounded-2xl group-hover:bg-green-500/20 transition-colors">
               <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-green-400"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
            </div>
            <h3 class="font-bold text-white text-lg mb-2">Aman & Privat</h3>
            <p class="text-slate-400 text-sm leading-relaxed">Tanpa pelacak. File Anda diproses di memori dan langsung dihapus.</p>
         </div>
      </div>

      <!-- Github Contributors -->
      <!-- <Collaborators owner="mel" repo="FastUnduh" /> -->

    </main>

    <!-- Footer -->
    <footer class="text-slate-600 text-xs py-8 w-full text-center mt-auto">
      &copy; 2025 FastUnduh. Dibuat dengan <span class="text-red-500">❤️</span> oleh Tim Developer
    </footer>

  </div>
</template>

<style scoped>
/* Scoped styles here because they are specific to App layout structure */
.feature-card {
  @apply bg-slate-800/40 backdrop-blur-md border border-white/5 rounded-3xl p-6 flex flex-col items-center hover:bg-slate-800/60 transition-all duration-300 hover:-translate-y-1 hover:border-white/10 hover:shadow-xl hover:shadow-black/20;
}

.animate-pulse-slow {
  animation: pulse 8s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 0.8; }
}
</style>
