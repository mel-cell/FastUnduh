<template>
<footer class="bg-white border-t border-gray-100 py-12" data-purpose="site-footer">
  <div class="max-w-7xl mx-auto px-4 flex flex-col md:flex-row justify-between items-center gap-8">
    <div class="flex flex-col gap-2">
      <span class="text-lg font-bold text-neutral-950">FastUnduh</span>
      <p class="text-xs text-gray-400">© {{ new Date().getFullYear() }} FastUnduh. Built for productivity.</p>
    </div>
    <div class="flex space-x-8 text-sm font-medium text-gray-500">
      <a class="hover:text-neutral-950 transition-colors" href="#">Terms</a>
      <a class="hover:text-neutral-950 transition-colors" href="#">Privacy</a>
      <a class="hover:text-neutral-950 transition-colors" href="#">Cookies</a>
      <a class="hover:text-neutral-950 transition-colors" href="#">Contact</a>
    </div>
  </div>
</footer>
</template>
