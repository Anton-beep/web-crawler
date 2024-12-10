import {useEffect, useRef, useState} from 'react';
import ForceGraph3D, {ForceGraph3DInstance} from '3d-force-graph';
//import {NodeObject} from "react-force-graph-3d";
import {GraphData} from '@/types/GraphData.ts';
import OpenSiteFromGraphCard from "@/components/OpenSiteFromGraphCard.tsx";
import {NodeObject} from 'three-forcegraph';


function getBordersForNodeOn2D(node: NodeObject, graph: ForceGraph3DInstance): {
    minX: number,
    minY: number,
    maxX: number,
    maxY: number
} {

    if (node.x === undefined || node.y === undefined || node.z === undefined) {
        console.error("Node has no x, y or z", node);
        return {minX: 0, minY: 0, maxX: 0, maxY: 0};
    }

    const coords = [
        graph.graph2ScreenCoords(node.x + 5, node.y, node.z),
        graph.graph2ScreenCoords(node.x - 5, node.y, node.z),
        graph.graph2ScreenCoords(node.x, node.y - 5, node.z),
        graph.graph2ScreenCoords(node.x, node.y + 5, node.z),
        graph.graph2ScreenCoords(node.x, node.y, node.z + 5),
        graph.graph2ScreenCoords(node.x, node.y, node.z - 5),
    ]

    let minX = coords[0].x;
    let maxX = coords[0].x;
    let minY = coords[0].y;
    let maxY = coords[0].y;

    for (let i = 1; i < coords.length; i++) {
        if (coords[i].x < minX) {
            minX = coords[i].x;
        }
        if (coords[i].x > maxX) {
            maxX = coords[i].x;
        }
        if (coords[i].y < minY) {
            minY = coords[i].y;
        }
        if (coords[i].y > maxY) {
            maxY = coords[i].y;
        }
    }


    return {minX, minY, maxX, maxY};
}

export default function SitesGraph({width, height, backgroundCol, data}: {
    width: number,
    height: number,
    backgroundCol: string,
    data: GraphData
}) {
    const graphRef = useRef<HTMLDivElement | null>(null);
    const [linkToOpen, setLinkToOpen] = useState("");

    useEffect(() => {
        if (graphRef.current) {

            const neighbors = new Map<string, Set<string>>();

            data.links.forEach(link => {
                if (!neighbors.has(link.source)) {
                    neighbors.set(link.source, new Set());
                }
                if (!neighbors.has(link.target)) {
                    neighbors.set(link.target, new Set());
                }

                neighbors.get(link.source)!.add(link.target);
                neighbors.get(link.target)!.add(link.source);
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

            Graph.onNodeClick((node: NodeObject, event: any) => {
                const borders = getBordersForNodeOn2D(node, Graph);
                if (event.layerX < borders.minX || event.layerX > borders.maxX || event.layerY < borders.minY || event.layerY > borders.maxY) {
                    return;
                }

                if (node.id === undefined) {
                    console.error("Node has no id", node);
                    return;
                }

                setLinkToOpen(node.id.toString());
            });

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

            Graph.nodeColor(node => {
                if (typeof node.id !== "string") {
                    console.error("node.id is not string ", node)
                    return "rgba(0,255,255,0.6)";
                }

                return highlightNodesID.has(node.id) ? node.id === hoverNode ? 'rgb(255,0,0,1)' : 'rgba(255,160,0,0.8)' : 'rgba(0,255,255,0.6)'
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
