import { useState } from 'react'
import { connectSocket, getSocket } from '../realtime/socket'
import { Message } from '@/packages/types'

export const useMessages = () => {
    const [messages, setMessages] = useState<Message[]>([])

    const socket = getSocket()

    socket?.on('chat.message', (msg: Message) => {
        setMessages(prev => [...prev, msg])
    })

    const sendMessage = (content: string) => {
        const msg: Message = {
            id: Math.random().toString(),
            senderId: 'me',
            content,
            channelId: '',
            createdAt: new Date().toISOString(),
            timestamp: new Date(),
            type: 'text',
        }
        socket?.emit('chat.message', msg)
        setMessages(prev => [...prev, msg])
    }

    return { messages, sendMessage }
}
