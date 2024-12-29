import { Loader2 } from "lucide-react";
import { Button } from "./ui/button";

export default function PendingButton() {
  return (
    <Button disabled>
      <Loader2 className="animate-spin" />
      请稍等
    </Button>
  );
}
