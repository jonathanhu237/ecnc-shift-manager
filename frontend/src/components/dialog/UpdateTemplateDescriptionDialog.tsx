import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { api, APIResponse } from "@/lib/api";
import { ScheduleTemplateMetaType } from "@/types/schedule-template";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import PendingButton from "../PendingButton";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  ssm?: ScheduleTemplateMetaType;
}

const formSchema = z.object({
  description: z.string(),
});

export default function UpdateTemplateDescriptionDialog({
  open,
  onOpenChange,
  ssm,
}: Props) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      description: "",
    },
  });

  const onSubmit = (value: z.infer<typeof formSchema>) => {
    mutation.mutate(value);
  };
  const queryClient = useQueryClient();

  const mutation = useMutation<
    APIResponse<ScheduleTemplateMetaType>,
    Error,
    { description: string }
  >({
    mutationFn: (data) =>
      api
        .post(`/schedule-template-meta/${ssm?.id}/update-description`, data)
        .then((res) => res.data),
    onError: (err) => {
      toast.error(err.message);
    },
    onSuccess: (res) => {
      queryClient.setQueryData(
        ["schedule-template-meta"],
        (data: ScheduleTemplateMetaType[]) => {
          return data.map((item) => {
            return item.id === ssm?.id ? res.data : item;
          });
        }
      );
      toast.success(res.message);
      form.reset();
      onOpenChange(false);
    },
  });

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>修改班表模板描述</AlertDialogHeader>
        <AlertDialogDescription>
          请在下方输入{ssm?.name}的新模板描述。
        </AlertDialogDescription>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>新的模板描述</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入新的模板描述" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <div className="flex justify-end space-x-2 mt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  onOpenChange(false);
                  form.reset();
                }}
              >
                取消
              </Button>
              {mutation.isPending ? (
                <PendingButton />
              ) : (
                <Button type="submit">确认修改</Button>
              )}
            </div>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
