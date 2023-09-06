declare module '*.md' {
    import type { ComponentOptions } from 'vue'
    const Component: ComponentOptions<any>
    export default Component
}

declare module 'vue' {
    import { CompatVue } from '@vue/runtime-dom'
    const Vue: CompatVue
    export default Vue
    export * from '@vue/runtime-dom'
    const { configureCompat } = Vue
    export { configureCompat }
  }
  