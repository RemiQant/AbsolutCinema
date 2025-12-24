import type { Movie, Showtime, Ticket } from './types';

export const movies: Movie[] = [
  {
    id: 1,
    title: "Interstellar",
    description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
    duration_minutes: 169,
    poster_url: "https://images.unsplash.com/photo-1536440136628-849c177e76a1?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 2,
    title: "Inception",
    description: "A thief who steals corporate secrets through dream-sharing technology is given the inverse task.",
    duration_minutes: 148,
    poster_url: "https://images.unsplash.com/photo-1440404653325-ab127d49abc1?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 3,
    title: "The Dark Knight",
    description: "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham.",
    duration_minutes: 152,
    poster_url: "https://images.unsplash.com/photo-1509347528160-9a9e33742cdb?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 4,
    title: "Avatar",
    description: "A paraplegic Marine dispatched to the moon Pandora on a unique mission.",
    duration_minutes: 162,
    poster_url: "https://images.unsplash.com/photo-1518676590629-3dcbd9c5a5c9?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 5,
    title: "Tenet",
    description: "Armed with only one word, Tenet, and fighting for the survival of the entire world.",
    duration_minutes: 150,
    poster_url: "https://images.unsplash.com/photo-1478720568477-152d9b164e26?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 6,
    title: "Forrest Gump",
    description: "The presidencies of Kennedy and Johnson, the Vietnam War, and other historical events unfold.",
    duration_minutes: 142,
    poster_url: "https://images.unsplash.com/photo-1485846234645-a62644f84728?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 7,
    title: "Dune",
    description: "Feature adaptation of Frank Herbert's science fiction novel about the son of a noble family.",
    duration_minutes: 155,
    poster_url: "https://images.unsplash.com/photo-1594908900066-3f47337549d8?w=300&h=450&fit=crop",
    rating: "PG-13"
  },
  {
    id: 8,
    title: "Parasite",
    description: "Greed and class discrimination threaten the newly formed symbiotic relationship.",
    duration_minutes: 132,
    poster_url: "https://images.unsplash.com/photo-1574267432644-f610fa866e72?w=300&h=450&fit=crop",
    rating: "R"
  }
];

export const showtimesData: Record<number, Showtime[]> = {
  1: [
    { id: 1, movie_id: 1, studio_id: 1, start_time: "2025-12-10T14:00:00", price: 50000, studio: { id: 1, name: "Studio 1", total_rows: 8, total_cols: 10 } },
    { id: 2, movie_id: 1, studio_id: 2, start_time: "2025-12-10T17:00:00", price: 60000, studio: { id: 2, name: "Studio 2", total_rows: 6, total_cols: 8 } },
    { id: 3, movie_id: 1, studio_id: 1, start_time: "2025-12-10T20:00:00", price: 75000, studio: { id: 1, name: "Studio 1", total_rows: 8, total_cols: 10 } }
  ],
  2: [
    { id: 4, movie_id: 2, studio_id: 1, start_time: "2025-12-10T13:30:00", price: 50000, studio: { id: 1, name: "Studio 1", total_rows: 8, total_cols: 10 } },
    { id: 5, movie_id: 2, studio_id: 3, start_time: "2025-12-10T19:00:00", price: 65000, studio: { id: 3, name: "Studio 3", total_rows: 10, total_cols: 12 } }
  ],
  3: [
    { id: 6, movie_id: 3, studio_id: 2, start_time: "2025-12-10T15:00:00", price: 55000, studio: { id: 2, name: "Studio 2", total_rows: 6, total_cols: 8 } }
  ],
  4: [
    { id: 7, movie_id: 4, studio_id: 1, start_time: "2025-12-10T18:30:00", price: 70000, studio: { id: 1, name: "Studio 1", total_rows: 8, total_cols: 10 } }
  ]
};

export const bookedTickets: Record<number, Ticket[]> = {
  1: [
    { showtime_id: 1, seat_number: "A3" },
    { showtime_id: 1, seat_number: "A4" },
    { showtime_id: 1, seat_number: "C5" }
  ]
};