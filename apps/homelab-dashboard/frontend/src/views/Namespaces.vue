<script setup lang="ts">
import { onMounted } from 'vue'
import { useClusterStore } from '../stores/cluster'

const store = useClusterStore()

onMounted(() => {
  store.fetchNamespaces()
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-foreground">Namespaces</h1>
      <p class="text-muted-foreground text-sm mt-1">Manage your Kubernetes namespaces</p>
    </div>

    <div v-if="store.loading" class="text-muted-foreground">Loading...</div>
    <div v-else-if="store.error" class="text-red-500">{{ store.error }}</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <RouterLink
        v-for="ns in store.namespaces"
        :key="ns.name"
        :to="`/namespaces/${ns.name}`"
        class="bg-card rounded-lg border p-4 hover:border-foreground/20 hover:shadow-sm transition-all"
      >
        <div class="flex items-center justify-between mb-2">
          <span class="font-medium text-foreground">{{ ns.name }}</span>
          <span
            class="px-2 py-0.5 rounded-full text-xs font-medium"
            :class="ns.status === 'Active' ? 'bg-emerald-50 text-emerald-700' : 'bg-secondary text-muted-foreground'"
          >
            {{ ns.status }}
          </span>
        </div>
        <div class="text-muted-foreground text-sm">
          Created: {{ new Date(ns.createdAt).toLocaleDateString() }}
        </div>
      </RouterLink>
    </div>
  </div>
</template>
