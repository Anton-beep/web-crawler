import {DialogContent, DialogDescription, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import { useNavigate } from 'react-router-dom';
import {ShortProject} from "@/types/ShortProject.ts";

export default function OpenProjectCard({projects}: {projects: ShortProject[]}) {
    const navigate = useNavigate();

    const getProjectsButtons = () => {
        if (projects.length === 0) {
            return (
                <div className="grid grid-cols-4 items-center gap-4">
                    <div className="col-span-4 text-accent text-center font-bold">
                        No projects
                    </div>
                </div>
            )
        }

        return projects.map((project) => {
            return (
                <div className="grid grid-cols-4 items-center gap-4" key={project.id}>
                    <Button variant="default" className="col-span-4 bg-blue-600 overflow-hidden text-ellipsis whitespace-nowrap" onClick={() => {
                        navigate(`/project/${project.id}`);
                    }}>
                        {project.name}
                    </Button>
                </div>
            )
        })
    }

    return (
        <DialogContent className="text-primary overflow-auto min-h-[25vh] max-h-[50vh]">
            <DialogHeader>
                <DialogTitle>Open Project</DialogTitle>
                <DialogDescription>
                    Click on a project which you want to open
                </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
                {getProjectsButtons()}
            </div>
        </DialogContent>
    )
}