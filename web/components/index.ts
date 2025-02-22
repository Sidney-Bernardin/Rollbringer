import type { AlpineComponent } from 'alpinejs'


export const defineComponent = <P, T>(fn: (params: P) => AlpineComponent<T>) => fn
