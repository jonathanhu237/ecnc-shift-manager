import { UserType } from "@/types/user";
import { DropdownMenu } from "@radix-ui/react-dropdown-menu";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { LogOut, Settings, User } from "lucide-react";
import { Link, useNavigate } from "react-router";
import {
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "../ui/sidebar";
import { api, APIResponse } from "@/lib/api";
import { AxiosResponse } from "axios";
import { toast } from "sonner";

export default function NavUser() {
    const queryClient = useQueryClient();
    const navigate = useNavigate();
    const myInfo: UserType | undefined = queryClient.getQueryData(["me"]);
    const mutation = useMutation({
        mutationFn: () => api.post("/auth/logout"),
        onSuccess: (res: AxiosResponse<APIResponse>) => {
            queryClient.clear();
            toast(res.data.message);
            navigate("/auth/login");
        },
        onError: (err) => {
            toast(err.message);
        },
    });

    return (
        <SidebarMenu>
            <SidebarMenuItem>
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <SidebarMenuButton
                            size="lg"
                            className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                        >
                            <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                                <User className="size-4" />
                            </div>
                            <div className="grid flex-1 text-left text-sm leading-tight">
                                <span className="truncate font-semibold">
                                    {myInfo?.full_name}
                                </span>
                                <span className="truncate text-xs">
                                    {myInfo?.username}({myInfo?.role})
                                </span>
                            </div>
                        </SidebarMenuButton>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent side="right" align="end">
                        <DropdownMenuGroup>
                            <Link to="/settings">
                                <DropdownMenuItem>
                                    <Settings />
                                    设置
                                </DropdownMenuItem>
                            </Link>
                        </DropdownMenuGroup>
                        <DropdownMenuSeparator />
                        <DropdownMenuGroup>
                            <DropdownMenuItem onClick={() => mutation.mutate()}>
                                <LogOut />
                                登出
                            </DropdownMenuItem>
                        </DropdownMenuGroup>
                    </DropdownMenuContent>
                </DropdownMenu>
            </SidebarMenuItem>
        </SidebarMenu>
    );
}
