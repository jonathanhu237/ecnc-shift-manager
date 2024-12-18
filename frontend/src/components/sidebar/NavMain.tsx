import { Home, Settings } from "lucide-react";
import {
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "../ui/sidebar";
import { Link, useLocation } from "react-router";

const items = [
    {
        title: "主页",
        url: "/",
        icon: Home,
    },
    {
        title: "设置",
        url: "/settings",
        icon: Settings,
    },
];

export default function NavMain() {
    const location = useLocation();

    return (
        <SidebarGroup>
            <SidebarGroupLabel>应用</SidebarGroupLabel>
            <SidebarGroupContent>
                <SidebarMenu>
                    {items.map((item) => (
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
