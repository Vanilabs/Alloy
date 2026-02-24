import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { apiFetch, withRetry } from "..";
import { AUTH_LINKS } from '@/packages/constants';

export async function requestMagicLinkApi(payload: { email: string }) {
    const res = await apiFetch(AUTH_LINKS.MAGIC_LINK, {
        method: 'POST',
        body: JSON.stringify(payload),
        retryOn401: false
    })
    return await res;
}

export async function verifyMagicLinkApi(payload: { token: string }) {
    const res = await apiFetch(AUTH_LINKS.MAGIC_LINK_VERIFY, {
        method: 'POST',
        body: JSON.stringify(payload),
        retryOn401: false,
    })
    return await res;
}

export async function acceptInvitationApi(payload: { token: string, email: string }) {
    const res = await apiFetch(AUTH_LINKS.AUTH_ACCEPT_INVITE, {
        method: 'PATCH',
        body: JSON.stringify(payload),
        retryOn401: false,
    })
    return await res;
}
