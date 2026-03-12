<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useClusterStore } from '../stores/cluster'

const route = useRoute()
const store = useClusterStore()
const activeTab = ref<'pods' | 'services'>('pods')

const namespace = ref(route.params.ns as string)

watch(() => route.params.ns, (ns) => {
  namespace.value = ns as string
  loadData()
})

function loadData() {
  store.fetchPods(namespace.value)
  store.fetchServices(namespace.value)
}

onMounted(loadData)
</script>

<template>
  <div>
    <div class="flex items-center gap-4 mb-6">
      <RouterLink to="/namespaces" class="text-gray-400 hover:text-white">
        ← Back
      </RouterLink>
      <h2 class="text-2xl font-bold">{{ namespace }}</h2>
    </div>

    <!-- Tabs -->
    <div class="flex gap-4 mb-6">
      <button
        @click="activeTab = 'pods'"
        class="px-4 py-2 rounded"
        :class="activeTab === 'pods' ? 'bg-blue-600' : 'bg-gray-700 hover:bg-gray-600'"
      >
        Pods ({{ store.pods.length }})
      </button>
      <button
        @click="activeTab = 'services'"
        class="px-4 py-2 rounded"
        :class="activeTab === 'services' ? 'bg-blue-600' : 'bg-gray-700 hover:bg-gray-600'"
      >
        Services ({{ store.services.length }})
      </button>
    </div>

    <div v-if="store.loading" class="text-gray-400">Loading...</div>
    <div v-else-if="store.error" class="text-red-400">{{ store.error }}</div>

    <!-- Pods Table -->
    <div v-else-if="activeTab === 'pods'" class="bg-gray-800 rounded-lg overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-700">
          <tr>
            <th class="px-4 py-3 text-left">Status</th>
            <th class="px-4 py-3 text-left">Name</th>
            <th class="px-4 py-3 text-left">Node</th>
            <th class="px-4 py-3 text-left">IP</th>
            <th class="px-4 py-3 text-left">Containers</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="pod in store.pods"
            :key="pod.name"
            class="border-t border-gray-700"
          >
            <td class="px-4 py-3">
              <span
                class="px-2 py-1 rounded text-sm"
                :class="{
                  'bg-green-900 text-green-300': pod.status === 'Running',
                  'bg-yellow-900 text-yellow-300': pod.status === 'Pending',
                  'bg-red-900 text-red-300': pod.status === 'Failed',
                  'bg-gray-700 text-gray-300': !['Running', 'Pending', 'Failed'].includes(pod.status)
                }"
              >
                {{ pod.status }}
              </span>
            </td>
            <td class="px-4 py-3 font-medium">{{ pod.name }}</td>
            <td class="px-4 py-3 text-gray-400">{{ pod.node }}</td>
            <td class="px-4 py-3 text-gray-400">{{ pod.ip || '-' }}</td>
            <td class="px-4 py-3">
              <div class="flex gap-1">
                <span
                  v-for="c in pod.containers"
                  :key="c.name"
                  class="w-3 h-3 rounded-full"
                  :class="c.ready ? 'bg-green-500' : 'bg-red-500'"
                  :title="c.name"
                ></span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Services Table -->
    <div v-else-if="activeTab === 'services'" class="bg-gray-800 rounded-lg overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-700">
          <tr>
            <th class="px-4 py-3 text-left">Name</th>
            <th class="px-4 py-3 text-left">Type</th>
            <th class="px-4 py-3 text-left">Cluster IP</th>
            <th class="px-4 py-3 text-left">Ports</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="svc in store.services"
            :key="svc.name"
            class="border-t border-gray-700"
          >
            <td class="px-4 py-3 font-medium">{{ svc.name }}</td>
            <td class="px-4 py-3">
              <span class="px-2 py-1 rounded text-sm bg-gray-700">
                {{ svc.type }}
              </span>
            </td>
            <td class="px-4 py-3 text-gray-400">{{ svc.clusterIP }}</td>
            <td class="px-4 py-3 text-gray-400">
              {{ svc.ports.map(p => `${p.port}/${p.protocol}`).join(', ') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
