import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { SidebarInset, SidebarTrigger } from "@/components/ui/sidebar";
import UsersTable from "@/components/UsersTable";

export default function UsersManagementPage() {
  return (
    <SidebarInset className="flex flex-col h-full">
      <header className="flex h-12 items-center gap-2">
        <SidebarTrigger className="ml-2" />
        <Separator orientation="vertical" className="h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbPage>用户管理</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <main className="px-2">
        <div className="container mx-auto">
          <UsersTable />
        </div>
      </main>
    </SidebarInset>
  );
}
