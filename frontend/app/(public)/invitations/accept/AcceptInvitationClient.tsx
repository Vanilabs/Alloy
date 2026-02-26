'use client'

import { useEffect, useMemo, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { CheckCircle2, XCircle } from 'lucide-react'
import { Button } from '@/packages/ui/components/Button'
import { cn } from '@/lib/utils'
import { apiFetch } from '@/packages/api'
import { AUTH_LINKS } from '@/packages/constants'
import { useAppDispatch } from '@/packages/store/hooks'
import { acceptInvitationThunk } from '@/packages/store/thunks'
import { safeErrorMessage, useLoopingTypedText } from '@/packages/shared/utils'

type Status = 'loading' | 'success' | 'error'

export default function AcceptInvitationPage() {
    const router = useRouter()
    const searchParams = useSearchParams()
    const dispatch = useAppDispatch();

    const token = useMemo(() => searchParams.get('token') ?? '', [searchParams])
    const email = useMemo(() => searchParams.get('email') ?? '', [searchParams])

    const [status, setStatus] = useState<Status>('loading')
    const [error, setError] = useState<string | null>(null)

    const typed = useLoopingTypedText('Accepting', 140, 650)

    useEffect(() => {
        let cancelled = false

        async function acceptInvite() {
            setStatus('loading')
            setError(null)

            if (!token || token.length < 10) {
                setStatus('error')
                setError('Invalid or missing invitation token.')
                return
            }

            try {
                const res = await dispatch(acceptInvitationThunk({ token, email }));

                if (!acceptInvitationThunk.fulfilled.match(res)) {
                    const msg = await safeErrorMessage(res as any)
                    const errorMsg = typeof res.payload === 'string' ? res.payload : res.error.message ?? msg;
                    throw new Error(errorMsg || 'Unable to accept invitation!')
                }

                if (cancelled) return
                setStatus('success')

                // redirect after 5 seconds
                setTimeout(() => {
                    router.replace('/')
                }, 5000)
            } catch (e: any) {
                if (cancelled) return
                setStatus('error')
                setError(e?.message || 'Error accepting invitation.')
            }
        }

        acceptInvite()

        return () => {
            cancelled = true
        }
    }, [token, router])

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

                    <h1 className="text-2xl font-bold text-foreground">Alloy Invitations</h1>
                    <p className="text-muted-foreground mt-2">
                        We’re processing your organization invite.
                    </p>
                </div>

                <div className="shadow-md rounded-2xl p-6 bg-white/60 backdrop-blur">
                    {status === 'loading' && (
                        <div className="flex flex-col items-center justify-center py-10">
                            <div
                                className={cn(
                                    'h-20 w-20 rounded-full border-4',
                                    'border-(--primary)/20 border-t-(--primary)',
                                    'animate-spin'
                                )}
                            />
                            <div className="mt-6 text-base font-semibold text-foreground">
                                {typed?typed:'. . .'}
                            </div>
                            <div className="mt-2 text-sm text-(--muted-foreground) text-center">
                                Please hold on — this will only take a moment.
                            </div>
                        </div>
                    )}

                    {status === 'success' && (
                        <div className="flex flex-col items-center justify-center py-10">
                            <CheckCircle2 className="h-20 w-20 text-green-600" />
                            <div className="mt-6 text-base font-semibold text-foreground">
                                Invitation accepted
                            </div>
                            <div className="mt-2 text-sm text-(--muted-foreground) text-center">
                                Redirecting you to login…
                            </div>

                            <div className="mt-6 w-full">
                                <Button
                                    className="w-full"
                                    onClick={() => router.push('/auth/login')}
                                >
                                    Go to login now
                                </Button>
                            </div>
                        </div>
                    )}

                    {status === 'error' && (
                        <div className="py-6">
                            <div className="mb-4 p-3 rounded-lg bg-(--destructive-foreground)/10 text-(--destructive) text-sm">
                                <div className="flex items-start gap-2">
                                    <XCircle className="mt-0.5 h-5 w-5" />
                                    <div>
                                        <div className="font-medium">Could not accept invite</div>
                                        <div className="mt-1">{error}</div>
                                    </div>
                                </div>
                            </div>

                            <div className="flex gap-3">
                                <Button className="flex-1" onClick={() => router.push('/auth/login')}>
                                    Go to login
                                </Button>
                                <Button
                                    //   variant="outline"
                                    className="flex-1"
                                    onClick={() => window.location.reload()}
                                >
                                    Retry
                                </Button>
                            </div>
                        </div>
                    )}
                </div>

                <p className="text-center text-sm text-(--muted-foreground) mt-6">
                    If this keeps failing, ask your admin to resend the invitation.
                </p>
            </div>
        </div>
    )
}
