import { DropdownMenu } from "@radix-ui/react-dropdown-menu";
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "../ui/sidebar";
import { DropdownMenuTrigger } from "../ui/dropdown-menu";
import { useQueryClient } from "@tanstack/react-query";
import { UserType } from "@/types/user";
import { User } from "lucide-react";

export default function NavUser() {
    const queryClient = useQueryClient();
    const myInfo: UserType | undefined = queryClient.getQueryData(["me"]);

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
                </DropdownMenu>
            </SidebarMenuItem>
        </SidebarMenu>
    );
}
