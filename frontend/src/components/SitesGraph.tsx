import { useEffect, useRef } from 'react';
import ForceGraph3D from '3d-force-graph';
import { GraphData } from '@/types/GraphData.ts';
// @ts-expect-error : Cannot find module 'three/src/nodes/tsl/TSLCore'
import {NodeObject} from "three/src/nodes/tsl/TSLCore";


export default function SitesGraph({width ,height, backgroundCol, data} : {width: number, height: number, backgroundCol: string, data: GraphData}) {
    const graphRef = useRef<HTMLDivElement | null>(null);

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
                // .linkDirectionalArrowLength(9)
                .linkDirectionalParticles(10)
                .linkDirectionalParticleSpeed(0.003);

            Graph.onNodeClick((node: NodeObject<number>) => {
                const url = node.id;
                if (url) {
                    window.open(url, '_blank');
                }
            });

            return () => {
                Graph._destructor(); // Clean up on unmount
            };
        }
    }, [data, height, width, backgroundCol]);

    return <div ref={graphRef} />;
}
