import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Film, Loader2, X } from 'lucide-react';
import api from '../../../src/api/axios';

interface Movie {
  id?: number;
  title: string;
  description: string;
  poster_url: string;
  duration_minutes: number;
  rating: string;
  genre: string;
}

const AdminMovies: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [movies, setMovies] = useState<Movie[]>([]);

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentMovie, setCurrentMovie] = useState<Movie | null>(null);
  const [formData, setFormData] = useState<Movie>({
    title: '',
    description: '',
    poster_url: '',
    duration_minutes: 0,
    rating: 'G',
    genre: ''
  });

  const fetchMovies = async () => {
    try {
      setLoading(true);
      const response = await api.get('admin/movies');
      setMovies(response.data?.data || []);
    } catch (err) {
      setError("Failed to fetch movies.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMovies();
  }, []);

  // --- HANDLERS ---

  const handleOpenModal = (movie: Movie | null = null) => {
    if (movie) {
      setCurrentMovie(movie);
      setFormData(movie);
    } else {
      setCurrentMovie(null);
      setFormData({ title: '', description: '', poster_url: '', duration_minutes: 0, rating: 'G', genre: '' });
    }
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (window.confirm("Are you sure you want to delete this movie?")) {
      try {
        await api.delete(`/movies/${id}`);
        setMovies(movies.filter(m => m.id !== id));
      } catch (err) {
        alert("Failed to delete movie.");
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const payload = {
      title: formData.title,
      description: formData.description,
      poster_url: formData.poster_url,
      duration_minutes: Number(formData.duration_minutes), 
      rating: formData.rating,
    };
  
    try {
      if (currentMovie) {
        // EDIT: /api/admin/movies/:id
        await api.put(`/admin/movies/${currentMovie.id}`, payload);
      } else {
        // ADD: /api/admin/movies
        await api.post('/admin/movies', payload);
      }
      
      setIsModalOpen(false);
      fetchMovies(); 
      alert("Movie saved successfully!");
    } catch (err: any) {
      console.error("Save Error:", err.response?.data);
      alert(`Failed: ${err.response?.data?.error || "Check console for details"}`);
    }
  };

  // --- RENDER ---

  if (loading) return <div className="flex h-64 items-center justify-center"><Loader2 className="animate-spin text-yellow-500" size={40} /></div>;

  return (
    <div className="relative">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-2xl font-bold text-yellow-500">Manage Movies</h1>
          <p className="text-gray-400 text-sm">Add, edit or remove movies</p>
        </div>
        <button 
          onClick={() => handleOpenModal()}
          className="flex items-center gap-2 bg-yellow-500 hover:bg-yellow-600 text-black px-4 py-2 rounded-lg font-bold transition-all shadow-lg shadow-yellow-500/20">
          <Plus size={20} /> Add Movie
        </button>
      </div>

      <div className="grid gap-4">
        {movies.map((movie) => (
          <div key={movie.id} className="bg-zinc-900 border border-yellow-600/20 rounded-xl p-4 flex items-center gap-6 hover:border-yellow-500/50 transition-all">
            <img src={movie.poster_url} alt={movie.title} className="w-20 h-28 object-cover rounded-lg shadow-md" />
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-1">
                <h3 className="text-xl font-bold text-white">{movie.title}</h3>
                <span className="text-xs px-2 py-0.5 bg-zinc-800 text-yellow-500 border border-yellow-500/30 rounded">{movie.rating}</span>
              </div>
              <p className="text-gray-400 text-sm line-clamp-2 max-w-2xl">{movie.description}</p>
              <div className="flex items-center gap-4 mt-2 text-xs text-zinc-500">
                 <span className="flex items-center gap-1"><Film size={14}/> {movie.duration_minutes} min</span>
              </div>
            </div>
            <div className="flex gap-2">
              <button onClick={() => handleOpenModal(movie)} className="p-2 hover:bg-blue-500/20 text-blue-400 rounded-lg transition-colors"><Edit size={20} /></button>
              <button onClick={() => handleDelete(movie.id!)} className="p-2 hover:bg-red-500/20 text-red-500 rounded-lg transition-colors"><Trash2 size={20} /></button>
            </div>
          </div>
        ))}
      </div>

      {/* MODAL COMPONENT */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="bg-zinc-900 border border-zinc-800 w-full max-w-lg rounded-2xl overflow-hidden shadow-2xl animate-in fade-in zoom-in duration-200">
            <div className="p-6 border-b border-zinc-800 flex justify-between items-center">
              <h2 className="text-xl font-bold text-yellow-500">{currentMovie ? 'Edit Movie' : 'Add New Movie'}</h2>
              <button onClick={() => setIsModalOpen(false)} className="text-zinc-500 hover:text-white"><X size={24} /></button>
            </div>
            
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="col-span-2">
                  <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Movie Title</label>
                  <input 
                    required 
                    type="text" 
                    value={formData.title} 
                    onChange={e => setFormData({...formData, title: e.target.value})}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-yellow-500" 
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Duration (Min)</label>
                  <input 
                    required 
                    type="number" 
                    value={formData.duration_minutes} 
                    onChange={e => setFormData({...formData, duration_minutes: parseInt(e.target.value)})}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-yellow-500" 
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Rating</label>
                  <select 
                    value={formData.rating} 
                    onChange={e => setFormData({...formData, rating: e.target.value})}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-yellow-500">
                    <option value="G">G</option>
                    <option value="PG">PG</option>
                    <option value="PG-13">PG-13</option>
                    <option value="R">R</option>
                    <option value="SU">SU</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Poster URL</label>
                <input 
                  required 
                  type="text" 
                  value={formData.poster_url} 
                  onChange={e => setFormData({...formData, poster_url: e.target.value})}
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-yellow-500" 
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase mb-1">Description</label>
                <textarea 
                  required 
                  rows={3}
                  value={formData.description} 
                  onChange={e => setFormData({...formData, description: e.target.value})}
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-yellow-500" 
                />
              </div>

              <div className="pt-4 flex gap-3">
                <button 
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="flex-1 bg-zinc-800 hover:bg-zinc-700 text-white font-bold py-3 rounded-xl transition-all">
                  Cancel
                </button>
                <button 
                  type="submit"
                  className="flex-1 bg-yellow-500 hover:bg-yellow-600 text-black font-bold py-3 rounded-xl transition-all">
                  {currentMovie ? 'Update Movie' : 'Save Movie'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminMovies;