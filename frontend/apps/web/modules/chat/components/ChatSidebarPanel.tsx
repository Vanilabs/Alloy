import { Search, ScanSearchIcon, Plus, ChevronDown } from "lucide-react";
import type { Conversation } from "@/packages/types";
import { ConversationItem } from "./ConversationItem";
import { cn } from "@/lib/utils";
import { useState } from "react";

interface ChatSidebarPanelProps {
    conversations: Conversation[];
    activeConversationId?: string;
    onSelectConversation: (id: string) => void;
    messageCount: number;
}

export function ChatSidebarPanel({
    conversations,
    activeConversationId,
    onSelectConversation,
    messageCount,
}: ChatSidebarPanelProps) {
    const [searchQuery, setSearchQuery] = useState("");

    const filteredConversations = conversations.filter((conv) =>
        conv.participant.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="w-80 bg-background border-r border-(--border) flex flex-col">
            <div className="p-4">
                <div className="flex items-center justify-between mb-4">
                    <div className="flex">
                        <button className="flex items-center gap-2 text-lg font-semibold">
                            Messages
                            <ChevronDown className="w-4 h-4" />
                        </button>

                        <span className="px-2 py-1 bg-(--tag-gray-bg) rounded-full text-sm">
                            {messageCount}
                        </span>
                    </div>

                    <Plus className="w-6 h-6 text-white bg-(--primary) rounded-full" />
                </div>

                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--input-search-icon)" />
                    <input
                        type="text"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        placeholder="Search messages"
                        className={cn(
                            "w-full pl-10 pr-4 py-2 rounded-lg bg-(--input-search) border-0",
                            "placeholder:text-(--muted-foreground)", "text-[14px]",
                            "focus:outline-none focus:ring-none focus:ring-primary/20"
                        )}
                    />
                </div>
            </div>

            <div className="flex-1 overflow-y-auto px-2">
                {filteredConversations.map((conversation) => (
                    <ConversationItem
                        key={conversation.id}
                        conversation={conversation}
                        isActive={conversation.id === activeConversationId}
                        onClick={() => onSelectConversation(conversation.id)}
                    />
                ))}
            </div>
        </div>
    );
}
