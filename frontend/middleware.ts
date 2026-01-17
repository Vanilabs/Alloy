import { NextRequest, NextResponse } from 'next/server'

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
    // const { pathname } = req.nextUrl

    // // allow next assets
    // if (
    //     pathname.startsWith('/_next') ||
    //     pathname.startsWith('/favicon.ico') ||
    //     pathname.startsWith('/public')
    // ) {
    //     return NextResponse.next()
    // }

    // // Admin can have separate auth later; for now allow it to run independently
    // if (pathname.startsWith(ADMIN_PREFIX)) {
    //     return NextResponse.next()
    // }

    // const session = req.cookies.get('alloy_session')?.value
    // const tenant = req.cookies.get('alloy_tenant')?.value

    // if (!session && !isPublic(pathname)) {
    //     return NextResponse.redirect(new URL('/login', req.url))
    // }

    // // logged in but no tenant picked yet → select tenant
    // if (session && !tenant && !isSystem(pathname) && !isPublic(pathname)) {
    //     return NextResponse.redirect(new URL('/select-tenant', req.url))
    // }

    // return NextResponse.next()
}

export const config = {
    matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}
