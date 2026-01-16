import { cn } from "@/lib/utils";
import { Avatar } from "./Avatar";
import { StatusTag } from "./StatusTag";
import { Conversation } from "@/packages/types";

interface ConversationItemProps {
    conversation: Conversation;
    isActive?: boolean;
    onClick: () => void;
}

function formatTime(date?: Date): string {
    if (!date) return '';
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return 'now';
    if (minutes < 60) return `${minutes}m`;
    if (hours < 24) return `${hours}h`;
    return `${days}d`;
}

export function ConversationItem({ conversation, isActive, onClick }: ConversationItemProps) {
    const { participant, lastMessage, lastMessageTime, tags } = conversation;

    return (
        <button
            onClick={onClick}
            className={cn(
                "w-full flex items-start gap-3 p-3 rounded-lg transition-colors text-left",
                "hover:bg-accent",
                isActive && "bg-(--chat-active-bg)"
            )}
        >
            <Avatar
                src={participant.avatar}
                name={participant.name}
                size="lg"
                status={participant.status}
            />

            <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between gap-2">
                    <span className="font-semibold text-(--gray-60) text-sm truncate">
                        {participant.name}
                    </span>
                    <span className="text-xs text-(--gray-1) shrink-0">
                        {formatTime(lastMessageTime)}
                    </span>
                </div>

                {lastMessage && (
                    <p className="text-xs text-(--text-muted) truncate mt-0.5">
                        {lastMessage}
                    </p>
                )}

                {tags && tags.length > 0 && (
                    <div className="flex flex-wrap gap-1 mt-2">
                        {tags.map((tag) => (
                            <StatusTag key={tag.id} tag={tag} />
                        ))}
                    </div>
                )}
            </div>
        </button>
    );
}
