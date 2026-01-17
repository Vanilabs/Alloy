import * as React from "react";
import * as AvatarPrimitive from "@radix-ui/react-avatar";

import { cn } from "@/lib/utils";

const sizeMap = {
    sm: "h-8 w-8 text-xs",
    md: "h-10 w-10 text-sm",
    lg: "h-14 w-14 text-base",
};

interface AvatarProps extends React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Root> {
    size?: keyof typeof sizeMap;
    src?: string;
    name?: string;
}

const Avatar = React.forwardRef<
    React.ElementRef<typeof AvatarPrimitive.Root>,
    AvatarProps
>(({ className, size = "md", ...props }, ref) => (
    <AvatarPrimitive.Root
        ref={ref}
        className={cn("relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full",
        sizeMap[size],
        className)}
        {...props}
    >
        {props.src && (
            <AvatarImage src={props.src} alt={props.name} className="aspect-square h-full w-full" />
        )}

        {!props.src && (
            <AvatarFallback>
                {props.name
                    ? props.name
                          .split(" ")
                          .map((n) => n[0])
                          .join("")
                          .toUpperCase()
                    : "?"}
            </AvatarFallback>
        )}
    </AvatarPrimitive.Root>
));
Avatar.displayName = AvatarPrimitive.Root.displayName;

const AvatarImage = React.forwardRef<
    React.ElementRef<typeof AvatarPrimitive.Image>,
    React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Image>
>(({ className, ...props }, ref) => (
    <AvatarPrimitive.Image ref={ref} className={cn("aspect-square h-full w-full", className)} {...props} />
));
AvatarImage.displayName = AvatarPrimitive.Image.displayName;

const AvatarFallback = React.forwardRef<
    React.ElementRef<typeof AvatarPrimitive.Fallback>,
    React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Fallback>
>(({ className, ...props }, ref) => (
    <AvatarPrimitive.Fallback
        ref={ref}
        className={cn("flex h-full w-full items-center justify-center rounded-full bg-muted", className)}
        {...props}
    />
));
AvatarFallback.displayName = AvatarPrimitive.Fallback.displayName;

export { Avatar, AvatarImage, AvatarFallback };
