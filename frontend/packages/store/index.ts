import { configureStore } from '@reduxjs/toolkit'
import { auth, chat, features, runtime, tenant } from './slices'
// import tenant from './slices/tenant.slice'
// import permissions from './slices/permissions.slice'
// import features from './slices/features.slice'

export const store = configureStore({
    reducer: {
        runtime: runtime.default,
        auth: auth.default,
        chat: chat.default,
        features: features.default,
        tenant: tenant.default,
    },
    devTools: process.env.NODE_ENV !== 'production',
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
