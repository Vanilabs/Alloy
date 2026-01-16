"use client";

import { createContext, useContext, useState } from "react";

const SidebarSlotContext = createContext<{
    slot: React.ReactNode;
    setSlot: (node: React.ReactNode) => void;
    showMobileChat: boolean;
    setShowMobileChat: (show: boolean) => void;
}>({
    slot: null,
    setSlot: () => { },
    showMobileChat: false,
    setShowMobileChat: () => { },
});

export function SidebarSlotProvider({ children }: { children: React.ReactNode }) {
    const [slot, setSlot] = useState<React.ReactNode>(null);
    const [showMobileChat, setShowMobileChat] = useState(false);

    return (
        <SidebarSlotContext.Provider value={{ slot, setSlot, showMobileChat, setShowMobileChat }}>
            {children}
        </SidebarSlotContext.Provider>
    );
}

export function useSidebarSlot() {
    return useContext(SidebarSlotContext);
}
