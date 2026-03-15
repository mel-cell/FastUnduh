<script setup>
import { ref, onUnmounted } from "vue";

const url = ref("");
const isLoading = ref(false);
const downloadStatus = ref("");
const downloadResult = ref(null);
const progressPercent = ref(0);
let progressInterval = null;

const startFakeProgress = () => {
  progressPercent.value = 0;
  progressInterval = setInterval(() => {
    if (progressPercent.value < 90) {
      // Simulate slow down towards the end
      const increment = Math.max(0.5, (95 - progressPercent.value) / 10);
      progressPercent.value = Math.min(95, progressPercent.value + increment);
    }
  }, 300);
};

const stopFakeProgress = () => {
  if (progressInterval) {
    clearInterval(progressInterval);
    progressInterval = null;
  }
};

onUnmounted(() => {
  stopFakeProgress();
});

const handleDownload = async () => {
  if (!url.value) return;
  isLoading.value = true;
  downloadStatus.value = "processing";
  downloadResult.value = null;
  
  startFakeProgress();

  try {
    const response = await fetch("/api/queue", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url: url.value }),
    });

    if (!response.ok) throw new Error("Failed to process the link.");

    const data = await response.json();
    pollStatus(data.job_id);
  } catch (error) {
    console.error(error);
    stopFakeProgress();
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
        stopFakeProgress();
        progressPercent.value = 100;
        
        // Add a tiny delay so user can see it hit 100%
        setTimeout(() => {
          isLoading.value = false;
          downloadStatus.value = "completed";
          downloadResult.value = {
            filename: data.filename,
            title: data.title || "Ready for download"
          };
        }, 400);

      } else if (status === "failed") {
        clearInterval(interval);
        stopFakeProgress();
        isLoading.value = false;
        downloadStatus.value = "failed";
        downloadResult.value = { error: "We couldn't extract media from this link." };
      }
    } catch {
      clearInterval(interval);
      stopFakeProgress();
      isLoading.value = false;
      downloadStatus.value = "failed";
      downloadResult.value = { error: "Connection lost. Please try again." };
    }
  }, 2000);
};

const startDownloadFile = () => {
  if (downloadResult.value && downloadResult.value.filename) {
    window.location.href = `/api/download/${downloadResult.value.filename}`;
  }
};

const resetUI = () => {
  url.value = "";
  downloadStatus.value = "";
  downloadResult.value = null;
  progressPercent.value = 0;
  stopFakeProgress();
};
</script>

<template>
  <section class="py-20 lg:py-32" data-purpose="hero-section">
    <div class="max-w-4xl mx-auto px-4 text-center">
      <h1 class="text-4xl md:text-6xl font-extrabold text-neutral-950 mb-6 leading-tight">
        Download high-quality content from any platform.
      </h1>
      <p class="text-lg md:text-xl text-gray-500 mb-12 max-w-2xl mx-auto">
        No watermarks. No complexity. Just paste the link and save your media in seconds.
      </p>
      
      <!-- Search/Input Container -->
      <div class="w-full max-w-2xl mx-auto" data-purpose="downloader-form">
        
        <form v-if="downloadStatus === ''" @submit.prevent="handleDownload" class="flex flex-col md:flex-row gap-3">
          <input 
            v-model="url"
            :disabled="isLoading"
            class="flex-grow px-5 py-4 bg-gray-50 border border-gray-200 focus:ring-2 focus:ring-brand focus:border-brand border-solid rounded-none text-base outline-none transition-all disabled:opacity-50" 
            placeholder="Paste link here (e.g., https://tiktok.com/@user/video/...)" 
            required 
            type="text"
          />
          <button 
            :disabled="!url"
            class="bg-brand hover:bg-brand-dark text-white px-8 py-4 font-bold transition-colors uppercase tracking-wider text-sm flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed" 
            style="min-width: 160px;" 
            type="submit"
          >
            Download
          </button>
        </form>

        <p v-if="downloadStatus === ''" class="mt-4 text-xs text-gray-400">
          By using our service you accept our <a class="underline" href="#">Terms of Service</a>.
        </p>

        <!-- Loading State: Clean Progress Bar -->
        <div v-if="isLoading" class="mt-2 p-8 bg-white border border-gray-100 flex flex-col items-center shadow-sm">
           <h3 class="text-lg font-bold text-neutral-950 mb-2">Analyzing Link...</h3>
           <p class="text-sm text-gray-400 mb-6">Fetching highest available quality from source servers</p>
           
           <!-- Progress bar track -->
           <div class="w-full max-w-md bg-gray-100 h-2 rounded-full overflow-hidden mb-2">
             <!-- Progress bar fill -->
             <div 
               class="bg-brand h-full transition-all duration-300 ease-out"
               :style="{ width: `${progressPercent}%` }"
             ></div>
           </div>
           
           <span class="text-xs font-semibold text-gray-500">{{ Math.round(progressPercent) }}% Processing</span>
        </div>

        <!-- Status Feedback -->
        <div v-if="downloadStatus === 'completed'" class="mt-2 p-8 bg-gray-50 border border-gray-200 flex flex-col items-center">
            <div class="w-12 h-12 rounded-full bg-green-100 text-green-600 flex items-center justify-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <p class="text-neutral-950 font-bold mb-2 text-xl">Ready to save!</p>
            <p class="text-sm text-gray-600 mb-8 truncate max-w-full italic px-4 font-medium">{{ downloadResult.title }}</p>
            
            <div class="flex flex-col sm:flex-row items-center justify-center gap-4 w-full max-w-md">
              <button @click="startDownloadFile" class="bg-brand hover:bg-brand-dark text-white font-bold py-4 px-8 text-sm uppercase tracking-wide w-full sm:w-auto transition-colors flex items-center justify-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
                Save File
              </button>
              <button @click="resetUI" class="text-gray-500 hover:text-neutral-950 font-bold py-4 px-6 text-sm uppercase tracking-wide w-full sm:w-auto transition-colors">Convert another</button>
            </div>
        </div>
        
        <!-- Error State -->
        <div v-if="downloadStatus === 'failed'" class="mt-2 p-8 border border-gray-200 flex flex-col items-center">
            <div class="w-12 h-12 rounded-full bg-red-50 text-red-600 flex items-center justify-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <p class="text-neutral-950 font-bold mb-2 text-xl">Download Failed</p>
            <p class="text-sm text-gray-500 mb-8 max-w-sm">{{ downloadResult.error || 'Failed to download video. Please check the URL or try again later.' }}</p>
            <button @click="resetUI" class="bg-neutral-950 hover:bg-neutral-800 text-white font-bold py-4 px-10 text-sm transition-colors uppercase tracking-wide">Try Again</button>
        </div>
        
      </div>
    </div>
  </section>
</template>
