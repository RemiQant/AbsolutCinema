import { useState } from 'react'
import { createBrowserRouter , RouterProvider , Navigate } from 'react-router-dom'
import { dashboardRoutes, successRoute } from '../Components/pages/Dashboard/dashboard.router'
import accountRoutes from '../Components/pages/Account/account.router'
import aboutRoutes from '../Components/pages/About/about.router'
import loginRoutes from '../Components/pages/Login/login.router'
import signUpRoutes from '../Components/pages/Signup/signup.router'
import adminRoutes from '../Components/pages/Admin/admin.router'



import './App.css'

const routes = createBrowserRouter([
  accountRoutes,
  aboutRoutes,
  dashboardRoutes,
  signUpRoutes,
  adminRoutes,
  successRoute,
  ...loginRoutes
])

function App() {
 
  return (
    <>
      <RouterProvider router={routes} />
    </>
  )
}

export default App
