import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogTrigger,
} from "@/components/ui/dialog"
import CreateProjectCard from "@/components/CreateProjectCard.tsx";
import OpenProjectCard from "@/components/OpenProjectCard.tsx";
import {useState} from "react";
import Api from "@/services/Api.ts";
import {ShortProject} from "@/types/ShortProject.ts";

export default function Dashboard() {
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [projects, setProjects] = useState<ShortProject[]>([]);

    const fetchProjects = () => {
        Api.getAllProjectsShort().then((response) => {
            setProjects(response.data);
        }).catch((e) => {
            console.error(e);
        });
    };

    return (
        <div className="mx-4">
            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                <DialogTrigger asChild>
                    <Button variant="jumping" className="bg-green-600 text-primary mx-4" onClick={() => setIsCreateDialogOpen(true)}>Create Project</Button>
                </DialogTrigger>
                <CreateProjectCard/>
            </Dialog>

            <Dialog onOpenChange={(isOpen) => {
                if (isOpen) {
                    fetchProjects();
                }
            }}>
                <DialogTrigger>
                    <Button variant="jumping" className="bg-blue-600 text-primary mx-4">Open Project</Button>
                </DialogTrigger>
                <OpenProjectCard projects={projects}/>
            </Dialog>

        </div>
    )
}
