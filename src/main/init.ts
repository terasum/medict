import { registerServices } from './init.mainsvc.register';
import { startResourceServer } from './init.resource.server';

// register main-process service for renderer
registerServices();
// start resource static server
startResourceServer();
