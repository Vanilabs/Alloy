import { cn } from "@/lib/utils";

interface AvatarProps {
    src?: string;
    name: string;
    size?: 'sm' | 'md' | 'lg';
    status?: 'online' | 'offline' | 'away';
    className?: string;
}

const sizeClasses = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-10 h-10 text-sm',
    lg: 'w-12 h-12 text-base',
};

const statusSizeClasses = {
    sm: 'w-2.5 h-2.5 right-0 bottom-0',
    md: 'w-3 h-3 right-0 bottom-0',
    lg: 'w-3.5 h-3.5 right-0.5 bottom-0.5',
};

const statusColors = {
    online: 'bg-success',
    offline: 'bg-muted-foreground',
    away: 'bg-warning',
};

export function Avatar({ src, name, size = 'md', status, className }: AvatarProps) {
    const initials = name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2);

    return (
        <div className={cn("relative inline-block", className)}>
            {src ? (
                <img
                    src={src}
                    alt={name}
                    className={cn(
                        "rounded-lg object-cover ring-none",
                        sizeClasses[size]
                    )}
                />
            ) : (
                <div
                    className={cn(
                        "rounded-full bg-primary flex items-center justify-center font-medium text-primary-foreground",
                        sizeClasses[size]
                    )}
                >
                    {initials}
                </div>
            )}
            {/* {status && (
                <span
                    className={cn(
                        "absolute rounded-full border-2 border-background",
                        statusSizeClasses[size],
                        statusColors[status]
                    )}
                />
            )} */}
        </div>
    );
}
