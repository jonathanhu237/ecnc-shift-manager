import { Home, UserCog } from "lucide-react";

export const routes = [
    {
        title: "主页",
        url: "/",
        icon: Home,
        levelRequired: 1,
    },
    {
        title: "用户管理",
        url: "/users-management",
        icon: UserCog,
        levelRequired: 3,
    },
];
