import { syncfn, asyncfn } from './service.manifest';
import {ipcMain} from 'electron';


function register(){
    for (const fnName in syncfn) {
        if (Object.prototype.hasOwnProperty.call(syncfn, fnName)) {
            const fn = syncfn[fnName];
            console.log("\n🔧 register main process service:", fnName);
            ipcMain.on(fnName, function(event, args) {
                console.log(`[main-rpc:sync]: ${fnName}| arg: ${args}`);
                const ret = fn(args);
                console.log(`[main-rpc:sync]: ${fnName}| ret: ${ret}`);
                event.returnValue = ret;
            });
        }
    }
    for (const fnName in asyncfn) {
        if (Object.prototype.hasOwnProperty.call(asyncfn, fnName)) {
            const fn = asyncfn[fnName];
            ipcMain.on(fnName, function(event, args) {
                console.log(`[main-rpc:asyn]: ${fnName} - ${args}`);
                event.returnValue = fn(event, args);
            });
        }
    }
}

// register main process service
register();