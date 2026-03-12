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
  // Auto-refresh nodes every 5 seconds
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
  {
    name: 'ArgoCD',
    url: 'https://argocd.seung.site',
    description: 'GitOps Continuous Delivery',
    icon: '🔄',
    color: 'from-orange-500 to-red-500'
  },
  {
    name: 'Traefik',
    url: 'https://traefik.seung.site/dashboard/',
    description: 'Ingress Controller Dashboard',
    icon: '🔀',
    color: 'from-cyan-500 to-blue-500'
  },
  {
    name: 'Longhorn',
    url: 'https://longhorn.seung.site',
    description: 'Distributed Storage',
    icon: '💾',
    color: 'from-red-500 to-orange-500'
  },
  {
    name: 'Grafana',
    url: 'https://grafana.seung.site',
    description: 'Monitoring & Dashboards',
    icon: '📊',
    color: 'from-orange-500 to-yellow-500'
  },
  {
    name: 'Harbor',
    url: 'https://harbor.seung.site',
    description: 'Container Registry',
    icon: '🐳',
    color: 'from-blue-500 to-indigo-500'
  }
]

function formatCPU(millicores: number): string {
  if (millicores >= 1000) {
    return `${(millicores / 1000).toFixed(1)} cores`
  }
  return `${millicores}m`
}

function formatMemory(bytes: number): string {
  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1) {
    return `${gb.toFixed(1)} GB`
  }
  const mb = bytes / (1024 * 1024)
  return `${mb.toFixed(0)} MB`
}

function getUsagePercent(usage: number, allocatable: number): number {
  if (allocatable === 0) return 0
  return Math.min(100, (usage / allocatable) * 100)
}

function getProgressColor(percent: number): string {
  if (percent >= 90) return 'bg-red-500'
  if (percent >= 70) return 'bg-yellow-500'
  return 'bg-green-500'
}
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Cluster Overview</h2>

    <!-- Stats cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-gray-800 rounded-lg p-6">
        <div class="text-gray-400 text-sm">Nodes</div>
        <div class="text-3xl font-bold text-green-400">
          {{ readyNodes }}/{{ totalNodes }}
        </div>
        <div class="text-gray-500 text-sm">Ready</div>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <div class="text-gray-400 text-sm">Namespaces</div>
        <div class="text-3xl font-bold text-blue-400">
          {{ store.namespaces.length }}
        </div>
        <div class="text-gray-500 text-sm">Active</div>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <div class="text-gray-400 text-sm">Cluster Status</div>
        <div class="text-3xl font-bold" :class="readyNodes === totalNodes ? 'text-green-400' : 'text-yellow-400'">
          {{ readyNodes === totalNodes ? 'Healthy' : 'Degraded' }}
        </div>
      </div>
    </div>

    <!-- Services -->
    <div class="mb-8">
      <h3 class="text-lg font-semibold mb-4">Services</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <a
          v-for="service in services"
          :key="service.name"
          :href="service.url"
          target="_blank"
          rel="noopener noreferrer"
          class="group relative overflow-hidden rounded-lg p-4 bg-gray-800 hover:scale-105 transition-all duration-200"
        >
          <div class="absolute inset-0 opacity-0 group-hover:opacity-20 transition-opacity bg-gradient-to-br" :class="service.color"></div>
          <div class="relative">
            <div class="text-3xl mb-2">{{ service.icon }}</div>
            <div class="font-semibold text-white group-hover:text-blue-400 transition-colors">
              {{ service.name }}
            </div>
            <div class="text-xs text-gray-500 mt-1">{{ service.description }}</div>
          </div>
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity text-gray-400">
            ↗
          </div>
        </a>
      </div>
    </div>

    <!-- Nodes -->
    <div>
      <div class="flex items-center gap-3 mb-4">
        <h3 class="text-lg font-semibold">Nodes</h3>
        <button
          @click="refreshNodes"
          class="p-1.5 rounded hover:bg-gray-700 transition-colors"
          title="Refresh"
        >
          <svg
            class="w-4 h-4 text-gray-400"
            :class="{ 'animate-spin': isRefreshing }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        </button>
        <span class="text-xs text-gray-500">auto-refresh 5s</span>
      </div>
      <div v-if="store.loading" class="text-gray-400">Loading...</div>
      <div v-else-if="store.error" class="text-red-400">{{ store.error }}</div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        <div
          v-for="node in store.nodes"
          :key="node.name"
          class="bg-gray-800 rounded-lg p-4"
        >
          <!-- Header -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span
                class="w-2.5 h-2.5 rounded-full"
                :class="node.status === 'Ready' ? 'bg-green-500' : 'bg-red-500'"
              ></span>
              <span class="font-semibold">{{ node.name }}</span>
            </div>
            <span class="text-xs px-2 py-0.5 rounded bg-gray-700 text-gray-400">
              {{ node.roles.join(', ') }}
            </span>
          </div>

          <!-- IP -->
          <div class="text-xs text-gray-500 mb-3">{{ node.ip }}</div>

          <!-- CPU -->
          <div class="mb-3">
            <div class="flex justify-between text-xs text-gray-400 mb-1">
              <span>CPU</span>
              <span>{{ formatCPU(node.cpu.usage) }} / {{ formatCPU(node.cpu.allocatable) }}</span>
            </div>
            <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="getProgressColor(getUsagePercent(node.cpu.usage, node.cpu.allocatable))"
                :style="{ width: `${getUsagePercent(node.cpu.usage, node.cpu.allocatable)}%` }"
              ></div>
            </div>
          </div>

          <!-- Memory -->
          <div class="mb-3">
            <div class="flex justify-between text-xs text-gray-400 mb-1">
              <span>Memory</span>
              <span>{{ formatMemory(node.memory.usage) }} / {{ formatMemory(node.memory.allocatable) }}</span>
            </div>
            <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="getProgressColor(getUsagePercent(node.memory.usage, node.memory.allocatable))"
                :style="{ width: `${getUsagePercent(node.memory.usage, node.memory.allocatable)}%` }"
              ></div>
            </div>
          </div>

          <!-- Storage -->
          <div>
            <div class="flex justify-between text-xs text-gray-400 mb-1">
              <span>Disk</span>
              <span>{{ formatMemory(node.storage.usage) }} / {{ formatMemory(node.storage.capacity) }}</span>
            </div>
            <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="getProgressColor(getUsagePercent(node.storage.usage, node.storage.capacity))"
                :style="{ width: `${getUsagePercent(node.storage.usage, node.storage.capacity)}%` }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
