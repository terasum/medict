import { getResourceRootPath } from '../config/config';
import koa_logger from 'koa-logger';
import serve from 'koa-static';
import Koa from 'koa';
import getPort from 'get-port';
import { logger } from '../utils/logger';

export let resourceServerPort = 0;

export async function startResourceServer() {
  const app = new Koa();
  app.use(
    koa_logger((str: string, args: object) => {
      logger.info(str);
    })
  );
  app.use(serve(getResourceRootPath()));
  resourceServerPort = await getPort({ port: getPort.makeRange(40000, 50000) });
  app.listen(resourceServerPort);
  logger.info(
    `⚙ static-server listening on port http://localhost:${resourceServerPort}`
  );
}

export default {
  startResourceServer,
  resourceServerPort,
};
