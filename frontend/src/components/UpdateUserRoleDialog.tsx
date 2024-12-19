import { UserType } from "@/types/user";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogHeader,
    AlertDialogTitle,
} from "./ui/alert-dialog";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel } from "./ui/form";
import { Select, SelectContent, SelectItem, SelectTrigger } from "./ui/select";
import { SelectValue } from "@radix-ui/react-select";
import { Button } from "./ui/button";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, APIResponse } from "@/lib/api";
import { AxiosResponse } from "axios";
import { toast } from "sonner";
import PendingButton from "./PendingButton";

interface Props {
    user: UserType | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const formSchema = z.object({
    role: z.string({
        required_error: "请提供用户的新身份",
    }),
});

export default function UpdateUserRoleDialog({
    user,
    open,
    onOpenChange: onOpenChange,
}: Props) {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            role: "普通助理",
        },
    });
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: (data: z.infer<typeof formSchema>) =>
            api.post(`/users/${user?.id}/update-role`, data),
        onSuccess: (res: AxiosResponse<APIResponse<UserType>>) => {
            const updatedUser = res.data.data;

            toast.success(res.data.message);
            queryClient.setQueryData(["users"], (data: UserType[]) => {
                return data.map((user) => {
                    return user.id === updatedUser.id ? updatedUser : user;
                });
            });
            onOpenChange(false);
            form.reset();
        },
        onError: (err) => {
            toast.error(err.message);
        },
    });

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>修改身份</AlertDialogTitle>
                    <AlertDialogDescription>
                        请选择{user?.fullName}({user?.username})的新身份。
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <Form {...form}>
                    <form
                        onSubmit={form.handleSubmit((data) =>
                            mutation.mutate(data)
                        )}
                    >
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
                        <div className="flex items-center justify-end gap-2 mt-4">
                            <Button
                                variant="outline"
                                type="button"
                                onClick={() => {
                                    onOpenChange(false);
                                    form.reset();
                                }}
                                disabled={mutation.isPending}
                            >
                                取消
                            </Button>
                            {mutation.isPending ? (
                                <PendingButton />
                            ) : (
                                <Button>确认</Button>
                            )}
                        </div>
                    </form>
                </Form>
            </AlertDialogContent>
        </AlertDialog>
    );
}
