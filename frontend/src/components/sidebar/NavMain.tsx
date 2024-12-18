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

export default function NavMain() {
    const location = useLocation();

    return (
        <SidebarGroup>
            <SidebarGroupLabel>应用</SidebarGroupLabel>
            <SidebarGroupContent>
                <SidebarMenu>
                    {routes.map((item) => (
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
