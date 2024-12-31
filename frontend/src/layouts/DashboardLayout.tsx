import AuthGuard from "@/components/auth/AuthGuard";
import AppSidebar from "@/components/sidebar/AppSidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { Outlet } from "react-router";

export default function DashboardLayout() {
  return (
    <AuthGuard>
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset className="p-4">
          <Outlet />
        </SidebarInset>
      </SidebarProvider>
    </AuthGuard>
  );
}
