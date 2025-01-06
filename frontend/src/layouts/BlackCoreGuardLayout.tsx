import BlackCoreGuard from "@/components/auth/BlackCoreGuard";
import { Outlet } from "react-router";

export default function BlackCoreGuardLayout() {
  return (
    <BlackCoreGuard>
      <Outlet />
    </BlackCoreGuard>
  );
}
