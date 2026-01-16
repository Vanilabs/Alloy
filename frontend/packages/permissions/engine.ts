import { RootState } from "../store"
import type { Permission } from "../types"

export function can(permissionSet: Permission[]) {
    return {
        perform(action: string) {
            return {
                on(resource: string) {
                    return permissionSet.some(p => p.action === action && p.resource === resource)
                },
            }
        },
    }
}

export function hasPermission(state: RootState, permission: string): boolean {
    const tenant = state.tenant.tenant
    if (!tenant) return false
    return tenant.permissions.includes(permission)
}
