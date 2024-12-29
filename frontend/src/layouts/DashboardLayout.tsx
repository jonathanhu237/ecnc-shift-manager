import AuthGuard from "@/components/auth/AuthGuard";
import AppSidebar from "@/components/sidebar/AppSidebar";
import { SidebarProvider } from "@/components/ui/sidebar";
import { Outlet } from "react-router";

export default function DashboardLayout() {
  return (
    <AuthGuard>
      <SidebarProvider>
        <AppSidebar />
        <Outlet />
      </SidebarProvider>
    </AuthGuard>
  );
}
