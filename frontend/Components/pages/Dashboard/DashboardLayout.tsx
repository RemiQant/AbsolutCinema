import React from 'react';
import { Outlet } from 'react-router-dom';
import Navbar from '../../Navbar/Navbar';

const DashboardLayout = () => {
  return (
    <div className="min-h-screen bg-black">
      <Navbar />
      {/* <Outlet /> is where the child routes (Movies, Showtimes, Seats) will render */}
      <Outlet />
    </div>
  );
};

export default DashboardLayout;