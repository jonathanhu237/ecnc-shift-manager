import AuthGuard from "@/components/auth/AuthGuard";
import { Outlet } from "react-router";

export default function DashboardLayout() {
    return (
        <AuthGuard>
            <Outlet />
        </AuthGuard>
    );
}
