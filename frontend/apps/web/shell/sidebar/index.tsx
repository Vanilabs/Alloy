"use client";
import Icon, {
    Bell,
    BarChart3,
    Users,
    Calendar,
    Zap,
    MessageCircle,
    Settings,
    Home as HomeIcon,
} from "lucide-react";
import { motion } from "framer-motion";
import { cn } from "@/lib/utils";
import type { User } from "@/packages/types";
import { Avatar } from "@/packages/ui";
import { useSidebarSlot } from "./sidebar-context";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { registry, selectEnabledNav } from "@/packages/plugins";
import { useAppSelector } from "@/packages/store/hooks";

// interface SidebarProps {
//     conversations: Conversation[];
//     activeConversationId?: string;
//     onSelectConversation: (id: string) => void;
//     currentUser: User;
//     messageCount?: number;
// }

interface SidebarProps {
    currentUser: User;
    children?: React.ReactNode; // injection point if needed
}

const navItems = [
    { icon: HomeIcon, label: "Home", href: "/" },
    // { icon: BarChart3, label: "Analytics", href: "/analytics" },
    // { icon: Users, label: "Team", href: "/team" },
    // { icon: Calendar, label: "Calendar", href: "/calendar" },
    // { icon: Zap, label: "Automation", href: "/automation" },
    // { icon: Bell, label: "Alerts", href: "/alerts" },
];

export function Sidebar({
    currentUser,
    children,
}: SidebarProps) {
    const pathname = usePathname();
    const navData = useAppSelector(selectEnabledNav)
    const [nav, setNav] = useState(registry.getNav())
    const { slot, showMobileChat } = useSidebarSlot();
    const [activeNav, setActiveNav] = useState<string>("Home");

    // useEffect(() => {
    //     if (pathname) setActiveNav(pathname);
    // }, [pathname]);

    // refresh if runtime registers later (optional)
    useEffect(() => {
        setNav(registry.getNav())
    }, [])

    return (
        <div className={`${showMobileChat ? 'hidden md:flex' : 'flex'} h-full`}>
            <div className="flex h-full">
                {/* Icon Navigation */}
                <div className="w-16 bg-background border-r border-(--border) flex flex-col items-center py-4">
                    <div className="w-10 h-10 rounded-xl bg-primary flex items-center justify-center text-white font-bold text-lg mb-6"
                        style={{
                            background: 'linear-gradient(135deg, hsla(243, 75%, 59%, 0.1) 0%, hsla(243, 75%, 59%, 1) 50%)',
                            boxShadow: '0px 4px 8px 0px rgba(79, 70, 229, 0.1)'
                        }}>
                        A
                    </div>

                    {/* Nav */}
                    <nav className="flex-1 flex flex-col items-center gap-2">
                        {[...navItems, ...nav].map(({ icon: Icon, href, label }) => {
                            // const isActive = href === "/"
                            //     ? pathname === "/"
                            //     : pathname.startsWith(href);
                            const isActive = label === activeNav;

                            return (
                                <Link
                                    key={href}
                                    href={href}
                                    onClick={() => setActiveNav(label)}
                                    aria-label={label}
                                    className={`
                                        w-10 h-10 rounded-lg flex items-center justify-center transition-colors
                                        ${isActive
                                            ? "bg-surfacetext-foreground bg-[hsl(var(--surface))] rounded-l-full rounded-r-full"
                                            : "text-muted-foreground hover:text-foreground hover:bg-accent"}`}
                                >
                                    <Icon className="w-5 h-5" />
                                </Link>
                            );
                        })}
                    </nav>

                    <div className="mt-auto flex flex-col items-center gap-3">
                        <button className="w-10 h-10 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-accent transition-colors">
                            <Settings className="w-5 h-5" />
                        </button>
                        <Avatar src={currentUser.avatar} name={currentUser.name} size="md" />
                    </div>
                </div>

                {/* Optional panel injected by pages */}
                {slot}
            </div>
        </div>
    );
}

export * from './sidebar-context';
