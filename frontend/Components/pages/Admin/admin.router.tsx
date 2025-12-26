import React from "react";
import AdminLayout from "./AdminLayout"; // The sidebar layout
import AdminMovies from "./AdminMovies"; // The movie management page
import AdminStudios from "./AdminStudios"; 
import AdminShowtimes from "./AdminShowtimes";
import { Navigate } from "react-router-dom";

export const adminRoutes = {
  path: '/admin',
  element: <AdminLayout />,
  children: [
    {
      index: true,
      element: <Navigate to="/admin/movies" replace /> // Redirect /admin to /admin/movies
    },
    {
      path: 'movies',
      element: <AdminMovies />
    },
    {
      path: 'studios',
      element: <AdminStudios />
    },
    {
      path: 'showtimes',
      element: <AdminShowtimes />
    }
  ]
};

export default adminRoutes;