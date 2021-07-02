export declare class Definition {
  keyText: string;
  definition: string;
}

export function NullDef(keyText: string) {
  return { keyText, definition: '' } as Definition;
}
