import React, {useEffect, useState} from 'react'
import CreateQuiz from './CreateQuiz'
import QuizList from './QuizList'
import Login from './Login'
import Register from './Register'
import QuizPlay from './QuizPlay'
import Results from './Results'

export default function App(){
  const [quizzes, setQuizzes] = useState([])
  const [view, setView] = useState('list') // list | play | results | login | register
  const [selectedQuiz, setSelectedQuiz] = useState(null)
  const [submissionId, setSubmissionId] = useState(null)
  const [user, setUser] = useState(null)

  const fetchQuizzes = async ()=>{
    const res = await fetch('http://localhost:8080/api/quizzes')
    const data = await res.json()
    setQuizzes(data)
  }

  useEffect(()=>{ fetchQuizzes() }, [])

  const handlePlay = (quizId)=>{ setSelectedQuiz(quizId); setView('play') }
  const handleSubmitted = (sid)=>{ setSubmissionId(sid); setView('results') }
  const handleBack = ()=>{ setView('list'); setSelectedQuiz(null); setSubmissionId(null); fetchQuizzes() }
  const handleLogin = (u, token)=>{ setUser(u); setView('list') }
  const handleRegistered = (u)=>{ alert('Registered. Please log in'); setView('login') }

  return (
    <div style={{padding:20,fontFamily:'Arial'}}>
      <h1>Quiz App</h1>
      {!user && view !== 'login' && view !== 'register' && (
        <div style={{marginBottom:12}}>
          <button onClick={()=>setView('login')} style={{marginRight:8}}>Log In</button>
          <button onClick={()=>setView('register')}>Sign Up</button>
        </div>
      )}

      {view === 'login' && <Login onLogin={handleLogin} />}
      {view === 'register' && <Register onRegistered={handleRegistered} />}
      {view === 'list' && (
        <>
          <CreateQuiz onCreated={fetchQuizzes} />
          <hr />
          <QuizList quizzes={quizzes} onPlay={handlePlay} />
        </>
      )}
      {view === 'play' && selectedQuiz && (
        <QuizPlay quizId={selectedQuiz} onSubmitted={handleSubmitted} onCancel={handleBack} />
      )}
      {view === 'results' && selectedQuiz && submissionId && (
        <Results quizId={selectedQuiz} submissionId={submissionId} onBack={handleBack} />
      )}
    </div>
  )
}
