export interface go {
  "apis": {
    "DictsAPI": {
		Destroy():Promise<void>
		Dicts():Promise<Array<PlainDictItem>>
		LookupDefinition(arg1:string,arg2:string,arg3:number):Promise<string|Error>
		Suggest(arg1:string):Promise<Array<WrappedWordItem>|Error>
    },
    "StaticInfos": {
		StaticSrvUrl():Promise<string>
    },
  }

  "main": {
    "App": {
		Greet(arg1:string):Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
