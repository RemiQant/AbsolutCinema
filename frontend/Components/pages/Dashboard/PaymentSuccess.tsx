import React from 'react';
import { Check, Ticket, Home, ArrowRight } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const PaymentSuccess: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-black text-white flex flex-col items-center justify-center px-6">
      
      {/* Container Utama dengan Glass Effect */}
      <div className="w-full max-w-sm p-8 rounded-3xl bg-zinc-900/40 border border-zinc-800 text-center">
        
        {/* Lingkaran Centang dengan Glow Halus */}
        <div className="w-20 h-20 bg-yellow-500 rounded-full flex items-center justify-center mx-auto mb-6 shadow-[0_0_30px_rgba(234,179,8,0.2)]">
          <Check size={40} className="text-black stroke-[3px]" />
        </div>

        <h1 className="text-2xl font-bold uppercase tracking-tight mb-2 ">
          Payment Absolute.
        </h1>
        <p className="text-zinc-500 text-sm mb-8">
          Transaksi berhasil. Sampai jumpa di bioskop!
        </p>

        {/* Detail Ringkas dengan Background Gelap */}
        <div className="bg-black/40 rounded-2xl p-6 mb-8 space-y-4 border border-zinc-800/50">
          <div className="flex justify-between items-center text-sm">
            <span className="text-zinc-500 uppercase tracking-widest text-[10px] font-bold">Movie</span>
            <span className="font-semibold">Interstellar</span>
          </div>
          <div className="flex justify-between items-center text-sm">
            <span className="text-zinc-500 uppercase tracking-widest text-[10px] font-bold">Seats</span>
            <span className="text-yellow-500 font-bold tracking-widest">A3, A4</span>
          </div>
          <div className="flex justify-between items-center text-sm">
            <span className="text-zinc-500 uppercase tracking-widest text-[10px] font-bold">Studio</span>
            <span className="font-semibold">Studio 1</span>
          </div>
        </div>

        {/* Action Button */}
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

      {/* Footer Note Kecil */}
      <p className="mt-8 text-zinc-700 text-[10px] uppercase tracking-[0.3em]">
        Invoice sent to your email
      </p>

    </div>
  );
};

export default PaymentSuccess;