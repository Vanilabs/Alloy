import { ChatRoutes } from './routes'
import ChatLayout from './components/ChatLayout'
import { ChatPermissions } from './permissions'
import { AlloyModule } from '@/packages/shared/types'
import React from 'react'
import {
    Bell,
    BarChart3,
    Users,
    Calendar,
    Zap,
    MessageCircle,
    Settings,
    Home as HomeIcon,
} from "lucide-react";

export const ChatModule: AlloyModule = {
    id: 'chat',
    name: 'Chat',
    description: 'Internal communication',
    routes: ChatRoutes,

    subscriptionKey: 'chat',
    featureFlagKey: 'chat',
    entryPermission: 'chat.access',

    permissions: Object.values(ChatPermissions),
    navigation: [
        { label: 'Chat', path: '/chat', icon: MessageCircle, href: '/chat' },
    ],
    init() {
        console.log('Chat module initialized')
    },
    destroy() {
        console.log('Chat module destroyed')
    },
    render(routePath: string): React.ReactNode {
        // For now, we maintain one page per module. Later: nested routes.
        if (routePath === '/chat') return React.createElement(ChatLayout)
        return React.createElement(ChatLayout)
    },
}
