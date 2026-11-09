import {api} from '../../../shared/api/client'
export const listJobs=()=>api('get','/jobs')
export const createJob=(dataset_id)=>api('post','/jobs',{dataset_id})
