"use client";

import { useState, useCallback, useEffect } from "react";
import { mockUsers, mockConversations, mockMessages } from "@/data/mockData";
import type { Message } from "@/packages/types";
import { ChatSidebarPanel } from "../../modules/chat/components/ChatSidebarPanel";
import { ChatWindow } from "../../modules/chat/components/ChatWindow";
import { Sidebar, useSidebarSlot } from "../../shell/sidebar";

export default function TenantHome() {
    const { setSlot, setShowMobileChat, showMobileChat } = useSidebarSlot();
    const [conversations] = useState(mockConversations);
    const [activeConversationId, setActiveConversationId] =
        useState<string | undefined>("conv-2");

    const [messages, setMessages] =
        useState<Record<string, Message[]>>(mockMessages);

    const currentUser = mockUsers["current-user"]!;
    const activeConversation = conversations.find(
        (c) => c.id === activeConversationId
    );

    const handleSelectConversation = useCallback((id: string) => {
        setActiveConversationId(id);
        setShowMobileChat(true);
    }, []);

    useEffect(() => {
        setSlot(
            <ChatSidebarPanel
                conversations={conversations}
                activeConversationId={activeConversationId}
                onSelectConversation={handleSelectConversation}
                messageCount={conversations.length}
            />
        );

        return () => setSlot(null); // cleanup on route change
    }, [conversations, activeConversationId, setSlot]);

    const handleSendMessage = useCallback(
        (content: string) => {
            if (!activeConversationId) return;

            const newMessage: Message = {
                id: `msg-${Date.now()}`,
                content,
                senderId: currentUser.id,
                timestamp: new Date(),
                type: "text",
            };

            setMessages((prev) => ({
                ...prev,
                [activeConversationId]: [
                    ...(prev[activeConversationId] || []),
                    newMessage,
                ],
            }));
        },
        [activeConversationId, currentUser.id]
    );

    const handleBack = useCallback(() => {
        setShowMobileChat(false);
    }, []);

    return (
        <div
            className={`flex-1 ${!showMobileChat && !activeConversation
                    ? "hidden md:flex"
                    : "flex"
                }`}
        >
            {activeConversation ? (
                <div className="flex-1">
                    <ChatWindow
                        participant={activeConversation.participant}
                        messages={messages[activeConversationId!] || []}
                        currentUserId={currentUser.id}
                        users={mockUsers}
                        onSendMessage={handleSendMessage}
                        onBack={handleBack}
                        showBackButton={showMobileChat}
                    />
                </div>
            ) : (
                <div className="flex-1 flex items-center justify-center text-muted-foreground">
                    <p className="text-lg">
                        Select a conversation to start messaging
                    </p>
                </div>
            )}
        </div>
    );
}
