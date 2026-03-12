import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api, type Node, type Namespace, type Pod, type Service } from '../api/client'

export const useClusterStore = defineStore('cluster', () => {
  const nodes = ref<Node[]>([])
  const namespaces = ref<Namespace[]>([])
  const pods = ref<Pod[]>([])
  const services = ref<Service[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchNodes() {
    // Only show loading on initial fetch
    if (nodes.value.length === 0) {
      loading.value = true
    }
    error.value = null
    try {
      nodes.value = await api.getNodes()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch nodes'
    } finally {
      loading.value = false
    }
  }

  async function fetchNamespaces() {
    loading.value = true
    error.value = null
    try {
      namespaces.value = await api.getNamespaces()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch namespaces'
    } finally {
      loading.value = false
    }
  }

  async function fetchPods(namespace: string) {
    loading.value = true
    error.value = null
    try {
      pods.value = await api.getPods(namespace)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch pods'
    } finally {
      loading.value = false
    }
  }

  async function fetchServices(namespace: string) {
    loading.value = true
    error.value = null
    try {
      services.value = await api.getServices(namespace)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch services'
    } finally {
      loading.value = false
    }
  }

  return {
    nodes,
    namespaces,
    pods,
    services,
    loading,
    error,
    fetchNodes,
    fetchNamespaces,
    fetchPods,
    fetchServices,
  }
})
