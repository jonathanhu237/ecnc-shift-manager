import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "../ui/sidebar";
import { Link, useLocation } from "react-router";
import { routes } from "@/lib/routes";
import { useQueryClient } from "@tanstack/react-query";
import { UserType } from "@/types/user";

export default function NavMain() {
  const location = useLocation();
  const queryClient = useQueryClient();
  const myInfo: UserType | undefined = queryClient.getQueryData(["me"]);

  if (!myInfo) {
    throw new Error(
      "NavMain component must be rendered after user is logged in"
    );
  }

  return (
    <SidebarGroup>
      <SidebarGroupLabel>应用</SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          {routes
            .filter((item) => item.levelRequired <= myInfo.level)
            .map((item) => (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton
                  asChild
                  isActive={location.pathname == item.url}
                >
                  <Link to={item.url}>
                    <item.icon />
                    <span>{item.title}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}
