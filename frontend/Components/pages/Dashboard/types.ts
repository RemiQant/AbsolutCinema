export interface Movie {
  id: number;
  title: string;
  description: string;
  duration_minutes: number;
  poster_url: string;
  rating: string;
}

export interface Studio {
  id: number;
  name: string;
  total_rows: number;
  total_cols: number;
}

export interface Showtime {
  id: number;
  movie_id: number;
  studio_id: number;
  start_time: string;
  price: number;
  studio: Studio;
}

export interface Seat {
  id: string;
  row: string;
  number: number;
  status: 'available' | 'selected' | 'booked';
}

export interface Ticket {
  showtime_id: number;
  seat_number: string;
}