import { DictService, EV_END_INDEXING, EV_START_INDEXING } from "../worksvc/Dictionary.svc";
import Configuration from "../worksvc/Configuration.svc";
import { UserConfigService } from "../worksvc/UserConfig.svc";
import { StorabeDictionary } from "../../model/StorableDictionary";
import fs  from "fs";

export class DictAPI {
    dictService: DictService;
    userConfigService: UserConfigService;
    config: Configuration;

    constructor() {
        this.config = Configuration.newInstance();
        this.userConfigService = new UserConfigService();

        this.dictService = new DictService(this.config.getResourceRootPath(), this.userConfigService.getDictBaseDir())
       
        this.dictService.eventEmitter.on(EV_START_INDEXING, (arg) => {
            console.log('===> BUILDING EVENT START', arg);
        })

        this.dictService.eventEmitter.on(EV_END_INDEXING, (arg) => {
            console.log('<=== BUILDING EVENT END', arg);
        })

        this.dictService.autoIndexing();
    }

    async reload() {
        if(this.userConfigService.getDictBaseDir() === this.dictService.dictsBaseDir) {
            console.log('[WORKER] dictionary directory is same, ignoring reload...')
            return;
        }
        console.log('=================== RELOAD Dictionaries ===================')
        this.dictService = new DictService(this.config.getResourceRootPath(), this.userConfigService.getDictBaseDir())
        this.dictService.eventEmitter.on(EV_START_INDEXING, (arg) => {
            console.log('===> BUILDING EVENT START', arg);
        })

        this.dictService.eventEmitter.on(EV_END_INDEXING, (arg) => {
            console.log('<=== BUILDING EVENT END', arg);
        })

        this.dictService.autoIndexing();
        return true;
    }

    getBaseDir() {
        return this.userConfigService.getDictBaseDir();
    }

    setBaseDir(dir: string) {
        if(!fs.existsSync(dir)) {
            return;
        }
        this.userConfigService.saveDictBaseDir(dir);
    }

    getDictInfo(dictid: string) {
        const loadedDict = this.dictService.loadDict(dictid);
        if (loadedDict && loadedDict.mdxDict && (loadedDict.mdxDict as any).header){
            const dictHeader = (loadedDict.mdxDict as any).header;
            // 转换为合法的可持久化词典返回
            return new StorabeDictionary(dictid, 
                dictHeader.Title, 
                dictHeader.Title,  
                loadedDict.mdxpath, 
                loadedDict.mddpath, 
                dictHeader.Description, 
                true)
        } else {
            const unloadedDict = this.dictService.findOne(dictid);
            return unloadedDict;
        }
    }

    hasIndexed(dictid: string) {
        const loadedDict = this.dictService.loadDict(dictid)
        if (!loadedDict) {
            return false
        }
        if (loadedDict.indexed) {
            return true
        }
        return false;
    }

    loadAllIndexed() {
        let indexed = this.dictService.listDicts();
        return indexed.map((dict) => {
            return {
                id: dict.id,
                alias: dict.alias,
                name: (dict.mdxDict as any).header!.Title || dict.name,
            };
        })
    }

    loadAllUnIndexed() {
        return this.dictService.findAll();
    }

    suggestWord(dictid: string, word: string) {
        return this.dictService.associate(dictid, word);
    }

    lookupDefinition(dictid: string, word: string, roffset: number) {
        console.log('[worker] lookupWordPrecisely', dictid, word, roffset);
        return this.dictService.lookupPrecisly(dictid, word, roffset)
    }

    postHandle(dictid: string, keyText: string, rawhtml: string) {
        return this.dictService.definitionReplace(dictid, keyText, rawhtml);
    }
}