import React, { useEffect, useState } from 'react';
import { Star, Info } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import api from '../../../src/api/axios'; // Make sure path matches your folder structure

interface Movie {
  id: number;
  title: string;
  poster_url: string;
  rating: number;
  age_rating: string;
}

const MovieGrid: React.FC = () => {
  const navigate = useNavigate();
  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        // Fetch Real Movies from Backend
        const response = await api.get('/movies');
        // The backend returns { data: [...] }
        setMovies(response.data.data);
      } catch (error) {
        console.error("Failed to fetch movies:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchMovies();
  }, []);

  if (loading) return <div className="text-white text-center mt-20">Loading Movies...</div>;

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
      {movies.map((movie) => (
        <div 
          key={movie.id}
          className="group relative bg-zinc-900 rounded-xl overflow-hidden hover:scale-105 transition-all duration-300 border border-zinc-800 hover:border-yellow-500/50 cursor-pointer shadow-lg hover:shadow-yellow-500/10"
          onClick={() => navigate(`/dashboard/movie/${movie.id}`)}
        >
          {/* Poster Image */}
          <div className="aspect-[2/3] relative">
            <img 
              src={movie.poster_url} 
              alt={movie.title}
              className="w-full h-full object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col justify-end p-4">
               <button className="bg-yellow-500 text-black font-bold py-2 px-4 rounded-lg flex items-center justify-center gap-2 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300">
                  <Info size={18} />
                  Details & Book
               </button>
            </div>
          </div>

          {/* Movie Info */}
          <div className="p-4">
            <h3 className="text-white font-bold truncate mb-1" title={movie.title}>{movie.title}</h3>
            <div className="flex justify-between items-center">
              <div className="flex items-center gap-1 text-yellow-500">
                <Star size={14} fill="currentColor" />
                <span className="text-sm font-medium">{movie.rating}</span>
              </div>
              <span className="text-xs font-bold text-zinc-500 border border-zinc-700 px-2 py-0.5 rounded">
                {movie.age_rating}
              </span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default MovieGrid;