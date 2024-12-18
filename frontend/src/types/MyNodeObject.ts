export interface MyNodeObject {
    id: string;
    neighbors?: number[];
    links?: { source: number; target: number }[];
    x?: number;
    y?: number;
    z?: number;
}
