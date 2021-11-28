import { DictService, EV_END_INDEXING, EV_START_INDEXING } from "../worksvc/Dictionary.svc";
import Configuration from "../worksvc/Configuration.svc";
import { SyncMainAPI } from '../../main/rpc.sync.main.reference';

export class DictAPI {
    dictService: DictService;
    config: Configuration;

    constructor() {
        this.config = Configuration.newInstance();

        this.dictService = new DictService(
            SyncMainAPI.syncGetResourceRootPath(),
            '/Users/chenquan/Workspace/nodejs/medict/testdict/testdict1',
            this.config.configJsonFilePath)

        this.dictService.autoIndexing();

        this.dictService.eventEmitter.on(EV_START_INDEXING, (ev, arg) => {
            console.log(ev, '==============', arg);
        })

        this.dictService.eventEmitter.on(EV_END_INDEXING, (ev, arg) => {
            console.log(ev, '==============', arg);
        })
    }

    loadAllIndexed() {
        let indexed = this.dictService.listDicts();
        return indexed.map((dict) => {
            return {
                id: dict.id,
                alias: dict.alias,
                name: dict.name,
            };
        })
    }

    loadAllUnIndexed() {
        return this.dictService.findAll();
    }

    suggestWord(dictid: string, word: string) {
        return this.dictService.associate(dictid, word);
    }

    lookupWordPrecisely(dictid: string, word: string, roffset: number) {
        console.log('[worker] lookupWordPrecisely', dictid, word, roffset);
        return this.dictService.lookupPrecisly(dictid, word, roffset)
    }

    postHandleDef(dictid: string, keyText: string, rawhtml: string) {
        return this.dictService.definitionReplace(dictid, keyText, rawhtml);
    }
}