<script setup lang="ts">
import { onMounted } from 'vue'
import { useClusterStore } from '../stores/cluster'

const store = useClusterStore()

onMounted(() => {
  store.fetchNamespaces()
})
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Namespaces</h2>

    <div v-if="store.loading" class="text-gray-400">Loading...</div>
    <div v-else-if="store.error" class="text-red-400">{{ store.error }}</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <RouterLink
        v-for="ns in store.namespaces"
        :key="ns.name"
        :to="`/namespaces/${ns.name}`"
        class="bg-gray-800 rounded-lg p-4 hover:bg-gray-750 transition cursor-pointer"
      >
        <div class="flex items-center justify-between mb-2">
          <span class="font-medium">{{ ns.name }}</span>
          <span
            class="px-2 py-1 rounded text-xs"
            :class="ns.status === 'Active' ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-300'"
          >
            {{ ns.status }}
          </span>
        </div>
        <div class="text-gray-500 text-sm">
          Created: {{ new Date(ns.createdAt).toLocaleDateString() }}
        </div>
      </RouterLink>
    </div>
  </div>
</template>
