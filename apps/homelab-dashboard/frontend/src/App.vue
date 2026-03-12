<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-gray-900 text-white flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-800 p-4">
      <h1 class="text-xl font-bold mb-8 text-blue-400">Homelab Dashboard</h1>
      <nav class="space-y-2">
        <RouterLink
          to="/"
          class="block px-4 py-2 rounded hover:bg-gray-700 transition"
          active-class="bg-blue-600"
        >
          Overview
        </RouterLink>
        <RouterLink
          to="/namespaces"
          class="block px-4 py-2 rounded hover:bg-gray-700 transition"
          active-class="bg-blue-600"
        >
          Namespaces
        </RouterLink>
      </nav>
    </aside>

    <!-- Main content -->
    <main class="flex-1 p-8">
      <RouterView />
    </main>
  </div>
</template>
