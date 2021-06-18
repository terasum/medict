import rpc from 'pauls-electron-rpc'
import { API_SERVICE_NAME, manifest } from '../service/service.manifest'

// import over the 'example-api' channel
const api = rpc.importAPI(API_SERVICE_NAME, manifest, { timeout: 30e3 })

// now use, as usual:
// api.readFileSync('/etc/hosts') // => '...'
api.sayHello().then((s:string) =>{
    console.log(s);
})
