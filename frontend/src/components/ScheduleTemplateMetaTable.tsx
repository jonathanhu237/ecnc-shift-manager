import { api, APIResponse } from "@/lib/api";
import { ScheduleTemplateMetaType } from "@/types/schedule-template";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";
import DataTable from "./DataTable";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { MoreHorizontal } from "lucide-react";

export default function ScheduleTemplateMetaTable() {
  const { data, isPending, isError, error } = useQuery({
    queryKey: ["schedule-template-meta"],
    queryFn: () =>
      api
        .get<APIResponse<ScheduleTemplateMetaType[]>>("/schedule-template-meta")
        .then((res) => res.data.data),
  });

  const columns: ColumnDef<ScheduleTemplateMetaType>[] = [
    {
      accessorKey: "name",
      header: "模板名",
    },
    {
      accessorKey: "description",
      header: "模板描述",
    },
    {
      accessorKey: "createdAt",
      header: "模板创建时间",
    },
    {
      id: "action",
      cell: () => {
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>操作</DropdownMenuLabel>
              <DropdownMenuItem>查看详情</DropdownMenuItem>
              <DropdownMenuItem>更改描述</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-destructive">
                删除模板
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
      <DataTable columns={columns} data={data} />
    </>
  );
}
