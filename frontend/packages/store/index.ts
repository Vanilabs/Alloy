import { configureStore } from '@reduxjs/toolkit'
// import storage from "redux-persist/lib/storage";
import { auth, chat, features, runtime, tenant } from './slices'
// import { persistReducer, persistStore } from "redux-persist";
import { rootReducer } from './reducer'
// import tenant from './slices/tenant.slice'
// import permissions from './slices/permissions.slice'
// import features from './slices/features.slice'

// const persistConfig = {
//     key: "alloy",
//     storage,
//     version: 1,
//     whitelist: ["auth", "tenant", "features"],
//     // blacklist: ["runtime", "chat"],
// };

// const persistedReducer = persistReducer(persistConfig, rootReducer);

export const store = configureStore({
    reducer: {
        runtime: runtime.default,
        auth: auth.default,
        chat: chat.default,
        features: features.default,
        tenant: tenant.default,
    },
    // devTools: process.env.NODE_ENV !== 'production',
    middleware: (getDefaultMiddleware) => getDefaultMiddleware({
        serializableCheck: false
    })
})

// export const persistor = persistStore(store);

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
