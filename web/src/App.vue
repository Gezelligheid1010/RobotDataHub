<template><div class="app"><aside class="sidebar"><h2>🤖 RobotDataHub</h2><p>Embodied AI Data Infra</p><router-link to="/datasets">数据集管理</router-link><router-link to="/jobs">质量任务中心</router-link><div class="token"><el-input v-model="token" type="password" show-password placeholder="API Token"/><el-button @click="saveToken">保存访问密钥</el-button></div></aside><main><div class="topbar"><span>RobotDataHub / 管理控制台</span><span class="pill">MVP</span></div><router-view/></main></div></template>
<script setup>
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
const token=ref(sessionStorage.getItem('robotdatahub-token')||'')
function saveToken(){sessionStorage.setItem('robotdatahub-token',token.value);ElMessage.success('已保存到当前标签页')}
</script>
