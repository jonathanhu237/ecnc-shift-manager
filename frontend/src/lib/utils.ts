import { ScheduleTemplateFormSchemaType } from "@/components/dialog/CreateShiftDialog";
import { clsx, type ClassValue } from "clsx";
import { isBefore, isEqual } from "date-fns";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function checkScheduleTemplateShiftConflict(
  scheduleTemplateShift: ScheduleTemplateFormSchemaType[],
  newShift: ScheduleTemplateFormSchemaType
) {
  for (const existShift of scheduleTemplateShift) {
    const existShiftStartTimeDate = new Date(
      `1970-01-01T${existShift.startTime}`
    );
    const existShiftEndTimeDate = new Date(`1970-01-01T${existShift.endTime}`);
    const newShiftStartTimeDate = new Date(`1970-01-01T${newShift.startTime}`);
    const newShiftEndTimeDate = new Date(`1970-01-01T${newShift.endTime}`);

    if (
      isBefore(existShiftEndTimeDate, newShiftStartTimeDate) ||
      isEqual(existShiftEndTimeDate, newShiftStartTimeDate)
    ) {
      continue;
    }
    if (
      isBefore(newShiftEndTimeDate, existShiftStartTimeDate) ||
      isEqual(newShiftEndTimeDate, existShiftStartTimeDate)
    ) {
      continue;
    }

    return true;
  }

  return false;
}
