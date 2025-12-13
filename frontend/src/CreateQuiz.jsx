import React, {useState} from 'react'

export default function CreateQuiz({onCreated}){
  const [title, setTitle] = useState('')
  const [questions, setQuestions] = useState([])

  const addQuestion = () => {
    setQuestions([...questions, { text: '', options: ['', ''], multiple: false }])
  }

  const updateQuestion = (idx, field, value) => {
    setQuestions(questions.map((q, i) => i === idx ? { ...q, [field]: value } : q))
  }

  const updateOption = (qIdx, optIdx, value) => {
    setQuestions(questions.map((q, i) => 
      i === qIdx ? { ...q, options: q.options.map((o, j) => j === optIdx ? value : o) } : q
    ))
  }

  const addOption = (qIdx) => {
    setQuestions(questions.map((q, i) => 
      i === qIdx ? { ...q, options: [...q.options, ''] } : q
    ))
  }

  const removeQuestion = (idx) => {
    setQuestions(questions.filter((_, i) => i !== idx))
  }

  const submit = async ()=>{
    if (!title.trim()) {
      alert('Please enter a quiz title')
      return
    }
    if (questions.length === 0) {
      alert('Please add at least one question')
      return
    }
    const payload = {
      title,
      questions: questions.map(q => ({
        text: q.text,
        options: q.options.filter(o => o.trim() !== ''),
        multiple: q.multiple,
        answers: []
      }))
    }
    await fetch('http://localhost:8080/api/quizzes', {
      method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)
    })
    setTitle('')
    setQuestions([])
    if(onCreated) onCreated()
  }

  return (
    <div style={{border: '1px solid #ccc', padding: '15px', borderRadius: '5px'}}>
      <h2>Create quiz</h2>
      <div>
        <label>Quiz Title</label><br />
        <input value={title} onChange={e=>setTitle(e.target.value)} style={{width: '100%', padding: '5px', marginBottom: '15px'}} />
      </div>

      <div>
        <h3>Questions</h3>
        {questions.length === 0 ? (
          <p style={{color: '#999'}}>No questions added yet</p>
        ) : (
          questions.map((q, qIdx) => (
            <div key={qIdx} style={{border: '1px solid #ddd', padding: '10px', marginBottom: '10px', borderRadius: '3px', backgroundColor: '#f9f9f9'}}>
              <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px'}}>
                <label><strong>Question {qIdx + 1}</strong></label>
                <button onClick={() => removeQuestion(qIdx)} style={{backgroundColor: '#ff6b6b', color: 'white', border: 'none', padding: '5px 10px', borderRadius: '3px', cursor: 'pointer'}}>Remove</button>
              </div>
              <input 
                value={q.text} 
                onChange={e=>updateQuestion(qIdx, 'text', e.target.value)}
                placeholder="Enter question text"
                style={{width: '100%', padding: '5px', marginBottom: '8px'}}
              />
              
              <div style={{marginBottom: '8px'}}>
                <label style={{marginRight: '10px'}}>
                  <input 
                    type="checkbox" 
                    checked={q.multiple}
                    onChange={e=>updateQuestion(qIdx, 'multiple', e.target.checked)}
                  />
                  Allow multiple answers
                </label>
              </div>

              <label>Options:</label>
              {q.options.map((o, optIdx) => (
                <div key={optIdx}>
                  <input 
                    value={o} 
                    onChange={e=>updateOption(qIdx, optIdx, e.target.value)}
                    placeholder={`Option ${optIdx + 1}`}
                    style={{width: '100%', padding: '5px', marginBottom: '5px'}}
                  />
                </div>
              ))}
              <button onClick={() => addOption(qIdx)} style={{marginTop: '5px', padding: '5px 10px', backgroundColor: '#4CAF50', color: 'white', border: 'none', borderRadius: '3px', cursor: 'pointer'}}>Add option</button>
            </div>
          ))
        )}
        <button onClick={addQuestion} style={{marginTop: '10px', padding: '8px 15px', backgroundColor: '#2196F3', color: 'white', border: 'none', borderRadius: '3px', cursor: 'pointer'}}>Add question</button>
      </div>

      <div style={{marginTop: '15px'}}>
        <button onClick={submit} style={{padding: '10px 20px', backgroundColor: '#4CAF50', color: 'white', border: 'none', borderRadius: '3px', cursor: 'pointer', fontSize: '16px'}}>Create Quiz</button>
      </div>
    </div>
  )
}
