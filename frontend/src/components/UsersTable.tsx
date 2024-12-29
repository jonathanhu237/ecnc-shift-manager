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
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import { useState } from "react";
import UpdateUserRoleDialog from "./UpdateUserRoleDialog";
import DeleteUserDialog from "./DeleteUserDialog";
import CreateUserDialog from "./CreateUserDialog";

export default function UsersTable() {
  const { data, isPending, isError, error } = useQuery({
    queryKey: ["users"],
    queryFn: () =>
      api.get<APIResponse<UserType[]>>("/users").then((res) => res.data.data),
  });
  const [updateRoleDialogOpen, setUpdateRoleDialogOpen] = useState(false);
  const [deleteUserDialogOpen, setDeleteUserDialogOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState<UserType | null>(null);
  const [createUserDialogOpen, setCreateUserDialogOpen] = useState(false);

  const columns: ColumnDef<UserType>[] = [
    {
      accessorKey: "username",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => {
              column.toggleSorting(column.getIsSorted() === "asc");
            }}
          >
            用户名
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
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
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => {
              column.toggleSorting(column.getIsSorted() === "asc");
            }}
          >
            身份
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
    },
    {
      accessorKey: "createdAt",
      header: "创建时间",
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
              <DropdownMenuItem
                onClick={() => {
                  setUpdateRoleDialogOpen(true);
                  setCurrentUser(user);
                }}
              >
                更改身份
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-destructive"
                onClick={() => {
                  setDeleteUserDialogOpen(true);
                  setCurrentUser(user);
                }}
              >
                删除用户
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  if (isPending) {
    return null;
  }
  if (isError) {
    toast.error(error.message);
    return null;
  }

  return (
    <>
      <div className="mb-4">
        <Button onClick={() => setCreateUserDialogOpen(true)}>添加用户</Button>
      </div>
      <DataTable columns={columns} data={data} />
      <UpdateUserRoleDialog
        user={currentUser}
        open={updateRoleDialogOpen}
        onOpenChange={setUpdateRoleDialogOpen}
      />
      <DeleteUserDialog
        user={currentUser}
        open={deleteUserDialogOpen}
        onOpenChange={setDeleteUserDialogOpen}
      />
      <CreateUserDialog
        open={createUserDialogOpen}
        onOpenChange={setCreateUserDialogOpen}
      />
    </>
  );
}
