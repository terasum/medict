export declare class Definition {
  keyText: string;
  definition: string;
  contentSize: number;
}

export function NullDef(keyText: string) {
  return { keyText, definition: '',contentSize:0 } as Definition;
}
