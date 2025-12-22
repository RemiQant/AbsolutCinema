import React from 'react';
import { Clock, ArrowLeft, Calendar } from 'lucide-react';
import { useParams, useNavigate, Link } from 'react-router-dom'; // New hooks
import { movies, showtimesData } from './mockData'; // Import data source

const ShowtimeSelection: React.FC = () => {
  const { movieId } = useParams(); // Get "1" from URL
  const navigate = useNavigate();

  // Find the movie (Simulating a DB fetch)
  // IMPORTANT: URL params are always strings, so we convert to Number
  const movie = movies.find(m => m.id === Number(movieId));
  const showtimes = showtimesData[Number(movieId)] || [];

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
  };

  if (!movie) {
    return <div className="text-white text-center mt-20">Movie not found! 😭</div>;
  }

  return (
    <div className="max-w-5xl mx-auto px-6 py-8 text-white">
      <button
        onClick={() => navigate('/dashboard')}
        className="flex items-center gap-2 text-gray-400 hover:text-yellow-500 mb-6 transition-colors cursor-pointer"
      >
        <ArrowLeft size={20} />
        Back to Movies
      </button>

      {/* ... (Keep Movie Info Header HTML) ... */}
       <div className="grid md:grid-cols-3 gap-6 mb-8">
            <div className="md:col-span-1">
              <img src={movie.poster_url} alt={movie.title} className="w-full rounded-lg border-2 border-yellow-600/20"/>
            </div>
            <div className="md:col-span-2">
               <h1 className="text-3xl font-bold text-yellow-500 mb-3">{movie.title}</h1>
               <p className="text-gray-300 leading-relaxed">{movie.description}</p>
            </div>
       </div>

      <h2 className="text-2xl font-bold text-yellow-500 mb-4">Select Showtime</h2>
      
      <div className="grid md:grid-cols-2 gap-4">
          {showtimes.map(showtime => (
            // WRAP IN LINK
            <Link
              to={`/dashboard/booking/${showtime.id}`}
              key={showtime.id}
              className="bg-zinc-900 hover:bg-zinc-800 rounded-lg p-6 border border-yellow-600/20 hover:border-yellow-500 transition-all text-left cursor-pointer"
            >
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2 text-yellow-500">
                  <Calendar size={18} />
                  <span className="font-bold text-lg">{formatTime(showtime.start_time)}</span>
                </div>
                <span className="text-yellow-500 font-bold">Rp {showtime.price.toLocaleString('id-ID')}</span>
              </div>
              <p className="text-gray-400 text-sm">{showtime.studio.name}</p>
            </Link>
          ))}
        </div>
    </div>
  );
};

export default ShowtimeSelection;