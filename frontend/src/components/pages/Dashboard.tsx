import {Button} from "@/components/ui/button"
import {Dialog, DialogTrigger,} from "@/components/ui/dialog"
import CreateProjectCard from "@/components/CreateProjectCard.tsx";
import {useEffect, useState} from "react";
import Api from "@/services/Api.ts";
import {ShortProject} from "@/types/ShortProject.ts";
import Cookies from "js-cookie";
import {HelpCard} from "@/components/HelpCard.tsx";
import SmallProjectCard from "@/components/SmallProjectCard.tsx";
import {BackgroundGradient} from "@/components/ui/background-gradient.tsx";

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
            <div className="mb-12">
            <BackgroundGradient className="rounded-[22px] p-4 sm:p-10 bg-background">
                <HelpCard setShowHelp={setShowHelp}/>
            </BackgroundGradient>
            </div>
        )
    }

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
                <SmallProjectCard projectName={project.name} projectId={project.id}/>
            )
        })
    }

    return (
        <div className="text-primary mx-20">
            {getHelp()}
            <div className="mb-12 text-xl font-extrabold">
                You have created <span className="text-accent font-normal"> {projects.length} </span> projects!
            </div>
            <div>
                <div className="mb-12">
                    <Dialog open={isCreateDialogOpen} onOpenChange={(arg) => {
                        fetchProjects();
                        setIsCreateDialogOpen(arg);
                    }}>
                        <DialogTrigger>
                            <Button variant="jumping" className="bg-green-600 text-primary"
                                    onClick={() => setIsCreateDialogOpen(true)}>Create Project</Button>
                        </DialogTrigger>
                        <CreateProjectCard/>
                    </Dialog>
                </div>

                <BackgroundGradient className="rounded-[22px] p-4 sm:p-10 bg-background">
                    <span className="text-xl font-bold">Your projects:</span>
                    <div className="overflow-auto content-center items-center text-center max-h-96">
                        {getProjectsButtons()}
                    </div>
                </BackgroundGradient>
            </div>
        </div>
    )
}
