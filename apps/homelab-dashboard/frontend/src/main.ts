import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './assets/main.css'

import Overview from './views/Overview.vue'
import Nodes from './views/Nodes.vue'
import Namespaces from './views/Namespaces.vue'

const routes = [
  { path: '/', name: 'overview', component: Overview },
  { path: '/nodes', name: 'nodes', component: Nodes },
  { path: '/namespaces', name: 'namespaces', component: Namespaces },
  { path: '/namespaces/:ns', name: 'namespace-detail', component: () => import('./views/NamespaceDetail.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
