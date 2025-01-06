import { Toaster } from "@/components/ui/sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import "./index.css";
import AuthLayout from "@/layouts/AuthLayout.tsx";
import DashboardLayout from "@/layouts/DashboardLayout.tsx";
import DashboardIndexPage from "@/pages/DashboardIndexPage";
import LoginPage from "@/pages/LoginPage.tsx";
import UsersManagementPage from "@/pages/UsersManagementPage.tsx";
import ScheduleTemplatesManagementPage from "@/pages/ShiftTemplatesManagementPage";
import CreateScheduleTemplatePage from "@/pages/CreateScheduleTemplatePage";
import ScheduleTemplateDetailsPage from "@/pages/ScheduleTemplateDetailsPage";
import BlackCoreGuardLayout from "@/layouts/BlackCoreGuardLayout";

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
            <Route index element={<DashboardIndexPage />} />
            <Route path="users-management" element={<BlackCoreGuardLayout />}>
              <Route index element={<UsersManagementPage />} />
            </Route>
            <Route
              path="schedule-templates-management"
              element={<BlackCoreGuardLayout />}
            >
              <Route index element={<ScheduleTemplatesManagementPage />} />
              <Route path=":id" element={<ScheduleTemplateDetailsPage />} />
              <Route path="create" element={<CreateScheduleTemplatePage />} />
            </Route>
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
