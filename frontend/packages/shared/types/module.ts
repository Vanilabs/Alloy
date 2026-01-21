export type RouteNode = {
  path: string // e.g. "/chat" or "/chat/channel/:channelId"
  // render is in module registry, not redux
}

export interface AlloyModule {
    id: string
    name: string
    description?: string
    routes?: any[]

    // gating keys
    subscriptionKey: string     // e.g. "chat"
    featureFlagKey: string      // e.g. "chat"
    entryPermission: string     // e.g. "chat.access"

    permissions?: string[]
    navigation?: NavigationItem[]
    init?(): void
    destroy?(): void

    // NOTE - IMPORTANT: how the module renders in the tenant router
    render?: (routePath: string) => React.ReactNode
}

export type FeatureFlags = Record<string, boolean>;

export type Permission = string;

export type NavigationItem = {
    label: string // e.g. "Chat"
    path: string // e.g. "/chat"
    href: string // optional full href if different from path
    icon: React.ComponentType<any> // optional icon component
}

export type AlloyPlugin = {
    id: string
    name: string
    init: (ctx: PluginContext) => void
}

export type PluginContext = {
    registerModule: (m: AlloyModule) => void
    registerPermission: (p: Permission) => void
    registerNav: (n: NavigationItem) => void
}
