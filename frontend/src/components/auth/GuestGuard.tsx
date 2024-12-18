import { api, APIResponse } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { PropsWithChildren } from "react";
import { Navigate } from "react-router";

export default function GuestGuard({ children }: PropsWithChildren) {
    const { isPending, isError } = useQuery({
        queryKey: ["me"],
        queryFn: () => api.get<APIResponse>("/me").then((res) => res.data),
    });

    if (isPending) return null;
    else if (isError) {
        return children;
    } else {
        return <Navigate to="/" />;
    }
}
