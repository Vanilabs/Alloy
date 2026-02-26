'use client'

import { useEffect, useMemo, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { cn } from '@/lib/utils'
import { Button } from '@/packages/ui/components/Button'
import { Mail, Lock, Eye, EyeOff, XCircle, CheckCircle2 } from "lucide-react";
import { useAppDispatch, useAppSelector } from '@/packages/store/hooks'
import { verifyMagicLinkThunk } from '@/packages/store/thunks'
// import { loginTenant } from '@/packages/auth/actions'
interface LoginPageProps {
    onLogin: (email: string, password: string) => Promise<void>;
}

export default function MagicLinkVerifyPage() {
    const router = useRouter()
    const searchParams = useSearchParams()
    const dispatch = useAppDispatch()

    const token = useMemo(() => searchParams.get('token') ?? '', [searchParams]);

    const { loading, error } = useAppSelector(s => s.auth);

    const [localError, setLocalError] = useState<string | null>(null);
    const [status, setStatus] = useState<'idle' | 'verifying' | 'success' | 'error'>('idle');

    async function verify() {
        setLocalError(null)

        if (!token || token.length < 10) {
            setStatus('error')
            setLocalError('Invalid or missing verification token. Please request a new magic link.')
            return
        }

        setStatus('verifying')

        setTimeout(() => {
            setStatus('verifying')
        }, 890000099090000)

        const res = await dispatch(verifyMagicLinkThunk({ token }))

        if (verifyMagicLinkThunk.fulfilled.match(res)) {
            setStatus('success')

            // Enterprise redirect strategy:
            // - If tenant is already selected via cookie -> go to '/'
            // - Else -> go to '/select-tenant'
            // We'll attempt '/select-tenant' first; middleware will redirect to '/' if tenant cookie exists.
            // setTimeout(() => {
            //     router.replace('/select-tenant')
            // }, 1200000900)

            return
        }

        setStatus('error')
        setLocalError(
            (res as any)?.error?.message ||
            'Verification failed. The link may have expired or already been used.'
        )
    }

    useEffect(() => {
        verify()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [token])

    useEffect(() => {
        if (!error) return
        const t = setTimeout(() => {
            // optional: dispatch(clearError())
        }, 4000)
        return () => clearTimeout(t)
    }, [error])

    const displayError = error || localError || null

    return (
        <div className="min-h-screen flex items-center justify-center bg-(--background) px-4">
            <div className="w-full max-w-md">
                <div className="text-center mb-8">
                    <div
                        className="w-16 h-16 rounded-2xl flex items-center justify-center text-(--accent) font-bold text-2xl mx-auto mb-4"
                        style={{
                            background:
                                'linear-gradient(135deg, hsla(243, 75%, 59%, 0.1) 0%, hsla(243, 75%, 59%, 1) 50%)',
                            boxShadow: '0px 4px 8px 0px rgba(79, 70, 229, 0.1)',
                        }}
                    >
                        A
                    </div>

                    <h1 className="text-2xl font-bold text-foreground">Verifying magic link</h1>
                    <p className="text-muted-foreground mt-2">
                        We&apos;re securely signing you into your organization workspace.
                    </p>
                </div>

                <div className="rounded-b-sm p-6 bg-white/60 backdrop-blur shadow-md">
                    {/* ERROR */}
                    {status === 'error' && displayError && (
                        <div className="mb-4 p-3 rounded-lg bg-(--destructive-foreground)/10 text-(--destructive) text-sm">
                            <div className="flex items-start gap-2">
                                <XCircle className="mt-0.5 h-5 w-5" />
                                <div>
                                    <div className="font-medium">Verification failed</div>
                                    <div className="mt-1">{displayError}</div>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* BIG LOADING */}
                    {status === 'verifying' && (
                        <div className="flex flex-col items-center justify-center py-10">
                            <div
                                className={cn(
                                    'h-20 w-20 rounded-full border-4',
                                    'border-(--primary)/20 border-t-(--primary)',
                                    'animate-spin'
                                )}
                            />
                            <div className="mt-6 text-base font-medium text-foreground">Verifying…</div>
                            <div className="mt-2 text-sm text-(--muted-foreground) text-center">
                                This usually takes a second. Don’t close this tab.
                            </div>
                        </div>
                    )}

                    {/* BIG SUCCESS */}
                    {status === 'success' && (
                        <div className="flex flex-col items-center justify-center py-10">
                            <CheckCircle2 className="h-20 w-20 text-green-600" />
                            <div className="mt-6 text-base font-semibold text-foreground">Verified!</div>
                            <div className="mt-2 text-sm text-(--muted-foreground) text-center">
                                Taking you to your workspace…
                            </div>
                        </div>
                    )}

                    {/* Idle fallback (rare) */}
                    {status === 'idle' && (
                        <div className="flex flex-col items-center justify-center py-10">
                            <div className="text-sm text-(--muted-foreground)">Preparing verification…</div>
                        </div>
                    )}

                    {/* Actions */}
                    <div className="mt-2 flex gap-3">
                        {status === 'error' && (
                            <>
                                <Button
                                    className="flex-1"
                                    onClick={verify}
                                    disabled={loading}
                                >
                                    {loading ? 'Retrying…' : 'Retry verification'}
                                </Button>
                            </>
                        )}

                        {status === 'success' && (
                            <Button className="w-full" onClick={() => router.replace('/')}>
                                Continue
                            </Button>
                        )}
                    </div>
                </div>

                <p className="text-center text-sm text-(--muted-foreground) mt-6">
                    If this keeps failing, request a new magic link from the login page.
                </p>
            </div>
        </div>
    )
}
