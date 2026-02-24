import { User } from '@/packages/types'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { acceptInvitationThunk, logoutThunk, requestMagicThunk, verifyMagicLinkThunk } from '../thunks';
import { IModuleKey } from '@/packages/shared/types';
import { maskEmail } from '@/packages/shared/utils';

// export type User = {
//     id: string
//     email: string
//     name: string
// }

export type ITokens = {
    accessToken: string;
    refreshToken: string;
    expiresIn?: number;
    tokenType?: string;
}

export type Session = {
    user: User;
    tokens: ITokens;
    features: IModuleKey[];
}

export interface AuthState {
    session: Session | null
    user: User | null
    isAuthenticated: boolean
    loading: boolean
    error?: string
    tokens?: ITokens | null
    magicLink: {
        requested: boolean
        verified: boolean
        emailMasked?: string
    }
}

const initialState: AuthState = {
    user: null,
    session: null,
    isAuthenticated: false,
    loading: false,
    error: undefined,
    tokens: null,
    magicLink: {
        requested: false,
        verified: false,
        emailMasked: undefined
    }
}

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        loginSuccess(state, action: PayloadAction<User>) {
            state.user = action.payload
            state.isAuthenticated = true
        },
        logout(state) {
            state.user = null
            state.isAuthenticated = false
        },
        setSession(state, action: PayloadAction<Session>) {
            state.user = action.payload.user;
            state.tokens = action.payload.tokens;
            state.session = action.payload;
            state.isAuthenticated = Boolean(action.payload)
        },
        setUser(state, action: PayloadAction<User>) {
            state.user = action.payload
            state.isAuthenticated = true
        },
        clearError(state) {
            state.error = undefined;
        },
        clearAuthTransient(state) {
            state.error = undefined;
            state.magicLink.requested = false;
            state.magicLink.verified = true;
            state.magicLink.emailMasked = undefined;
        }
        // setSession(state, action: PayloadAction<Session>) {
        //     state.user = action.payload.user
        //     state.isAuthenticated = true
        // },
    },
    extraReducers: builder => {
        // ---- REQUEST MAGIC LINK ----
        builder
            .addCase(requestMagicThunk.pending, (state, action) => {
                state.loading = true
                state.error = undefined
            })
            .addCase(requestMagicThunk.fulfilled, (state, action) => {
                state.loading = false
                state.magicLink.requested = true
                state.magicLink.emailMasked = maskEmail(action.payload.toString())
            })
            .addCase(requestMagicThunk.rejected, (state, action) => {
                state.loading = false
                state.error = typeof action.payload === 'string' && action.payload  || 'Unable to send link. Try again.'
            })

        // ---- VERIFY MAGIC LINK ----
        builder
            .addCase(verifyMagicLinkThunk.pending, state => {
                state.loading = true
                state.error = undefined
            })
            .addCase(verifyMagicLinkThunk.fulfilled, (state, action) => {
                state.loading = false
                state.magicLink.verified = true
                // If API returns a session (optional), store it
                const session = (action.payload as any)?.session as Session | undefined
                if (session) {
                    state.user = session.user
                    state.tokens = session.tokens;
                    state.session = session;
                    state.isAuthenticated = true
                    // Dispatch to features
                } else {
                    // At minimum, cookies should now be set.
                    // We mark authenticated true only if you also have /me endpoint.
                    // If you want, set to true here:
                    state.isAuthenticated = true
                }

                // update tenants if fetched successfully
                // if (action.payload.tenants?.length) {
                //     state.tenants = action.payload.tenants
                // }
            })
            .addCase(verifyMagicLinkThunk.rejected, (state, action) => {
                state.loading = false
                state.error =
                    typeof action.payload === 'string' && action.payload  || action.error.message || 'Verification failed. The link may have expired.'
            })

        // ---- ACCEPT INVITATION ----
        builder
            .addCase(acceptInvitationThunk.pending, state => {
                state.loading = true
                state.error = undefined
            })
            .addCase(acceptInvitationThunk.fulfilled, (state, action) => {
                state.loading = false
                state.magicLink.verified = true
                // If API returns a session (optional), store it
                const session = (action.payload as any)?.session as Session | undefined
                if (session) {
                    state.user = session.user
                    state.tokens = session.tokens;
                    state.session = session;
                    state.isAuthenticated = true
                    // Dispatch to features
                } else {
                    // At minimum, cookies should now be set.
                    // We mark authenticated true only if you also have /me endpoint.
                    // If you want, set to true here:
                    state.isAuthenticated = true
                }

                // update tenants if fetched successfully
                // if (action.payload.tenants?.length) {
                //     state.tenants = action.payload.tenants
                // }
            })
            .addCase(acceptInvitationThunk.rejected, (state, action) => {
                state.loading = false
                state.error =
                    typeof action.payload === 'string' && action.payload  || action.error.message || 'Accept invitation failed.'
            })

        // ---- LOGOUT ----
        builder
            .addCase(logoutThunk.pending, state => {
                state.loading = true
                state.error = undefined
            })
            .addCase(logoutThunk.fulfilled, state => {
                state.loading = false
                state.user = null
                state.session = null
                state.tokens = null
                state.isAuthenticated = false
                state.magicLink.requested = false
                state.magicLink.verified = false
                state.magicLink.emailMasked = undefined
            })
            .addCase(logoutThunk.rejected, (state, action) => {
                state.loading = false
                state.error = action.error.message || 'Logout failed'
            })
    }
})

export const { loginSuccess, logout, setSession, setUser, clearAuthTransient, clearError } = authSlice.actions
export default authSlice.reducer
