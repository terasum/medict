import { Definition } from '../model/Definition';

export type ResourceFn = (key: string, withPayload?: boolean) => Definition;
export type LookupFn = (key: string) => Definition;

export interface Replacer {
  replace(
    dictid: string,
    keyText: string,
    html: string,
    lookfn: LookupFn,
    resourceFn: ResourceFn
  ): {keyText: string, definition: string};
}
