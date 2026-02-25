'use client'

import { Button } from '@/packages/ui/components/Button';
import { useRouter } from 'next/navigation'
import { useEffect, useMemo, useRef, useState } from 'react';
import { Mail, Lock, Eye, EyeOff } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppDispatch, useAppSelector } from '@/packages/store/hooks';
import { requestMagicThunk } from '@/packages/store/thunks';
import { clearAuthTransient, clearError } from '@/packages/store/slices/auth.slice';
import { useToast } from '@/packages/ui/use-toast';
// import { loginTenant } from '@/packages/auth/actions'
interface LoginPageProps {
    onLogin: (email: string, password: string) => Promise<void>;
}

export default function LoginPage() {
    const router = useRouter()
    const dispatch = useAppDispatch();
    // const { toast } = useToast();

    const { loading, error, magicLink } = useAppSelector(s => s.auth)

    const [email, setEmail] = useState("");
    const [localError, setLocalError] = useState<string | null>(null);

    const canSubmit = useMemo(() => {
        const trimmed = email.trim();
        return trimmed.length > 3 && trimmed.includes('@')
    }, [email]);

    async function submit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        setLocalError(null);

        // clear any previous timers
        if (clearLocalErrorTimerRef.current) {
            window.clearTimeout(clearLocalErrorTimerRef.current);
            clearLocalErrorTimerRef.current = null;
        }

        setLocalError(null);
        dispatch(clearError());

        if (!canSubmit) {
            setLocalError('Please enter a valid email.')
            clearLocalErrorTimerRef.current = window.setTimeout(() => {
                setLocalError(null);
            }, 3000);
            return
        }

        const res = await dispatch(requestMagicThunk({ email: email.trim() }))

        if (requestMagicThunk.fulfilled.match(res)) {
            // Stay on page and show success state
            // NOTE: Do not reveal if email exists
        }
    }

    const clearLocalErrorTimerRef = useRef<number | null>(null);

    useEffect(() => {
        if (error) dispatch(clearError());
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [email]);

    useEffect(() => {
        return () => {
            if (clearLocalErrorTimerRef.current) {
                window.clearTimeout(clearLocalErrorTimerRef.current);
            }
        };
    }, []);

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

                <form onSubmit={submit} className="space-y-4">
                    {(localError || error) && (
                        <div className="p-3 rounded-lg bg-(--destructive-foreground)/10 text-(--destructive) text-sm animate-fade-in">
                            {localError || error}
                        </div>
                    )}

                    {/* Success state (generic, enterprise safe) */}
                    {magicLink.requested && (
                        <div className="mb-4 p-3 rounded-lg bg-emerald-50 text-emerald-700 text-sm">
                            If an account exists for <span className="font-medium">{magicLink.emailMasked}</span>,
                            you&apos;ll receive a sign-in link shortly.
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
                        disabled={loading || !canSubmit}
                    >
                        {loading ? (
                            <span className="flex items-center gap-2">
                                <span className="w-4 h-4 border-2 border-(--primary-foreground)/30 border-t-(--primary-foreground) rounded-full animate-spin" />
                                Signing in...
                            </span>
                        ) : (
                            "Sign in"
                        )}
                    </Button>

                    {/* resend */}
                    {magicLink.requested && (
                        <div className="mt-4 text-center text-sm text-(--muted-foreground)">
                            Didn&apos;t get it?{' '}
                                <button
                                type='submit'
                                className="text-(--primary) hover:underline"
                                onClick={(e) => {
                                    e.preventDefault();
                                    submit(e as any)
                                }}
                                // onClick={submit}
                                disabled={loading}
                            >
                                Resend
                            </button>
                        </div>
                    )}
                </form>

                <p className="text-center text-sm text-(--muted-foreground) mt-6">
                    Don't have an account?{" "}
                    <span className="text-(--primary)">
                        Contact your organization admin for an invite.
                    </span>
                </p>

                <p className="text-center text-sm text-(--muted-foreground) mt-6">
                    For security, the link expires quickly and can be used only once.
                </p>
            </div>
        </div>
    )
}
