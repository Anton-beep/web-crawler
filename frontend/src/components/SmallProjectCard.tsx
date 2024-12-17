import {Button} from "@/components/ui/button.tsx";
import {useNavigate} from "react-router-dom";

export default function SmallProjectCard({projectName, projectId}: { projectName: string, projectId: string }) {
    const navigate = useNavigate();

    return (
        <div key={projectId}>
            <Button variant="jumpingLink" className="col-span-4 overflow-hidden text-ellipsis whitespace-nowrap text-lg"
                    onClick={() => {
                        navigate(`/project/${projectId}`);
                    }}>
                <span>{projectName}</span>
            </Button>
        </div>
    )
}