export function formatDeltaLabel(delta: number): string {
  return `${delta > 0 ? '+' : ''}${delta}%`
}
