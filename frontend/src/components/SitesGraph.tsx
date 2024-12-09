import {useEffect, useRef, useState} from 'react';
import ForceGraph3D, {ForceGraph3DInstance} from '3d-force-graph';
import {GraphData} from '@/types/GraphData.ts';
// @ts-expect-error : Cannot find module 'three/src/nodes/tsl/TSLCore'
import {NodeObject} from "three/src/nodes/tsl/TSLCore";
import OpenSiteFromGraphCard from "@/components/OpenSiteFromGraphCard.tsx";

function getBordersForNodeOn2D(node: NodeObject<number>, graph: ForceGraph3DInstance): {minX: number, minY : number, maxX: number, maxY: number} {
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
            const Graph = ForceGraph3D()(graphRef.current)
                .backgroundColor(backgroundCol)
                .graphData(data)
                .nodeLabel('id')
                .linkWidth(1)
                .nodeAutoColorBy('id')
                .width(width)
                .height(height)
                .linkDirectionalParticles(10)
                .linkDirectionalParticleSpeed(0.003)
                .nodeRelSize(5);

            Graph.onNodeClick((node: NodeObject<number>, event: any) => {
                const borders = getBordersForNodeOn2D(node, Graph);
                if (event.layerX < borders.minX || event.layerX > borders.maxX || event.layerY < borders.minY || event.layerY > borders.maxY) {
                    return;
                }
                setLinkToOpen(node.id);
            });

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
