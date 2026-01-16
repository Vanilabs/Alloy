'use client'

import { useAppSelector } from "../store/hooks"

export default function RequirePermission({
    permission,
    children,
}: {
    permission: string
    children: React.ReactNode
}) {
    const allowed = useAppSelector(s => s.tenant.tenant?.permissions?.includes(permission) ?? false)
    if (!allowed) return <div className="p-6">Access denied</div>
    return <>{children}</>
}
