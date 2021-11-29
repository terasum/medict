import { SyncMainAPI } from "../../main/rpc.sync.main.reference";
import Configuration from "./Configuration.svc"
import { LowSync, JSONFileSync } from 'lowdb';
import { Config } from '../../model/Config';

export class UserConfigService {
    userConfigFilePath: string;
    rawconfig: Configuration;
    db: LowSync<Config>;

    constructor() {
        this.userConfigFilePath = SyncMainAPI.syncGetConfigJsonPath();
        this.rawconfig = Configuration.newInstance();
        this.db = new LowSync(new JSONFileSync<Config>(this.userConfigFilePath));
        this.db.read();
        this.db.data ||= new Config();
    }

    getDictBaseDir(): string {
        return this.db.data!.dictBaseDir;
    }

    saveDictBaseDir(dir: string) {
        this.db.data!.dictBaseDir = dir;
        this.db.write();
    }

    getBaiduApiConfig(): { appid: string, appkey: string } {
        return this.db.data!.translateApis.baidu;
    }

    getYoudaoApiConfig(): { appid: string, appkey: string } {
        return this.db.data!.translateApis.youdao;
    }

    saveBaiduApiConfig(args: { appid: string; appkey: string }) {
        this.db.data!.translateApis.baidu = args;
        this.db.write();
    }
    saveYoudaoApiConfig(args: { appid: string; appkey: string }) {
        this.db.data!.translateApis.youdao = args;
        this.db.write();
    }
}
