import React, {useEffect, useState} from 'react'
import { fetchWithAuth } from './api'

export default function QuizPlay({quizId, onSubmitted, onCancel}){
  const [quiz, setQuiz] = useState(null)
  const [answers, setAnswers] = useState({})
  const [submitted, setSubmitted] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(()=>{
    let cancelled = false
    fetch(`http://localhost:8080/api/quizzes/${quizId}`)
      .then(r=>r.json()).then(data=>{ if(!cancelled) { setQuiz(data); setLoading(false) } })
    return ()=>{ cancelled = true }
  },[quizId])

  if(loading) return <div>Loading...</div>
  if(!quiz) return <div>Quiz not found</div>

  const toggleAnswer = (qid, idx, multiple)=>{
    setAnswers(prev=>{
      const cur = prev[qid] || []
      if(multiple){
        const set = new Set(cur)
        if(set.has(idx)) set.delete(idx)
        else set.add(idx)
        return {...prev, [qid]: Array.from(set)}
      } else {
        return {...prev, [qid]: [idx]}
      }
    })
  }

  const submit = async ()=>{
    const payload = { answers: Object.keys(answers).map(k=>({ questionId: Number(k), selected: answers[k] })) }
    try{
      const res = await fetchWithAuth(`http://localhost:8080/api/quizzes/${quizId}/submit`, { method: 'POST', body: JSON.stringify(payload) })
      const data = await res.json()
      // show thank you message and provide button to go back
      setSubmitted({ id: data.submissionId, score: data.score })
      // also call onSubmitted so caller can navigate to results if desired
      if(onSubmitted) onSubmitted(data.submissionId)
    }catch(err){
      let msg = 'Submit failed'
      try{
        const j = await err.json()
        msg = j.error || JSON.stringify(j)
      }catch(e){ }
      alert(msg)
    }
  }

  return (
    <div style={{maxWidth:760}}>
      <h2>{quiz.title}</h2>
      {submitted && (
        <div style={{padding:12, border:'1px solid #4caf50', background:'#e8f5e9', marginBottom:12}}>
          <div style={{fontWeight:700}}>Thank you for answering!</div>
          <div style={{marginTop:8}}>
            <button onClick={onCancel} style={{background:'#1e66d0',color:'#fff',padding:'8px 12px',border:'none'}}>Go to Main Page</button>
          </div>
        </div>
      )}
      {quiz.questions.map(q=> (
        <div key={q.id} style={{border:'1px solid #eee', padding:12, marginBottom:12}}>
          <div style={{fontWeight:600}}>{q.text}</div>
          <div style={{marginTop:8}}>
            {q.options.map((op, idx)=>{
              const text = typeof op === 'string' ? op : (op.text || op.Text)
              const sel = (answers[q.id] || []).includes(idx)
              return (
                <label key={idx} style={{display:'block', marginBottom:6}}>
                  <input type={q.multiple? 'checkbox':'radio'} checked={sel} onChange={()=>toggleAnswer(q.id, idx, q.multiple)} /> {text}
                </label>
              )
            })}
          </div>
        </div>
      ))}
      <div>
        <button onClick={onCancel} style={{marginRight:8}}>Cancel</button>
        <button onClick={submit} style={{background:'#1e66d0',color:'#fff',padding:'8px 12px',border:'none'}}>Submit</button>
      </div>
    </div>
  )
}
