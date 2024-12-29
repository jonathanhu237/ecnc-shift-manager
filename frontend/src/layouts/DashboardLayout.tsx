import AuthGuard from "@/components/auth/AuthGuard";
import AppSidebar from "@/components/sidebar/AppSidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { routes } from "@/lib/routes";
import { Outlet, useLocation } from "react-router";

export default function DashboardLayout() {
  const location = useLocation();

  return (
    <AuthGuard>
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset className="p-4">
          <header className="mb-4 flex items-center font-bold">
            {routes.find((route) => route.url === location.pathname)?.title}
          </header>
          <Outlet />
        </SidebarInset>
      </SidebarProvider>
    </AuthGuard>
  );
}
