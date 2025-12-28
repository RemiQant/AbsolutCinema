import React from "react";
import DashboardLayout from "./DashboardLayout";
import MovieGrid from "./MovieGrid";
import ShowtimeSelection from "./ShowtimeSelection";
import SeatSelection from "./SeatSelection";
import PaymentSuccess from "./PaymentSuccess";


// 1. Rute Dashboard (Tetap seperti biasa)
export const dashboardRoutes = {
  path: '/dashboard',
  element: <DashboardLayout />,
  children: [
    {
      index: true,
      element: <MovieGrid />
    },
    {
      path: 'movie/:movieId',
      element: <ShowtimeSelection />
    },
    {
      path: 'booking/:showtimeId',
      element: <SeatSelection />
    }
  ]
};

export const successRoute = {
  path: 'booking/success', 
  element: <PaymentSuccess />
};

