import React, { useState, useEffect } from 'react';
import { ArrowLeft } from 'lucide-react';
import { useParams, useNavigate } from 'react-router-dom';
import { showtimesData, movies, bookedTickets as bookedTicketsDB } from './mockData';
import { type Seat } from './types';

const SeatSelection: React.FC = () => {
  const { showtimeId } = useParams();
  const navigate = useNavigate();
  const [seats, setSeats] = useState<Seat[]>([]);

  // 1. Find the Showtime
  // In a real app, you'd iterate all showtimes or fetch by ID. 
  // Since our data structure is weird (grouped by movie), we flat search:
  const allShowtimes = Object.values(showtimesData).flat();
  const showtime = allShowtimes.find(s => s.id === Number(showtimeId));
  
  // 2. Find the Movie (for the title)
  const movie = showtime ? movies.find(m => m.id === showtime.movie_id) : null;
  
  // 3. Find Booked Tickets
  const bookedTickets = showtime ? (bookedTicketsDB[showtime.id] || []) : [];

  useEffect(() => {
    if (!showtime) return;

    const rows = Array.from({ length: showtime.studio.total_rows }, (_, i) => String.fromCharCode(65 + i));
    const booked = bookedTickets.map(t => t.seat_number);

    const generatedSeats = rows.flatMap(row =>
      Array.from({ length: showtime.studio.total_cols }, (_, i) => ({
        id: `${row}${i + 1}`,
        row,
        number: i + 1,
        status: booked.includes(`${row}${i + 1}`) ? 'booked' : 'available'
      } as Seat))
    );
    setSeats(generatedSeats);
  }, [showtime]); // Re-run if showtime changes

  const handleSeatClick = (seatId: string) => {
    // ... (Keep existing logic) ...
    setSeats(prev => prev.map(seat => {
      if (seat.id === seatId && seat.status !== 'booked') {
        return { ...seat, status: seat.status === 'selected' ? 'available' : 'selected' };
      }
      return seat;
    }));
  };

  if (!showtime || !movie) return <div>Loading or Not Found...</div>;

  const selectedSeats = seats.filter(s => s.status === 'selected');
  const totalPrice = selectedSeats.length * showtime.price;
  const rowLetters = Array.from({ length: showtime.studio.total_rows }, (_, i) => String.fromCharCode(65 + i));

  return (
    <div className="max-w-5xl mx-auto px-6 py-8 text-white">
      <button
        onClick={() => navigate(`/dashboard/movie/${movie.id}`)} // Go back to specific movie
        className="flex items-center gap-2 text-gray-400 hover:text-yellow-500 mb-6 transition-colors cursor-pointer"
      >
        <ArrowLeft size={20} />
        Back to Showtimes
      </button>

      {/* ... (Paste the rest of your Seat UI here, it's mostly visual) ... */}
      
       <div className="bg-zinc-900 rounded-lg p-6 mb-8 border border-yellow-600/20">
        <h2 className="text-2xl font-bold mb-2 text-yellow-500">{movie.title}</h2>
        {/* ... Info ... */}
      </div>

       {/* ... Screen & Seats Map ... */}
            <div className="flex flex-col items-center overflow-x-auto pb-4">
            {/* 1. TOP HEADER (Seat Numbers) */}
            <div className="flex items-center gap-2 mb-4">
                {/* Empty space to align with the Row Letters column */}
                <div className="w-10" /> 
                {Array.from({ length: showtime.studio.total_cols }, (_, i) => (
                <div key={i} className="w-10 text-center text-xs font-bold text-zinc-500">
                    {i + 1}
                </div>
                ))}
            </div>

            {/* 2. SEAT ROWS */}
            <div className="space-y-3">
                {rowLetters.map(row => (
                <div key={row} className="flex items-center justify-center gap-2">
                    {/* ROW LETTER (Left Column) */}
                    <div className="w-10 flex items-center justify-center text-sm font-bold text-yellow-500/60 uppercase">
                    {row}
                    </div>

                    {/* INDIVIDUAL SEATS */}
                    {seats.filter(s => s.row === row).map(seat => (
                    <button 
                        key={seat.id} 
                        onClick={() => handleSeatClick(seat.id)} 
                        disabled={seat.status === 'booked'}
                        className={`
                        w-10 h-10
                        flex-shrink-0
                        flex items-center justify-center 
                        text-xs font-bold 
                        rounded-t-lg
                        transition-all
                        ${seat.status === 'selected' ? 'bg-yellow-500 text-black scale-110 shadow-lg shadow-yellow-500/20' : 
                            seat.status === 'booked' ? 'bg-zinc-800 text-zinc-600 cursor-not-allowed' : 
                            'bg-zinc-700 hover:bg-zinc-600 text-white'}
                        `}
                    >
                        {/* Keeping the number inside is helpful for mobile, but the grid now handles the labels */}
                        {seat.number}
                    </button>
                    ))}
                </div>
                ))}
            </div>

            {/* 3. SCREEN INDICATOR (Optional but recommended) */}
            <div className="mt-12 w-full max-w-md mx-auto">
                <div className="h-1 bg-gradient-to-r from-transparent via-yellow-500 to-transparent rounded-full shadow-[0_4px_20px_rgba(234,179,8,0.5)]" />
                <p className="text-center text-xs text-zinc-500 mt-2 tracking-[0.2em] uppercase">Screen</p>
            </div>
        </div>
      {/* ... Total Price & Booking Button ... */}
    </div>
  );
};

export default SeatSelection;