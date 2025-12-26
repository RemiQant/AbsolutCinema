import React, { useEffect, useState } from 'react';
import { Check, Home, ArrowRight } from 'lucide-react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import api from '../../../src/api/axios'; 

interface BookingData {
  id: string;
  tickets: Array<{
    seat_number: string;
    showtime: {
        movie: { title: string };
        studio: { name: string };
    }
  }>;
}

const PaymentSuccess: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  // Xendit might return ?id=... or ?external_id=... check your Xendit settings
  // Assuming Xendit redirects with external_id (which is our booking ID)
  const bookingId = searchParams.get('external_id') || searchParams.get('id');

  const [booking, setBooking] = useState<BookingData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchBooking = async () => {
      if (!bookingId) {
          setLoading(false);
          return;
      }
      try {
        const res = await api.get(`/bookings/${bookingId}`);
        setBooking(res.data.data);
      } catch (err) {
        console.error("Failed to fetch receipt", err);
      } finally {
        setLoading(false);
      }
    };

    fetchBooking();
  }, [bookingId]);

  if (loading) return <div className="text-white text-center p-10">Verifying Payment...</div>;

  // Derive display data from the first ticket (since all tickets in a booking are same movie/studio)
  const firstTicket = booking?.tickets[0];
  const movieTitle = firstTicket?.showtime.movie.title || "Unknown Movie";
  const studioName = firstTicket?.showtime.studio.name || "Unknown Studio";
  const seatList = booking?.tickets.map(t => t.seat_number).join(", ") || "N/A";

  return (
    <div className="min-h-screen bg-black text-white flex flex-col items-center justify-center px-6">
      <div className="w-full max-w-sm p-8 rounded-3xl bg-zinc-900/40 border border-zinc-800 text-center">
        
        <div className="w-20 h-20 bg-yellow-500 rounded-full flex items-center justify-center mx-auto mb-6 shadow-[0_0_30px_rgba(234,179,8,0.2)]">
          <Check size={40} className="text-black stroke-[3px]" />
        </div>

        <h1 className="text-2xl font-bold uppercase tracking-tight mb-2 ">Payment Absolute.</h1>
        <p className="text-zinc-500 text-sm mb-8">Transaksi berhasil. Sampai jumpa di bioskop!</p>

        {/* DYNAMIC RECEIPT */}
        <div className="bg-black/40 rounded-2xl p-6 mb-8 space-y-4 border border-zinc-800/50">
          <div className="flex justify-between items-center text-sm">
            <span className="text-zinc-500 uppercase tracking-widest text-[10px] font-bold">Movie</span>
            <span className="font-semibold">{movieTitle}</span>
          </div>
          <div className="flex justify-between items-center text-sm">
            <span className="text-zinc-500 uppercase tracking-widest text-[10px] font-bold">Seats</span>
            <span className="text-yellow-500 font-bold tracking-widest">{seatList}</span>
          </div>
          <div className="flex justify-between items-center text-sm">
            <span className="text-zinc-500 uppercase tracking-widest text-[10px] font-bold">Studio</span>
            <span className="font-semibold">{studioName}</span>
          </div>
        </div>

        <div className="space-y-4">
          <button 
            onClick={() => navigate('/dashboard')}
            className="w-full bg-yellow-500 hover:bg-yellow-600 text-black font-bold py-4 rounded-xl transition-all flex items-center justify-center gap-2"
          >
            <Home size={18} />
            BACK TO HOME
          </button>
        </div>
      </div>
    </div>
  );
};

export default PaymentSuccess;