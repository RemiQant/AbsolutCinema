import { useState } from "react";
import { User, X, Menu } from "lucide-react";
import { NavLink } from "react-router-dom";



const activePage = ({ isActive } : {isActive : boolean}) =>
  `transition-colors ${
    isActive
      ? "text-yellow-500"
      : "text-gray-200 hover:text-yellow-500"
  }`;

const Navbar = () => {
  const [open, setOpen] = useState(false);

  return (
    <nav className="bg-black border-b-3 border-yellow-500">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">

          <div className="flex text-xl font-bold text-yellow-500 gap-4">
            AbsolutCinema <img src="../public/Logo/LogoAbsolutCinema.png" className="w-8 rounded" />
          </div>

          <div className="hidden md:flex space-x-6">
            <NavLink to="/dashboard" className= "text-gray-200 hover:text-yellow-500 active:text-yellow-500">
              Home
            </NavLink>

            <NavLink to="/about" className={activePage}>
              About
            </NavLink>

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
              Account
            </NavLink>
          </div>

          <button
            onClick={() => setOpen(!open)}
            className="md:hidden text-gray-200 transition duration-300 hover:text-yellow-500"
          >
            {open ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>
      </div>

      <div
        className={`
          md:hidden overflow-hidden px-4
          transition-all duration-300 ease-in-out
          ${open ? "max-h-40 opacity-100 pb-4" : "max-h-0 opacity-0"}
        `}
      >
        <NavLink
            to="/"
            className={({ isActive }) =>
                `block py-2 ${activePage({ isActive })}`
            }
        >
        Home
        </NavLink>

        <NavLink
            to="/about"
            className={({ isActive }) =>
                `block py-2 ${activePage({ isActive })}`
            }
        >
        About Us
        </NavLink>

        <NavLink
          to="/account"
          className={({ isActive }) =>
            `flex items-center gap-2 py-2 transition-colors ${
              isActive
                ? "text-yellow-500"
                : "text-gray-200 hover:text-yellow-500"
            }`
          }
        >
          <User size={20} />
          Account
        </NavLink>
      </div>
    </nav>
  );
};

export default Navbar;
