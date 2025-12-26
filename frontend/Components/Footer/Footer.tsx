import React from 'react'
import { Mail, Phone, MapPin } from 'lucide-react';

const Footer = () => {
  return (
    <footer className="bg-black border-t border-zinc-900">
      <div className="max-w-5xl mx-auto px-6 py-16">
        <div className="grid md:grid-cols-3 gap-12 text-center md:text-left">
          
          <div className="space-y-4 group">
            <div className="flex items-center gap-3 justify-center md:justify-start text-yellow-500">
              <MapPin size={20} className="group-hover:scale-110 transition-transform" />
              <span className="font-bold uppercase text-xs tracking-widest">Address</span>
            </div>
            <p className="text-gray-400 text-sm leading-relaxed">
              Jl. Cinema Boulevard No. 123,<br /> 
              Jakarta Selatan, Indonesia
            </p>
          </div>

          {/* Email */}
          <div className="space-y-4 group">
            <div className="flex items-center gap-3 justify-center md:justify-start text-yellow-500">
              <Mail size={20} className="group-hover:scale-110 transition-transform" />
              <span className="font-bold uppercase text-xs tracking-widest">Email</span>
            </div>
            <p className="text-gray-400 text-sm">hello@absolutcinema.com</p>
          </div>

          {/* Telepon */}
          <div className="space-y-4 group">
            <div className="flex items-center gap-3 justify-center md:justify-start text-yellow-500">
              <Phone size={20} className="group-hover:scale-110 transition-transform" />
              <span className="font-bold uppercase text-xs tracking-widest">Phone</span>
            </div>
            <p className="text-gray-400 text-sm">+62 21 888 999 10</p>
          </div>

        </div>

        <div className="mt-16 pt-8 border-t border-zinc-900/50 text-center">
          <p className="text-zinc-600 text-[10px] uppercase tracking-[0.2em]">
            © 2025 AbsolutCinema. All Rights Reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}

export default Footer