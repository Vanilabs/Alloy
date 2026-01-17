import { useState } from "react";
import {
    Bell,
    BarChart3,
    Users,
    Calendar,
    Zap,
    Settings,
    ChevronDown,
    Plus,
    Search,
    Home as HomeIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Avatar } from "./Avatar";
import { ConversationItem } from "./ConversationItem";
import { Button } from "@/packages/ui";
import { Conversation, User } from "@/packages/types";

interface SidebarProps {
    conversations: Conversation[];
    activeConversationId?: string;
    onSelectConversation: (id: string) => void;
    currentUser: User;
    messageCount?: number;
}

const navItems = [
    { icon: HomeIcon, label: "Notifications" },
    { icon: BarChart3, label: "Analytics" },
    { icon: Users, label: "Team" },
    { icon: Calendar, label: "Calendar" },
    { icon: Zap, label: "Automation" },
    { icon: Bell, label: "Alerts" },
];

export function Sidebar({
    conversations,
    activeConversationId,
    onSelectConversation,
    currentUser,
    messageCount = 0,
}: SidebarProps) {
    const [searchQuery, setSearchQuery] = useState("");

    const filteredConversations = conversations.filter((conv) =>
        conv.participant.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="flex h-full">
            {/* Icon Navigation */}
            <div className="w-16 bg-background border-r border-(--border) flex flex-col items-center py-4">
                <div className="w-10 h-10 rounded-xl bg-primary flex items-center justify-center text-primary-foreground font-bold text-lg mb-6">
                    A
                </div>

                <nav className="flex-1 flex flex-col items-center gap-2">
                    {navItems.map((item, index) => (
                        <button
                            key={index}
                            className={cn(
                                "w-10 h-10 rounded-lg flex items-center justify-center transition-colors",
                                "text-muted-foreground hover:text-foreground hover:bg-accent"
                            )}
                        >
                            <item.icon className="w-5 h-5" />
                        </button>
                    ))}
                </nav>

                <div className="mt-auto flex flex-col items-center gap-3">
                    <button className="w-10 h-10 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-accent transition-colors">
                        <Settings className="w-5 h-5" />
                    </button>
                    <Avatar src={currentUser.avatar} name={currentUser.name} size="md" />
                </div>
            </div>

            {/* Conversation List */}
            <div className="w-80 bg-background border-r border-(--border) flex flex-col">
                <div className="p-4">
                    <div className="flex items-center justify-between mb-4">
                        <button className="flex items-center gap-2 text-lg font-semibold">
                            Messages
                            <ChevronDown className="w-4 h-4" />
                        </button>
                        <div className="flex items-center gap-2">
                            <span className="px-2 py-0.5 bg-muted rounded-full text-sm font-medium">
                                {messageCount}
                            </span>
                            <Button size="icon" className="rounded-full">
                                <Plus className="w-4 h-4" />
                            </Button>
                        </div>
                    </div>

                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                        <input
                            type="text"
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            placeholder="Search messages"
                            className={cn(
                                "w-full pl-10 pr-4 py-2 rounded-lg bg-muted border-0",
                                "placeholder:text-muted-foreground",
                                "focus:outline-none focus:ring-2 focus:ring-primary/20"
                            )}
                        />
                    </div>
                </div>

                <div className="flex-1 overflow-y-auto scrollbar-thin px-2">
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
        </div>
    );
}
