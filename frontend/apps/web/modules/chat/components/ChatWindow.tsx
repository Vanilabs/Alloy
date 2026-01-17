import { useRef, useEffect } from "react";
import { ChatHeader } from "./ChatHeader";
import { MessageBubble } from "./MessageBubble";
import { MessageInput } from "./MessageInput";
import { Message, User } from "@/packages/types";

interface ChatWindowProps {
    participant: User;
    messages: Message[];
    currentUserId: string;
    users: Record<string, User>;
    onSendMessage: (content: string) => void;
    onBack?: () => void;
    showBackButton?: boolean;
}

export function ChatWindow({
    participant,
    messages,
    currentUserId,
    users,
    onSendMessage,
    onBack,
    showBackButton,
}: ChatWindowProps) {
    const messagesEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages]);

    const shouldShowAvatar = (index: number, message: Message) => {
        if (index === 0) return true;
        const prevMessage = messages[index - 1]!;
        return prevMessage.senderId !== message.senderId;
    };

    return (
        <div className="flex flex-col h-full bg-background">
            <ChatHeader
                user={participant}
                onBack={onBack}
                showBackButton={showBackButton}
            />

            <div className="flex-1 overflow-y-auto p-4 space-y-3 scrollbar-thin">
                {messages.map((message, index) => {
                    const sender = users[message.senderId]!;
                    const isOwn = message.senderId === currentUserId;

                    return (
                        <MessageBubble
                            key={message.id}
                            message={message}
                            sender={sender}
                            isOwn={isOwn}
                            showAvatar={shouldShowAvatar(index, message)}
                        />
                    );
                })}
                <div ref={messagesEndRef} />
            </div>

            <MessageInput onSend={onSendMessage} />
        </div>
    );
}
