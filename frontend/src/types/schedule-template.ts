export type ScheduleTemplateMetaType = {
  id: number;
  name: string;
  description: string;
  createdAt: string;
};

export type ScheduleTemplateShiftType = {
  id: number;
  startTime: string;
  endTime: string;
  requiredAssistants: number;
  applicableDays: number[];
};
