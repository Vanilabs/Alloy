import { Metadata } from 'next'
import "../global.css"
import { Providers } from './(public)/provider';

export const metadata: Metadata = {
    title: ". : : Alloy System : : .",
    description: "Organizational Operating System for the Modern Team",
    icons: {
        icon: '/favicon.ico',
    }
};

export default function PublicRootLayout({ children }: { children: React.ReactNode }) {
    return (
        <html lang="en" className="light">
            <body>
                <Providers>
                    {children}
                </Providers>
            </body>
        </html>
    )
}
