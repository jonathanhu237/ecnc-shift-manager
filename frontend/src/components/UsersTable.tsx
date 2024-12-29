import { api, APIResponse } from "@/lib/api";
import { UserType } from "@/types/user";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import CreateUserDialog from "./CreateUserDialog";
import DataTable from "./DataTable";
import DeleteUserDialog from "./DeleteUserDialog";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import UpdateUserRoleDialog from "./UpdateUserRoleDialog";
import { format } from "date-fns";

export default function UsersTable() {
  const { data, isPending, isError, error } = useQuery({
    queryKey: ["users"],
    queryFn: () =>
      api.get<APIResponse<UserType[]>>("/users").then((res) => res.data.data),
  });
  const [globalFilter, setGlobalFilter] = useState("");
  const [updateRoleDialogOpen, setUpdateRoleDialogOpen] = useState(false);
  const [deleteUserDialogOpen, setDeleteUserDialogOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState<UserType | null>(null);
  const [createUserDialogOpen, setCreateUserDialogOpen] = useState(false);

  const columns: ColumnDef<UserType>[] = [
    {
      accessorKey: "username",
      header: () => <div className="text-center">用户名</div>,
      cell: ({ row }) => (
        <div className="text-center">{row.original.username}</div>
      ),
    },
    {
      accessorKey: "fullName",
      header: () => <div className="text-center">姓名</div>,
      cell: ({ row }) => (
        <div className="text-center">{row.original.fullName}</div>
      ),
    },
    {
      accessorKey: "email",
      header: () => <div className="text-center">邮箱</div>,
      cell: ({ row }) => (
        <div className="text-center">{row.original.email}</div>
      ),
    },
    {
      accessorKey: "role",
      header: () => <div className="text-center">身份</div>,
      cell: ({ row }) => <div className="text-center">{row.original.role}</div>,
    },
    {
      accessorKey: "createdAt",
      header: () => <div className="text-center">创建时间</div>,
      cell: ({ row }) => {
        const createdAt = row.original.createdAt;
        const formattedDate = format(
          new Date(createdAt),
          "yyyy-MM-dd HH:mm:ss"
        );

        return <div className="text-center">{formattedDate}</div>;
      },
    },
    {
      id: "action",
      header: () => <div className="text-center">操作</div>,
      cell: ({ row }) => {
        const user = row.original;

        return (
          <div className="flex justify-center">
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
          </div>
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
      <DataTable
        columns={columns}
        data={data}
        globalFilter={globalFilter}
        setGlobalFilter={setGlobalFilter}
        actions={
          <Button onClick={() => setCreateUserDialogOpen(true)}>
            添加用户
          </Button>
        }
      />
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
