import { Phone, MoreVertical, ArrowLeft } from "lucide-react";
import { Avatar } from "./Avatar";
import { User } from "@/packages/types";
import { Button } from "@/packages/ui";

interface ChatHeaderProps {
    user: User;
    onBack?: () => void;
    onCall?: () => void;
    showBackButton?: boolean;
}

export function ChatHeader({ user, onBack, onCall, showBackButton }: ChatHeaderProps) {
    return (
        <div className="flex items-center justify-between px-4 py-3 border-b border-(--border) bg-background">
            <div className="flex items-center gap-3">
                {showBackButton && (
                    <Button variant="ghost" size="icon" onClick={onBack} className="md:hidden">
                        <ArrowLeft className="w-5 h-5" />
                    </Button>
                )}
                <Avatar src={user.avatar} name={user.name} size="md" status={user.status} />
                <div>
                    <h2 className="font-semibold text-foreground">{user.name}</h2>
                    <p className="text-xs text-success flex items-center gap-1">
                        <span className={`w-2 h-2 rounded-full ${user.status.toLocaleLowerCase() === 'online' ? 'bg-(--tag-green)' : user.status === 'away' ? 'bg-(--tag-orange)' : 'bg-(--tag-gray)'} inline-block`} />
                        {user.status === 'online' ? 'Online' : user.status === 'away' ? 'Away' : 'Offline'}
                    </p>
                </div>
            </div>

            <div className="flex items-center gap-2 bg-(--primary-foreground) text-(--primary) rounded-sm">
                <Button
                    // variant="outline"
                    size="sm"
                    onClick={onCall}
                    className="text-(--primary) hover:bg-(--primary-foreground)"
                >
                    <Phone className="w-4 h-4 mr-2" />
                    Call
                </Button>
                {/* <Button variant="ghost" size="icon">
                    <MoreVertical className="w-5 h-5" />
                </Button> */}
            </div>
        </div>
    );
}
