import React from "react";
import { Navigate, Outlet } from "react-router-dom";
import Cookies from "js-cookie";
import AdminLayout from "./AdminLayout";
import AdminMovies from "./AdminMovies";
import AdminStudios from "./AdminStudios";
import AdminShowtimes from "./AdminShowtimes";

/**
 * Gatekeeper: Menjaga agar Customer tidak bisa masuk ke /admin
 * Kita mengecek role dari cookie.
 */
const AdminGuard = () => {
  // Ambil role dari cookie bernama 'role'
  // Pastikan saat login, backend/frontend kamu menset cookie ini
  const role = Cookies.get('role');

  // Jika role bukan admin, tendang ke dashboard
  if (role !== 'admin') {
    console.warn("Akses ditolak: Anda bukan admin.");
    return <Navigate to="/dashboard" replace />;
  }

  // Jika admin, izinkan akses ke semua children (AdminLayout & Isinya)
  return <Outlet />;
};

export const adminRoutes = {
  path: '/admin',
  element: <AdminGuard />, // Pelindung utama di level paling atas
  children: [
    {
      element: <AdminLayout />, // Berisi Sidebar dan Header Admin
      children: [
        {
          index: true,
          element: <Navigate to="movies" replace />
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
    }
  ]
};

export default adminRoutes;