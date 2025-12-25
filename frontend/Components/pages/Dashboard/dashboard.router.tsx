import React from "react";
import DashboardLayout from "./DashboardLayout"; // Make sure you renamed the file!
import MovieGrid from "./MovieGrid";
import ShowtimeSelection from "./ShowtimeSelection";
import SeatSelection from "./SeatSelection";
import PaymentSuccess from "./PaymentSuccess";

export const dashboardRoutes = {
  path: '/dashboard',
  element: <DashboardLayout />,
  children: [
    {
      index: true, // This means /dashboard shows this
      element: <MovieGrid />
    },
    {
      path: 'movie/:movieId', // /dashboard/movie/1
      element: <ShowtimeSelection />
    },
    {
      path: 'booking/:showtimeId', // /dashboard/booking/55
      element: <SeatSelection />
    },
    {
      path: 'payment-success',
      element: <PaymentSuccess />
    }
  ]
};

export default dashboardRoutes;