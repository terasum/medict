class TranslateApis {
  youdao = { appid: "", appkey: "" }
  baidu = { appid: "", appkey: "" }
}

export class Config {
  dictBaseDir: string;
  translateApis: TranslateApis;

  constructor() {
    this.dictBaseDir = "";
    this.translateApis = new TranslateApis();
  }
}
