import React from 'react';
import { Plus, Trash2, Calendar, Clock, DollarSign } from 'lucide-react';

const AdminShowtimes: React.FC = () => {
  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold text-yellow-500">Showtime Schedule</h1>
        <button className="bg-yellow-500 text-black px-4 py-2 rounded-lg font-bold flex items-center gap-2">
          <Plus size={20}/> Create Showtime
        </button>
      </div>

      <div className="bg-zinc-900 border border-yellow-600/20 rounded-xl overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-zinc-800/50 text-yellow-500 text-sm uppercase tracking-wider">
              <th className="p-4 font-bold">Movie</th>
              <th className="p-4 font-bold">Studio</th>
              <th className="p-4 font-bold">Date & Time</th>
              <th className="p-4 font-bold">Price</th>
              <th className="p-4 font-bold text-center">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-800">
            {/* Example Row */}
            <tr className="hover:bg-zinc-800/30 transition-colors">
              <td className="p-4 text-white font-medium">Interstellar</td>
              <td className="p-4 text-gray-400 text-sm">Studio 1</td>
              <td className="p-4">
                <div className="flex flex-col">
                  <span className="text-white text-sm flex items-center gap-1"><Calendar size={12}/> Dec 25, 2025</span>
                  <span className="text-gray-500 text-xs flex items-center gap-1"><Clock size={12}/> 19:00 PM</span>
                </div>
              </td>
              <td className="p-4 text-yellow-500 font-bold">Rp 50.000</td>
              <td className="p-4">
                <div className="flex justify-center">
                  <button className="text-red-500 hover:bg-red-500/10 p-2 rounded-lg transition-all">
                    <Trash2 size={18} />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AdminShowtimes;