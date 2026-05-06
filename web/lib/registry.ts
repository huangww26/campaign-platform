import { ComponentType } from 'react'

export type ActivityComponent = ComponentType<{
  props: Record<string, unknown>
  ctx: { lang: string }
}>

const registry = new Map<string, ActivityComponent>()

export function register(name: string, comp: ActivityComponent) {
  registry.set(name, comp)
}

export function get(name: string): ActivityComponent | undefined {
  return registry.get(name)
}

export function has(name: string): boolean {
  return registry.has(name)
}
