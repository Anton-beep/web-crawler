import {useParams} from 'react-router-dom';
import SitesGraph from "@/components/SitesGraph.tsx";
import {useCallback, useEffect, useState} from "react";
import Api from "@/services/Api.ts";
import ReactMarkdown from 'react-markdown';
import {YourDataIsLoading} from "@/components/YourDataIsLoading.tsx";

export default function Project() {
    const {projectId} = useParams();
    const [graphData, setGraphData] = useState<{ nodes: [], links: [] } | null>(null);
    const [mainIdeas, setMainIdeas] = useState("");
    const [keyWords, setKeyWords] = useState("");
    const [projectName, setProjectName] = useState("");
    const [loadingGraph, setLoadingGraphGraphData] = useState(true);
    const [error, setError] = useState("");
    const [dimensions, setDimensions] = useState({width: window.innerWidth, height: window.innerHeight});
    const [loadingKeyWords, setLoadingKeyWords] = useState(true);
    const [loadingMainIdeas, setLoadingMainIdeas] = useState(true);

    useEffect(() => {
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
    }, [loadingGraph, error]);

    const fetchData = useCallback(() => {
        Api.getProject(projectId as string).then((response) => {
            console.log(response.data);

            if (response.data.name !== "") {
                setProjectName(response.data.name);
            }

            if (graphData === null && response.data.web_graph !== "") {
                try {
                    setGraphData(JSON.parse(response.data.web_graph));
                    console.log("data", JSON.parse(response.data.web_graph))
                    setLoadingGraphGraphData(false);
                } catch (e) {
                    setError("Error parsing project data");
                    setLoadingGraphGraphData(false);
                    console.error(e);
                }
            }

            if (response.data.key_words !== "") {
                setLoadingKeyWords(false);
                setKeyWords(response.data.key_words);
            }

            if (response.data.main_ideas !== "") {
                setLoadingMainIdeas(false);
                setMainIdeas(response.data.main_ideas);
            }
        }).catch((e) => {
            setError("Project not found");
            setLoadingGraphGraphData(false);
            console.error(e);
        });
    }, [projectId, graphData]);

    const isAnyAnalysisLoading = useCallback(() => {
        return loadingKeyWords || loadingMainIdeas;
    }, [loadingKeyWords, loadingMainIdeas]);

    useEffect(() => {
        fetchData();

        const interval = setInterval(() => {
            if (isAnyAnalysisLoading()) {
                fetchData();
            }
        }, 3000); // 1000ms = 1 second

        return () => clearInterval(interval);
    }, [fetchData, loadingKeyWords, loadingMainIdeas, isAnyAnalysisLoading]);

    const getGraphContent = () => {
        if (loadingGraph) {
            return (
                <div className="text-2xl font-extrabold">
                    <YourDataIsLoading/>
                </div>
            )
                ;
        } else if (error) {
            return <p className="text-error">{error}</p>;
        } else {
            if (graphData === null) {
                return <p className="text-error">No data</p>;
            }
            return (
                // don't ask me about dimensions.width * 0.1
                <SitesGraph width={dimensions.width - dimensions.width * 0.2} height={dimensions.height}
                            backgroundCol={"#1a1a1c"} data={graphData}/>
            );
        }
    }

    return (
        <div className="items-center justify-center py-20 h-full md:h-auto relative w-full">
            <div className="text-primary font-bold text-center text-xl mb-16">
                Name of the project: <span className="text-accent">{projectName}</span>
            </div>
            <div className="flex justify-center mb-32">
                {getGraphContent()}
            </div>
            <div>
                <div className="border-2 border-gray-700 p-4 rounded-lg mb-16">
                    <div className="text-primary font-bold italic text-center text-xl mb-8">
                        Key words (Be careful, AI generated)
                    </div>
                    <div className="text-primary text-center">
                        {loadingKeyWords ? <div className="font-bold text-xl">
                            <YourDataIsLoading/></div> : <ReactMarkdown>{keyWords}</ReactMarkdown>}
                    </div>
                </div>
                <div className="border-2 border-gray-700 p-4 rounded-lg mb-16">
                    <div className="text-primary font-bold italic text-center text-xl mb-8">
                        Main idea (Be careful, AI generated)
                    </div>
                    <div className="text-primary text-center">
                        {loadingMainIdeas ? <div className="font-bold text-xl"><YourDataIsLoading/></div> :
                            <ReactMarkdown>{mainIdeas}</ReactMarkdown>}
                    </div>
                </div>
            </div>
        </div>
    )
}