import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const PUBLIC_PATHS = ['/login']
const SYSTEM_PATHS = ['/select-tenant']
const ADMIN_PREFIX = '/admin'

function isPublic(pathname: string) {
    return PUBLIC_PATHS.some(p => pathname === p || pathname.startsWith(p + '/'))
}

function isSystem(pathname: string) {
    return SYSTEM_PATHS.some(p => pathname === p || pathname.startsWith(p + '/'))
}

export function middleware(req: NextRequest) {
    const { pathname, searchParams } = req.nextUrl

    // Redirect legacy /invite to new accept page
    if (pathname === '/auth/invite') {
        const token = searchParams.get('token') ?? ''

        const url = req.nextUrl.clone()
        url.pathname = '/invitations/accept'
        url.search = token ? `?token=${encodeURIComponent(token)}` : ''
        return NextResponse.redirect(url)
    }

    return NextResponse.next()
}

export const config = {
    matcher: ['/((?!_next/static|_next/image|favicon.ico).*)', '/invite'],
}
