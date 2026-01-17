import { cn } from "@/lib/utils";
import { Avatar } from "./Avatar";
import { Message, User } from "@/packages/types";

interface MessageBubbleProps {
    message: Message;
    sender: User;
    isOwn: boolean;
    showAvatar?: boolean;
}

export function MessageBubble({ message, sender, isOwn, showAvatar = true }: MessageBubbleProps) {
    return (
        <div
            className={cn(
                "flex gap-2 max-w-[80%] animate-fade-in",
                isOwn ? "ml-auto flex-row-reverse" : "mr-auto"
            )}
        >
            {showAvatar && !isOwn && (
                <Avatar src={sender.avatar} name={sender.name} size="sm" />
            )}
            {!showAvatar && !isOwn && <div className="w-8" />}

            <div
                className={cn(
                    "px-4 py-2.5 rounded-2xl",
                    isOwn
                        ? "bg-(--chat-sent-foreground) text-(--chat-sent) rounded-br-md"
                        : "bg-(--chat-received-foreground) text-(--chat-received) rounded-bl-md"
                )}
            >
                {!isOwn && showAvatar && (
                    <p className="text-xs font-medium text-primary mb-1">{sender.name}</p>
                )}
                <p className="text-sm leading-relaxed whitespace-pre-wrap">{message.content}</p>
            </div>

            {showAvatar && isOwn && (
                <Avatar src={sender.avatar} name={sender.name} size="sm" />
            )}
            {!showAvatar && isOwn && <div className="w-8" />}
        </div>
    );
}
