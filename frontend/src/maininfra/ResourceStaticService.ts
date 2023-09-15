/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import {getResourceRootPath} from './BasicConfig';
import koa_logger from 'koa-logger';
import serve from 'koa-static';
import Koa from 'koa';
import cors from '@koa/cors';
import getPort from 'get-port';
const logger = {
  info: function(...arg :any) {console.log(arg)},
  error: function(...arg :any) {console.error(arg)}
}

export let resourceServerPort = 0;

export async function startResourceServer() {
  const app = new Koa();
  app.use(
    koa_logger((str: string, args: object) => {
      logger.info(str);
    })
  );
  app.use(cors());
  app.use(serve(getResourceRootPath()));
  resourceServerPort = await getPort({ port: getPort.makeRange(40000, 50000) });
  app.listen(resourceServerPort);
  logger.info(
    `⚙ static-server listening on port http://localhost:${resourceServerPort}`
  );
}

export function staticServerPort() {
  return resourceServerPort;
}

export default {
  startResourceServer,
  resourceServerPort,
};
