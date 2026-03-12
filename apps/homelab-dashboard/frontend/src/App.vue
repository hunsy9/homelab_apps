<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'

const sidebarOpen = ref(false)

function closeSidebar() {
  sidebarOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-background flex">
    <!-- Mobile header -->
    <header class="md:hidden fixed top-0 left-0 right-0 h-14 bg-card border-b flex items-center px-4 z-40">
      <button @click="sidebarOpen = true" class="p-2 -ml-2 rounded-md hover:bg-accent">
        <svg class="w-5 h-5 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <h1 class="text-lg font-semibold text-foreground ml-3">Homelab</h1>
    </header>

    <!-- Mobile overlay -->
    <div
      v-if="sidebarOpen"
      class="md:hidden fixed inset-0 bg-black/50 z-40"
      @click="closeSidebar"
    ></div>

    <!-- Sidebar -->
    <aside
      class="fixed md:static inset-y-0 left-0 w-64 border-r bg-card p-6 z-50 transform transition-transform duration-200 ease-in-out md:transform-none"
      :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
    >
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-xl font-semibold text-foreground">Homelab</h1>
        <button @click="closeSidebar" class="md:hidden p-1 rounded-md hover:bg-accent">
          <svg class="w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <nav class="space-y-1">
        <RouterLink
          to="/"
          @click="closeSidebar"
          class="flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
          active-class="bg-secondary text-foreground"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
          Overview
        </RouterLink>
        <RouterLink
          to="/namespaces"
          @click="closeSidebar"
          class="flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
          active-class="bg-secondary text-foreground"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          Namespaces
        </RouterLink>
      </nav>
    </aside>

    <!-- Main content -->
    <main class="flex-1 p-4 md:p-8 bg-muted/30 pt-18 md:pt-8">
      <RouterView />
    </main>
  </div>
</template>
