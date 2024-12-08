import {useParams} from 'react-router-dom';
import SitesGraph from "@/components/SitesGraph.tsx";
import {useEffect, useState} from "react";
import Api from "@/services/Api.ts";

export default function Project() {
    const {projectId} = useParams();
    const [data, setData] = useState({nodes: [], links: []});
    const [isInProcess, setIsInProcess] = useState(false);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [dimensions, setDimensions] = useState({width: window.innerWidth, height: window.innerHeight});

    useEffect(() => {
        Api.getProject(projectId as string).then((response) => {
            if (response.data.processing) {
                setLoading(false);
                setIsInProcess(true);
                return;
            }

            try {
                setData(JSON.parse(response.data.web_graph));
            } catch (e) {
                setError("Error parsing project data");
                setLoading(false);
                console.error(e);
            }

            if (response.data.web_graph === "") {
                setError("No data for project");
                setLoading(false);
                return;
            }

            setLoading(false);
        }).catch((e) => {
            setError("Project not found");
            setLoading(false);
            console.error(e);
        });
    }, [projectId]);

    useEffect(() => {
        const handleResize = () => {
            setDimensions({
                width: window.innerWidth,
                height: window.innerHeight,
            });
        };
        window.addEventListener("resize", handleResize);

        return () => window.removeEventListener("resize", handleResize);
    }, [loading, isInProcess, error]);

    const getContent = () => {
        if (loading) {
            return (
                <h2 className="text-2xl font-bold">
                    Loading...</h2>
            );
        } else if (isInProcess) {
            return (
                <h2 className="text-2xl font-bold">
                    Your project will be available shortly...</h2>
            );
        } else if (error) {
            return <p className="text-error">{error}</p>;
        } else {
            return (
                // don't ask me about dimensions.width * 0.1
                <SitesGraph width={dimensions.width - dimensions.width * 0.05} height={dimensions.height} backgroundCol={"#18181b"} data={data}/>
            );
        }
    }

    return (
        <div className="text-primary">
            <h1>Project: {projectId}</h1>
            {getContent()}
        </div>
    )
}