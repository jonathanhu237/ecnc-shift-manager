import { z } from "zod";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogHeader,
    AlertDialogTitle,
} from "./ui/alert-dialog";
import { Button } from "./ui/button";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "./ui/form";
import { Input } from "./ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "./ui/select";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, APIResponse } from "@/lib/api";
import { UserType } from "@/types/user";
import { toast } from "sonner";
import PendingButton from "./PendingButton";

interface Props {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const formSchema = z.object({
    username: z.string().min(1, { message: "用户名不能为空" }),
    email: z
        .string()
        .min(1, { message: "邮箱不能为空" })
        .email({ message: "邮箱格式不正确" }),
    fullName: z.string().min(1, { message: "姓名不能为空" }),
    role: z.string(),
});

export default function CreateUserDialog({ open, onOpenChange }: Props) {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            username: "",
            email: "",
            fullName: "",
            role: "普通助理",
        },
    });
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationFn: (data: z.infer<typeof formSchema>) =>
            api.post<APIResponse<UserType>>("/users", data),
        onSuccess: (res) => {
            onOpenChange(false);
            form.reset();
            toast.success(res.data.message);
            queryClient.setQueryData(["users"], (data: UserType[]) => [
                ...data,
                res.data.data,
            ]);
        },
        onError: (err) => {
            toast.error(err.message);
        },
    });

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>添加用户</AlertDialogTitle>
                    <AlertDialogDescription>
                        成功添加用户后，将会给新用户的邮箱发送一封邮件，其中包含用户名以及随机生成的密码。
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <Form {...form}>
                    <form
                        onSubmit={form.handleSubmit((data) =>
                            mutation.mutate(data)
                        )}
                        className="space-y-2"
                    >
                        <FormField
                            control={form.control}
                            name="username"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>用户名(NetID)</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="请输入用户名"
                                            {...field}
                                            disabled={mutation.isPending}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>邮箱</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="请输入邮箱"
                                            {...field}
                                            disabled={mutation.isPending}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="fullName"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>姓名</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="请输入姓名"
                                            {...field}
                                            disabled={mutation.isPending}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="role"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>身份</FormLabel>
                                    <Select
                                        onValueChange={field.onChange}
                                        value={field.value}
                                        disabled={mutation.isPending}
                                    >
                                        <FormControl>
                                            <SelectTrigger>
                                                <SelectValue />
                                            </SelectTrigger>
                                        </FormControl>
                                        <SelectContent>
                                            <SelectItem value="普通助理">
                                                普通助理
                                            </SelectItem>
                                            <SelectItem value="资深助理">
                                                资深助理
                                            </SelectItem>
                                            <SelectItem value="黑心">
                                                黑心
                                            </SelectItem>
                                        </SelectContent>
                                    </Select>
                                </FormItem>
                            )}
                        />
                        <div className="flex justify-end items-center gap-2 pt-2">
                            <Button
                                variant="outline"
                                onClick={() => {
                                    onOpenChange(false);
                                    form.reset();
                                }}
                                type="button"
                                disabled={mutation.isPending}
                            >
                                取消
                            </Button>
                            {mutation.isPending ? (
                                <PendingButton />
                            ) : (
                                <Button type="submit">确定</Button>
                            )}
                        </div>
                    </form>
                </Form>
            </AlertDialogContent>
        </AlertDialog>
    );
}
