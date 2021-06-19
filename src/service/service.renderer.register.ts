import { syncfn, asyncfn } from './service.manifest'
import { ipcRenderer } from 'electron'
// ensure this defination before than register function
const rendererAPI = {}

// register automatically
register();

function register() {
  for (const fnName in syncfn) {
    if (Object.prototype.hasOwnProperty.call(syncfn, fnName)) {
        console.log(`🔨 register renderer sync function: ${fnName}`);
        rendererAPI[fnName] = (args: any) => {
        return ipcRenderer.sendSync(fnName, args);
      }
    }
  }

  for (const fnName in asyncfn) {
    if (Object.prototype.hasOwnProperty.call(asyncfn, fnName)) {
        console.log(`🔨 register renderer async function: ${fnName}`);
        rendererAPI[fnName] = (args: any) => {
        ipcRenderer.send(fnName, args)
      }
    }
  }
}

export default rendererAPI;