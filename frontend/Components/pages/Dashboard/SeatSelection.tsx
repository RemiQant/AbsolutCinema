import React, { useState, useEffect } from 'react';
import { ArrowLeft } from 'lucide-react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../../src/api/axios';
import { type Seat } from './types';

interface ShowtimeDetails {
  id: number;
  price: number;
  start_time: string;
  end_time: string;
  studio: {
    name: string;
    total_rows: number;
    total_cols: number;
  };
  movie: {
    id: number;
    title: string;
  };
}

const SeatSelection: React.FC = () => {
  const { showtimeId } = useParams();
  const navigate = useNavigate();
  
  const [showtime, setShowtime] = useState<ShowtimeDetails | null>(null);
  const [seats, setSeats] = useState<Seat[]>([]);
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!showtimeId) return;

        // 1. Get Showtime Details
        const showtimeRes = await api.get(`/showtimes/${showtimeId}`);
        const showtimeData = showtimeRes.data.data;
        setShowtime(showtimeData);

        // 2. Get Occupied Seats (FIXED CRASH LOGIC)
        const occupiedRes = await api.get(`/showtimes/${showtimeId}/seats`);
        const rawData = occupiedRes.data.data;
        
        // SAFEGUARD: Force it to be an array so .includes() never crashes
        let occupiedSeats: string[] = [];
        if (Array.isArray(rawData)) {
            occupiedSeats = rawData;
        } else if (rawData && Array.isArray(rawData.occupied_seats)) {
            occupiedSeats = rawData.occupied_seats;
        }

        // 3. Generate Seat Grid
        const rows = Array.from({ length: showtimeData.studio.total_rows }, (_, i) => String.fromCharCode(65 + i));
        
        const generatedSeats = rows.flatMap(row =>
          Array.from({ length: showtimeData.studio.total_cols }, (_, i) => {
            const seatId = `${row}${i + 1}`; // <--- BACKTICKS ARE IMPORTANT HERE
            return {
              id: seatId,
              row,
              number: i + 1,
              status: occupiedSeats.includes(seatId) ? 'booked' : 'available'
            } as Seat;
          })
        );
        setSeats(generatedSeats);
        setLoading(false);

      } catch (error) {
        console.error("Failed to fetch data:", error);
        alert("Failed to load seat data. Check console!");
        // Don't navigate away automatically so you can debug
      }
    };

    fetchData();
  }, [showtimeId]); // Removed 'navigate' from deps to prevent loop

  const handleSeatClick = (seatId: string) => {
    setSeats(prev => prev.map(seat => {
      if (seat.id === seatId && seat.status !== 'booked') {
        return { ...seat, status: seat.status === 'selected' ? 'available' : 'selected' };
      }
      return seat;
    }));
  };

  const handleCheckout = async () => {
    const selectedSeats = seats.filter(s => s.status === 'selected');
    if (selectedSeats.length === 0) return;

    setProcessing(true);
    try {
      const payload = {
        showtime_id: Number(showtimeId),
        seat_numbers: selectedSeats.map(s => s.id)
      };

      const response = await api.post('/bookings', payload);
      const { payment_url } = response.data;

      if (payment_url) {
        window.location.href = payment_url;
      } else {
        alert("Booking created, but no payment link returned?");
      }

    } catch (error: any) {
      console.error("Booking failed:", error);
      if (error.response?.status === 401) {
        alert("You need to login first!");
        navigate('/login');
      } else if (error.response?.status === 409) {
        alert("Someone just stole your seat! Refreshing...");
        window.location.reload();
      } else {
        alert("Booking failed: " + (error.response?.data?.error || "Unknown error"));
      }
    } finally {
      setProcessing(false);
    }
  };

  if (loading || !showtime) return <div className="text-white p-10 text-center">Loading Cinema...</div>;

  const selectedSeats = seats.filter(s => s.status === 'selected');
  const totalPrice = selectedSeats.length * showtime.price;
  const rowLetters = Array.from({ length: showtime.studio.total_rows }, (_, i) => String.fromCharCode(65 + i));

  return (
    <div className="max-w-5xl mx-auto px-6 py-8 text-white">
      <button
        onClick={() => navigate(`/dashboard/movie/${showtime.movie.id}`)}
        className="flex items-center gap-2 text-gray-400 hover:text-yellow-500 mb-6 transition-colors cursor-pointer"
      >
        <ArrowLeft size={20} />
        Back to Movie
      </button>

      {/* HEADER */}
      <div className="bg-zinc-900 rounded-lg p-6 mb-8 border border-yellow-600/20">
        <h2 className="text-2xl font-bold mb-2 text-yellow-500">{showtime.movie.title}</h2>
        <div className="flex gap-6 text-sm text-gray-400">
          <p>Studio: <span className="text-white">{showtime.studio.name}</span></p>
          <p>Time: <span className="text-white">{new Date(showtime.start_time).toLocaleTimeString()}</span></p>
          <p>Price: <span className="text-white">Rp {showtime.price.toLocaleString('id-ID')}</span></p>
        </div>
      </div>

      {/* SCREEN */}
      <div className="flex flex-col items-center overflow-x-auto pb-4">
        <div className="flex items-center gap-2 mb-4">
          <div className="w-10" /> 
          {Array.from({ length: showtime.studio.total_cols }, (_, i) => (
            <div key={i} className="w-10 text-center text-xs font-bold text-zinc-500">{i + 1}</div>
          ))}
        </div>

        <div className="space-y-3">
          {rowLetters.map(row => (
            <div key={row} className="flex items-center justify-center gap-2">
              <div className="w-10 flex items-center justify-center text-sm font-bold text-yellow-500/60 uppercase">{row}</div>
              {seats.filter(s => s.row === row).map(seat => (
                <button 
                  key={seat.id} 
                  onClick={() => handleSeatClick(seat.id)} 
                  disabled={seat.status === 'booked'}
                  className={`
                    w-10 h-10 flex-shrink-0 flex items-center justify-center text-xs font-bold rounded-t-lg transition-all
                    ${seat.status === 'selected' ? 'bg-yellow-500 text-black scale-110 shadow-lg shadow-yellow-500/20' : 
                      seat.status === 'booked' ? 'bg-zinc-800 text-zinc-600 cursor-not-allowed' : 
                      'bg-zinc-700 hover:bg-zinc-600 text-white'}
                  `}
                >
                  {seat.number}
                </button>
              ))}
            </div>
          ))}
        </div>
        
        <div className="mt-12 w-full max-w-md mx-auto">
           <div className="h-1 bg-gradient-to-r from-transparent via-yellow-500 to-transparent rounded-full shadow-[0_4px_20px_rgba(234,179,8,0.5)]" />
           <p className="text-center text-xs text-zinc-500 mt-2 tracking-[0.2em] uppercase">Screen</p>
        </div>
      </div>

      {/* FOOTER */}
      <div className="mt-8 bg-zinc-900 rounded-xl p-6 border border-zinc-800 flex justify-between items-center sticky bottom-4 shadow-2xl">
        <div>
            <p className="text-gray-400 text-sm mb-1">Total Price</p>
            <p className="text-2xl font-bold text-yellow-500">Rp {totalPrice.toLocaleString('id-ID')}</p>
            <p className="text-xs text-zinc-500">{selectedSeats.length} seats selected</p>
        </div>
        <button
            onClick={handleCheckout}
            disabled={selectedSeats.length === 0 || processing}
            className={`
                px-8 py-3 rounded-lg font-bold transition-all
                ${selectedSeats.length > 0 
                    ? 'bg-yellow-500 hover:bg-yellow-600 text-black cursor-pointer' 
                    : 'bg-zinc-800 text-zinc-600 cursor-not-allowed'}
            `}
        >
            {processing ? 'Processing...' : 'Confirm Booking'}
        </button>
      </div>
    </div>
  );  
};

export default SeatSelection;