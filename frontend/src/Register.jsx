import React, {useState} from 'react'

export default function Register({onRegistered}){
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [err, setErr] = useState(null)

  const submit = async (e)=>{
    e.preventDefault()
    try{
      const res = await fetch('http://localhost:8080/api/register', {
        method: 'POST', headers: {'Content-Type':'application/json'},
        body: JSON.stringify({email, password, displayName})
      })
      if(!res.ok){ setErr('Registration failed'); return }
      const data = await res.json()
      onRegistered(data)
    }catch(err){ setErr('Network error') }
  }

  return (
    <div style={{maxWidth:360, padding:12, border:'1px solid #ddd'}}>
      <h3>Register</h3>
      <form onSubmit={submit}>
        <div><input placeholder="Email" value={email} onChange={e=>setEmail(e.target.value)} style={{width:'100%',padding:8,marginBottom:8}}/></div>
        <div><input placeholder="Display name" value={displayName} onChange={e=>setDisplayName(e.target.value)} style={{width:'100%',padding:8,marginBottom:8}}/></div>
        <div><input type="password" placeholder="Password" value={password} onChange={e=>setPassword(e.target.value)} style={{width:'100%',padding:8,marginBottom:8}}/></div>
        <div><button style={{background:'#1e66d0',color:'#fff',padding:'8px 12px',border:'none'}}>Sign Up</button></div>
        {err && <div style={{color:'red',marginTop:8}}>{err}</div>}
      </form>
    </div>
  )
}
