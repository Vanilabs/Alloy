"use client";

import { notFound, useRouter } from 'next/navigation'
import { registry, selectEnabledModules } from '@/packages/plugins'
import { useMemo } from 'react'
import { useAppSelector } from '@/packages/store/hooks'

export default function DynamicTenantRoute({
    params,
}: {
    params: { path?: string[] }
}) {
    const router = useRouter()
    const enabledModules = useAppSelector(selectEnabledModules)


    const path = '/' + (params.path?.join('/') ?? '')
    const enabledModuleIds = useMemo(() => new Set(enabledModules.map(m => m.id)), [enabledModules])

    const mod = registry.getModuleByPath(path)
    // if (!mod || !mod.render) return notFound()

    if (!mod || !enabledModuleIds.has(mod.id) || !mod.render) return notFound()

    return <>{mod.render(path)}</>
}
