import React, { useState, useEffect } from 'react';
import { Plus, Trash2, Calendar, Clock, Loader2, X } from 'lucide-react';
import api from '../../../src/api/axios';

interface Showtime {
  id: number;
  movie_id: number;
  studio_id: number;
  start_time: string;
  price: number;
  movie: { title: string };
  studio: { name: string };
}

interface Movie { id: number; title: string; }
interface Studio { id: number; name: string; }

const AdminShowtimes: React.FC = () => {
  const [showtimes, setShowtimes] = useState<Showtime[]>([]);
  const [movies, setMovies] = useState<Movie[]>([]);
  const [studios, setStudios] = useState<Studio[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Form State
  const [formData, setFormData] = useState({
    movie_id: '',
    studio_id: '',
    start_time: '',
    price: 0
  });

  const fetchData = async () => {
    try {
      setLoading(true);
      const [stRes, mvRes, sdRes] = await Promise.all([
        api.get('/admin/showtimes'),
        api.get('/admin/movies'),
        api.get('/admin/studios')
      ]);
      setShowtimes(stRes.data?.data || []);
      setMovies(mvRes.data?.data || []);
      setStudios(sdRes.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch data", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, []);

  const handleDelete = async (id: number) => {
    if (!window.confirm("Hapus jadwal tayang ini?")) return;
    try {
      await api.delete(`/admin/showtimes/${id}`);
      setShowtimes(prev => prev.filter(s => s.id !== id));
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal menghapus");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Pastikan format ISO String untuk Golang time.Time
      const payload = {
        ...formData,
        movie_id: Number(formData.movie_id),
        studio_id: Number(formData.studio_id),
        start_time: new Date(formData.start_time).toISOString(),
        price: Number(formData.price)
      };

      await api.post('/admin/showtimes', payload);
      setIsModalOpen(false);
      fetchData(); // Refresh list
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal membuat showtime");
    }
  };

  if (loading) return <div className="flex h-64 items-center justify-center"><Loader2 className="animate-spin text-yellow-500" size={40} /></div>;

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold text-yellow-500">Showtime Schedule</h1>
        <button 
          onClick={() => setIsModalOpen(true)}
          className="bg-yellow-500 hover:bg-yellow-600 text-black px-4 py-2 rounded-lg font-bold flex items-center gap-2 transition-all">
          <Plus size={20}/> Create Showtime
        </button>
      </div>

      <div className="bg-zinc-900 border border-yellow-600/20 rounded-xl overflow-hidden shadow-xl">
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
            {showtimes.map((st) => (
              <tr key={st.id} className="hover:bg-zinc-800/30 transition-colors text-sm">
                <td className="p-4 text-white font-medium">{st.movie?.title}</td>
                <td className="p-4 text-gray-400">{st.studio?.name}</td>
                <td className="p-4">
                  <div className="flex flex-col">
                    <span className="text-white flex items-center gap-1">
                      <Calendar size={12}/> {new Date(st.start_time).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })}
                    </span>
                    <span className="text-gray-500 flex items-center gap-1">
                      <Clock size={12}/> {new Date(st.start_time).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })}
                    </span>
                  </div>
                </td>
                <td className="p-4 text-yellow-500 font-bold">Rp {st.price.toLocaleString('id-ID')}</td>
                <td className="p-4">
                  <div className="flex justify-center">
                    <button onClick={() => handleDelete(st.id)} className="text-red-500 hover:bg-red-500/10 p-2 rounded-lg transition-all">
                      <Trash2 size={18} />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* CREATE MODAL */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="bg-zinc-900 border border-zinc-800 w-full max-w-md rounded-2xl shadow-2xl animate-in zoom-in duration-200">
            <div className="p-6 border-b border-zinc-800 flex justify-between items-center">
              <h2 className="text-xl font-bold text-yellow-500">New Showtime</h2>
              <button onClick={() => setIsModalOpen(false)} className="text-zinc-500 hover:text-white"><X size={24} /></button>
            </div>
            
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Select Movie</label>
                <select 
                  required
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white outline-none focus:border-yellow-500"
                  value={formData.movie_id}
                  onChange={e => setFormData({...formData, movie_id: e.target.value})}
                >
                  <option value="">Choose Movie...</option>
                  {movies.map(m => <option key={m.id} value={m.id}>{m.title}</option>)}
                </select>
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Select Studio</label>
                <select 
                  required
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white outline-none focus:border-yellow-500"
                  value={formData.studio_id}
                  onChange={e => setFormData({...formData, studio_id: e.target.value})}
                >
                  <option value="">Choose Studio...</option>
                  {studios.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                </select>
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Start Time</label>
                <input 
                  required
                  type="datetime-local"
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white outline-none focus:border-yellow-500"
                  value={formData.start_time}
                  onChange={e => setFormData({...formData, start_time: e.target.value})}
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Price (IDR)</label>
                <input 
                  required
                  type="number"
                  placeholder="50000"
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white outline-none focus:border-yellow-500"
                  value={formData.price}
                  onChange={e => setFormData({...formData, price: parseInt(e.target.value)})}
                />
              </div>

              <button type="submit" className="w-full bg-yellow-500 hover:bg-yellow-600 text-black font-bold py-3 rounded-xl transition-all shadow-lg shadow-yellow-500/20">
                Save Schedule
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminShowtimes;