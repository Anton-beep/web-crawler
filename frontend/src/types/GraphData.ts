import { GraphLink } from './GraphLink.ts'
import { GraphNode } from './GraphNode'

export interface GraphData {
    nodes: GraphNode[]
    links: GraphLink[]
}