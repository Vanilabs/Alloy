import { getTenants, getTenant as loadTenantContext } from '@/apps/web/app/api'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface TenantContext {
    id: string
    slug: string
    subscriptions: string[]
    role: string
    featureFlags?: Record<string, boolean>
    permissions: string[]
}

export type RuntimeState = {
    initialized: boolean
    tenant: TenantContext | null
    loading: boolean
    error?: string
}

// const initialState: TenantContext = { id: '', slug: '', role: '', subscriptions: [], permissions: [] }

const initialState: RuntimeState = {
    initialized: false,
    tenant: null,
    loading: false,
}

const tenantSlice = createSlice({
    name: 'tenant',
    initialState,
    reducers: {
        markInitialized(state) {
                    state.initialized = true
                },
        setTenant(state, action: PayloadAction<TenantContext>) {
            state.tenant = action.payload;
        },
        resetRuntime(state) {
            Object.assign(state, initialState)
        },
        clearTenant(state) {
            Object.assign(state, initialState)
        },
        // get tenant info
        // getTenant(state, action: PayloadAction<TenantContext>) {
        //     return Object.assign(state, action.payload)
        // },
    },
    extraReducers: builder => {
            builder
                .addCase(loadTenantContext.pending, state => {
                    state.loading = true
                    state.error = undefined
                })
                .addCase(loadTenantContext.fulfilled, (state, action) => {
                    state.loading = false
                    state.tenant = action.payload
                })
                .addCase(loadTenantContext.rejected, (state, action) => {
                    state.loading = false
                    state.error = action.error.message || 'Unknown error'
                })
        },
})

export const { setTenant, clearTenant, markInitialized, resetRuntime } = tenantSlice.actions
export default tenantSlice.reducer
