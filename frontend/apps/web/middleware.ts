import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(req: NextRequest) {
    const token = req.cookies.get('accessToken')
    if (!token) return NextResponse.redirect(new URL('/auth/login', req.url))

    const tenantId = req.headers.get('x-tenant-id') || null
    if (!tenantId) return NextResponse.redirect(new URL('/select-tenant', req.url))

    return NextResponse.next()
}

export const config = {
    matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}
