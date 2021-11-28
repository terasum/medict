
import {registerAPIs} from './rpc.sync.main.register';
import { startResourceServer } from './init.resource.server';

// register main-sync API
registerAPIs();
// start resource static server
startResourceServer();
