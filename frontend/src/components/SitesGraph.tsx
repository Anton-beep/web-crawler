import {useEffect, useRef, useState} from 'react';
import ForceGraph3D from '3d-force-graph';
import {GraphData} from '@/types/GraphData.ts';
import OpenSiteFromGraphCard from "@/components/OpenSiteFromGraphCard.tsx";
import {NodeObject} from 'three-forcegraph';
import {getHoveredNodeColor, getSelectedNodeColor, getNumberedColor} from "@/utils/nodeColors.ts";

export default function SitesGraph({width, height, backgroundCol, data}: {
    width: number,
    height: number,
    backgroundCol: string,
    data: GraphData
}) {
    const graphRef = useRef<HTMLDivElement | null>(null);
    const [linkToOpen, setLinkToOpen] = useState("");

    useEffect(() => {
        console.log("data", data);
        if (graphRef.current) {
            const neighbors = new Map<string, Set<string>>();

            data.links.forEach((link: any) => {
                let src;
                if (typeof link.source === "string") {
                    src = link.source;
                } else {
                    src = link.source?.id;
                    if (src === undefined) {
                        console.error("link.source.id is undefined", link);
                        return;
                    }
                }

                let target;
                if (typeof link.target === "string") {
                    target = link.target;
                } else {
                    target = link.target?.id;
                    if (target === undefined) {
                        console.error("link.target.id is undefined", link);
                        return;
                    }
                }

                if (!neighbors.has(src)) {
                    neighbors.set(src, new Set());
                }
                if (!neighbors.has(target)) {
                    neighbors.set(target, new Set());
                }

                neighbors.get(src)!.add(target);
                neighbors.get(target)!.add(src);
            });

            const highlightNodesID = new Set<string>();
            const highlightLinks = new Map<string, Set<string>>();
            let hoverNode: string | null = null;

            // Docs to graph: https://github.com/vasturiano/3d-force-graph
            const Graph = new ForceGraph3D(graphRef.current)
                .backgroundColor(backgroundCol)
                .graphData(data)
                .nodeLabel('id')
                .linkWidth(1)
                .nodeAutoColorBy('id')
                .width(width)
                .height(height)
                .linkDirectionalParticles(5)
                .linkDirectionalParticleWidth(4)
                .linkDirectionalParticleSpeed(0.003)
                .nodeRelSize(5);

            // make first node bigger
            Graph.nodeVal((node: NodeObject) => node === data.nodes[0] ? 1000 : 5);

            const handleNodeClick = (node: NodeObject) => {
                if (typeof node.id !== "string") {
                    console.error("node.id is not string", node);
                    return;
                }

                if (hoverNode === node.id) {
                    setLinkToOpen(node.id);
                }

                hoverNode = null;
                highlightNodesID.clear();
                highlightLinks.clear();
                updateHighlight();
            }

            Graph.onNodeClick(handleNodeClick);

            // highlight nodes and links on hover

            function updateHighlight() {
                // trigger update of highlighted objects in scene
                Graph
                    .nodeColor(Graph.nodeColor())
                    .linkWidth(Graph.linkWidth())
                    .linkDirectionalParticles(Graph.linkDirectionalParticles());
            }

            Graph.onNodeHover((node => {
                if (node === hoverNode) {
                    return;
                }

                if (node === null) {
                    if (hoverNode === null) {
                        return;
                    } else {
                        hoverNode = null;
                        highlightNodesID.clear();
                        highlightLinks.clear();
                        updateHighlight();
                        return;
                    }
                }

                highlightNodesID.clear();
                highlightLinks.clear();

                if (typeof node.id !== "string") {
                    console.error("Node has no id", node);
                    return;
                }

                highlightNodesID.add(node.id);

                const nodeNeighbors = neighbors.get(node.id);
                if (nodeNeighbors === undefined) {
                    console.error("Node has no neighbors", node);
                    return;
                }


                nodeNeighbors.forEach((neighbor: string) => {
                    highlightNodesID.add(neighbor);
                });

                highlightLinks.set(node.id, nodeNeighbors)

                hoverNode = node.id

                updateHighlight();
            }));

            Graph.onLinkHover((link => {
                highlightNodesID.clear();
                highlightLinks.clear();

                if (link) {
                    if (typeof link.source !== "object" || typeof link.target !== "object") {
                        console.error("link.source or link.target is not object ", link);
                        return;
                    }

                    if (typeof link.source.id !== "string" || typeof link.target.id !== "string") {
                        console.error("link.source.id or link.target.id is not string");
                        return;
                    }

                    const setWithNeighbor = new Set<string>();
                    setWithNeighbor.add(link.target.id);
                    highlightLinks.set(link.source.id, setWithNeighbor);

                    highlightNodesID.add(link.source.id);
                    highlightNodesID.add(link.target.id);
                }

                updateHighlight();
            }));

            Graph.nodeColor((node: any) => {
                if (typeof node.id !== "string") {
                    console.error("node.id is not string ", node)
                    return "rgba(0,255,255,0.6)";
                }

                let colorNum = node?.val;
                if (colorNum === undefined) {
                    console.error("node.val is undefined", node);
                    return "rgba(0,255,255,0.6)";
                }

                return highlightNodesID.has(node.id) ? node.id === hoverNode ? getHoveredNodeColor() : getSelectedNodeColor() : getNumberedColor(colorNum);
            })

            Graph.linkWidth(link => {
                if (link.source === undefined || link.target === undefined) {
                    console.error("link.source or link.target is undefined")
                    return 1;
                }

                if (typeof link.source !== "object" || typeof link.target !== "object") {
                    if (typeof link.source === "string" && typeof link.target === "string") {
                        return highlightLinks.get(link.source)?.has(link.target) || highlightLinks.get(link.target)?.has(link.source) ? 4 : 1
                    }

                    console.error("link.source or link.target is not an object")
                    return 1;
                }

                if (typeof link.source.id !== "string" || typeof link.target.id !== "string") {
                    console.error("link.source.id or link.target.id is not a string")
                    return 1;
                }

                return highlightLinks.get(link.source.id)?.has(link.target.id) || highlightLinks.get(link.target.id)?.has(link.source.id) ? 4 : 1
            })

            Graph.linkDirectionalParticles(link => {
                if (link.source === undefined || link.target === undefined) {
                    console.error("link.source or link.target is undefined")
                    return 1;
                }

                if (typeof link.source !== "object" || typeof link.target !== "object") {
                    if (typeof link.source === "string" && typeof link.target === "string") {
                        return highlightLinks.has(link.source + link.target) || highlightLinks.has(link.target + link.source) ? 4 : 0
                    }

                    console.error("link.source or link.target is not an object")
                    return 1;
                }

                if (typeof link.source.id !== "string" || typeof link.target.id !== "string") {
                    console.error("link.source.id or link.target.id is not a string")
                    return 1;
                }

                return highlightLinks?.get(link.source.id)?.has(link.target.id) || highlightLinks.get(link.target.id)?.has(link.source.id) ? 4 : 0
            })

            return () => {
                Graph._destructor(); // Clean up on unmount
            };
        }
    }, [data, height, width, backgroundCol]);

    return (
        <>
            {linkToOpen !== "" && (
                <div>
                    <OpenSiteFromGraphCard url={linkToOpen} setUrl={setLinkToOpen}/>
                </div>
            )}
            <div ref={graphRef}/>
        </>
    );
}
