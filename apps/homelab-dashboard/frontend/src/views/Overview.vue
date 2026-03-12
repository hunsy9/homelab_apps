<script setup lang="ts">
import { onMounted, onUnmounted, computed, ref } from 'vue'
import { useClusterStore } from '../stores/cluster'

const store = useClusterStore()
const refreshInterval = ref<number | null>(null)
const isRefreshing = ref(false)

async function refreshNodes() {
  isRefreshing.value = true
  await store.fetchNodes()
  setTimeout(() => { isRefreshing.value = false }, 300)
}

onMounted(() => {
  store.fetchNodes()
  store.fetchNamespaces()
  refreshInterval.value = window.setInterval(() => {
    store.fetchNodes()
  }, 5000)
})

onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})

const readyNodes = computed(() => store.nodes.filter(n => n.status === 'Ready').length)
const totalNodes = computed(() => store.nodes.length)

const services = [
  { name: 'ArgoCD', url: 'https://argocd.seung.site', description: 'GitOps CD' },
  { name: 'Traefik', url: 'https://traefik.seung.site/dashboard/', description: 'Ingress' },
  { name: 'Longhorn', url: 'https://longhorn.seung.site', description: 'Storage' },
  { name: 'Grafana', url: 'https://grafana.seung.site', description: 'Monitoring' },
  { name: 'Harbor', url: 'https://harbor.seung.site', description: 'Registry' },
]

function formatCPU(millicores: number): string {
  if (millicores >= 1000) return `${(millicores / 1000).toFixed(1)}`
  return `${(millicores / 1000).toFixed(2)}`
}

function formatMemory(bytes: number): string {
  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1) return `${gb.toFixed(1)} GB`
  const mb = bytes / (1024 * 1024)
  return `${mb.toFixed(0)} MB`
}

function getUsagePercent(usage: number, allocatable: number): number {
  if (allocatable === 0) return 0
  return Math.min(100, (usage / allocatable) * 100)
}

function getProgressColor(percent: number): string {
  if (percent >= 90) return 'bg-red-500'
  if (percent >= 70) return 'bg-amber-500'
  return 'bg-emerald-500'
}
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-semibold text-foreground">Cluster Overview</h1>
      <p class="text-muted-foreground text-sm mt-1">Monitor your Kubernetes cluster resources</p>
    </div>

    <!-- Stats cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-card rounded-lg border p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-muted-foreground">Nodes</p>
          <svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" />
          </svg>
        </div>
        <p class="text-3xl font-bold text-foreground mt-2">{{ readyNodes }}<span class="text-muted-foreground text-lg font-normal">/{{ totalNodes }}</span></p>
        <p class="text-xs text-muted-foreground mt-1">Ready nodes</p>
      </div>

      <div class="bg-card rounded-lg border p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-muted-foreground">Namespaces</p>
          <svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <p class="text-3xl font-bold text-foreground mt-2">{{ store.namespaces.length }}</p>
        <p class="text-xs text-muted-foreground mt-1">Active namespaces</p>
      </div>

      <div class="bg-card rounded-lg border p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-muted-foreground">Status</p>
          <span
            class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium"
            :class="readyNodes === totalNodes ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'"
          >
            {{ readyNodes === totalNodes ? 'Healthy' : 'Degraded' }}
          </span>
        </div>
        <p class="text-3xl font-bold text-foreground mt-2">{{ readyNodes === totalNodes ? '100%' : Math.round(readyNodes / totalNodes * 100) + '%' }}</p>
        <p class="text-xs text-muted-foreground mt-1">Cluster health</p>
      </div>
    </div>

    <!-- Services -->
    <div>
      <h2 class="text-lg font-semibold text-foreground mb-4">Services</h2>
      <div class="grid grid-cols-2 md:grid-cols-5 gap-3">
        <a
          v-for="service in services"
          :key="service.name"
          :href="service.url"
          target="_blank"
          rel="noopener noreferrer"
          class="group bg-card rounded-lg border p-4 hover:border-foreground/20 hover:shadow-sm transition-all"
        >
          <p class="font-medium text-foreground group-hover:text-primary">{{ service.name }}</p>
          <p class="text-xs text-muted-foreground mt-1">{{ service.description }}</p>
        </a>
      </div>
    </div>

    <!-- Nodes -->
    <div>
      <div class="flex items-center gap-3 mb-4">
        <h2 class="text-lg font-semibold text-foreground">Nodes</h2>
        <button
          @click="refreshNodes"
          class="p-1.5 rounded-md hover:bg-accent transition-colors"
          title="Refresh"
        >
          <svg
            class="w-4 h-4 text-muted-foreground"
            :class="{ 'animate-spin': isRefreshing }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
        <span class="text-xs text-muted-foreground">auto-refresh 5s</span>
      </div>

      <div v-if="store.loading" class="text-muted-foreground">Loading...</div>
      <div v-else-if="store.error" class="text-red-500">{{ store.error }}</div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        <div
          v-for="node in store.nodes"
          :key="node.name"
          class="bg-card rounded-lg border p-4"
        >
          <!-- Header -->
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <span
                class="w-2 h-2 rounded-full"
                :class="node.status === 'Ready' ? 'bg-emerald-500' : 'bg-red-500'"
              ></span>
              <span class="font-medium text-foreground">{{ node.name }}</span>
            </div>
            <span class="text-xs px-2 py-0.5 rounded-full bg-secondary text-muted-foreground">
              {{ node.roles.join(', ') }}
            </span>
          </div>

          <!-- IP -->
          <p class="text-xs text-muted-foreground mb-4">{{ node.ip }}</p>

          <!-- Resources -->
          <div class="space-y-3">
            <!-- CPU -->
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-muted-foreground">CPU</span>
                <span class="text-foreground font-medium">{{ formatCPU(node.cpu.usage) }} / {{ formatCPU(node.cpu.allocatable) }} cores</span>
              </div>
              <div class="h-1.5 bg-secondary rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="getProgressColor(getUsagePercent(node.cpu.usage, node.cpu.allocatable))"
                  :style="{ width: `${getUsagePercent(node.cpu.usage, node.cpu.allocatable)}%` }"
                ></div>
              </div>
            </div>

            <!-- Memory -->
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-muted-foreground">Memory</span>
                <span class="text-foreground font-medium">{{ formatMemory(node.memory.usage) }} / {{ formatMemory(node.memory.allocatable) }}</span>
              </div>
              <div class="h-1.5 bg-secondary rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="getProgressColor(getUsagePercent(node.memory.usage, node.memory.allocatable))"
                  :style="{ width: `${getUsagePercent(node.memory.usage, node.memory.allocatable)}%` }"
                ></div>
              </div>
            </div>

            <!-- Disk -->
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-muted-foreground">Disk</span>
                <span class="text-foreground font-medium">{{ formatMemory(node.storage.usage) }} / {{ formatMemory(node.storage.capacity) }}</span>
              </div>
              <div class="h-1.5 bg-secondary rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="getProgressColor(getUsagePercent(node.storage.usage, node.storage.capacity))"
                  :style="{ width: `${getUsagePercent(node.storage.usage, node.storage.capacity)}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
