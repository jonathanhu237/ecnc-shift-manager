import { api, APIResponse } from "@/lib/api";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { AxiosResponse } from "axios";
import { AlertCircle, Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const formSchema = z
  .object({
    oldPassword: z.string().min(0, {
      message: "旧密码不能为空",
    }),
    newPassword: z.string().min(8, {
      message: "密码不得少于 8 位",
    }),
    confirmPassword: z.string(),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "确认密码与新密码不一致",
    path: ["confirmPassword"],
  });

export default function UpdatePasswordDialog({ open, onOpenChange }: Props) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
  });
  const mutation = useMutation<
    AxiosResponse<APIResponse<null>>,
    Error,
    { oldPassword: string; newPassword: string }
  >({
    mutationFn: (data) => api.post("/me/update-password", data),
    onSuccess: (res) => {
      onOpenChange(false);
      toast.success(res.data.message);
      setErrorMessage("");
      form.reset();
    },
    onError: (err) => {
      setErrorMessage(err.message);
    },
  });
  const onSubmit = (data: z.infer<typeof formSchema>) => {
    mutation.mutate({
      oldPassword: data.oldPassword,
      newPassword: data.newPassword,
    });
  };
  const [errorMessage, setErrorMessage] = useState<string>("");

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>修改密码</AlertDialogTitle>
          <AlertDialogDescription>
            在这里修改你的密码，请勿使用过于简单的密码。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="oldPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>旧密码</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="请输入你的旧密码"
                      {...field}
                      type="password"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="newPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>新密码</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="请输入你的新密码"
                      {...field}
                      type="password"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>确认密码</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="请确认你的新密码"
                      {...field}
                      type="password"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {errorMessage && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>修改密码失败</AlertTitle>
                <AlertDescription>{errorMessage}</AlertDescription>
              </Alert>
            )}
            <div className="flex gap-2 justify-end">
              <Button
                type="button"
                onClick={() => {
                  onOpenChange(false);
                  form.reset();
                  setErrorMessage("");
                }}
                variant="outline"
              >
                取消
              </Button>

              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending ? (
                  <>
                    <Loader2 className="animate-spin" />
                    请稍等
                  </>
                ) : (
                  "确认"
                )}
              </Button>
            </div>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
