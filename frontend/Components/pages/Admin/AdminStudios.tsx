import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Grid3X3, Loader2, X } from 'lucide-react';
import api from '../../../src/api/axios';

interface Studio {
  id?: number;
  name: string;
  total_rows: number;
  total_cols: number;
}

const AdminStudios: React.FC = () => {
  const [studios, setStudios] = useState<Studio[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentStudio, setCurrentStudio] = useState<Studio | null>(null);

  const [formData, setFormData] = useState<Studio>({
    name: '',
    total_rows: 0,
    total_cols: 0
  });

  const fetchData = async () => {
    try {
      setLoading(true);
      const response = await api.get('/admin/studios');
      setStudios(response.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch studios");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  // --- HANDLERS ---

  const handleOpenModal = (studio: Studio | null = null) => {
    if (studio) {
      setCurrentStudio(studio);
      setFormData(studio);
    } else {
      setCurrentStudio(null);
      setFormData({ name: '', total_rows: 0, total_cols: 0 });
    }
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Hapus studio ini? Semua jadwal terkait mungkin akan terpengaruh.")) return;
    try {
      await api.delete(`/admin/studios/${id}`);
      setStudios(prev => prev.filter(s => s.id !== id));
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal menghapus studio");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const payload = {
      ...formData,
      total_rows: Number(formData.total_rows),
      total_cols: Number(formData.total_cols)
    };

    try {
      if (currentStudio?.id) {
        await api.put(`/admin/studios/${currentStudio.id}`, payload);
      } else {
        await api.post('/admin/studios', payload);
      }
      setIsModalOpen(false);
      fetchData();
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal menyimpan studio");
    }
  };

  if (loading) return <div className="flex h-64 items-center justify-center"><Loader2 className="animate-spin text-yellow-500" size={40} /></div>;

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-2xl font-bold text-yellow-500">Theater Studios</h1>
          <p className="text-gray-400 text-sm">Manage studio capacity and seat layouts</p>
        </div>
        <button 
          onClick={() => handleOpenModal()}
          className="flex items-center gap-2 bg-yellow-500 hover:bg-yellow-600 text-black px-4 py-2 rounded-lg font-bold transition-all shadow-lg shadow-yellow-500/20">
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
                <span className="text-yellow-500 font-bold">{studio.total_rows * studio.total_cols} Seats</span>
              </div>
              <div className="flex justify-between text-xs text-zinc-500 italic">
                <span>Layout</span>
                <span>{studio.total_rows} Rows x {studio.total_cols} Columns</span>
              </div>
            </div>
            <div className="flex gap-3">
              <button 
                onClick={() => handleOpenModal(studio)}
                className="flex-1 flex items-center justify-center gap-2 py-2 bg-zinc-800 hover:bg-zinc-700 rounded-lg text-sm transition-colors text-white">
                <Edit size={16}/> Edit
              </button>
              <button 
                onClick={() => handleDelete(studio.id!)}
                className="p-2 text-red-500 hover:bg-red-500/10 rounded-lg transition-colors">
                <Trash2 size={18}/>
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* MODAL FORM */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="bg-zinc-900 border border-zinc-800 w-full max-w-md rounded-2xl shadow-2xl animate-in fade-in zoom-in duration-200">
            <div className="p-6 border-b border-zinc-800 flex justify-between items-center">
              <h2 className="text-xl font-bold text-yellow-500">{currentStudio ? 'Edit Studio' : 'Create New Studio'}</h2>
              <button onClick={() => setIsModalOpen(false)} className="text-zinc-500 hover:text-white"><X size={24} /></button>
            </div>
            
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Studio Name</label>
                <input 
                  required
                  type="text"
                  placeholder="e.g. Studio 1"
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:border-yellow-500 outline-none"
                  value={formData.name}
                  onChange={e => setFormData({...formData, name: e.target.value})}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Total Rows</label>
                  <input 
                    required
                    type="number"
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:border-yellow-500 outline-none"
                    value={formData.total_rows}
                    onChange={e => setFormData({...formData, total_rows: parseInt(e.target.value) || 0})}
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Total Columns</label>
                  <input 
                    required
                    type="number"
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:border-yellow-500 outline-none"
                    value={formData.total_cols}
                    onChange={e => setFormData({...formData, total_cols: parseInt(e.target.value) || 0})}
                  />
                </div>
              </div>

              <div className="bg-yellow-500/5 p-4 rounded-lg border border-yellow-500/10 mb-4">
                <p className="text-xs text-yellow-500/70 text-center italic">
                  Preview Capacity: <span className="font-bold">{formData.total_rows * formData.total_cols} seats</span>
                </p>
              </div>

              <button type="submit" className="w-full bg-yellow-500 hover:bg-yellow-600 text-black font-bold py-3 rounded-xl transition-all shadow-lg shadow-yellow-500/20">
                {currentStudio ? 'Update Studio' : 'Create Studio'}
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminStudios;