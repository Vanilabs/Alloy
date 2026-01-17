import type { AlloyModule } from '../shared/types'

export interface AlloyPlugin {
    id: string
    name: string
    modules?: AlloyModule[]
    permissions?: string[]
    init(ctx: PluginContext): void
}

export interface PluginContext {
    addRoute(route: any): void
    addSidebarItem(item: any): void
    registerPermission(permission: string): void
    registerCommand(cmd: any): void
}
