import { useState } from "react";
import { User, X, Menu, LogOut, ShieldCheck } from "lucide-react"; // Tambah ShieldCheck
import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../../src/context/AuthContext"

const activePage = ({ isActive } : {isActive : boolean}) =>
  `transition-colors ${
    isActive
      ? "text-yellow-500"
      : "text-gray-200 hover:text-yellow-500"
  }`;

const Navbar = () => {
  const [open, setOpen] = useState(false);
  const {user, logout, loading} = useAuth(); // Rapikan pengambilan loading
  const navigate = useNavigate();

  const handleLogout = async () => {
      await logout();
      navigate('/login');
  }

  if (loading) {
    return <nav className="bg-black h-16 border-b-3 border-yellow-500"></nav>; 
  }

  return (
    <nav className="bg-black border-b-3 border-yellow-500">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">

          <div className="flex text-xl font-bold text-yellow-500 gap-4 items-center">
            AbsolutCinema <img src="/Logo/LogoAbsolutCinema.png" className="w-8 rounded" />
          </div>

          <div className="hidden md:flex space-x-6 items-center">
            <NavLink to="/dashboard" className={activePage}>
              Home
            </NavLink>

            <NavLink to="/about" className={activePage}>
              About
            </NavLink>

            {user && user.role === 'admin' && (
              <NavLink to="/admin" className={activePage}>
                <div className="flex items-center gap-1">
                  <ShieldCheck size={18} />
                  Admin
                </div>
              </NavLink>
            )}

            {user ? (
                <>
                    <NavLink
                      to="/account"
                      className={({ isActive }) =>
                        `flex items-center gap-2 transition-colors ${
                          isActive
                            ? "text-yellow-500"
                            : "text-gray-200 hover:text-yellow-500"
                        }`
                      }
                    >
                        <User size={20} />
                    <span>
                        Hi, {user.username}!
                    </span>
                    </NavLink>
                    
                    <button onClick={handleLogout} className="text-gray-200 hover:text-red-500 transition-colors">
                        <LogOut size={20} />
                    </button>
                </>
            ) : (
                <NavLink to="/login" className="text-gray-200 hover:text-yellow-500">
                    Login
                </NavLink>
            )}
          </div>

          {/* Hamburger Menu untuk Mobile */}
          <button
            onClick={() => setOpen(!open)}
            className="md:hidden text-gray-200 transition duration-300 hover:text-yellow-500"
          >
            {open ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>
      </div>

      {/* MOBILE MENU */}
      <div
        className={`
          md:hidden overflow-hidden px-4
          transition-all duration-300 ease-in-out
          ${open ? "max-h-64 opacity-100 pb-4" : "max-h-0 opacity-0"}
        `}
      >
        <NavLink to="/dashboard" className={({ isActive }) => `block py-2 ${activePage({ isActive })}`}>
          Home
        </NavLink>

        <NavLink to="/about" className={({ isActive }) => `block py-2 ${activePage({ isActive })}`}>
          About Us
        </NavLink>

        {/* ADMIN LINK MOBILE */}
        {user && user.role === 'admin' && (
          <NavLink to="/admin" className={({ isActive }) => `block py-2 ${activePage({ isActive })}`}>
            Admin Panel
          </NavLink>
        )}

        {user && (
          <NavLink
            to="/account"
            className={({ isActive }) =>
              `flex items-center gap-2 py-2 transition-colors ${
                isActive ? "text-yellow-500" : "text-gray-200 hover:text-yellow-500"
              }`
            }
          >
            <User size={20} />
            Account
          </NavLink>
        )}
      </div>
    </nav>
  );
};

export default Navbar;