'use client';

import React from 'react'
import { Provider } from 'react-redux'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { store } from '../store'
import { useInitRuntime } from '@/apps/web/bootstrap/init-runtime';

function RuntimeBootstrap() {
    useInitRuntime()
    return null
}

export function AppProvider({ children }: { children: React.ReactNode }) {
    const [queryClient] = React.useState(() => new QueryClient())

    return (
        <Provider store={store}>
            <QueryClientProvider client={queryClient}>
                <RuntimeBootstrap />
                {children}
            </QueryClientProvider>
        </Provider>
    )
}
