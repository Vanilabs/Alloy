import { User } from '@/packages/types'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

// export type User = {
//     id: string
//     email: string
//     name: string
// }

export type Session = {
    user: User
    tokens: {
        accessToken: string
        refreshToken: string
    }
}

export interface AuthState {
    user: User | null
    isAuthenticated: boolean
}

const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
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
        setUser(state, action: PayloadAction<User>) {
            state.user = action.payload
            state.isAuthenticated = true
        },
        setSession(state, action: PayloadAction<Session>) {
            state.user = action.payload.user
            state.isAuthenticated = true
        },
    },
})

export const { loginSuccess, logout, setSession, setUser } = authSlice.actions
export default authSlice.reducer
