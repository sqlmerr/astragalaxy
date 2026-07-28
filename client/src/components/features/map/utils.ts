export function getOrbitRadius(orbit: number) {
  return (orbit + 1) * 80 + (30 * (orbit + 1)) ** 1.25
}
