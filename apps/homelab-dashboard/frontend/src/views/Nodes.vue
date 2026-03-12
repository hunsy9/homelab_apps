<script setup lang="ts">
import { onMounted } from 'vue'
import { useClusterStore } from '../stores/cluster'

const store = useClusterStore()

onMounted(() => {
  store.fetchNodes()
})

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
    <h2 class="text-2xl font-bold mb-6">Nodes</h2>

    <div v-if="store.loading" class="text-gray-400">Loading...</div>
    <div v-else-if="store.error" class="text-red-400">{{ store.error }}</div>
    <div v-else class="bg-gray-800 rounded-lg overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-700">
          <tr>
            <th class="px-4 py-3 text-left">Status</th>
            <th class="px-4 py-3 text-left">Name</th>
            <th class="px-4 py-3 text-left">Roles</th>
            <th class="px-4 py-3 text-left">IP</th>
            <th class="px-4 py-3 text-left">CPU</th>
            <th class="px-4 py-3 text-left">Memory</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="node in store.nodes"
            :key="node.name"
            class="border-t border-gray-700 hover:bg-gray-750"
          >
            <td class="px-4 py-3">
              <span
                class="px-2 py-1 rounded text-sm"
                :class="node.status === 'Ready' ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300'"
              >
                {{ node.status }}
              </span>
            </td>
            <td class="px-4 py-3 font-medium">{{ node.name }}</td>
            <td class="px-4 py-3 text-gray-400">{{ node.roles.join(', ') }}</td>
            <td class="px-4 py-3 text-gray-400">{{ node.ip }}</td>
            <td class="px-4 py-3">
              <div class="w-32">
                <div class="flex justify-between text-xs text-gray-400 mb-1">
                  <span>{{ formatCPU(node.cpu.usage) }}</span>
                  <span>{{ formatCPU(node.cpu.allocatable) }}</span>
                </div>
                <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all"
                    :class="getProgressColor(getUsagePercent(node.cpu.usage, node.cpu.allocatable))"
                    :style="{ width: `${getUsagePercent(node.cpu.usage, node.cpu.allocatable)}%` }"
                  ></div>
                </div>
                <div class="text-xs text-gray-500 mt-1 text-center">
                  {{ getUsagePercent(node.cpu.usage, node.cpu.allocatable).toFixed(0) }}%
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <div class="w-32">
                <div class="flex justify-between text-xs text-gray-400 mb-1">
                  <span>{{ formatMemory(node.memory.usage) }}</span>
                  <span>{{ formatMemory(node.memory.allocatable) }}</span>
                </div>
                <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all"
                    :class="getProgressColor(getUsagePercent(node.memory.usage, node.memory.allocatable))"
                    :style="{ width: `${getUsagePercent(node.memory.usage, node.memory.allocatable)}%` }"
                  ></div>
                </div>
                <div class="text-xs text-gray-500 mt-1 text-center">
                  {{ getUsagePercent(node.memory.usage, node.memory.allocatable).toFixed(0) }}%
                </div>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
