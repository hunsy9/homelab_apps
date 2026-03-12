<script setup lang="ts">
import { onMounted } from 'vue'
import { useClusterStore } from '../stores/cluster'

const store = useClusterStore()

onMounted(() => {
  store.fetchNodes()
})
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
            <td class="px-4 py-3 text-gray-400">{{ node.cpu.allocatable }}</td>
            <td class="px-4 py-3 text-gray-400">{{ node.memory.allocatable }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
