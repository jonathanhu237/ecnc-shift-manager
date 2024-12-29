import { Toaster } from "@/components/ui/sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import "./index.css";
import AuthLayout from "./layouts/AuthLayout.tsx";
import DashboardLayout from "./layouts/DashboardLayout.tsx";
import DashboardIndex from "./pages/DashboardIndex.tsx";
import LoginPage from "./pages/LoginPage.tsx";
import UsersManagementPage from "./pages/UsersManagementPage.tsx";
import BlackCoreGuard from "./components/auth/BlackCoreGuard.tsx";
import ScheduleTemplatesManagement from "./pages/ShiftTemplatesManagement.tsx";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: Infinity,
      retry: false,
    },
  },
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<DashboardLayout />}>
            <Route index element={<DashboardIndex />} />
            <Route
              path="users-management"
              element={
                <BlackCoreGuard>
                  <UsersManagementPage />
                </BlackCoreGuard>
              }
            />
            <Route
              path="shift-templates-management"
              element={
                <BlackCoreGuard>
                  <ScheduleTemplatesManagement />
                </BlackCoreGuard>
              }
            />
          </Route>
          <Route path="auth" element={<AuthLayout />}>
            <Route path="login" element={<LoginPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
      <ReactQueryDevtools initialIsOpen={false} buttonPosition="top-right" />
      <Toaster />
    </QueryClientProvider>
  </StrictMode>
);
