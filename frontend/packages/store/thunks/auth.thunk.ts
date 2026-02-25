import { acceptInvitationApi, requestMagicLinkApi, verifyMagicLinkApi } from "@/packages/api";
import { createAsyncThunk } from "@reduxjs/toolkit";
import { Session } from "../slices/auth.slice";
import { exit } from "process";
import { mapSubscriptionsToFlags } from "@/packages/shared/utils";
import { setFeatures } from "../slices/features.slice";

// Request Magic Link
export const requestMagicThunk = createAsyncThunk(
    'auth/requestMagicLink',
    async (payload: { email: string; }, { dispatch, rejectWithValue, fulfillWithValue }) => {
        try {
            await requestMagicLinkApi(payload)
            return fulfillWithValue(payload.email);
        } catch (error: any) {
            return rejectWithValue(error?.message || 'Failed to request magic link');
        }
    }
)

/**
 *  Verify Magic Link
 *  - CONSIDERATION: set session cookies (recommended)
 *  - May return session tokens/user; we treat it as a "login event"
 *  - Then fetch tenants (optional) and/or tenant context later after select tenant
 */
export const verifyMagicLinkThunk = createAsyncThunk(
    'auth/verifyMagicLink',
    async (payload: { token: string }, { dispatch, rejectWithValue, fulfillWithValue }) => {
        try {
            const result = await verifyMagicLinkApi(payload)

            if (result instanceof Error) {
                throw result
            }

            const data: Session = {
                user: result?.user,
                tokens: {
                    accessToken: result?.auth?.access_token,
                    refreshToken: result?.auth?.refresh_token,
                    tokenType: result?.auth?.token_type,
                    expiresIn: Number(result?.auth?.expires_in),
                },
                features: ['chat', 'hrm', 'meetings', 'tms']
            }

            const feature = await mapSubscriptionsToFlags(data.features);

            // dispatch({ type: 'featureFlags/setFlags', payload: feature })
            dispatch(setFeatures(feature));

            return data;
        } catch (error: any) {
            return rejectWithValue(error?.message || 'Failed to verify magic link');
        }
    }
)

export const acceptInvitationThunk = createAsyncThunk(
    'auth/acceptInvitation',
    async (payload: { token: string, email: string }, { dispatch, rejectWithValue, fulfillWithValue }) => {
        try {
            const result = await acceptInvitationApi(payload)

            if (result instanceof Error) {
                throw new Error(result?.message)
            }

            const data: Session = {
                user: result?.user,
                tokens: {
                    accessToken: result?.auth?.access_token,
                    refreshToken: result?.auth?.refresh_token,
                    tokenType: result?.auth?.token_type,
                    expiresIn: Number(result?.auth?.expires_in),
                },
                features: ['chat', 'hrm', 'meetings', 'tms']
            }

            const feature = await mapSubscriptionsToFlags(data.features);

            // dispatch({ type: 'featureFlags/setFlags', payload: feature })
            dispatch(setFeatures(feature));

            return data;
        } catch (error: any) {
            return rejectWithValue(error?.message || 'Failed to accept invitation.');
        }
    }
)

export const logoutThunk = createAsyncThunk('auth/logout', async (_, { dispatch }) => {
    // await logoutApi()
    // dispatch(clearUserSession())
    return true
})
