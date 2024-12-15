import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { BrowserRouter, Route, Routes } from "react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import AuthLayout from "./layouts/AuthLayout.tsx";
import LoginPage from "./pages/LoginPage.tsx";
import { Toaster } from "@/components/ui/sonner";
import DashboardLayout from "./layouts/DashboardLayout.tsx";
import DashboardIndex from "./pages/DashboardIndex.tsx";

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
                    </Route>
                    <Route path="auth" element={<AuthLayout />}>
                        <Route path="login" element={<LoginPage />} />
                    </Route>
                </Routes>
            </BrowserRouter>
            <ReactQueryDevtools initialIsOpen={false} />
            <Toaster />
        </QueryClientProvider>
    </StrictMode>
);
