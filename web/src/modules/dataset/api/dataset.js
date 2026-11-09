import {api} from '../../../shared/api/client'
export const listDatasets=()=>api('get','/datasets')
export const createDataset=(data)=>api('post','/datasets',data)
export const deleteDataset=(id)=>api('delete',`/datasets/${id}`)
export const uploadFile=(id,file,path)=>{const data=new FormData();data.append('file',file);data.append('path',path||file.name);return api('post',`/datasets/${id}/files`,data)}
export const importManifest=(id,body)=>api('post',`/datasets/${id}/import-manifest`,body)
export const fileUrl=(id,path)=>api('get',`/datasets/${id}/file-url`,undefined,{params:{path}})
