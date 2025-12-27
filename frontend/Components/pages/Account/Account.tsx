import React from 'react';
import { User, Mail, Award, LogOut, Zap, Calendar } from 'lucide-react';
import { useAuth } from "../../../src/context/AuthContext";
import { useNavigate } from 'react-router-dom';
import Navbar from '../../Navbar/Navbar';

const AccountPage: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  if (!user) return null;

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-black text-white font-sans selection:bg-yellow-500/30">
        <div className="max-w-2xl mx-auto px-6 py-20">
          
          {/* Section Profil Utama */}
          <div className="flex flex-col items-center text-center mb-16">
            <div className="w-24 h-24 bg-zinc-900 border-2 border-yellow-500/20 rounded-full flex items-center justify-center mb-6 shadow-[0_0_40px_rgba(234,179,8,0.05)]">
              <User size={48} className="text-yellow-500" />
            </div>
            <h1 className="text-3xl font-bold tracking-tight mb-2">{user.username}</h1>
            <div className="flex items-center gap-2">
               <Zap size={14} className="text-yellow-500 animate-pulse" />
               <span className="text-zinc-500 text-[10px] font-black uppercase tracking-[0.2em]">
                 Status: Active
               </span>
            </div>
          </div>

          <div className="space-y-1 border-t border-zinc-900">
            
            <div className="flex items-center justify-between py-6 border-b border-zinc-900">
              <div className="flex items-center gap-4 text-zinc-500">
                <Mail size={18} />
                <span className="text-sm font-semibold uppercase tracking-wider text-[11px]">Email</span>
              </div>
              <span className="text-sm text-zinc-200">{user.email}</span>
            </div>


            {/* Join Date (Ganti dari Security) */}
            <div className="flex items-center justify-between py-6 border-b border-zinc-900">
              <div className="flex items-center gap-4 text-zinc-500">
                <Calendar size={18} />
                <span className="text-sm font-semibold uppercase tracking-wider text-[11px]">Joined Since</span>
              </div>
              <span className="text-sm text-zinc-200">December 2025</span>
            </div>
          </div>

          {/* Tombol Logout */}
          <div className="mt-16 text-center">
            <button 
              onClick={handleLogout}
              className="group flex items-center gap-3 mx-auto text-zinc-600 hover:text-red-500 transition-all duration-300 cursor-pointer"
            >
              <LogOut size={18} className="group-hover:-translate-x-1 transition-transform " />
              <span className="text-xs font-black uppercase tracking-[0.2em] cursor-pointer">Sign Out</span>
            </button>
          </div>

        </div>
      </div>
    </>
  );
};

export default AccountPage;