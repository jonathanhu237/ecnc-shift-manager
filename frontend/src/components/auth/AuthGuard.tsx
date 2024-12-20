import { api, APIResponse } from "@/lib/api";
import { UserType } from "@/types/user";
import { useQuery } from "@tanstack/react-query";
import { PropsWithChildren } from "react";
import { useNavigate } from "react-router";

export default function AuthGuard({ children }: PropsWithChildren) {
    const { isPending, isError } = useQuery({
        queryKey: ["me"],
        queryFn: () =>
            api.get<APIResponse<UserType>>("/me").then((res) => res.data.data),
    });

    const navigate = useNavigate();

    if (isPending) return null;
    else if (isError) {
        navigate("/auth/login");
    } else {
        return children;
    }
}
