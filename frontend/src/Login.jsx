import React, {useState} from 'react'
import { fetchWithAuth } from './api'

export default function Login({onLogin}){
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [err, setErr] = useState(null)

  const submit = async (e)=>{
    e.preventDefault()
    try{
      const res = await fetch('http://localhost:8080/api/login', {
        method: 'POST', headers: {'Content-Type':'application/json'},
        body: JSON.stringify({email, password})
      })
      if(!res.ok){ setErr('Login failed'); return }
      const data = await res.json()
      localStorage.setItem('qa_token', data.token)
      onLogin(data.user, data.token)
    }catch(err){ setErr('Network error') }
  }

  return (
    <div style={{maxWidth:360, padding:12, border:'1px solid #ddd'}}>
      <h3>Log In</h3>
      <form onSubmit={submit}>
        <div><input placeholder="Email" value={email} onChange={e=>setEmail(e.target.value)} style={{width:'100%',padding:8,marginBottom:8}}/></div>
        <div><input type="password" placeholder="Password" value={password} onChange={e=>setPassword(e.target.value)} style={{width:'100%',padding:8,marginBottom:8}}/></div>
        <div><button style={{background:'#1e66d0',color:'#fff',padding:'8px 12px',border:'none'}}>Log In</button></div>
        {err && <div style={{color:'red',marginTop:8}}>{err}</div>}
      </form>
    </div>
  )
}
