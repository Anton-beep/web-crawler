import { useEffect, useRef } from 'react';
import ForceGraph3D from '3d-force-graph';
import SpriteText from 'three-spritetext';
import { GraphData } from '@/types/GraphData.ts';
import './SitesGraph.css'
import { UnrealBloomPass } from 'three/examples/jsm/postprocessing/UnrealBloomPass';
import getDomainFromUrl from "../../utils/getDomainFromUrl.ts";
import {NodeObject} from "three/src/nodes/tsl/TSLCore";


function ForceGraph({width ,height, backgroundCol, data} : {width: number, height: number, backgroundCol: string, data: GraphData}) {
    const graphRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (graphRef.current) {
            const Graph = ForceGraph3D()(graphRef.current)
                .backgroundColor(backgroundCol)
                .graphData(data)
                .nodeLabel((node: NodeObject<number>) => getDomainFromUrl(node.id))
                .linkWidth(1)
                .nodeAutoColorBy('id')
                .width(width)
                .height(height)
                // .linkDirectionalArrowLength(9)
                .linkDirectionalParticles(10)
                .linkDirectionalParticleSpeed(0.003);

            Graph.nodeThreeObject((node: NodeObject<number>) => {
                const sprite = new SpriteText(node.id) as any;
                sprite.material.depthWrite = false;
                sprite.color = node.color;
                sprite.textHeight = 8;
                return sprite;
            });

            const bloomPass = new UnrealBloomPass();
            bloomPass.strength = 1.4;
            bloomPass.radius = 0;
            bloomPass.threshold = 0;
            Graph.postProcessingComposer().addPass(bloomPass);

            return () => {
                Graph._destructor(); // Clean up on unmount
            };
        }
    }, [data, height, width, backgroundCol]);

    return <div ref={graphRef} />;
};

export default ForceGraph;
