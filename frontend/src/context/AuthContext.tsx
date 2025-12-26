import React, { createContext, useContext, useState, useEffect } from 'react';
import api from '../api/axios';

interface User {
  id: string;
  username: string;
  email: string;
  role: string;
}

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

  // Fungsi utama untuk memverifikasi session
  const verifyAuth = async () => {
    try {
      // 1. Coba ambil data user (Backend validasi access token)
      const response = await api.get('/auth/me');
      setUser(response.data.user);
    } catch (err: any) {
      // 2. ELSE (401): Jika access token tidak valid/expired
      if (err.response?.status === 401) {
        try {
          // 3. Frontend triggers refresh flow
          await api.post('/auth/refresh');
          
          // 4. Retry original request (Ambil data user lagi setelah refresh)
          const retryResponse = await api.get('/auth/me');
          setUser(retryResponse.data.user);
        } catch (refreshErr) {
          // Jika refresh token juga expired/gagal
          setUser(null);
        }
      } else {
        setUser(null);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    verifyAuth();
  }, []);

  const login = (userData: User) => {
    setUser(userData);
  };

  const logout = async () => {
    try {
      await api.post('/auth/logout');
    } catch (err) {
      console.error("Logout failed", err);
    } finally {
      setUser(null);
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {/* Tahan render UI sampai verifikasi selesai agar tidak ada 'flicker' */}
      {!loading ? (
        children
      ) : (
        <div className="flex h-screen items-center justify-center bg-black text-yellow-500">
          <p className="animate-pulse font-bold">VERIFYING SESSION...</p>
        </div>
      )}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within an AuthProvider');
  return context;
};