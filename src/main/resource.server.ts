import { getResourceRootPath } from '../config/config';
import logger from 'koa-logger';
import serve from 'koa-static';
import Koa from 'koa';
import getPort from 'get-port';


export let resourceServerPort = 0;

export async function  startServer() {
    const app = new Koa();
    app.use(logger())
    app.use(serve(getResourceRootPath()));
    resourceServerPort = await getPort({ port: getPort.makeRange(40000, 50000) });
    app.listen(resourceServerPort);
    console.log(`⚙ static-server listening on port http://localhost:${resourceServerPort}`);
}

export default {
    startServer,
    resourceServerPort
}