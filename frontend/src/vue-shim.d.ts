declare module '*.vue' {
  import Vue from 'vue';
  export default Vue;

}

declare module 'vue-friendly-iframe'

declare module '*.md' {
  import type { ComponentOptions } from 'vue'
  const Component: ComponentOptions<any>
  export default Component
}