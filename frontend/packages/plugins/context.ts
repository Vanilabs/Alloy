import type { PluginContext } from '@/packages/shared/types'
import { registry } from './registry'

export const pluginContext: PluginContext = {
    registerModule: registry.registerModule,
    registerPermission: registry.registerPermission,
    registerNav: registry.registerNav,
}
