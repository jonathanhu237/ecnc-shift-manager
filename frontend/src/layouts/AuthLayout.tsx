import { Outlet } from "react-router";

export default function AuthLayout() {
    return (
        <div className="h-screen flex items-center justify-center">
            <Outlet />
        </div>
    );
}
