<template><div class="titlebar"><div><h1>质量任务中心</h1><p class="sub">Python Worker 校验任务与执行结果</p></div><el-button @click="refresh">刷新</el-button></div><div class="panel"><el-table :data="items" v-loading="loading" stripe><el-table-column prop="id" label="Job ID" min-width="210"/><el-table-column prop="dataset_id" label="Dataset ID" min-width="200"/><el-table-column label="状态" width="120"><template #default="{row}"><el-tag :type="tag(row.state)">{{row.state}}</el-tag></template></el-table-column><el-table-column label="尝试次数" width="120"><template #default="{row}">{{row.attempts}} / {{row.max_attempts}}</template></el-table-column><el-table-column label="结果" width="140"><template #default="{row}"><el-button link @click="show(row)">查看报告</el-button></template></el-table-column></el-table><el-empty v-if="!items.length&&!loading" description="暂无质量检查任务"/></div><el-dialog v-model="detail" title="质量校验报告" width="680"><pre style="white-space:pre-wrap;max-height:55vh;overflow:auto">{{selected}}</pre></el-dialog></template>
<script setup>
import {ref,onMounted} from 'vue';import {ElMessage} from 'element-plus';import {listJobs} from '../api/job'
const items=ref([]),loading=ref(false),detail=ref(false),selected=ref('')
const tag=s=>({succeeded:'success',failed:'danger',running:'warning',pending:'info'}[s]||'info')
async function refresh(){loading.value=true;try{items.value=await listJobs()}catch(e){ElMessage.error(e.response?.data?.error||e.message)}finally{loading.value=false}}
function show(row){selected.value=JSON.stringify({error:row.error,result:row.result},null,2);detail.value=true}
onMounted(refresh)
</script>
