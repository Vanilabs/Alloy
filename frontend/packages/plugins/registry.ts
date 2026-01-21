import type { AlloyModule, NavigationItem, Permission } from '@/packages/shared/types'

export type RegistryState = {
    modules: Record<string, AlloyModule>
    nav: NavigationItem[]
    permissions: Set<Permission>
    initialized: boolean
}

const state: RegistryState = {
    modules: {},
    nav: [],
    permissions: new Set(),
    initialized: false,
}

export const registry = {
    isInitialized: () => state.initialized,

    markInitialized: () => {
        state.initialized = true
    },

    registerModule: (m: AlloyModule) => {
        state.modules[m.id] = m
        console.log('registry.registerModule', m)
        m.navigation?.forEach(n => registry.registerNav(n))
        m.permissions?.forEach(p => registry.registerPermission(p))
    },

    registerNav: (n: NavigationItem) => {
        // avoid dup
        console.log('registry.registerNav', n)
        if (!state.nav.some(x => x.path === n.path)) state.nav.push(n)
    },

    registerPermission: (p: Permission) => {
        state.permissions.add(p)
    },

    // isModuleEnabled: (m: AlloyModule): boolean => {
    //     const ctx = tenantStore.get()

    //     const subscribed = ctx.subscriptions.includes(m.subscriptionKey)
    //     const featureOn = isFeatureEnabled(m.featureFlagKey)
    //     const permitted = hasPermission(m.entryPermission)

    //     return subscribed && featureOn && permitted
    // },

    getModules: () => Object.values(state.modules),
    getNav: (): NavigationItem[] => {
        console.log('registry.getNav', state)
        const nav: NavigationItem[] = []
        for (const m of Object.values(state.modules)) {
            //   if (!registry.isModuleEnabled(m)) continue
            if (m.navigation) for (const n of m.navigation) nav.push(n)
        }
        return nav
    },
    // getNav: () => state.nav.slice(),
    getModuleByPath: (path: string): AlloyModule | undefined => {
        const mod = Object.values(state.modules).find(m =>
            m.navigation?.some(n => n.path === path)
        )
        if (!mod) return undefined
        // if (!registry.isModuleEnabled(mod)) return undefined
        console.log('registry.getModuleByPath', path, '->', mod)
        return mod
    },
    // getModuleByPath: (path: string) => {
    //     // find module whose navigation includes this path
    //     const all = Object.values(state.modules)
    //     return all.find(m => m.navigation?.some(n => n.path === path))
    // },
}
