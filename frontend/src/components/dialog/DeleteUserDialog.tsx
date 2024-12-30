import { api, APIResponse } from "@/lib/api";
import { UserType } from "@/types/user";
import { AlertDialogTitle } from "@radix-ui/react-alert-dialog";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import PendingButton from "@/components/PendingButton";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";

interface Props {
  user: UserType | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export default function DeleteUserDialog({ user, open, onOpenChange }: Props) {
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: () => api.delete<APIResponse<null>>(`/users/${user?.id}`),
    onSuccess: (res) => {
      queryClient.setQueryData(["users"], (data: UserType[]) => {
        return data.filter((u) => u.id !== user?.id);
      });
      toast.success(res.data.message);
      onOpenChange(false);
    },
    onError: (err) => {
      toast.error(err.message);
      onOpenChange(false);
    },
  });

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>删除用户</AlertDialogTitle>
          <AlertDialogDescription>
            你确定要删除{user?.fullName}({user?.username}
            )吗？删除用户之后，与其相关的所有数据都会被一起删除。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div className="flex items-center justify-end gap-2 mt-4">
          <Button
            variant="outline"
            disabled={mutation.isPending}
            onClick={() => onOpenChange(false)}
          >
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
