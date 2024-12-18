import { UserType } from "@/types/user";
import { useQueryClient } from "@tanstack/react-query";
import { PropsWithChildren } from "react";

export default function BlackCoreGuard({ children }: PropsWithChildren) {
    const queryClient = useQueryClient();
    const myInfo: UserType | undefined = queryClient.getQueryData(["me"]);

    if (!myInfo) {
        throw new Error(
            "BlackCoreGuard component must be rendered after user is logged in"
        );
    }

    if (myInfo.level < 3) {
        return null;
    }

    return children;
}
