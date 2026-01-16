import { Conversation, Message, Tag, User } from "@/packages/types";

export const mockUsers: Record<string, User> = {
    "current-user": {
        id: "current-user",
        name: "Lucas Miller",
        avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop",
        status: "online",
        role: "admin",
    },
    "user-1": {
        id: "user-1",
        name: "Ibim Victor",
        avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop",
        status: "online",
        role: "member",
    },
    "user-2": {
        id: "user-2",
        name: "Samuel Samuel",
        avatar: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&h=100&fit=crop",
        status: "online",
        role: "member",
    },
    "user-3": {
        id: "user-3",
        name: "Joshua Joshua",
        avatar: "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=100&h=100&fit=crop",
        status: "away",
        role: "member",
    },
    "user-4": {
        id: "user-4",
        name: "Titus Kitamura",
        avatar: "https://images.unsplash.com/photo-1519345182560-3f2917c472ef?w=100&h=100&fit=crop",
        status: "offline",
        role: "viewer",
    },
    "user-5": {
        id: "user-5",
        name: "Geoffrey Mott",
        avatar: "https://images.unsplash.com/photo-1463453091185-61582044d556?w=100&h=100&fit=crop",
        status: "offline",
        role: "member",
    },
    "user-6": {
        id: "user-6",
        name: "Alfonzo Schuessler",
        avatar: "https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=100&h=100&fit=crop",
        status: "online",
        role: "member",
    },
    "user-7": {
        id: "user-7",
        name: "Brooke Davis",
        avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100&h=100&fit=crop",
        status: "online",
        role: "member",
    },
};

const tags: Record<string, Tag> = {
    "rough-day": { id: "rough-day", label: "Rough Day", color: "orange" },
    "help-wanted": { id: "help-wanted", label: "Help wanted", color: "green" },
    "follow-up": { id: "follow-up", label: "Follow up", color: "gray" },
    "hacktoberfest": { id: "hacktoberfest", label: "Hacktoberfest", color: "purple" },
    "resting": { id: "resting", label: "Resting", color: "orange" },
    "some-content": { id: "some-content", label: "Some content", color: "gray" },
    "request": { id: "request", label: "Request", color: "orange" },
};

export const mockConversations: Conversation[] = [
    {
        id: "conv-1",
        participant: mockUsers["user-1"]!,
        lastMessage: "Haha oh man 🔥",
        lastMessageTime: new Date(Date.now() - 12 * 60000),
        tags: [tags["rough-day"]!, tags["help-wanted"]!],
    },
    {
        id: "conv-2",
        participant: mockUsers["user-2"]!,
        lastMessage: "woohoooo",
        lastMessageTime: new Date(Date.now() - 24 * 60000),
        tags: [tags["follow-up"]!],
    },
    {
        id: "conv-3",
        participant: mockUsers["user-3"]!,
        lastMessage: "Haha that's terrifying 😆",
        lastMessageTime: new Date(Date.now() - 60 * 60000),
        tags: [tags["rough-day"]!, tags["hacktoberfest"]!],
    },
    {
        id: "conv-4",
        participant: mockUsers["user-4"]!,
        lastMessage: "omg, this is amazing",
        lastMessageTime: new Date(Date.now() - 5 * 60 * 60000),
        tags: [tags["resting"]!, tags["some-content"]!],
    },
    {
        id: "conv-5",
        participant: mockUsers["user-5"]!,
        lastMessage: "aww 😍",
        lastMessageTime: new Date(Date.now() - 2 * 24 * 60 * 60000),
        tags: [tags["request"]!],
    },
    {
        id: "conv-6",
        participant: mockUsers["user-6"]!,
        lastMessage: "perfect!",
        lastMessageTime: new Date(Date.now() - 60000),
        tags: [tags["follow-up"]!],
    },
    {
        id: "conv-7",
        participant: mockUsers["user-7"]!,
        lastMessage: "You're the best!",
        lastMessageTime: new Date(Date.now() - 30 * 60000),
        tags: [],
    },
];

export const mockMessages: Record<string, Message[]> = {
    "conv-2": [
        {
            id: "msg-1",
            content: "omg, this is amazing",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 30 * 60000),
            type: "text",
        },
        {
            id: "msg-2",
            content: "perfect! ✅",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 29 * 60000),
            type: "text",
        },
        {
            id: "msg-3",
            content: "Wow, this is really epic",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 28 * 60000),
            type: "text",
        },
        {
            id: "msg-4",
            content: "How are you?",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 27 * 60000),
            type: "text",
        },
        {
            id: "msg-5",
            content: "just ideas for next time",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 26 * 60000),
            type: "text",
        },
        {
            id: "msg-6",
            content: "I'll be there in 2 mins 🍑",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 25 * 60000),
            type: "text",
        },
        {
            id: "msg-7",
            content: "woohoooo",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 24 * 60000),
            type: "text",
        },
        {
            id: "msg-8",
            content: "Haha oh man",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 23 * 60000),
            type: "text",
        },
        {
            id: "msg-9",
            content: "Haha that's terrifying 😱",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 22 * 60000),
            type: "text",
        },
        {
            id: "msg-10",
            content: "aww",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 21 * 60000),
            type: "text",
        },
        {
            id: "msg-11",
            content: "omg, this is amazing",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 20 * 60000),
            type: "text",
        },
        {
            id: "msg-12",
            content: "woohoooo 🔥",
            senderId: "user-2",
            timestamp: new Date(Date.now() - 19 * 60000),
            type: "text",
        },
    ],
    "conv-7": [
        {
            id: "msg-b-1",
            content: "Hey Lucas!",
            senderId: "user-7",
            timestamp: new Date(Date.now() - 60 * 60000),
            type: "text",
        },
        {
            id: "msg-b-2",
            content: "How's your project going?",
            senderId: "user-7",
            timestamp: new Date(Date.now() - 59 * 60000),
            type: "text",
        },
        {
            id: "msg-b-3",
            content: "Hi Brooke!",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 55 * 60000),
            type: "text",
        },
        {
            id: "msg-b-4",
            content: "It's going well. Thanks for asking!",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 54 * 60000),
            type: "text",
        },
        {
            id: "msg-b-5",
            content: "No worries. Let me know if you need any help 😊",
            senderId: "user-7",
            timestamp: new Date(Date.now() - 50 * 60000),
            type: "text",
        },
        {
            id: "msg-b-6",
            content: "You're the best!",
            senderId: "current-user",
            timestamp: new Date(Date.now() - 30 * 60000),
            type: "text",
        },
    ],
};

// Generate default empty messages for other conversations
Object.keys(mockConversations).forEach((_, index) => {
    const convId = `conv-${index + 1}`;
    if (!mockMessages[convId]) {
        mockMessages[convId] = [];
    }
});
