import React, {useEffect, useState} from 'react'
import CreateQuiz from './CreateQuiz'
import QuizList from './QuizList'

export default function App(){
  const [quizzes, setQuizzes] = useState([])

  const fetchQuizzes = async ()=>{
    const res = await fetch('http://localhost:8080/api/quizzes')
    const data = await res.json()
    setQuizzes(data)
  }

  useEffect(()=>{ fetchQuizzes() }, [])

  return (
    <div style={{padding:20,fontFamily:'Arial'}}>
      <h1>Quiz App</h1>
      <CreateQuiz onCreated={fetchQuizzes} />
      <hr />
      <QuizList quizzes={quizzes} />
    </div>
  )
}
