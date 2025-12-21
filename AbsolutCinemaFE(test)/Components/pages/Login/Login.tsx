import React from "react";
import { NavLink  , useNavigate} from "react-router-dom";
import { useState , useEffect} from "react";

const Landing = () => {

  const [user , setUser] = useState<string>("");
  const [pwd , setPwd] = useState<string>("");

  const navigate = useNavigate();

  return (
    <div
      className="relative min-h-screen bg-cover bg-center"
      style={{
        backgroundImage:
          "url('https://images.unsplash.com/photo-1517604931442-7e0c8ed2963c?auto=format&fit=crop&w=1920')",
      }}
    >
      <div className="absolute inset-0 bg-black/80"></div>

      <div className="absolute top-8 left-10 z-10">
        <h1 className="flex gap-4 text-3xl font-extrabold tracking-wide text-yellow-500">
          AbsolutCinema  <img src="../public/Logo/LogoAbsolutCinema.png" className="w-12 rounded" />
        </h1>
      </div>

      <div className="relative z-10 min-h-screen flex items-center justify-center px-4">
        <div className="w-full max-w-md rounded-2xl backdrop-blur-md border border-white/20 p-10 shadow-2xl">
          <h1 className="text-3xl font-bold text-center text-yellow-500 mb-2">
            Welcome Back
          </h1>
          <p className="text-center text-gray-300 mb-8">
            Login to continue your cinematic journey
          </p>

          <form className="space-y-5">
            <div>
              <label className="block text-sm text-gray-300 mb-1">
                Email
              </label>
              <input
                type="email"
                placeholder="you@example.com"
                className="w-full bg-black/40 text-white border border-gray-600 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-yellow-500"
                onChange = {(e) => setUser(e.target.value)}
                value = {user}
                required
              />
            </div>

            <div>
              <label className="block text-sm text-gray-300 mb-1">
                Password
              </label>
              <input
                type="password"
                placeholder="••••••••"
                className="w-full bg-black/40 text-white border border-gray-600 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-yellow-500"
                onChange = {(e) => setPwd(e.target.value)}
                value = {pwd}
                required
              />
            </div>

            <button
              type="submit"
              className="w-full bg-yellow-600 text-black font-semibold py-2.5 rounded-md hover:bg-yellow-500 transition cursor-pointer"
              onClick = {() => navigate('/dashboard')}
            >
              Login
            </button>
          </form>
          <p className = "text-center mt-5 text-gray-400">Don't have an account yet? <NavLink to = "/signup" className = "text-yellow-500">Sign Up </NavLink></p>
          <p className="text-center text-sm text-gray-400 mt-8">
            © 2025 AbsolutCinema. All rights reserved.
          </p>
        </div>
      </div>
    </div>
  );
};

export default Landing;
