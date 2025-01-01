import * as React from "react";
import { TimePickerInput } from "./ui/time-picker-input";

interface TimePickerDemoProps {
  date: Date | undefined;
  setDate: (date: Date | undefined) => void;
}

export function TimePicker({ date, setDate }: TimePickerDemoProps) {
  const minuteRef = React.useRef<HTMLInputElement>(null);
  const hourRef = React.useRef<HTMLInputElement>(null);
  const secondRef = React.useRef<HTMLInputElement>(null);

  return (
    <div className="flex items-center gap-2">
      <TimePickerInput
        picker="hours"
        date={date}
        setDate={setDate}
        ref={hourRef}
        onRightFocus={() => minuteRef.current?.focus()}
      />
      <span>时</span>
      <TimePickerInput
        picker="minutes"
        date={date}
        setDate={setDate}
        ref={minuteRef}
        onLeftFocus={() => hourRef.current?.focus()}
        onRightFocus={() => secondRef.current?.focus()}
      />
      <span>分</span>
      <TimePickerInput
        picker="seconds"
        date={date}
        setDate={setDate}
        ref={secondRef}
        onLeftFocus={() => minuteRef.current?.focus()}
      />
      <span>秒</span>
    </div>
  );
}
