import { TimePicker } from "@/components/time-picker";
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
import { checkScheduleTemplateShiftConflict } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { isBefore } from "date-fns";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  scheduleTemplateShifts: ScheduleTemplateFormSchemaType[];
  setScheduleTemplateShifts: (sts: ScheduleTemplateFormSchemaType[]) => void;
}

export const ScheduleTemplateFormSchema = z
  .object({
    startTime: z.string(),
    endTime: z.string(),
    requiredAssistants: z
      .number({ message: "请输入所需要的助理人数" })
      .int()
      .positive("请输入有效的助理人数"),
    applicableDays: z.array(z.number()).default([]), // if you don't add default([]) here, zod will not work
  })
  .refine(
    (data) => {
      const startTimeDate = new Date(`1970-01-01T${data.startTime}`);
      const endTimeDate = new Date(`1970-01-01T${data.endTime}`);
      return isBefore(startTimeDate, endTimeDate);
    },
    {
      message: "结束时间必须要在开始时间之后",
      path: ["endTime"],
    }
  );

export type ScheduleTemplateFormSchemaType = z.infer<
  typeof ScheduleTemplateFormSchema
>;

export default function CreateScheduleTemplateDialog({
  open,
  onOpenChange,
  scheduleTemplateShifts,
  setScheduleTemplateShifts,
}: Props) {
  const form = useForm<ScheduleTemplateFormSchemaType>({
    resolver: zodResolver(ScheduleTemplateFormSchema),
    defaultValues: {
      startTime: "00:00:00",
      endTime: "00:00:00",
      requiredAssistants: 0,
      applicableDays: [],
    },
  });

  const onSubmit = (data: ScheduleTemplateFormSchemaType) => {
    if (checkScheduleTemplateShiftConflict(scheduleTemplateShifts, data)) {
      toast.error("新班次与已有班次的时间冲突");
      return;
    }

    setScheduleTemplateShifts([...scheduleTemplateShifts, data]);
    toast.success("创建班次成功");
    onOpenChange(false);
    form.reset();
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>添加班次</AlertDialogTitle>
          <AlertDialogDescription>
            请输入班次的信息以创建班次。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="startTime"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>班次开始时间</FormLabel>
                  <FormControl>
                    <TimePicker
                      setDate={(date) => {
                        const hour = date?.getHours();
                        const minute = date?.getMinutes();
                        const second = date?.getSeconds();

                        field.onChange(
                          `${hour?.toString().padStart(2, "0")}:${minute
                            ?.toString()
                            .padStart(2, "0")}:${second
                            ?.toString()
                            .padStart(2, "0")}`
                        );
                      }}
                      date={
                        field.value
                          ? new Date(`1970-01-01T${field.value}`)
                          : undefined
                      }
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="endTime"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>班次结束时间</FormLabel>
                  <FormControl>
                    <TimePicker
                      setDate={(date) => {
                        const hour = date?.getHours();
                        const minute = date?.getMinutes();
                        const second = date?.getSeconds();

                        field.onChange(
                          `${hour?.toString().padStart(2, "0")}:${minute
                            ?.toString()
                            .padStart(2, "0")}:${second
                            ?.toString()
                            .padStart(2, "0")}`
                        );
                      }}
                      date={
                        field.value
                          ? new Date(`1970-01-01T${field.value}`)
                          : undefined
                      }
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="requiredAssistants"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>所需要助理人数</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      {...field}
                      onChange={(e) => field.onChange(parseInt(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex justify-end gap-2">
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
              <Button type="submit">添加</Button>
            </div>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
