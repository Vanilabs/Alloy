export interface User {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    status: 'online' | 'offline' | 'away';
    role?: 'admin' | 'member' | 'viewer';
}

export interface Message {
    id: string;
    content: string;
    senderId: string;
    timestamp: Date;
    createdAt?: string;
    channelId?: string;
    type: 'text' | 'file' | 'image';
}

export interface Conversation {
    id: string;
    participant: User;
    lastMessage?: string;
    lastMessageTime?: Date;
    unreadCount?: number;
    tags?: Tag[];
}

export interface Tag {
    id: string;
    label: string;
    color: 'orange' | 'green' | 'purple' | 'gray';
}

export interface Organization {
    id: string;
    name: string;
    logo?: string;
    members: User[];
}
