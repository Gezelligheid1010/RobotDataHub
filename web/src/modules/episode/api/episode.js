import {api} from '../../../shared/api/client'
export const listEpisodes=(id)=>api('get',`/datasets/${id}/episodes`)
