import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { ScheduleTemplateMetaType } from "@/types/schedule-template";
import { Button } from "@/components/ui/button";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, APIResponse } from "@/lib/api";
import PendingButton from "../PendingButton";
import { toast } from "sonner";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  ssm?: ScheduleTemplateMetaType;
}

export default function DeleteScheduleTemplateDialog({
  open,
  onOpenChange,
  ssm,
}: Props) {
  const queryClient = useQueryClient();

  const mutation = useMutation<APIResponse<null>>({
    mutationFn: () =>
      api.delete(`/schedule-templates/${ssm?.id}`).then((res) => res.data),
    onSuccess: (res) => {
      toast.success(res.message);
      queryClient.setQueryData(
        ["schedule-template-meta"],
        (data: ScheduleTemplateMetaType[]) => {
          return data.filter((d) => d.id !== ssm?.id);
        }
      );
      onOpenChange(false);
    },
    onError: (err) => {
      toast.error(err.message);
    },
  });

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>删除模板</AlertDialogTitle>
          <AlertDialogDescription>
            你确定要删除{ssm?.name}吗？（你无法删除已经被应用了的模板）
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div className="flex justify-end mt-2 space-x-2">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            取消
          </Button>
          {mutation.isPending ? (
            <PendingButton />
          ) : (
            <Button onClick={() => mutation.mutate()}>确定</Button>
          )}
        </div>
      </AlertDialogContent>
    </AlertDialog>
  );
}
