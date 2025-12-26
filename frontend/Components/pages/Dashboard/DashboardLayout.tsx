import React from 'react';
import { Outlet } from 'react-router-dom';
import Navbar from '../../Navbar/Navbar';
import Footer from '../../Footer/Footer';

const DashboardLayout = () => {
  return (
    <div className="min-h-screen bg-black">
      <Navbar />
      {/* <Outlet /> is where the child routes (Movies, Showtimes, Seats) will render */}
      <Outlet />
      <Footer />
    </div>
  );
};

export default DashboardLayout;