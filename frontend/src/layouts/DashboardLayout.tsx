import AuthGuard from "@/components/auth/AuthGuard";
import AppSidebar from "@/components/sidebar/AppSidebar";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import {
    SidebarInset,
    SidebarProvider,
    SidebarTrigger,
} from "@/components/ui/sidebar";
import { routes } from "@/lib/routes";
import { Link, Outlet, useLocation } from "react-router";

export default function DashboardLayout() {
    const location = useLocation();

    return (
        <AuthGuard>
            <SidebarProvider>
                <AppSidebar />
                <SidebarInset className="flex flex-col h-full">
                    <header className="flex h-12 items-center gap-2">
                        <SidebarTrigger className="ml-2" />
                        <Separator orientation="vertical" className="h-4" />
                        <Breadcrumb>
                            <BreadcrumbList>
                                <BreadcrumbItem>
                                    <BreadcrumbLink asChild>
                                        <Link to="/">ECNC 假勤系统</Link>
                                    </BreadcrumbLink>
                                </BreadcrumbItem>
                                <BreadcrumbSeparator />
                                <BreadcrumbItem>
                                    <BreadcrumbPage>
                                        {
                                            routes.find(
                                                (item) =>
                                                    item.url ==
                                                    location.pathname
                                            )?.title
                                        }
                                    </BreadcrumbPage>
                                </BreadcrumbItem>
                            </BreadcrumbList>
                        </Breadcrumb>
                    </header>
                    <main className="flex-1 flex">
                        <Outlet />
                    </main>
                </SidebarInset>
            </SidebarProvider>
        </AuthGuard>
    );
}
