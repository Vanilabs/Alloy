import ChatLayout from './components/ChatLayout'

export const ChatRoutes = [
    {
        path: '/chat',
        name: 'Chat',
        element: <ChatLayout />,
        permissions: ['chat.read'],
    },
]
