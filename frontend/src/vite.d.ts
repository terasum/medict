declare module '*.md' {
    import type { ComponentOptions } from 'vue'
    const Component: ComponentOptions<any>
    export default Component
}