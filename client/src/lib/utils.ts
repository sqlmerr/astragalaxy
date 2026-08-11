import { clsx  } from "clsx"
import type {ClassValue} from "clsx";
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function randomNumberBetween(a: number, b: number): number {
  const min = Math.min(a, b)
  const max = Math.max(a, b)

  return Math.random() * (max - min) + min
}

export function degreesToRadians(degrees: number): number {
  return (degrees * Math.PI) / 180
}
