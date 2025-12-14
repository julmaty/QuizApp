import React from 'react'

export default function QuizList({quizzes, onPlay}){
  if(!quizzes || quizzes.length===0) return <div>No quizzes yet</div>
  return (
    <div>
      <h2>Quizzes</h2>
      {quizzes.map(q=> (
        <div key={q.id} style={{border:'1px solid #ddd', padding:8, marginBottom:8}}>
          <strong>{q.title}</strong>
          <div style={{fontSize:12,color:'#666'}}>{new Date(q.createdAt).toLocaleString()}</div>
          <ul>
            {q.questions.map((qq, idx)=> (
              <li key={idx}>{qq.text}
                <ul>
                  {qq.options.map((op, i)=> (
                    <li key={i}>{typeof op === 'string' ? op : (op.text ?? op.Text ?? JSON.stringify(op))}</li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>
          <div style={{marginTop:8}}>
            <button onClick={()=> onPlay && onPlay(q.id)} style={{background:'#1e66d0',color:'#fff',border:'none',padding:'6px 10px'}}>Play</button>
          </div>
        </div>
      ))}
    </div>
  )
}
