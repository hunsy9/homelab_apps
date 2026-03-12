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
    color: 'from-orange-500 to-red-500'
  },
  {
    name: 'Traefik',
    url: 'https://traefik.seung.site/dashboard/',
    description: 'Ingress Controller Dashboard',
    color: 'from-cyan-500 to-blue-500'
  },
  {
    name: 'Longhorn',
    url: 'https://longhorn.seung.site',
    description: 'Distributed Storage',
    color: 'from-red-500 to-orange-500'
  },
  {
    name: 'Grafana',
    url: 'https://grafana.seung.site',
    description: 'Monitoring & Dashboards',
    color: 'from-orange-500 to-yellow-500'
  },
  {
    name: 'Harbor',
    url: 'https://harbor.seung.site',
    description: 'Container Registry',
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
            <div class="w-8 h-8 mb-2">
              <!-- ArgoCD -->
              <svg v-if="service.name === 'ArgoCD'" viewBox="0 0 24 24" fill="currentColor" class="text-orange-500">
                <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.6 0 12 0zm0 2.2c5.4 0 9.8 4.4 9.8 9.8s-4.4 9.8-9.8 9.8S2.2 17.4 2.2 12 6.6 2.2 12 2.2zm-1.3 4.5L7 12l3.7 5.3h2.6L10.2 12l3.1-5.3h-2.6zm4 0L11 12l3.7 5.3h2.6L14.2 12l3.1-5.3h-2.6z"/>
              </svg>
              <!-- Traefik -->
              <svg v-else-if="service.name === 'Traefik'" viewBox="0 0 24 24" fill="currentColor" class="text-cyan-400">
                <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.6 0 12 0zm-.5 3.5l7 5.5-7 5.5v-3.3H6V8.8h5.5V3.5zm1 10l7 5.5-7 5.5v-3.3H6v-4.4h6.5V13.5z"/>
              </svg>
              <!-- Longhorn -->
              <svg v-else-if="service.name === 'Longhorn'" viewBox="0 0 24 24" fill="currentColor" class="text-red-500">
                <path d="M21 3H3v18h18V3zM10 17H7v-4h3v4zm0-6H7V7h3v4zm7 6h-5v-4h5v4zm0-6h-5V7h5v4z"/>
              </svg>
              <!-- Grafana -->
              <svg v-else-if="service.name === 'Grafana'" viewBox="0 0 24 24" fill="currentColor" class="text-orange-400">
                <path d="M22.6 10.9c-.1-.5-.2-1.1-.4-1.6-.1-.3-.4-.5-.7-.5h-1.3c-.2-.5-.4-.9-.7-1.4l.9-.9c.2-.2.3-.6.1-.9-.3-.5-.7-.9-1.1-1.3-.2-.2-.3-.3-.5-.4-.3-.2-.6-.1-.9.1l-.9.9c-.4-.3-.9-.5-1.4-.7V2.9c0-.3-.2-.6-.5-.7-.5-.2-1.1-.3-1.6-.4-.3 0-.6.1-.8.4L12 3.4l-.8-1.2c-.2-.3-.5-.4-.8-.4-.5.1-1.1.2-1.6.4-.3.1-.5.4-.5.7v1.3c-.5.2-.9.4-1.4.7l-.9-.9c-.2-.2-.6-.3-.9-.1-.5.3-.9.7-1.3 1.1-.2.2-.3.3-.4.5-.2.3-.1.6.1.9l.9.9c-.3.4-.5.9-.7 1.4H2.4c-.3 0-.6.2-.7.5-.2.5-.3 1.1-.4 1.6 0 .3.1.6.4.8l1.2.8-1.2.8c-.3.2-.4.5-.4.8.1.5.2 1.1.4 1.6.1.3.4.5.7.5h1.3c.2.5.4.9.7 1.4l-.9.9c-.2.2-.3.6-.1.9.3.5.7.9 1.1 1.3.2.2.3.3.5.4.3.2.6.1.9-.1l.9-.9c.4.3.9.5 1.4.7v1.3c0 .3.2.6.5.7.5.2 1.1.3 1.6.4.3 0 .6-.1.8-.4l.8-1.2.8 1.2c.2.3.5.4.8.4.5-.1 1.1-.2 1.6-.4.3-.1.5-.4.5-.7v-1.3c.5-.2.9-.4 1.4-.7l.9.9c.2.2.6.3.9.1.5-.3.9-.7 1.3-1.1.2-.2.3-.3.4-.5.2-.3.1-.6-.1-.9l-.9-.9c.3-.4.5-.9.7-1.4h1.3c.3 0 .6-.2.7-.5.2-.5.3-1.1.4-1.6 0-.3-.1-.6-.4-.8l-1.2-.8 1.2-.8c.3-.2.4-.5.4-.8zM12 16c-2.2 0-4-1.8-4-4s1.8-4 4-4 4 1.8 4 4-1.8 4-4 4z"/>
              </svg>
              <!-- Harbor -->
              <svg v-else-if="service.name === 'Harbor'" viewBox="0 0 24 24" fill="currentColor" class="text-blue-400">
                <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2zm0 18c-4.4 0-8-3.6-8-8s3.6-8 8-8 8 3.6 8 8-3.6 8-8 8zm-1-13h2v6h-2V7zm0 8h2v2h-2v-2z"/>
                <path d="M8 10h2v5H8v-5zm6 0h2v5h-2v-5z"/>
              </svg>
            </div>
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
