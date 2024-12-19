import { api, APIResponse } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { UserType } from "@/types/user";
import { ColumnDef } from "@tanstack/react-table";
import DataTable from "./DataTable";
import { toast } from "sonner";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Button } from "./ui/button";
import { MoreHorizontal } from "lucide-react";

const columns: ColumnDef<UserType>[] = [
    {
        accessorKey: "username",
        header: "用户名",
    },
    {
        accessorKey: "fullName",
        header: "姓名",
    },
    {
        accessorKey: "email",
        header: "邮箱",
    },
    {
        accessorKey: "role",
        header: "身份",
    },
    {
        id: "action",
        cell: ({ row }) => {
            const user = row.original;

            return (
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                            <MoreHorizontal className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuLabel>操作</DropdownMenuLabel>
                        <DropdownMenuItem>更改身份</DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem className="text-destructive">
                            删除用户
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            );
        },
    },
];

export default function UsersTable() {
    const { data, isPending, isError, error } = useQuery({
        queryKey: ["users"],
        queryFn: () =>
            api
                .get<APIResponse<UserType[]>>("/users")
                .then((res) => res.data.data),
    });

    if (isPending) {
        return null;
    }
    if (isError) {
        toast(error.message);
        return null;
    }

    return <DataTable columns={columns} data={data} />;
}
