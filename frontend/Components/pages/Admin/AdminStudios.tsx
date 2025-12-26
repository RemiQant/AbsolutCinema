import React from 'react';
import { Plus, Edit, Trash2, Grid3X3 } from 'lucide-react';

const AdminStudios: React.FC = () => {
  // Replace with actual data import
  const studios = [
    { id: 1, name: "Studio 1", rows: 8, cols: 10 },
    { id: 2, name: "Studio 2", rows: 6, cols: 8 },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold text-yellow-500">Theater Studios</h1>
        <button className="flex items-center gap-2 bg-yellow-500 text-black px-4 py-2 rounded-lg font-bold">
          <Plus size={20} /> New Studio
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {studios.map(studio => (
          <div key={studio.id} className="bg-zinc-900 border border-yellow-600/20 rounded-xl p-6 relative overflow-hidden group">
            <div className="absolute top-0 right-0 p-3 opacity-10 group-hover:opacity-20 transition-opacity">
              <Grid3X3 size={80} />
            </div>
            <h3 className="text-xl font-bold text-white mb-4">{studio.name}</h3>
            <div className="space-y-2 mb-6">
              <div className="flex justify-between text-sm text-gray-400">
                <span>Total Capacity</span>
                <span className="text-yellow-500 font-bold">{studio.rows * studio.cols} Seats</span>
              </div>
              <div className="flex justify-between text-xs text-zinc-500 italic">
                <span>Layout</span>
                <span>{studio.rows} Rows x {studio.cols} Columns</span>
              </div>
            </div>
            <div className="flex gap-3">
              <button className="flex-1 flex items-center justify-center gap-2 py-2 bg-zinc-800 hover:bg-zinc-700 rounded-lg text-sm transition-colors">
                <Edit size={16}/> Edit
              </button>
              <button className="p-2 text-red-500 hover:bg-red-500/10 rounded-lg">
                <Trash2 size={18}/>
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default AdminStudios;