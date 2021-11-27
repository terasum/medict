import { DictService,EV_END_INDEXING, EV_START_INDEXING } from "../worksvc/Dictionary.svc";
import Configuration from "../worksvc/Configuration.svc";

export default class DictAPI {
    dictService: DictService;
    config : Configuration;
    constructor() {
        this.config = Configuration.newInstance();

        this.dictService = new DictService(
            '/Users/chenquan/Workspace/nodejs/medict/testdict/testrscroot', 
            '/Users/chenquan/Workspace/nodejs/medict/testdict/testdict1', 
            this.config.configJsonFilePath) 

        this.dictService.autoIndexing();

        this.dictService.eventEmitter.on(EV_START_INDEXING,(ev, arg) =>{
            console.log(ev, arg);
        })

        this.dictService.eventEmitter.on(EV_END_INDEXING,(ev, arg) =>{
            console.log(ev, arg);
        })
    }
}