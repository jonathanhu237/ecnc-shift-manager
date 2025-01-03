export type ScheduleTemplateMetaType = {
  id: number;
  name: string;
  description: string;
  createdAt: string;
  version: number;
};

export type ScheduleTemplateShiftType = {
  id: number;
  startTime: string;
  endTime: string;
  requiredAssistants: number;
  applicableDays: number[];
};

export type ScheduleTemplateType = ScheduleTemplateMetaType & {
  shifts: ScheduleTemplateShiftType[];
};
