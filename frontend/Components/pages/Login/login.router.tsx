import { Navigate } from "react-router-dom";
import Login from "./Login";

const loginRoutes = [
  {
    path: '/',
    element: <Login />
  },
  {
    path: '/login',
    element: <Navigate to="/" replace />
  }
];

export default loginRoutes;
