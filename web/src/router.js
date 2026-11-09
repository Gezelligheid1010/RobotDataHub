import {createRouter,createWebHistory} from 'vue-router'
import Dataset from './modules/dataset/views/Index.vue'
import Episode from './modules/episode/views/Index.vue'
import Job from './modules/job/views/Index.vue'
export default createRouter({history:createWebHistory(),routes:[{path:'/',redirect:'/datasets'},{path:'/datasets',component:Dataset},{path:'/datasets/:id/episodes',component:Episode},{path:'/jobs',component:Job}]})
