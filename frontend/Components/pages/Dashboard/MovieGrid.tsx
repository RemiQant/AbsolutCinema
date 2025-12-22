import React from 'react';
import { Clock } from 'lucide-react';
import { Link } from 'react-router-dom'; // Import Link
import { movies } from './mockData'; // Import data directly here now

const MovieGrid: React.FC = () => {
  // No more props needed! We fetch data directly (or from API later)
  
  return (
    <div className="max-w-7xl mx-auto px-6 py-8">
      {/* ... (Keep your Hero Section HTML here) ... */}

      <h2 className="text-yellow-500 text-2xl font-bold mb-6 tracking-wide">CURRENTLY SHOWING</h2>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {movies.map(movie  => (
          // Change div to Link
          <Link
            to={`/dashboard/movie/${movie.id}`} 
            key={movie.id}
            className="group cursor-pointer"
          >
            {/* ... (Keep exactly the same card inner HTML) ... */}
            <div className="relative rounded-lg overflow-hidden mb-3 border-2 border-transparent hover:border-yellow-500 transition-all">
                {/* ... existing image/rating code ... */}
                <img src={movie.poster_url} alt={movie.title} className="w-full h-96 object-cover group-hover:scale-105 transition-transform duration-300" />
            </div>
            <h3 className="text-white font-bold text-center mb-2 group-hover:text-yellow-500 transition-colors">{movie.title}</h3>
             {/* ... */}
          </Link>
        ))}
      </div>
    </div>
  );
};

export default MovieGrid;