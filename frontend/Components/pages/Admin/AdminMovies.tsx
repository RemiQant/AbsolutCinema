import React, { useState } from 'react';
import { Plus, Edit, Trash2, Film } from 'lucide-react';
import { movies as mockMovies } from '../Dashboard/mockData';

const AdminMovies: React.FC = () => {
  const [movies, setMovies] = useState(mockMovies);

  const [nowPlayingMovies, setNowPlayingMovies] = useState();

  const fetchMovies = () => {
    
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-2xl font-bold text-yellow-500">Manage Movies</h1>
          <p className="text-gray-400 text-sm">Add, edit or remove movies from the catalog</p>
        </div>
        <button 
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
              <button className="p-2 hover:bg-blue-500/20 text-blue-400 rounded-lg transition-colors"><Edit size={20} /></button>
              <button className="p-2 hover:bg-red-500/20 text-red-500 rounded-lg transition-colors"><Trash2 size={20} /></button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default AdminMovies;