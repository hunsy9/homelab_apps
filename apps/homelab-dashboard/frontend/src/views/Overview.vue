<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useClusterStore } from '../stores/cluster'

const store = useClusterStore()

onMounted(() => {
  store.fetchNodes()
  store.fetchNamespaces()
})

const readyNodes = computed(() => store.nodes.filter(n => n.status === 'Ready').length)
const totalNodes = computed(() => store.nodes.length)

const services = [
  {
    name: 'Traefik',
    url: 'https://traefik.seung.site',
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
  },
  {
    name: 'Hello World',
    url: 'https://hello.seung.site',
    description: 'Test Application',
    icon: '👋',
    color: 'from-green-500 to-teal-500'
  }
]
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

    <!-- Node list preview -->
    <div class="bg-gray-800 rounded-lg p-6">
      <h3 class="text-lg font-semibold mb-4">Nodes</h3>
      <div v-if="store.loading" class="text-gray-400">Loading...</div>
      <div v-else-if="store.error" class="text-red-400">{{ store.error }}</div>
      <div v-else class="space-y-2">
        <div
          v-for="node in store.nodes"
          :key="node.name"
          class="flex items-center justify-between p-3 bg-gray-700 rounded"
        >
          <div class="flex items-center gap-3">
            <span
              class="w-3 h-3 rounded-full"
              :class="node.status === 'Ready' ? 'bg-green-500' : 'bg-red-500'"
            ></span>
            <span class="font-medium">{{ node.name }}</span>
            <span class="text-gray-400 text-sm">{{ node.roles.join(', ') }}</span>
          </div>
          <span class="text-gray-400">{{ node.ip }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
