import {useParams} from 'react-router-dom';
import SitesGraph from "@/components/SitesGraph.tsx";
import {useCallback, useEffect, useState} from "react";
import Api from "@/services/Api.ts";
import ReactMarkdown from 'react-markdown';
import {YourDataIsLoading} from "@/components/YourDataIsLoading.tsx";
import {BackgroundGradient} from "@/components/ui/background-gradient.tsx";
import {Button} from "@/components/ui/button.tsx";
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from "@/components/ui/dialog.tsx";
import { useNavigate } from 'react-router-dom';

export default function Project() {
    const {projectId} = useParams();
    const [graphData, setGraphData] = useState<{ nodes: [], links: [] } | null>(null);
    const [deadSites, setDeadSites] = useState<string[]>([]);
    const [mainIdeas, setMainIdeas] = useState("");
    const [keyWords, setKeyWords] = useState("");
    const [projectName, setProjectName] = useState("");
    const [projectStartUrl, setProjectStartUrl] = useState("");
    const [loadingGraph, setLoadingGraphGraphData] = useState(true);
    const [error, setError] = useState("");
    const [dimensions, setDimensions] = useState({width: window.innerWidth, height: window.innerHeight});
    const [loadingKeyWords, setLoadingKeyWords] = useState(true);
    const [loadingMainIdeas, setLoadingMainIdeas] = useState(true);
    const navigate = useNavigate();

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
            if (response.data.name !== "") {
                setProjectName(response.data.name);
            }

            if (response.data.start_url !== "") {
                setProjectStartUrl(response.data.start_url);
            }

            if (graphData === null && !response.data.processing) {
                try {
                    setGraphData(JSON.parse(response.data.web_graph));
                    setLoadingGraphGraphData(false);
                } catch (e) {
                    setError("Error parsing project data");
                    setLoadingGraphGraphData(false);
                    console.error(e);
                }

                setDeadSites(response.data.dlq_sites);
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
                            backgroundCol={"#18181b"} data={graphData}/>
            );
        }
    }

    const getDeadLinksContent = () => {
        if (deadSites === null || deadSites.length === 0) {
            return
        } else {
            return (
                <div className="mb-16 text-error text-center">
                    <span className="text-lg font-bold">Dead links: </span>
                    {deadSites.join(", ")}
                </div>
            );
        }
    }

    const downloadGraphData = () => {
        if (graphData) {
            const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(graphData));
            const downloadAnchorNode = document.createElement('a');
            downloadAnchorNode.setAttribute("href", dataStr);
            downloadAnchorNode.setAttribute("download", `web_graph_${projectId}.json`);
            document.body.appendChild(downloadAnchorNode);
            downloadAnchorNode.click();
            downloadAnchorNode.remove();
        }
    }

    const deleteProject = () => {
        Api.deleteProject(projectId as string).then(() => {
            navigate("/dashboard");
        }).catch((e) => {
            console.error(e);
        });
    }

    return (
        <div className="items-center justify-center py-20 h-full md:h-auto relative w-full">
            <div className="flex justify-center mb-10">
                <div className="text-primary font-bold text-center text-xl mx-10">
                    Name of the project: <span className="text-accent">{projectName}</span>
                </div>
                <div className="text-primary font-bold text-center text-xl mx-10">
                    Start link: <span className="text-accent">{projectStartUrl}</span>
                </div>
            </div>
            <div>
                {getDeadLinksContent()}
            </div>
            <div className="flex justify-center mb-16">
                <BackgroundGradient className="rounded-[22px] p-4 sm:p-10 bg-background">
                    {getGraphContent()}
                </BackgroundGradient>
            </div>
            <div className="flex justify-center mb-10">
                <Button onClick={downloadGraphData} variant="default">Download Graph Data</Button>
            </div>
            <div className="flex justify-center mb-10">
                <div className="mx-10 w-full h-full ">
                    <BackgroundGradient className="rounded-[22px] p-4 sm:p-10 bg-background">
                        <div className="text-primary font-bold italic text-center text-xl mb-8">
                            Key words <span className="text-warning">(Be careful, AI generated)</span>
                        </div>
                        <div className="text-primary text-left max-h-96 overflow-auto">
                            {loadingKeyWords ? <div className="font-bold text-xl">
                                <YourDataIsLoading/></div> : <ReactMarkdown>{keyWords}</ReactMarkdown>}
                        </div>
                    </BackgroundGradient>
                </div>
                <div className="mx-10 w-full h-full">
                    <BackgroundGradient className="rounded-[22px] p-4 sm:p-10 bg-background">
                        <div className="text-primary font-bold italic text-center text-xl mb-8">
                            Main idea <span className="text-warning">(Be careful, AI generated)</span>
                        </div>
                        <div className="text-primary text-left max-h-96 overflow-auto">
                            {loadingMainIdeas ? <div className="font-bold text-xl"><YourDataIsLoading/></div> :
                                <ReactMarkdown>{mainIdeas}</ReactMarkdown>}
                        </div>
                    </BackgroundGradient>
                </div>
            </div>
            <div>
                <Dialog>
                    <DialogTrigger>
                        <Button variant="default" className="bg-error text-primary">
                            Delete project
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle className="text-warning">Are you sure you want to delete this project?</DialogTitle>
                        </DialogHeader>
                        <DialogFooter>
                            <Button variant="default" className="bg-error text-primary" onClick={deleteProject}>Delete Permanently</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>
        </div>
    )
}