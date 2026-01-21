"use client";

import * as React from "react";
import { Toaster } from '@/packages/ui/toaster';

export function Providers({ children }: { children: React.ReactNode }) {
    return (
        <>
            {children}
            <Toaster />
        </>
    );
}
