import { io, Socket } from 'socket.io-client'

let socket: Socket | null = null

export const connectSocket = (tenantId: string, userId: string) => {
    socket = io(`${process.env.NEXT_PUBLIC_WS_URL}`, {
        query: { tenantId, userId },
    })
    return socket
}

export const getSocket = () => socket
