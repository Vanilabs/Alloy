import { cn } from "@/lib/utils";
import { Tag } from "@/packages/types";

interface StatusTagProps {
    tag: Tag;
    className?: string;
}

const colorClasses = {
    orange: 'bg-(--tag-orange-bg) text-(--tag-orange)',
    green: 'bg-(--tag-green-bg) text-(--tag-green)',
    purple: 'bg-(--tag-purple-bg) text-(--tag-purple)',
    gray: 'border-2 text-(--tag-gray)',
};

export function StatusTag({ tag, className }: StatusTagProps) {
    return (
        <span
            className={cn(
                "inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium",
                colorClasses[tag.color],
                className
            )}
        >
            {tag.label}
        </span>
    );
}
