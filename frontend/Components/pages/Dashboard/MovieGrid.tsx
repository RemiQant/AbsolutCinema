import React from 'react';
import { Loader2, Star } from 'lucide-react'; // Menambah icon Star untuk estetika jika diperlukan
import { Link } from 'react-router-dom';  
import { useState, useEffect } from 'react';
import api from "../../../src/api/axios";

interface Movie {
  id: number;
  title: string;
  description: string;
  duration_minutes: number;
  poster_url: string;
  rating: string;
}

const MovieGrid: React.FC = () => {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, serError] = useState<string | null>(null);
  
  useEffect(() => {
    const fetchMovies = async () => {
      try {
        setLoading(true);
        const response = await api.get('/movies');
        // Pastikan akses ke response.data.data aman
        const movieData = response.data?.data || [];
        setMovies(movieData);
      } catch (err) {
        serError("Failed to fetch movies. Please try again later.");
      } finally {
        setLoading(false);
      }
    }
    fetchMovies();
  }, [])

  if (loading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="animate-spin text-yellow-500" size={40} />
      </div>
    );
  }

  if (error) {
    return <div className="text-red-500 text-center p-10">{error}</div>;
  }

  return (
    <div className="max-w-7xl mx-auto px-6 py-8">
      <h2 className="text-yellow-500 text-2xl font-bold mb-6 tracking-wide uppercase">Currently Showing</h2>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-10">
        {movies.map(movie => (
          /* PEMBUNGKUS UTAMA PER CARD */
          <div key={movie.id} className="bg-[#1a1a1a] p-3 rounded-xl shadow-2xl hover:bg-[#252525] transition-colors duration-300">
            
            <Link
              to={`/dashboard/movie/${movie.id}`} 
              className="group cursor-pointer flex flex-col"
            >
              {/* Container Gambar */}
              <div className="relative aspect-[2/3] rounded-lg overflow-hidden mb-4 border border-gray-800 group-hover:border-yellow-500 transition-all duration-300">
                  
                  {/* Badge Rating */}
                  <div className="absolute top-2 left-2 z-10 bg-yellow-500 text-black text-[10px] font-black px-2 py-0.5 rounded uppercase">
                    {movie.rating}
                  </div>

                  <img 
                    src={movie.poster_url} 
                    alt={movie.title} 
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
                  />
              </div>

              {/* Bagian Teks dalam Pembungkus */}
              <div className="px-1 pb-2">
                <h3 className="text-white font-bold text-base line-clamp-1 group-hover:text-yellow-500 transition-colors">
                  {movie.title}
                </h3>
                
                <div className="flex items-center justify-between mt-2">
                  <span className="text-gray-500 text-xs font-medium">
                    {movie.duration_minutes} min
                  </span>
                  <span className="text-[10px] text-yellow-500/50 border border-yellow-500/30 px-1.5 rounded">
                    HD
                  </span>
                </div>
              </div>
            </Link>

          </div>
        ))}
      </div>
    </div>
  );
};

export default MovieGrid;