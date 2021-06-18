import rpc from 'pauls-electron-rpc'
import {API_SERVICE_NAME, manifest} from '../service/service.manifest';

const api = rpc.exportAPI(API_SERVICE_NAME, manifest, {
  sayHello: () => {
    return Promise.resolve('hello!')
  }
})

// log any errors
api.on('error', console.log);
