import ScheduleTemplateMetaTable from "@/components/table/ScheduleTemplateMetaTable";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";

export default function ScheduleTemplatesManagementPage() {
  return (
    <>
      <Breadcrumb className="mb-4">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbPage className="font-bold">班表模板管理</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <ScheduleTemplateMetaTable />
    </>
  );
}
