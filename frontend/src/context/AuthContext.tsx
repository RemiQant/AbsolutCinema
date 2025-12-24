import React, { createContext, useContext, useState, useEffect } from 'react';
import api from '../api/axios';

// Define what our User looks like (matching the backend response)
interface User {
  id: string;
  username: string;
  email: string;
  role: string;
}

// Define what the Context looks like
interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (userData: User) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Check if user is logged in when the app starts
  useEffect(() => {
    const checkAuth = async () => {
      try {
        // This endpoint requires the cookie we set earlier!
        const response = await api.get('/auth/me'); 
        setUser(response.data.user);
      } catch (err) {
        // If 401 Unauthorized, just set user to null (not logged in)
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = (userData: User) => {
    setUser(userData);
  };

  const logout = async () => {
    try {
      await api.post('/auth/logout');
      setUser(null);
    } catch (err) {
      console.error("Logout failed", err);
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook to use the context easily
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};