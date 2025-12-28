import React, { useEffect, useState } from 'react';
import { Check, Home, Ticket } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const PaymentSuccess: React.FC = () => {
    const navigate = useNavigate();
    const [booking, setBooking] = useState<any>(null);

    useEffect(() => {
        // Ambil snapshot data yang kita simpan di SeatSelection
        const savedData = localStorage.getItem('last_successful_booking');
        if (savedData) {
            setBooking(JSON.parse(savedData));
        }
    }, []);

    // Gunakan data snapshot, jika tidak ada baru pakai Unknown
    const movieTitle = booking?.movieTitle || "Movie Confirmed";
    const studioName = booking?.studioName || "Cinema Studio";
    const seatList = booking?.seats || "Confirmed";
    const bookingId = booking?.bookingId || "-------";
    
    const formattedDate = booking?.startTime 
        ? new Date(booking.startTime).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
        : "Today";
    const formattedTime = booking?.startTime 
        ? new Date(booking.startTime).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
        : "--:--";

    return (
        <div className="min-h-screen bg-black text-white flex flex-col items-center justify-center p-6">
            <div className="w-full max-w-md bg-zinc-900/50 border border-zinc-800 rounded-[2.5rem] overflow-hidden shadow-2xl">
                <div className="p-8 pb-4 text-center">
                    <div className="w-20 h-20 bg-yellow-500 rounded-full flex items-center justify-center mx-auto mb-6 shadow-[0_0_40px_rgba(234,179,8,0.3)]">
                        <Check size={40} className="text-black stroke-[3px]" />
                    </div>
                    <h1 className="text-2xl font-black uppercase italic tracking-tighter text-yellow-500">Success!</h1>
                    <p className="text-zinc-500 text-sm mt-2 font-medium italic">Enjoy your movie absolute.</p>
                </div>

                <div className="px-8 pb-8">
                    <div className="bg-black/40 border border-zinc-800/50 rounded-3xl p-6 space-y-5 relative overflow-hidden">
                        <div className="absolute top-0 right-0 w-32 h-32 -mr-8 -mt-8 opacity-10">
                            <Ticket size={120} className="rotate-12 text-yellow-500" />
                        </div>

                        <div className="space-y-1 relative z-10">
                            <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Movie Title</span>
                            <p className="text-lg font-bold leading-tight uppercase ">{movieTitle}</p>
                        </div>

                        <div className="grid grid-cols-2 gap-4 relative z-10">
                            <div className="space-y-1">
                                <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Studio</span>
                                <p className="text-sm font-semibold">{studioName}</p>
                            </div>
                            <div className="space-y-1">
                                <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Seats</span>
                                <p className="text-sm font-bold text-yellow-500 uppercase">{seatList}</p>
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4 relative z-10">
                            <div className="space-y-1">
                                <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Date</span>
                                <p className="text-sm font-semibold">{formattedDate}</p>
                            </div>
                            <div className="space-y-1">
                                <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Time</span>
                                <p className="text-sm font-semibold">{formattedTime}</p>
                            </div>
                        </div>

                        <div className="pt-4 border-t border-zinc-800 flex justify-between items-center opacity-50 italic">
                            <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Ref ID</span>
                            <span className="text-xs font-mono">#{bookingId.slice(0, 8).toUpperCase()}</span>
                        </div>
                    </div>

                    <button 
                        onClick={() => {
                            localStorage.removeItem('last_successful_booking');
                            navigate('/dashboard');
                        }}
                        className="w-full mt-8 bg-white hover:bg-yellow-500 text-black font-black py-4 rounded-2xl transition-all flex items-center justify-center gap-3"
                    >
                        <Home size={20} />
                        BACK TO HOME
                    </button>
                </div>
            </div>
        </div>
    );
};

export default PaymentSuccess;