import Configuration from './worksvc/Configuration.svc';
import './service.worker.manifest';
console.log('[WORKER] ======= Worker services initing... =======');


// init load, this function will load during startup process (main-process)
// so, if load dictionary failed, it will block main-process
// we should try-catch the errors
// synchronized load dictionaries  process

(function () {
    console.log('----------- loading... -------------')
    Configuration.newInstance();
}());

