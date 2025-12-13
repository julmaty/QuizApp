import React, {useState} from 'react'

export default function CreateQuiz({onCreated}){
  const [title, setTitle] = useState('')
  const [question, setQuestion] = useState('')
  const [options, setOptions] = useState(['',''])

  const addOption = ()=> setOptions([...options,''])
  const setOption = (i,v)=> setOptions(options.map((o,idx)=> idx===i? v: o))

  const submit = async ()=>{
    const payload = {
      title,
      questions: [
        { text: question, options, multiple:false, answers: [] }
      ]
    }
    await fetch('http://localhost:8080/api/quizzes', {
      method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)
    })
    setTitle('')
    setQuestion('')
    setOptions(['',''])
    if(onCreated) onCreated()
  }

  return (
    <div>
      <h2>Create quiz</h2>
      <div>
        <label>Title</label><br />
        <input value={title} onChange={e=>setTitle(e.target.value)} />
      </div>
      <div style={{marginTop:8}}>
        <label>Question</label><br />
        <input value={question} onChange={e=>setQuestion(e.target.value)} />
      </div>
      <div style={{marginTop:8}}>
        <label>Options</label>
        {options.map((o,i)=> (
          <div key={i}>
            <input value={o} onChange={e=>setOption(i,e.target.value)} />
          </div>
        ))}
        <button onClick={addOption} style={{marginTop:6}}>Add option</button>
      </div>
      <div style={{marginTop:8}}>
        <button onClick={submit}>Create</button>
      </div>
    </div>
  )
}
