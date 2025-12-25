import React from 'react';
import { Users, Award, Star, ShieldCheck, Armchair, Volume2 } from 'lucide-react';
import Navbar from '../../Navbar/Navbar';
import Footer from '../../Footer/Footer';

const About: React.FC = () => {
  return (
    <>
      <Navbar />
      <div className="bg-black min-h-screen text-white font-sans selection:bg-yellow-500/30">
        
        {/* Hero Section */}
        <div className="py-24 md:py-32 text-center border-b border-zinc-900">
          <h1 className="text-5xl md:text-8xl font-black tracking-tighter text-yellow-500 uppercase mb-4 italic">
            AbsolutCinema
          </h1>
          <p className="text-zinc-500 tracking-[0.4em] uppercase text-xs font-semibold ">
            The Absolute Standard of Cinematic Excellence
          </p>
        </div>

        <div className="max-w-5xl mx-auto px-6 py-20">
          
          {/* Main Story Section */}
          <div className="grid md:grid-cols-2 gap-16 items-center mb-32">
            <div className="space-y-6">
              <div className="inline-block px-3 py-1 bg-yellow-500/10 border border-yellow-500/20 rounded-full">
                <span className="text-yellow-500 text-xs font-bold uppercase tracking-widest">Since 2025</span>
              </div>
              <h2 className="text-4xl font-bold leading-tight">
                Menghadirkan Kembali <br />
                <span className="text-yellow-500">Esensi Menonton Film.</span>
              </h2>
              <p className="text-zinc-400 leading-relaxed text-lg text-justify">
                Absolut Cinema hadir untuk menghidupkan kembali pengalaman menonton Anda.
                Kami memadukan kenyamanan eksklusif dengan teknologi proyeksi dan sound system terkini
                guna menghadirkan pengalaman menonton yang berkelas, imersif, dan tak terlupakan
                bagi setiap pecinta film.
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="bg-zinc-900/40 p-6 rounded-2xl border border-zinc-800 text-center shadow-[0_0_15px_rgba(234,179,8,0.05)] hover:shadow-[0_0_20px_rgba(234,179,8,0.15)] transition-shadow">
                <Volume2 className="text-yellow-500 mx-auto mb-3" size={24} />
                <p className="text-2xl font-bold uppercase tracking-tighter">DOLBY</p>
                <p className="text-zinc-500 text-xs uppercase font-medium">Atmos Sound</p>
              </div>
              <div className="bg-zinc-900/40 p-6 rounded-2xl border border-zinc-800 text-center shadow-[0_0_15px_rgba(234,179,8,0.05)] hover:shadow-[0_0_20px_rgba(234,179,8,0.15)] transition-shadow">
                <Star className="text-yellow-500 mx-auto mb-3" size={24} />
                <p className="text-2xl font-bold">4K</p>
                <p className="text-zinc-500 text-xs uppercase font-medium">Resolution</p>
              </div>
              <div className="bg-zinc-900/40 p-6 rounded-2xl border border-zinc-800 text-center col-span-2 shadow-[0_0_15px_rgba(234,179,8,0.05)]">
                <ShieldCheck className="text-yellow-500 mx-auto mb-2" size={24} />
                <p className="text-zinc-400 text-sm font-medium">Certified Premium Experience</p>
              </div>
            </div>
          </div>

          <div className="bg-zinc-900/20 p-12 rounded-[2rem] border border-zinc-800/50 shadow-[0_0_50px_rgba(234,179,8,0.03)]">
            <div className="text-center mb-16">
              <h3 className="text-yellow-500 uppercase tracking-[0.3em] text-xs font-bold">Absolute Standard</h3>
            </div>
            
            <div className="grid md:grid-cols-3 gap-12">
              {[
                { 
                  icon: <Award size={32} />, 
                  title: "4K Laser Technology", 
                  desc: "Ketajaman gambar tak tertandingi dengan reproduksi warna yang sangat akurat di setiap frame." 
                },
                { 
                  icon: <Volume2 size={32} />, 
                  title: "Dolby Atmos", 
                  desc: "Sistem audio yang mengelilingi Anda, memberikan dimensi suara yang terasa nyata." 
                },
                { 
                  icon: <Armchair size={32} />, 
                  title: "Ultimate Comfort", 
                  desc: "Kursi premium yang dirancang khusus untuk kenyamanan menonton durasi panjang." 
                }
              ].map((item, idx) => (
                <div key={idx} className="group relative p-6 rounded-2xl hover:bg-zinc-800/30 transition-all duration-500 hover:shadow-[0_10px_30px_rgba(234,179,8,0.05)]">
                  <div className="text-yellow-500 mb-6 group-hover:scale-110 transition-transform duration-500">
                    {item.icon}
                  </div>
                  <h4 className="font-bold text-xl mb-4 text-zinc-100">{item.title}</h4>
                  <p className="text-zinc-500 text-sm leading-relaxed">{item.desc}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Quote Section */}
          <div className="mt-32 text-center">
             <div className="w-12 h-1 bg-yellow-500 mx-auto mb-8 opacity-50"></div>
             <p className="text-zinc-500 italic font-serif text-2xl max-w-2xl mx-auto leading-relaxed">
               "Cinema is a matter of what's in the frame and what's out." 
             </p>
             <p className="mt-4 text-zinc-700 text-xs uppercase tracking-widest">— Martin Scorsese</p>
          </div>

        </div>
      </div>
      <Footer />
    </>
  );
};

export default About;