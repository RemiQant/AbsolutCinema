import React from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { Film, LayoutDashboard, DoorOpen, CalendarDays , ArrowLeft } from 'lucide-react';

const AdminLayout: React.FC = () => {
  const location = useLocation();

  const menuItems = [
    { path: '/admin/movies', icon: <Film size={20} />, label: 'Movies' },
    { path: '/admin/studios', icon: <DoorOpen size={20} />, label: 'Studios' },
    { path: '/admin/showtimes', icon: <CalendarDays size={20} />, label: 'Showtimes' },
  ];

  return (
    <div className="flex min-h-screen bg-black">
      <aside className="w-64 border-r border-yellow-600/20 bg-zinc-950 p-6">
        <div className="flex items-center gap-3 mb-10 px-2">
          <LayoutDashboard className="text-yellow-500" />
          <span className="text-yellow-500 font-bold tracking-wider">ADMIN PANEL</span>
        </div>
        
        <nav className="space-y-2">
          {menuItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                location.pathname === item.path 
                ? 'bg-yellow-500 text-black font-bold' 
                : 'text-gray-400 hover:bg-zinc-900 hover:text-yellow-500'
              }`}
            >
              {item.icon}
              {item.label}
            </Link>
          ))}
          <hr className="border-zinc-800 my-4" />
          <Link to="/dashboard" className="flex items-center gap-3 px-4 py-3 text-gray-500 hover:text-white">
             <ArrowLeft size={18}/> Exit to User View
          </Link>
        </nav>
      </aside>

      <main className="flex-1 p-8 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  );
};

export default AdminLayout;