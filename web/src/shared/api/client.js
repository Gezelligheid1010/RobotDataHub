import axios from 'axios'
export const client=axios.create({baseURL:'/api/v1',timeout:30000})
client.interceptors.request.use(config=>{config.headers['X-API-Key']=sessionStorage.getItem('robotdatahub-token')||'';return config})
export async function api(method,url,data,opts={}){return (await client({method,url,data,...opts})).data.data}
