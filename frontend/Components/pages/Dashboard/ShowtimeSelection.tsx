import React, { useEffect, useState } from 'react';
import { ArrowLeft, Calendar } from 'lucide-react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../../src/api/axios';

interface MovieDetail {
  id: number;
  title: string;
  description: string;
  poster_url: string;
  duration_minutes: number;
  genre: string;
}

interface Showtime {
  id: number;
  start_time: string;
  price: number;
  studio: {
    name: string;
  };
}

const ShowtimeSelection: React.FC = () => {
  const { movieId } = useParams();
  const navigate = useNavigate();

  const [movie, setMovie] = useState<MovieDetail | null>(null);
  const [showtimes, setShowtimes] = useState<Showtime[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!movieId) return;

        // 1. Get Movie Details (For the header)
        const movieRes = await api.get(`/movies/${movieId}`);
        setMovie(movieRes.data.data);

        // 2. Get Showtimes for this Movie
        // Backend route matches: GET /api/showtimes?movie_id=X
        const showtimeRes = await api.get(`/showtimes?movie_id=${movieId}`);
        setShowtimes(showtimeRes.data.data || []);
        
      } catch (error) {
        console.error("Failed to load data", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [movieId]);

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
  };

  if (loading) return <div className="text-white text-center mt-20">Loading Showtimes...</div>;
  if (!movie) return <div className="text-white text-center mt-20">Movie not found! 😭</div>;

  return (
    <div className="max-w-5xl mx-auto px-6 py-8 text-white">
      <button
        onClick={() => navigate('/dashboard')}
        className="flex items-center gap-2 text-gray-400 hover:text-yellow-500 mb-6 transition-colors cursor-pointer"
      >
        <ArrowLeft size={20} />
        Back to Movies
      </button>

      {/* Movie Info Header */}
       <div className="grid md:grid-cols-3 gap-6 mb-8">
            <div className="md:col-span-1">
              <img src={movie.poster_url} alt={movie.title} className="w-full rounded-lg border-2 border-yellow-600/20 shadow-2xl"/>
            </div>
            <div className="md:col-span-2">
               <h1 className="text-4xl font-bold text-yellow-500 mb-4">{movie.title}</h1>
               <div className="flex gap-3 mb-6">
                 <span className="px-3 py-1 bg-zinc-800 rounded-full text-xs font-bold text-zinc-400">{movie.genre}</span>
                 <span className="px-3 py-1 bg-zinc-800 rounded-full text-xs font-bold text-zinc-400">{movie.duration_minutes} Mins</span>
               </div>
               <p className="text-gray-300 leading-relaxed text-lg">{movie.description}</p>
            </div>
       </div>

      <h2 className="text-2xl font-bold text-yellow-500 mb-6 border-b border-zinc-800 pb-2">Select Showtime</h2>
      
      {showtimes.length === 0 ? (
        <div className="text-zinc-500 italic">No showtimes available for this movie yet. Tell Admin to work! 😡</div>
      ) : (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
          {showtimes.map(showtime => (
            <button
              key={showtime.id}
              onClick={() => navigate(`/dashboard/booking/${showtime.id}`)}
              className="bg-zinc-900 hover:bg-zinc-800 rounded-xl p-6 border border-zinc-800 hover:border-yellow-500 transition-all text-left group relative overflow-hidden"
            >
              <div className="absolute inset-0 bg-gradient-to-r from-yellow-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
              
              <div className="relative z-10">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2 text-yellow-500">
                    <Calendar size={20} />
                    <span className="font-bold text-xl">{formatTime(showtime.start_time)}</span>
                  </div>
                  <span className="text-white font-bold bg-zinc-800 px-3 py-1 rounded-lg">
                    Rp {showtime.price.toLocaleString('id-ID')}
                  </span>
                </div>
                <p className="text-gray-400 text-sm flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-green-500"></span>
                  {showtime.studio.name}
                </p>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default ShowtimeSelection;