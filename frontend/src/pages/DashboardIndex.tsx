import {
    Breadcrumb,
    BreadcrumbList,
    BreadcrumbItem,
    BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { SidebarInset, SidebarTrigger } from "@/components/ui/sidebar";

export default function DashboardIndex() {
    return (
        <SidebarInset className="flex flex-col h-full">
            <header className="flex h-12 items-center gap-2">
                <SidebarTrigger className="ml-2" />
                <Separator orientation="vertical" className="h-4" />
                <Breadcrumb>
                    <BreadcrumbList>
                        <BreadcrumbItem>
                            <BreadcrumbPage>主页</BreadcrumbPage>
                        </BreadcrumbItem>
                    </BreadcrumbList>
                </Breadcrumb>
            </header>
            <main className="px-2 flex-1 flex">
                <div className="flex items-center justify-center flex-1">
                    <span>
                        如果你看到这里一片空白，不用担心，目前主页不会展示任何内容。
                    </span>
                </div>
            </main>
        </SidebarInset>
    );
}
