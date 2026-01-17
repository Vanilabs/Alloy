'use client'

import { Button } from '@/packages/ui/components/Button';
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react';
import { Mail, Lock, Eye, EyeOff } from "lucide-react";
import { cn } from "@/lib/utils";
// import { loginTenant } from '@/packages/auth/actions'
interface LoginPageProps {
    onLogin: (email: string, password: string) => Promise<void>;
}

export default function LoginPage() {
    const router = useRouter()

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setIsLoading(true);

        try {
            // await onLogin(email, password);
            await handleLogin();
        } catch (err) {
            setError(err instanceof Error ? err.message : "Login failed. Please try again.");
        } finally {
            setIsLoading(false);
        }
    };

    async function handleLogin() {
        await loginTenant()
        router.replace('/')
    }

    useEffect(() => {
        if (!error) return;

        const timer = setTimeout(() => {
            setError(null);
        }, 3000);

        return () => clearTimeout(timer);
    }, [error]);

    return (
        <div className="min-h-screen flex items-center justify-center bg-(--background) px-4">
            <div className="w-full max-w-md">
                <div className="text-center mb-8">
                    <div className="w-16 h-16 rounded-2xl flex items-center justify-center text-(--accent) font-bold text-2xl mx-auto mb-4"
                        style={{
                            background: 'linear-gradient(135deg, hsla(243, 75%, 59%, 0.1) 0%, hsla(243, 75%, 59%, 1) 50%)',
                            boxShadow: '0px 4px 8px 0px rgba(79, 70, 229, 0.1)'
                        }}
                    >
                        A
                    </div>
                    <h1 className="text-2xl font-bold text-foreground">Welcome back</h1>
                    <p className="text-muted-foreground mt-2">
                        Sign in to your organization account
                    </p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {error && (
                        <div className="p-3 rounded-lg bg-(--destructive-foreground)/10 text-(--destructive) text-sm animate-fade-in">
                            {error}
                        </div>
                    )}

                    <div className="space-y-2">
                        <label htmlFor="email" className="text-sm font-medium text-foreground">
                            Email
                        </label>
                        <div className="relative">
                            <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-(--muted-foreground)" />
                            <input
                                id="email"
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Enter your email"
                                required
                                className={cn(
                                    "w-full pl-10 pr-4 py-3 rounded-lg bg-(--muted)",
                                    "placeholder:text-(--muted-foreground) outline-none",
                                    "focus:outline-none focus:ring-2 focus:ring-(--primary) focus:text-(--primary) focus:border-transparent",
                                    "transition-all"
                                )}
                            />
                        </div>
                    </div>

                    {/* <div className="space-y-2">
                        <label htmlFor="password" className="text-sm font-medium text-foreground">
                            Password
                        </label>
                        <div className="relative">
                            <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-(--muted-foreground)" />
                            <input
                                id="password"
                                type={showPassword ? "text" : "password"}
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Enter your password"
                                required
                                className={cn(
                                    "w-full pl-10 pr-12 py-3 rounded-lg bg-(--muted) border border-input",
                                    "placeholder:text-(--muted-foreground)",
                                    "focus:outline-none focus:ring-2 focus:ring-(--primary) focus:text-(--primary) focus:border-transparent",
                                    "transition-all"
                                )}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-(--muted-foreground) hover:text-foreground transition-colors"
                            >
                                {showPassword ? (
                                    <EyeOff className="w-5 h-5" />
                                ) : (
                                    <Eye className="w-5 h-5" />
                                )}
                            </button>
                        </div>
                    </div> */}

                    {/* <div className="flex items-center justify-between">
                        <label className="flex items-center gap-2 cursor-pointer">
                            <input
                                type="checkbox"
                                className="w-4 h-4 rounded border-input text-(--primary) focus:ring-(--primary)"
                            />
                            <span className="text-sm text-(--muted-foreground)">Remember me</span>
                        </label>
                        <button
                            type="button"
                            className="text-sm text-(--primary) hover:underline"
                        >
                            Forgot password?
                        </button>
                    </div> */}

                    <Button
                        type="submit"
                        className="w-full py-3 text-base text-(--accent) font-medium bg-(--primary) hover:bg-(--primary-hover) rounded-lg transition-colors"
                        disabled={isLoading}
                    >
                        {isLoading ? (
                            <span className="flex items-center gap-2">
                                <span className="w-4 h-4 border-2 border-(--primary-foreground)/30 border-t-(--primary-foreground) rounded-full animate-spin" />
                                Signing in...
                            </span>
                        ) : (
                            "Sign in"
                        )}
                    </Button>
                </form>

                <p className="text-center text-sm text-(--muted-foreground) mt-6">
                    Don't have an account?{" "}
                    <span className="text-(--primary)">
                        Contact your organization admin for an invite.
                    </span>
                </p>
            </div>
        </div>
    )
}

function loginTenant() {
    throw new Error('Function not implemented.')
}

