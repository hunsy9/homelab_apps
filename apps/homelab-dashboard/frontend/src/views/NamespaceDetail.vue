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
  <div class="space-y-6">
    <div class="flex items-center gap-4">
      <RouterLink to="/namespaces" class="text-muted-foreground hover:text-foreground transition-colors">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </RouterLink>
      <div>
        <h1 class="text-2xl font-semibold text-foreground">{{ namespace }}</h1>
        <p class="text-muted-foreground text-sm">Namespace details</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 p-1 bg-secondary rounded-lg w-fit">
      <button
        @click="activeTab = 'pods'"
        class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
        :class="activeTab === 'pods' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'"
      >
        Pods ({{ store.pods.length }})
      </button>
      <button
        @click="activeTab = 'services'"
        class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
        :class="activeTab === 'services' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'"
      >
        Services ({{ store.services.length }})
      </button>
    </div>

    <div v-if="store.loading" class="text-muted-foreground">Loading...</div>
    <div v-else-if="store.error" class="text-red-500">{{ store.error }}</div>

    <!-- Pods Table -->
    <div v-else-if="activeTab === 'pods'" class="bg-card rounded-lg border overflow-hidden overflow-x-auto">
      <table class="w-full min-w-[600px]">
        <thead class="bg-muted/50">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Status</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Name</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Node</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">IP</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Containers</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-for="pod in store.pods" :key="pod.name" class="hover:bg-muted/30">
            <td class="px-4 py-3">
              <span
                class="px-2 py-0.5 rounded-full text-xs font-medium"
                :class="{
                  'bg-emerald-50 text-emerald-700': pod.status === 'Running',
                  'bg-amber-50 text-amber-700': pod.status === 'Pending',
                  'bg-red-50 text-red-700': pod.status === 'Failed',
                  'bg-secondary text-muted-foreground': !['Running', 'Pending', 'Failed'].includes(pod.status)
                }"
              >
                {{ pod.status }}
              </span>
            </td>
            <td class="px-4 py-3 font-medium text-foreground">{{ pod.name }}</td>
            <td class="px-4 py-3 text-muted-foreground">{{ pod.node }}</td>
            <td class="px-4 py-3 text-muted-foreground font-mono text-sm">{{ pod.ip || '-' }}</td>
            <td class="px-4 py-3">
              <div class="flex gap-1">
                <span
                  v-for="c in pod.containers"
                  :key="c.name"
                  class="w-2.5 h-2.5 rounded-full"
                  :class="c.ready ? 'bg-emerald-500' : 'bg-red-500'"
                  :title="c.name"
                ></span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Services Table -->
    <div v-else-if="activeTab === 'services'" class="bg-card rounded-lg border overflow-hidden overflow-x-auto">
      <table class="w-full min-w-[500px]">
        <thead class="bg-muted/50">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Name</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Type</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Cluster IP</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Ports</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-for="svc in store.services" :key="svc.name" class="hover:bg-muted/30">
            <td class="px-4 py-3 font-medium text-foreground">{{ svc.name }}</td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-secondary text-muted-foreground">
                {{ svc.type }}
              </span>
            </td>
            <td class="px-4 py-3 text-muted-foreground font-mono text-sm">{{ svc.clusterIP }}</td>
            <td class="px-4 py-3 text-muted-foreground">
              {{ svc.ports.map(p => `${p.port}/${p.protocol}`).join(', ') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
