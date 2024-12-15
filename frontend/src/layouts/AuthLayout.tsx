import GuestGuard from "@/components/auth/GuestGuard";
import { Outlet } from "react-router";

export default function AuthLayout() {
    return (
        <GuestGuard>
            <div className="h-screen flex items-center justify-center">
                <Outlet />
            </div>
        </GuestGuard>
    );
}
