import { api, APIResponse } from "@/lib/api";
import { ScheduleTemplateMetaType } from "@/types/schedule-template";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";
import DataTable from "@/components/table/DataTable";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { MoreHorizontal } from "lucide-react";
import { useState } from "react";
import { format } from "date-fns";
import UpdateTemplateDescriptionDialog from "@/components/dialog/UpdateTemplateDescriptionDialog";

export default function ScheduleTemplateMetaTable() {
  const { data, isPending, isError, error } = useQuery({
    queryKey: ["schedule-template-meta"],
    queryFn: () =>
      api
        .get<APIResponse<ScheduleTemplateMetaType[]>>("/schedule-template-meta")
        .then((res) => res.data.data),
  });
  const [globalFilter, setGlobalFilter] = useState("");
  const [updateDescriptionDialogOpen, setUpdateDescriptionDialogOpen] =
    useState(false);
  const [currentssm, setCurrentssm] = useState<
    ScheduleTemplateMetaType | undefined
  >(undefined);

  const columns: ColumnDef<ScheduleTemplateMetaType>[] = [
    {
      accessorKey: "name",
      header: () => <div className="text-center">模板名</div>,
      cell: ({ row }) => (
        <div className="text-center text-nowrap">{row.original.name}</div>
      ),
    },
    {
      accessorKey: "description",
      header: () => <div className="text-center">模板描述</div>,
      cell: ({ row }) => (
        <div className="text-center">{row.original.description}</div>
      ),
    },
    {
      accessorKey: "createdAt",
      header: () => <div className="text-center">模板创建时间</div>,
      cell: ({ row }) => {
        const createdAt = row.original.createdAt;
        const formattedDate = format(
          new Date(createdAt),
          "yyyy-MM-dd HH:mm:ss"
        );

        return <div className="text-center text-nowrap">{formattedDate}</div>;
      },
    },
    {
      id: "action",
      header: () => <div className="text-center">操作</div>,
      cell: ({ row }) => {
        const ssm = row.original;

        return (
          <div className="flex items-center justify-center">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>操作</DropdownMenuLabel>
                <DropdownMenuItem>查看详情</DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => {
                    setCurrentssm(ssm);
                    setUpdateDescriptionDialogOpen(true);
                  }}
                >
                  更改描述
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="text-destructive">
                  删除模板
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
        actions={<Button>创建班表模板</Button>}
      />
      <UpdateTemplateDescriptionDialog
        open={updateDescriptionDialogOpen}
        onOpenChange={setUpdateDescriptionDialogOpen}
        ssm={currentssm}
      />
    </>
  );
}
