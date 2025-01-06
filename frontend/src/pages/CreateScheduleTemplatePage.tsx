import CreateScheduleTemplateDialog, {
  ScheduleTemplateFormSchema,
} from "@/components/dialog/CreateShiftDialog";
import PendingButton from "@/components/PendingButton";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { api, APIResponse } from "@/lib/api";
import { DayOfWeek } from "@/lib/const";
import {
  ScheduleTemplateMetaType,
  ScheduleTemplateType,
} from "@/types/schedule-template";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CirclePlusIcon, Trash2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

const formSchema = z.object({
  name: z.string().min(1, "请输入模板名称"),
  description: z.string(),
  shifts: z.array(ScheduleTemplateFormSchema),
});

type FormSchemaType = z.infer<typeof formSchema>;

export default function CreateScheduleTemplatePage() {
  const [createShiftDialogOpen, setCreateShiftDialogOpen] = useState(false);
  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      description: "",
      shifts: [],
    },
  });

  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const mutation = useMutation<
    APIResponse<ScheduleTemplateType>,
    Error,
    FormSchemaType
  >({
    mutationFn: (data: FormSchemaType) =>
      api.post("/schedule-templates", data).then((res) => res.data),
    onSuccess: (res) => {
      // set query cache
      queryClient.setQueryData(["schedule-template", res.data.id], res.data);
      queryClient.setQueryData(
        ["schedule-template-meta"],
        (prev: ScheduleTemplateMetaType[]) => [
          ...prev,
          {
            id: res.data.id,
            name: res.data.name,
            description: res.data.description,
            createdAt: res.data.createdAt,
            version: res.data.version,
          },
        ]
      );
      // toast
      toast.success(res.message);
      // reset form
      form.reset();
      // navigate to the table
      navigate("/schedule-templates-management", { replace: true });
    },
    onError: (err) => {
      toast.error(err.message);
    },
  });

  const onSubmit = (data: FormSchemaType) => {
    mutation.mutate(data);
  };

  return (
    <>
      {/* header */}
      <Breadcrumb className="mb-4">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/schedule-templates-management">班表模板管理</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage className="font-bold">创建班表模板</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      {/* form */}
      <div className="flex justify-center items-center flex-1">
        <Card>
          <CardHeader>
            <CardTitle>创建班表模板</CardTitle>
            <CardDescription>
              请在下方填入新班表模板的信息，注意班表模板名字不能重复。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form className="space-y-4">
                {/* the name of the template */}
                <FormField
                  name="name"
                  control={form.control}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>模板名称</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="请输入模板名称"
                          {...field}
                          disabled={mutation.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                {/* the description of the template */}
                <FormField
                  name="description"
                  control={form.control}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>模板描述</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="（可选）请输入模板描述"
                          {...field}
                          disabled={mutation.isPending}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                {/* the shifts of the template */}
                <FormField
                  name="shifts"
                  control={form.control}
                  render={({ field }) => (
                    <>
                      <FormItem>
                        {/* label */}
                        <div className="flex justify-between items-center">
                          <FormLabel>班次</FormLabel>
                          <Button
                            variant="ghost"
                            type="button"
                            size="icon"
                            onClick={() => setCreateShiftDialogOpen(true)}
                            disabled={mutation.isPending}
                          >
                            <CirclePlusIcon />
                          </Button>
                        </div>
                        {/* control */}
                        <FormControl>
                          {/* tabs */}
                          <Tabs defaultValue="mon" className="mt-2">
                            {/* tabsList */}
                            <TabsList className="grid grid-cols-7">
                              {DayOfWeek.map((item) => (
                                <TabsTrigger key={item.key} value={item.key}>
                                  {item.name}
                                </TabsTrigger>
                              ))}
                            </TabsList>
                            {/* tabsContent */}
                            {DayOfWeek.map((item) => (
                              <TabsContent key={item.key} value={item.key}>
                                <div className="space-y-2">
                                  {field.value.map((shift, index) => (
                                    <Card key={index}>
                                      <CardContent className="flex justify-between items-center text-sm p-2 pl-4">
                                        <span>
                                          {shift.startTime}~{shift.endTime} (
                                          {shift.requiredAssistants}名助理)
                                        </span>
                                        <div className="flex items-center gap-2">
                                          <Switch
                                            disabled={mutation.isPending}
                                            checked={field.value.some(
                                              (shiftField) =>
                                                shift === shiftField &&
                                                shift.applicableDays.includes(
                                                  item.id
                                                )
                                            )}
                                            onCheckedChange={(checked) => {
                                              field.onChange(
                                                field.value.map((s) =>
                                                  s === shift
                                                    ? {
                                                        ...s,
                                                        applicableDays: checked
                                                          ? [
                                                              ...s.applicableDays,
                                                              item.id,
                                                            ]
                                                          : s.applicableDays.filter(
                                                              (day) =>
                                                                day !== item.id
                                                            ),
                                                      }
                                                    : s
                                                )
                                              );
                                            }}
                                          />
                                          <Button
                                            variant="ghost"
                                            type="button"
                                            size="icon"
                                            onClick={() => {
                                              field.onChange(
                                                field.value.filter(
                                                  (s) => s !== shift
                                                )
                                              );
                                            }}
                                            disabled={mutation.isPending}
                                          >
                                            <Trash2 />
                                          </Button>
                                        </div>
                                      </CardContent>
                                    </Card>
                                  ))}
                                </div>
                              </TabsContent>
                            ))}
                          </Tabs>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                      {/* dialog for creating shift */}
                      <CreateScheduleTemplateDialog
                        open={createShiftDialogOpen}
                        onOpenChange={setCreateShiftDialogOpen}
                        scheduleTemplateShifts={field.value}
                        setScheduleTemplateShifts={field.onChange}
                      />
                    </>
                  )}
                />
                {/* button */}
                <div className="flex justify-end">
                  {mutation.isPending ? (
                    <PendingButton />
                  ) : (
                    <Button type="button" onClick={form.handleSubmit(onSubmit)}>
                      提交
                    </Button>
                  )}
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>
    </>
  );
}
