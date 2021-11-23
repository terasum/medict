import Configuration from './Configuration.svc';
import { DictService } from './Dictionary.svc';
console.log('[WORKER] ======= Worker services initing... =======');

// init load, this function will load during startup process (main-process)
// so, if load dictionary failed, it will block main-process
// we should try-catch the errors
// synchronized load dictionaries  process

(function () {
    Configuration.newInstance();
    new Promise((resolve) => {
        let dict = new DictService();
        resolve(dict);
    }).then(() => {
        console.log('[WORKER] ======= Worker dictionary reloaded ... =======');
    }).catch(error => {
        console.error('[WORKER] ======= Worker dictionary reloaded error =======');
        console.error(error)
    })
})();

