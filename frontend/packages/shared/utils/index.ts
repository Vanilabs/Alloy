import { Features } from "@/packages/store/slices/features.slice"
import { IModuleKey } from "../types"
import { useEffect, useState } from "react"

const defaultModules: Record<IModuleKey, boolean> = {
    chat: false,
    hrm: false,
    tms: false,
    meetings: false,
}

export function mapSubscriptionsToFlags(
    subscriptions: string[]
): Features {
    const result = { ...defaultModules }

    for (const sub of subscriptions) {
        if (sub in result) {
            result[sub as IModuleKey] = true
        }
    }

    return result
}

export function maskEmail(email: string) {
    const [name, domain] = email.split('@')
    const n = String(name);
    if (!domain) return email
    const maskedName =
        n.length <= 2 ? `${n[0] ?? ''}*` : `${n.slice(0, 2)}***`
    return `${maskedName}@${domain}`
}

export function useLoopingTypedText(text: string, speedMs = 140, pauseMs = 700) {
    const [value, setValue] = useState('')

    useEffect(() => {
        let mounted = true
        let i = 0
        let timeout: any

        const tick = () => {
            if (!mounted) return
            i += 1
            setValue(text.slice(0, i))

            if (i >= text.length) {
                timeout = setTimeout(() => {
                    if (!mounted) return
                    i = 0
                    setValue('')
                    timeout = setTimeout(tick, speedMs)
                }, pauseMs)
            } else {
                timeout = setTimeout(tick, speedMs)
            }
        }

        tick()

        return () => {
            mounted = false
            clearTimeout(timeout)
        }
    }, [text, speedMs, pauseMs])

    return value
}

export async function safeErrorMessage(res: Response) {
    try {
        const ct = res.headers.get('content-type') || ''
        if (ct.includes('application/json')) {
            const j = await res.json()
            return j?.message || j?.error
        }
        const t = await res.text()
        return t?.slice(0, 200)
    } catch {
        return null
    }
}
