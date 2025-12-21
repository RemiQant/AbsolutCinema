import { useState } from 'react'
import reactLogo from './assets/react.svg'
import { createBrowserRouter , RouterProvider , Navigate } from 'react-router-dom'
import { dashboardRoutes } from '../Components/pages/Dashboard/dashboard.router'
import accountRoutes from '../Components/pages/Account/account.router'
import aboutRoutes from '../Components/pages/About/about.router'
import loginRoutes from '../Components/pages/Login/login.router'
import signUpRoutes from '../Components/pages/Signup/signup.router'

import './App.css'
import { sign } from 'crypto'

function App() {
 
  const routes = createBrowserRouter([
    accountRoutes,
    aboutRoutes,
    dashboardRoutes,
    signUpRoutes,
    loginRoutes
  ])

  return (
    <>
      <RouterProvider router={routes} />

    </>
  )
}

export default App
