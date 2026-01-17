import type { AlloyPlugin } from '@/packages/shared/types'
import { pluginContext } from '@/packages/plugins'

const plugins: AlloyPlugin[] = [
    // Example: future plugins go here
]

export async function loadPlugins() {
    for (const p of plugins) {
        p.init(pluginContext)
    }
}
