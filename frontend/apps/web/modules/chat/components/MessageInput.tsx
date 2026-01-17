import { useState, useRef } from "react";
import { Send, Paperclip, Image, Plus } from "lucide-react";
import { cn } from "@/lib/utils";
import {
    Button,
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/packages/ui";

interface MessageInputProps {
    onSend: (content: string) => void;
    onAttach?: (type: 'file' | 'screenshot') => void;
    placeholder?: string;
    disabled?: boolean;
}

export function MessageInput({ onSend, onAttach, placeholder = "Type a message", disabled }: MessageInputProps) {
    const [message, setMessage] = useState("");
    const [showAttach, setShowAttach] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);

    const handleSend = () => {
        if (message.trim()) {
            onSend(message.trim());
            setMessage("");
            inputRef.current?.focus();
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    return (
        <div className="flex items-center gap-2 p-4 shadow-2xl bg-background">
            <Popover open={showAttach} onOpenChange={setShowAttach}>
                <PopoverTrigger asChild>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="shrink-0 text-white bg-(--primary)/10 rounded-full"
                    >
                        <Paperclip className="w-5 h-5" />
                    </Button>
                </PopoverTrigger>
                <PopoverContent
                    side="top"
                    align="start"
                    className="w-66 p-2 bg-white border-(--border)"
                >
                    <div className="flex flex-col">
                        <button
                            onClick={() => {
                                onAttach?.('file');
                                setShowAttach(false);
                            }}
                            className="flex items-center gap-3 px-3 py-2 rounded-lg text-(--muted-foreground) hover:bg-(--primary) hover:text-white transition-colors text-left"
                        >
                            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
                                <Paperclip className="w-4 h-4" />
                            </div>
                            <span className="font-medium">Send File</span>
                        </button>
                        <button
                            onClick={() => {
                                onAttach?.('screenshot');
                                setShowAttach(false);
                            }}
                            className="flex items-center gap-3 px-3 py-2 rounded-lg text-(--muted-foreground) hover:bg-(--primary) hover:text-white transition-colors text-left"
                        >
                            <div className="w-8 h-8 rounded-lg bg-accent flex items-center justify-center">
                                <Image className="w-4 h-4" />
                            </div>
                            <span className="font-medium">Attach a screenshot</span>
                        </button>
                    </div>
                </PopoverContent>
            </Popover>

            <div className="flex-1 relative">
                <input
                    ref={inputRef}
                    type="text"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder={placeholder}
                    disabled={disabled}
                    className={cn(
                        "w-full px-4 py-2.5 rounded-md bg-white border border-(--border) focus:border-(--primary) transition-colors",
                        "placeholder:text-(--muted-foreground)",
                        "focus:outline-none focus:ring-0 focus:ring-(--primary)/20",
                        "disabled:opacity-50 disabled:cursor-not-allowed"
                    )}
                />

                <Button
                    variant="ghost"
                    size="icon"
                    onClick={handleSend}
                    disabled={!message.trim() || disabled}
                    className={cn(
                        "absolute right-3 top-1/2 -translate-y-1/2 shrink-0 transition-colors",
                        message.trim()
                            ? "text-(--primary) hover:text-(--primary)/10"
                            : "text-(--muted-foreground)"
                    )}
                >
                    <Send className="w-5 h-5" />
                </Button>
            </div>
        </div>
    );
}
