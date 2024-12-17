import {Button} from "@/components/ui/button"
import {Dialog, DialogTrigger,} from "@/components/ui/dialog"
import CreateProjectCard from "@/components/CreateProjectCard.tsx";
import OpenProjectCard from "@/components/OpenProjectCard.tsx";
import {useEffect, useState} from "react";
import Api from "@/services/Api.ts";
import {ShortProject} from "@/types/ShortProject.ts";
import Cookies from "js-cookie";
import {HelpCard} from "@/components/HelpCard.tsx";

export default function Dashboard() {
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [projects, setProjects] = useState<ShortProject[]>([]);
    const [showHelp, setShowHelp] = useState(true);

    useEffect(() => {
        setShowHelp(Cookies.get("showHelp") !== "false");
        fetchProjects();
    }, [showHelp]);

    const fetchProjects = () => {
        Api.getAllProjectsShort().then((response) => {
            setProjects(response.data);
        }).catch((e) => {
            console.error(e);
        });
    };

    const getHelp = () => {
        if (!showHelp) {
            return null;
        }
        return (
            <HelpCard setShowHelp={setShowHelp}/>
        )
    }

    return (
        <div className="text-primary mx-20">
            {getHelp()}
            <div className="mb-20 text-lg font-extrabold">
                You have created <span className="text-accent font-normal"> {projects.length} </span> projects!
            </div>
            <div>
                <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                    <DialogTrigger>
                        <Button variant="jumping" className="bg-green-600 text-primary"
                                onClick={() => setIsCreateDialogOpen(true)}>Create Project</Button>
                    </DialogTrigger>
                    <CreateProjectCard/>
                </Dialog>

                <Dialog onOpenChange={(isOpen) => {
                    if (isOpen) {
                        fetchProjects();
                    }
                }}>
                    <DialogTrigger>
                        <Button variant="jumping" className="bg-blue-600 text-primary mx-4 h-10">Open
                            Project</Button>
                    </DialogTrigger>
                    <OpenProjectCard projects={projects}/>
                </Dialog>
            </div>
        </div>
    )
}
