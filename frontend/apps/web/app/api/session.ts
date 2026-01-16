import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Session } from "../../../../packages/store/slices/auth.slice"
import { Features } from "../../../../packages/store/slices/features.slice"
import { TenantContext } from "../../../../packages/store/slices/tenant.slice"

export const getSession = createAsyncThunk(
    'session/getSession',
    async (): Promise<Session | null> => {
        try {
            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/session`, {
                credentials: 'include',
            })
            // for demo return sample session
            // Check if in development mode
            if (process.env.NODE_ENV === 'development')
                return {
                    user: {
                        id: 'user-123',
                        name: 'Demo User',
                        email: 'demo@example.com',
                        status: 'online'
                    },
                    tokens: {
                        accessToken: 'demo-access-token',
                        refreshToken: 'demo-refresh-token',
                    },
                }
            else
                if (!res.ok) return null
            return res.json()
        } catch {
            return null
        }
    }
);

export const getTenants = createAsyncThunk(
    'session/getTenants',
    async (token: string): Promise<TenantContext[]> => {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/tenants`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        // for demo return sample tenants
        if (process.env.NODE_ENV === 'development')
            return [
                {
                    id: 'tenant-1', slug: 'tenant-one', role: 'admin',
                    subscriptions: ['chat', 'hrm', 'tms'],
                    permissions: []
                },
                {
                    id: 'tenant-2', slug: 'tenant-two', role: 'member',
                    subscriptions: ['chat', 'meetings'],
                    permissions: []
                },
            ]
        else
            return res.ok ? res.json() : []
    }
)

export const getTenant = createAsyncThunk(
    'session/getTenant',
    async (): Promise<TenantContext> => {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/tenant`, {
            cache: 'no-store',
        })
        // for demo return sample tenants
        if (process.env.NODE_ENV === 'development')
            return {
                    id: 'tenant-1', slug: 'tenant-one', role: 'admin',
                    subscriptions: ['chat', 'hrm', 'tms'],
                    permissions: []
                }
        else
            return res.ok ? res.json() : { id: '', slug: '', role: '', subscriptions: [], permissions: []  }
    }
)

export const getFeatureFlags = createAsyncThunk(
    'session/getFeatureFlags',
    async (args: { tenantId: string; token: string }): Promise<Features> => {
        const res = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/tenants/${args.tenantId}/features`,
            {
                headers: { Authorization: `Bearer ${args.token}` },
            }
        )
        // for demo return sample feature flags
        if (process.env.NODE_ENV === 'development')
            return {
                chat: true,
                hrm: true,
                tms: false,
                meetings: true,
            }
        else
            return res.ok ? res.json() : {}
    }
)
