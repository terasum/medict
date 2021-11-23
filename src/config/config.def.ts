export declare class DictItem {
    id: string;
    name: string;
    mdxpath: string;
    mddpath: [string];
    alias: string;
    description: string;
    resourceBaseDir: string;
}

export class DictConfig {
    translateApis: {
        baidu: {
            appid: string;
            appkey: string;
        },
        youdao: {
            appkey: string,
            appid: string
        }
    };
    dicts: DictItem[];
    constructor() {
        this.dicts = [];
        this.translateApis = { baidu: { appid: '', appkey: '' }, youdao: { appid: '', appkey: '' }, };
    }
}
