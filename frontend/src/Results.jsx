import React, {useEffect, useState} from 'react'

export default function Results({quizId, submissionId, onBack}){
  const [data, setData] = useState(null)
  useEffect(()=>{
    let cancelled = false
    fetch(`http://localhost:8080/api/quizzes/${quizId}/results/${submissionId}`)
      .then(r=>r.json()).then(d=>{ if(!cancelled) setData(d) })
    return ()=>{ cancelled = true }
  },[quizId, submissionId])

  if(!data) return <div>Loading results...</div>
  return (
    <div style={{maxWidth:760}}>
      <h2>Results</h2>
      <div style={{fontSize:22,fontWeight:700}}>Score: {data.score}</div>
      <div style={{marginTop:12}}>
        {Array.isArray(data.perQuestion) && data.perQuestion.map((pq,i)=> (
          <div key={i} style={{border:'1px solid #eee', padding:10, marginBottom:8}}>
            <div><strong>Question ID {pq.questionId}</strong></div>
            <div>Selected: {Array.isArray(pq.selected) ? pq.selected.join(', ') : ''}</div>
            <div>Correct: {Array.isArray(pq.correct) ? pq.correct.join(', ') : ''}</div>
            <div style={{color: pq.correctBool ? 'green' : 'red'}}>{pq.correctBool ? 'Correct' : 'Incorrect'}</div>
          </div>
        ))}
      </div>
      <div><button onClick={onBack}>Back</button></div>
    </div>
  )
}
