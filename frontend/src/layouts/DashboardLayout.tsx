import AuthGuard from "@/components/auth/AuthGuard";
import AppSidebar from "@/components/sidebar/AppSidebar";
import {
    SidebarInset,
    SidebarProvider,
    SidebarTrigger,
} from "@/components/ui/sidebar";
import { Outlet } from "react-router";

export default function DashboardLayout() {
    return (
        <AuthGuard>
            <SidebarProvider>
                <AppSidebar />
                <SidebarInset>
                    <main>
                        <SidebarTrigger />
                        <Outlet />
                    </main>
                </SidebarInset>
            </SidebarProvider>
        </AuthGuard>
    );
}
