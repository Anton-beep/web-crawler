// const SOME_NICE_COLORS = ["#b91c1c", "#a16207", "#4d7c0f", "#047857", "#0e7490", "#1d4ed8", "#6d28d9", "#a21caf", "#be123c"]
const SOME_NICE_COLORS = ["#ef4444", "#eab308", "#22c55e", "#3b82f6", "#a855f7", "#f43f5e"]

export function getNumberedColor(number: number): string {
    return SOME_NICE_COLORS[number % SOME_NICE_COLORS.length];
}

export function getHoveredNodeColor(): string {
    return "#a21caf";
}

export function getSelectedNodeColor(): string {
    return "#c2410c";
}