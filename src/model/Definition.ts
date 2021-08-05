export declare class Definition {
  keyText: string;
  definition: string;
  contentSize: number;
  payload?: string
}

export function NullDef(keyText: string) {
  return { keyText, definition: '',contentSize:0 } as Definition;
}
