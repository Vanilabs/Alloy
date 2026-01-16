import { registry } from '@/packages/plugins'
import { ChatModule } from '../modules/chat/register'

export const registerModules = () => {
    const modules = [ChatModule /*, HRMModule, TMSModule */]
    // modules.forEach(m => m.init?.())
    modules.forEach(m => registry.registerModule(m))
}
