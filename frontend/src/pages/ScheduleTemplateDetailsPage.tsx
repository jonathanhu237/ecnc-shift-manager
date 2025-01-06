import ScheduleTemplateDetailsCard from "@/components/ScheduleTemplateDetailsCard";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { api, APIResponse } from "@/lib/api";
import { ScheduleTemplateType } from "@/types/schedule-template";
import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router";
import { toast } from "sonner";

export default function ScheduleTemplateDetailsPage() {
  const { id: scheduleTemplateID } = useParams();
  const {
    data: res,
    isPending,
    isError,
    error,
  } = useQuery<APIResponse<ScheduleTemplateType>>({
    queryKey: ["schedule-template", scheduleTemplateID],
    queryFn: () =>
      api
        .get(`/schedule-templates/${scheduleTemplateID}`)
        .then((res) => res.data),
  });

  if (isPending) return null;
  if (isError) {
    toast.error(error.message);
    return null;
  }

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
            <BreadcrumbPage className="font-bold">
              {res.data.name}
            </BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      {/* schedule template details */}
      <div className="mx-auto flex items-center justify-center flex-1 max-w-xl">
        <ScheduleTemplateDetailsCard scheduleTemplate={res.data} />
      </div>
    </>
  );
}
