export interface GraphData {
    nodes: {
        id: string,
        group?: string,
        neighbors?: string[],
    }[]
    links: {
        source: string,
        target: string
    }[]
}