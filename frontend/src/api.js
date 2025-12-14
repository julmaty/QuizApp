export function getAuthToken(){
  return localStorage.getItem('qa_token')
}

export async function fetchWithAuth(url, opts={}){
  const token = getAuthToken()
  const headers = opts.headers || {}
  headers['Content-Type'] = headers['Content-Type'] || 'application/json'
  if(token) headers['Authorization'] = `Bearer ${token}`
  const res = await fetch(url, {...opts, headers})
  if(!res.ok) throw res
  return res
}
