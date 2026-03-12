const API_BASE = import.meta.env.VITE_API_URL || '/api/v1'

async function fetchApi<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`)
  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`)
  }
  return response.json()
}

export interface Node {
  name: string
  status: string
  roles: string[]
  ip: string
  cpu: { capacity: string; allocatable: string }
  memory: { capacity: string; allocatable: string }
  labels: Record<string, string>
  createdAt: string
}

export interface Namespace {
  name: string
  status: string
  createdAt: string
}

export interface Pod {
  name: string
  namespace: string
  status: string
  node: string
  ip: string
  containers: { name: string; image: string; ready: boolean }[]
  labels: Record<string, string>
  createdAt: string
}

export interface Service {
  name: string
  namespace: string
  type: string
  clusterIP: string
  ports: { name: string; port: number; targetPort: string; protocol: string }[]
  labels: Record<string, string>
  createdAt: string
}

export const api = {
  getNodes: () => fetchApi<Node[]>('/cluster/nodes'),
  getNamespaces: () => fetchApi<Namespace[]>('/cluster/namespaces'),
  getPods: (ns: string) => fetchApi<Pod[]>(`/namespaces/${ns}/pods`),
  getServices: (ns: string) => fetchApi<Service[]>(`/namespaces/${ns}/services`),
}
