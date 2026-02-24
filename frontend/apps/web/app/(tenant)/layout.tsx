import { mockUsers } from '@/data/mockData'
import { Metadata } from 'next'
import { AppProvider } from '@/packages/providers/AppProvider'
import { Sidebar, SidebarSlotProvider } from '@/apps/web/shell/sidebar'
import '../../../../global.css'

export const metadata: Metadata = {
    title: ". : : Alloy System : : .",
    description: "Organizational Operating System for the Modern Team",
    icons: {
        icon: '/favicon.ico',
    }
};

export default function RootLayout({ children }: { children: React.ReactNode }) {

    const currentUser = mockUsers["current-user"]!;

    return (
        <html lang="en" className="light">
            <body>
                <AppProvider>
                    {!currentUser?.id ? <>{children}</>: <SidebarSlotProvider>
                        <div className="h-screen flex overflow-hidden">
                            {/* Sidebar - hidden on mobile when chat is open */}
                            {currentUser?.id && <Sidebar currentUser={currentUser} />}
                            {children}
                        </div>
                    </SidebarSlotProvider>}
                </AppProvider>
            </body>
        </html>
    )
}
