import axios from 'axios';

type ApiFetchOptions = RequestInit & {
    retryOn401?: boolean,
}

class ApiError extends Error {
    status: number;
    data: any;
    constructor(message: string, status: number, data: any) {
        super(message);
        this.name = "ApiError";
        this.status = status;
        this.data = data;
    }
}

async function safeParseJson(res: Response) {
    const text = await res.text();
    if (!text) return null;
    try {
        return JSON.parse(text);
    } catch {
        return text;
    }
}

function pickErrorMessage(data: any, fallback: string) {
    if (!data) return fallback;
    if (typeof data === "string") return data;
    return data.message || data.error || data.details || fallback;
}

let refreshing: Promise<void> | null = null;

function buildHeaders(options?: RequestInit) {
    const headers = new Headers(options?.headers)

    // Only set JSON content-type if body is not FormData
    const isFormData = typeof FormData !== 'undefined' && options?.body instanceof FormData
    if (!isFormData && !headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
    }

    // ngrok bypass
    if (typeof window === "undefined") {
        headers.set("ngrok-skip-browser-warning", "true");
    }

    return headers
}

async function refreshTokenOnce(arg?: Record<string, any>) {
    try {
        if (!refreshing) {
            refreshing = (async () => {
                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}`, {
                    method: 'POST',
                    credentials: 'include',
                    // body
                })
                if (!res.ok) throw new Error('Token refresh failed');
            })().finally(() => {
                refreshing = null
            })
        }

        return refreshing;
    } catch (error) {
        throw error;
    }
}

export async function apiFetch(path: string, options: ApiFetchOptions) {
    const doRequest = async (withCreds: boolean) => {
        const res = await fetch(path, {
            ...options,
            credentials: withCreds ? "include" : options.credentials,
            headers: buildHeaders(options),
        });

        const data = await safeParseJson(res);

        if (res.ok) return data;

        if (res.status === 401) {
            throw new ApiError("Unauthorized", 401, data);
        }

        throw new ApiError(
            pickErrorMessage(data, `Request failed (${res.status})`),
            res.status,
            data
        );
    };
    try {
        return await doRequest(false);
    } catch (err: any) {
        const shouldRetry =
            err?.name === "ApiError" &&
            err?.status === 401 &&
            options.retryOn401 !== false;

        if (!shouldRetry) throw err;

        await refreshTokenOnce();
        return await doRequest(true);
    }
}

// exponential backoff helper
export async function withRetry<T>(fn: () => Promise<T>, attempts = 3, baseDelayMs = 200) {
    let delay = baseDelayMs
    for (let i = 0; i < attempts; i++) {
        try {
            return await fn()
        } catch (e) {
            if (i === attempts - 1) throw e
            await new Promise(r => setTimeout(r, delay))
            delay *= 2
        }
    }

    throw new Error('Resource exhausted')
}
