import { Metadata } from 'next'
import "../global.css"

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
                {children}
            </body>
        </html>
    )
}
