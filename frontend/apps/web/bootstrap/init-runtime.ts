'use client'

import { useEffect } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { useDispatch } from 'react-redux'
import { Session, setSession, setUser } from '../../../packages/store/slices/auth.slice'
import { markInitialized, setTenant, TenantContext } from '../../../packages/store/slices/tenant.slice'
import { getSession, getTenants, getFeatureFlags } from '../app/api'
import { registerModules } from './register-modules'
import { loadPlugins } from './load-plugins'
import { registry } from '@/packages/plugins'
import { store } from '@/packages/store'
import { Features } from '@/packages/store/slices/features.slice'

let runtimeBootstrapped = false // survives strict-mode remount in dev

export const useInitRuntime = () => {
    const router = useRouter()
    const pathname = usePathname()
    const dispatch = useDispatch()

    useEffect(() => {
        if (registry.isInitialized()) return
        if (runtimeBootstrapped) return
        runtimeBootstrapped = true

        const run = async () => {
            try {
                console.log('🚀 Bootstrapping Alloy OS runtime...')
                // const session = await getSession()

                const session = await store.dispatch(getSession());

                if (!session) {
                    if (pathname !== '/login') router.replace('/login')
                    return
                }

                const payload = session.payload as Session;

                store.dispatch(setUser(payload?.user))

                const tenants = await store.dispatch(getTenants(payload.tokens.accessToken));

                const tenantsPayload = tenants.payload as TenantContext[];

                if (!tenantsPayload?.length) {
                    if (pathname !== '/system/select-tenant') router.replace('/system/select-tenant')
                    return
                }

                const tenantIdFromUrl = pathname.split('/')[2]
                const tenant = tenantsPayload.find((t: any) => t.id === tenantIdFromUrl) ?? tenantsPayload[0]
                if (!tenant) {
                    router.replace('/system/select-tenant')
                    return
                }

                store.dispatch(setTenant(tenant))

                const features = await store.dispatch(getFeatureFlags({ tenantId: tenant.id, token: payload.tokens.accessToken}));

                const featureFlags = (features.payload as Features) || {}

                dispatch({ type: 'featureFlags/setFlags', payload: featureFlags })

                if (registry.isInitialized()) return

                await loadPlugins() //NOTE - load plugins (can register modules/nav/permissions) loaded before modules

                registerModules()

                for (const mod of registry.getModules()) {
                    mod.init?.()
                }

                store.dispatch(markInitialized())

                registry.markInitialized()

                // Only route if you are not already on a valid tenant route
                // (avoid forcing reload)
                if (pathname === '/login' || pathname === '/bootstrap' || pathname === '/system/select-tenant') {
                    router.replace('/')
                }

                console.log('✅ Alloy OS runtime initialized')
            } catch (err) {
                console.error('Runtime bootstrap failed:', err)
                router.replace('/login')
            }
        }

        run()
    }, [dispatch, pathname, router])
}
