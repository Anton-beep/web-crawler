import {DialogContent, DialogDescription, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import { useNavigate } from 'react-router-dom';
import {ShortProject} from "@/types/ShortProject.ts";

export default function OpenProjectCard({projects}: {projects: ShortProject[]}) {
    const navigate = useNavigate();

    return (
        <DialogContent className="text-primary">
            <DialogHeader>
                <DialogTitle>Open Project</DialogTitle>
                <DialogDescription>
                    Click on a project which you want to open
                </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
                {projects.map((project) => {
                    return (
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Button className="col-span-4 bg-blue-600 text-primary" onClick={() => {
                                navigate(`/project/${project.id}`);
                            }}>
                                {project.name}
                            </Button>
                        </div>
                    )
                })}
            </div>
        </DialogContent>
    )
}