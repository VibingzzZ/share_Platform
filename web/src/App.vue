<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import CustomizeDrawer from './components/CustomizeDrawer.vue'
import { api, type Dashboard } from './api/client'
import { defaultLayout, moveModule, type ModuleKey } from './lib/layout'
const fallback: Dashboard = { resources: [{ title: '产品需求模板', type: 'template', description: '跨职能协作的需求模板。' }, { title: '服务发布手册', type: 'document', description: '从检查到回滚的标准流程。' }], posts: [{ title: '从零搭建组织内容空间', summary: '内容沉淀应该服务于每天的协作节奏。', publishedAt: '2026-08-28' }], aiProducts: [{ name: 'Dify 工作流实验室', summary: '将内部知识和模型能力组合成可复用工作流。', url: '#', tags: ['Dify'] }] }
const dashboard = ref<Dashboard>(fallback); const layout = ref(structuredClone(defaultLayout)); const search = ref(''); const drawer = ref(false)
const visible = (key: ModuleKey) => !layout.value.hidden.includes(key); const contains = (text: string) => text.toLowerCase().includes(search.value.toLowerCase())
function toggle(key: ModuleKey) { layout.value.hidden = layout.value.hidden.includes(key) ? layout.value.hidden.filter(item => item !== key) : [...layout.value.hidden, key] }
async function save() { localStorage.setItem('share-layout', JSON.stringify(layout.value)); try { await api.saveLayout(layout.value) } catch {} drawer.value = false }
watch(layout, value => { document.body.className = `${value.theme} ${value.density}` }, { deep: true })
onMounted(async () => { try { dashboard.value = await api.dashboard() } catch {} try { layout.value = await api.layout() } catch { layout.value = JSON.parse(localStorage.getItem('share-layout') || 'null') || layout.value } })
</script>
<template>
  <header class="topbar"><a class="brand">SHARE <span>Platform</span></a><nav><a href="#resources">资源</a><a href="#posts">开发日志</a><a href="#ai-lab">AI 实验室</a></nav><div class="top-actions"><input v-model="search" type="search" placeholder="搜索资源与文章"><button class="action" @click="drawer = true">调整布局</button></div></header>
  <main class="shell"><section class="intro"><div><p class="eyebrow">INTERNAL KNOWLEDGE HUB</p><h1>组织的知识，在这里持续生长。</h1><p>资源、开发过程与 AI 实验统一沉淀，按自己的工作习惯组织工作台。</p></div><div class="pulse"><b>{{ dashboard.resources.length }}</b><span>已沉淀资源</span><b>{{ dashboard.posts.length }}</b><span>开发记录</span></div></section>
    <section v-if="visible('overview')" class="module overview"><div><p class="section-label">本周概览</p><h2>保持同步，快速进入工作。</h2></div><div class="quick"><a href="#resources">浏览资源库 →</a><a href="#posts">查看最近发布 →</a><a href="#ai-lab">探索 AI 产品 →</a></div></section>
    <section v-if="visible('resources')" id="resources" class="module"><header class="section-head"><div><p class="section-label">资源库</p><h2>最近加入</h2></div></header><div class="resource-grid"><article v-for="item in dashboard.resources.filter(x => contains(x.title + x.description))" :key="item.title" class="card"><span class="tag">{{ item.type }}</span><h3>{{ item.title }}</h3><p>{{ item.description }}</p></article></div></section>
    <section v-if="visible('posts')" id="posts" class="module split"><div><header class="section-head"><div><p class="section-label">开发日志</p><h2>正在发生的事</h2></div></header><article v-for="item in dashboard.posts.filter(x => contains(x.title + x.summary))" :key="item.title" class="post"><small>{{ item.publishedAt?.slice(0, 10) }}</small><h3>{{ item.title }}</h3><p>{{ item.summary }}</p></article></div><aside class="notice"><p class="section-label">贡献入口</p><h3>有值得分享的内容？</h3><p>把开发经验、可复用资源和有趣实验沉淀到组织空间。</p><button>提交内容</button></aside></section>
    <section v-if="visible('ai-lab')" id="ai-lab" class="module"><header class="section-head"><div><p class="section-label">AI 实验室</p><h2>组织正在试用的产品</h2></div></header><div class="ai-grid"><article v-for="item in dashboard.aiProducts.filter(x => contains(x.name + x.summary))" :key="item.name" class="card"><span class="tag">{{ item.tags?.[0] || 'AI' }}</span><h3>{{ item.name }}</h3><p>{{ item.summary }}</p><a :href="item.url">打开产品 →</a></article></div></section>
  </main><CustomizeDrawer :open="drawer" :layout="layout" @close="drawer = false" @save="save" @toggle="toggle" @move="(key, offset) => layout = moveModule(layout, key, offset)" />
</template>
