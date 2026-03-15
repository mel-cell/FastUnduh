<script setup>
import SiteNavbar from './components/SiteNavbar.vue';
import HeroDownloader from './components/HeroDownloader.vue';
import SupportedPlatforms from './components/SupportedPlatforms.vue';
import AppFeatures from './components/AppFeatures.vue';
import HowToSteps from './components/HowToSteps.vue';
import SiteFooter from './components/SiteFooter.vue';
</script>

<template>
  <div class="bg-white text-neutral-950 font-sans antialiased min-h-screen flex flex-col">
    <!-- BEGIN: Navigation -->
    <SiteNavbar />
    <!-- END: Navigation -->

    <main class="flex-grow">
      <!-- BEGIN: Hero Section (with Downloader Form) -->
      <HeroDownloader />
      <!-- END: Hero Section -->

      <!-- BEGIN: Supported Platforms Section -->
      <SupportedPlatforms />
      <!-- END: Supported Platforms Section -->

      <!-- BEGIN: Feature List Section -->
      <AppFeatures />
      <!-- END: Feature List Section -->

      <!-- BEGIN: How-To Steps -->
      <HowToSteps />
      <!-- END: How-To Steps -->
    </main>

    <!-- BEGIN: Footer -->
    <SiteFooter />
    <!-- END: Footer -->
  </div>
</template>

<style>
/* Global Layout Adjustments specific to this layout */
.grid-border {
  border: 1px solid #e5e7eb;
}
.input-shadow-none {
  box-shadow: none !important;
}
</style>
