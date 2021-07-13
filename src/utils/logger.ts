// import pino from 'pino';
// import prettifier from 'pino-pretty';

import { getLoggerFilePath } from '../config/config';

import log from 'electron-log';

log.transports.file.resolvePath = () => getLoggerFilePath();

export const logger = log.scope('main-process');

// export const logger = pino({
//   prettyPrint: {
//     levelFirst: true,
//     ignore: 'pid,hostname', // --ignore
//     crlf: false, // --crlf
//     errorLikeObjectKeys: ['err', 'error'], // --errorLikeObjectKeys
//     errorProps: '', // --errorProps
//     messageKey: 'msg', // --messageKey
//     levelKey: 'level', // --levelKey
//     messageFormat: false, // --messageFormat
//     timestampKey: 'time', // --timestampKey
//     translateTime: true, // --translateTime
//     hideObject: false, // --hideObject
//     singleLine: false, // --singleLine
//   },
//   prettifier: prettifier,
// });
