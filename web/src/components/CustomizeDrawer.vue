<script setup lang="ts">
import type { Layout, ModuleKey } from '../lib/layout'
import { moduleLabels } from '../lib/layout'
defineProps<{ open: boolean; layout: Layout }>()
const emit = defineEmits<{ close: []; save: []; move: [ModuleKey, number]; toggle: [ModuleKey] }>()
</script>
<template>
  <dialog :open="open">
    <div class="drawer-head"><div><p class="section-label">PAGE DIY</p><h2>调整你的工作台</h2></div><button @click="emit('close')">×</button></div>
    <div class="settings"><div v-for="(key, index) in layout.order" :key="key" class="setting"><span>{{ moduleLabels[key] }}</span><span><button :disabled="index === 0" @click="emit('move', key, -1)">↑</button><button :disabled="index === layout.order.length - 1" @click="emit('move', key, 1)">↓</button><input type="checkbox" :checked="!layout.hidden.includes(key)" @change="emit('toggle', key)"></span></div></div>
    <label>内容密度<select v-model="layout.density"><option value="comfortable">舒适</option><option value="compact">紧凑</option></select></label>
    <label>界面主题<select v-model="layout.theme"><option value="light">浅色</option><option value="dark">深色</option></select></label>
    <menu><button class="save" @click="emit('save')">保存布局</button></menu>
  </dialog>
</template>
