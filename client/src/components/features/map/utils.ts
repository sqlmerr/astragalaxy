export function getOrbitRadius(orbit: number) {
  return (orbit + 1) * 80 + (30 * (orbit + 1)) ** 1.25
}

export interface Point {
  x: number
  y: number
}

export function regularPolygon(
  cx: number,
  cy: number,
  radius: number,
  sides: number,
  rotation = -Math.PI / 2
): Point[] {
  if (sides < 3) {
    throw new Error("A polygon must have at least 3 sides.")
  }

  const points: Point[] = []

  const step = (Math.PI * 2) / sides

  for (let i = 0; i < sides; i++) {
    const angle = rotation + i * step

    points.push({
      x: cx + radius * Math.cos(angle),
      y: cy + radius * Math.sin(angle),
    })
  }

  return points
}
