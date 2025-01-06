import { ScheduleTemplateType } from "@/types/schedule-template";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import { DayOfWeek } from "@/lib/const";
import { Clock, UserRound, UsersRound } from "lucide-react";

interface Props {
  scheduleTemplate: ScheduleTemplateType;
}

export default function ScheduleTemplateDetails({ scheduleTemplate }: Props) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{scheduleTemplate.name}</CardTitle>
        <CardDescription>{scheduleTemplate.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue={DayOfWeek[0].key}>
          <TabsList className={cn("grid", `grid-cols-${DayOfWeek.length}`)}>
            {DayOfWeek.map((day) => (
              <TabsTrigger key={day.key} value={day.key}>
                {day.name}
              </TabsTrigger>
            ))}
          </TabsList>
          {DayOfWeek.map((day) => (
            <TabsContent key={day.key} value={day.key} className="space-y-2">
              {scheduleTemplate.shifts.map(
                (shift) =>
                  shift.applicableDays.includes(day.id) && (
                    <Card>
                      <CardContent className="p-4">
                        <div className="flex justify-between">
                          <div className="flex justify-center items-center space-x-2">
                            <Clock size={16} />
                            <span>
                              {shift.startTime}~{shift.endTime}
                            </span>
                          </div>
                          <div className="flex justify-center items-center space-x-2">
                            <UsersRound size={16} />
                            <span>{shift.requiredAssistants} 名助理</span>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  )
              )}
            </TabsContent>
          ))}
        </Tabs>
      </CardContent>
    </Card>
  );
}
