import { api, APIResponse } from "@/lib/api";
import { UserType } from "@/types/user";
import { DropdownMenu } from "@radix-ui/react-dropdown-menu";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosResponse } from "axios";
import { KeyRound, LogOut, User } from "lucide-react";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import {
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "../ui/sidebar";
import { useState } from "react";
import UpdatePasswordDialog from "@/components/dialog/UpdatePasswordDialog";

export default function NavUser() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const myInfo: UserType | undefined = queryClient.getQueryData(["me"]);
  const mutation = useMutation({
    mutationFn: () => api.post("/auth/logout"),
    onSuccess: (res: AxiosResponse<APIResponse<UserType>>) => {
      queryClient.clear();
      toast.success(res.data.message);
      navigate("/auth/login");
    },
    onError: (err) => {
      toast.error(err.message);
    },
  });
  const [dialogOpen, setDialogOpen] = useState(false);

  return (
    <>
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
                    {myInfo?.fullName}
                  </span>
                  <span className="truncate text-xs">
                    {myInfo?.username}({myInfo?.role})
                  </span>
                </div>
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent side="right" align="end">
              <DropdownMenuGroup>
                <DropdownMenuItem onClick={() => setDialogOpen(true)}>
                  <KeyRound />
                  修改密码
                </DropdownMenuItem>
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
      <UpdatePasswordDialog open={dialogOpen} onOpenChange={setDialogOpen} />
    </>
  );
}
