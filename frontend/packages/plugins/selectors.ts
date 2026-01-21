import { createSelector } from '@reduxjs/toolkit'
import { registry } from './registry'
import type { AlloyModule, NavigationItem } from '@/packages/shared/types'
import { RootState } from '../store'

const selectTenant = (s: RootState) => s.tenant.tenant;

console.log('selectTenant', selectTenant);

export const selectEnabledModules = createSelector([selectTenant], (tenant): AlloyModule[] => {
    const all = registry.getModules()
    console.log('selectEnabledModules - all modules:', all);
    console.log('selectEnabledModules - tenant:', tenant);
    if (!tenant) return []

    return all.filter(m => {
        const subscribed = tenant.subscriptions.includes(m.subscriptionKey)
        const flagOn = Boolean(tenant.featureFlags?.[m.featureFlagKey])
        const permitted = tenant.permissions.includes(m.entryPermission)
        return subscribed && flagOn && permitted
    })
})

export const selectEnabledNav = createSelector([selectEnabledModules], (mods): NavigationItem[] => {
    const nav: NavigationItem[] = []
    for (const m of mods) {
        if (m.navigation) for (const n of m.navigation) nav.push(n)
    }
    return nav
})
