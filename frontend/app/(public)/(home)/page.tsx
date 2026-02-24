"use client"
import { useState } from 'react';
import { ArrowRight, Calendar, ChartNoAxesCombined, Home, Play, User } from 'lucide-react';
import { toast } from '@/packages/ui/use-toast';
// import { showAlert } from '@/store/alertSlice';
// import { useAppDispatch } from '@/packages/store/hooks';

export default function LandingPage() {
    const [email, setEmail] = useState('');
    // const dispatch = useAppDispatch();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (email) {
            toast({
                title: 'Request Submitted!',
                status: 'success',
                description: "We'll be in touch soon.",
                duration: 4000,
            });
            setEmail('');
        } else {
            toast({
                title: 'Email Required',
                status: 'error',
                description: 'Please enter your email address.',
                duration: 4000,
            });
        }
    };

    return (
        <div className="min-h-screen h-screen bg-[#0D0D12] relative overflow-hidden">
            {/* Gradient Orbs */}
            {/* Top Right Gradient */}
            <div
                className="absolute -top-45 -right-45 w-90 h-90 rounded-full opacity-50 blur-2xl"
                style={{
                    background:
                        "radial-gradient(circle, hsla(235, 88%, 67%) 0%)",
                }}
            />
            {/* <div
                className="absolute top-0 left-1/4 w-150 h-150 rounded-full opacity-30 blur-[120px]"
                style={{
                    background: 'radial-gradient(circle, hsl(239, 84%, 67%) 0%, transparent 70%)',
                }}
            /> */}

            {/* Bottom Left Gradient */}
            <div
                className="absolute -bottom-40 -left-40 w-130 h-105 rounded-full opacity-70 blur-[90px]"
                style={{
                    background:
                        "radial-gradient(circle, hsla(235, 88%, 67%) 0%)",
                }}
            />
            {/* <div
                className="absolute bottom-6 right-1/4 w-150 h-125 rounded-full opacity-20 blur-[100px]"
                style={{
                    background: 'radial-gradient(circle, hsl(280, 87%, 65%) 0%, transparent 70%)',
                }}
            /> */}

            {/* Header */}
            <header className="relative z-10 flex items-center justify-between px-8 py-6 max-w-7xl mx-auto">
                <div className="flex items-center gap-3">
                    <div
                        className="w-9 h-9 rounded-xl flex items-center justify-center text-white font-bold text-base shadow-lg"
                        style={{
                            background: 'linear-gradient(135deg, hsl(239, 84%, 67%) 0%, hsl(280, 87%, 65%) 100%)',
                            boxShadow: '0px 4px 8px 0px rgba(79, 70, 229, 0.1)'
                        }}
                    >
                        A
                    </div>
                    <span className="text-white font-semibold text-lg">Alloy</span>
                </div>

                <button
                    onClick={() => location.replace('/')}
                    className="flex items-center gap-2 px-5 py-2.5 bg-black hover:bg-black/50 rounded-lg text-white text-sm font-medium transition-colors backdrop-blur-sm"
                >
                    <Play className="w-4 h-4 fill-current text-black bg-amber-50 rounded-lg p-0.5" />
                    View Demo
                </button>
            </header>

            {/* Hero Section */}
            <main className="relative z-10 flex flex-col items-center justify-center px-8 pt-20 pb-32">
                <h1 className="text-5xl md:text-6xl lg:text-7xl font-bold text-white text-center max-w-4xl leading-tight">
                    All your essentials,
                    <br />
                    unified in one app.
                </h1>

                <p className="mt-6 text-lg text-gray-400 text-center max-w-2xl">
                    Designed to remove complexity, save time, and give you a smoother, more focused experience from start to finish.
                </p>

                {/* Email Input */}
                <form onSubmit={handleSubmit} className="mt-12 w-full max-w-md">
                    <div className="relative group">
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="your email"
                            className="w-full bg-transparent border-b border-gray-600 focus:border-primary py-4 pr-12 text-white placeholder:text-gray-500 text-lg outline-none transition-colors"
                        />
                        <button
                            type="submit"
                            className="absolute right-0 top-1/2 -translate-y-1/2 p-2 text-gray-400 hover:text-white transition-colors"
                        >
                            <ArrowRight className="w-6 h-6" />
                        </button>
                    </div>
                </form>

                {/* App Preview */}
                <div className="mt-20 w-full max-w-5xl">
                    <div
                        className="relative rounded-2xl overflow-hidden shadow-2xl"
                        style={{
                            background: 'linear-gradient(180deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.02) 100%)',
                            border: '1px solid rgba(255,255,255,0.1)',
                        }}
                    >
                        {/* Mock App Window */}
                        <div className="bg-[#1A1A24] rounded-2xl p-1">
                            <div className="flex items-center gap-2 px-4 py-3 border-b border-white/5">
                                <div className="flex gap-2">
                                    <div className="w-3 h-3 rounded-full bg-red-500/80" />
                                    <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
                                    <div className="w-3 h-3 rounded-full bg-green-500/80" />
                                </div>
                            </div>

                            <div className="flex h-100">
                                {/* Sidebar Preview */}
                                <div className="w-16 bg-[#0D0D12] border-r border-white/5 flex flex-col items-center py-4">
                                    <div
                                        className="w-8 h-8 rounded-lg flex items-center justify-center text-white font-bold text-sm mb-4"
                                        style={{
                                            background: 'linear-gradient(135deg, hsl(239, 84%, 67%) 0%, hsl(280, 87%, 65%) 100%)',
                                        }}
                                    >
                                        A
                                    </div>
                                    <div className="space-y-4 mt-4">
                                        {[{i:1, icon: Home}, {i:2, icon: ChartNoAxesCombined}, {i:3, icon: User}, {i:4, icon: Play}, {i:5, icon: Calendar}].map((i) => (
                                            <div key={i.i} className={`flex justify-center align-middle p-1.5 rounded-full text-white ${i.i === 1 ? 'bg-white/10' : 'hover:bg-white/5 cursor-pointer'}`}>
                                                <i.icon className="w-4 h-5" />
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                {/* Conversation List Preview */}
                                <div className="w-64 bg-[#12121A] border-r border-white/5 p-4">
                                    <div className="flex items-center justify-between mb-4">
                                        <span className="text-white/80 text-sm font-medium">Messages</span>
                                        <div className="w-6 h-6 rounded-full bg-(--primary) flex items-center justify-center">
                                            <span className="text-xs text-white">+</span>
                                        </div>
                                    </div>
                                    <div className="space-y-3">
                                        {[1, 2, 3, 4].map((i) => (
                                            <div key={i} className="flex items-center gap-3 p-2 rounded-lg bg-white/5">
                                                <div className="w-10 h-10 rounded-full bg-white/10" />
                                                <div className="flex-1">
                                                    <div className="h-3 w-20 bg-white/20 rounded" />
                                                    <div className="h-2 w-24 bg-white/10 rounded mt-1.5" />
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                {/* Chat Preview */}
                                <div className="flex-1 bg-[#0D0D12] p-6 flex flex-col">
                                    <div className="flex-1 space-y-4">
                                        <div className="flex gap-3">
                                            <div className="w-8 h-8 rounded-full bg-white/10" />
                                            <div className="space-y-2">
                                                <div className="px-4 py-2 rounded-2xl bg-[#2A2A3A] text-white/80 text-sm">
                                                    omg, this is amazing
                                                </div>
                                                <div className="px-4 py-2 rounded-2xl bg-[#2A2A3A] text-white/80 text-sm w-fit">
                                                    perfect! ✅
                                                </div>
                                            </div>
                                        </div>
                                        <div className="flex justify-end">
                                            <div className="px-4 py-2 rounded-2xl bg-primary text-white text-sm">
                                                How are you?
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};
