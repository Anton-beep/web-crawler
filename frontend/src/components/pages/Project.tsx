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
            console.error(e);
        });
    }, [projectId]);

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
            return <SitesGraph width={1600} height={1100} backgroundCol={"#18181b"} data={data}/>;
        }
    }

    return (
        <div className="text-primary">
            <h1>Project: {projectId}</h1>
            {getContent()}
        </div>
    )
}