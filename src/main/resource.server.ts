// import http from 'http';
import { getResourceRootPath } from '../config/config'
import logger from 'koa-logger';

// import url from 'url';

// import fs from 'fs';
// import path from 'path';
// const port = process.argv[2] || "9001";


// http.createServer(function (req, res) {
//   console.log(`${req.method} ${req.url}`);
//   // parse URL
//   const parsedUrl = url.parse(req.url!);

//   // extract URL path
//   let pathname = `.${parsedUrl.pathname}`;
//   // based on the URL path, extract the file extention. e.g. .js, .doc, ...
//   const ext = path.parse(pathname).ext;
//   // maps file extention to MIME typere
//   const map = {
//     '.ico': 'image/x-icon',
//     '.html': 'text/html',
//     '.js': 'text/javascript',
//     '.json': 'application/json',
//     '.css': 'text/css',
//     '.png': 'image/png',
//     '.jpg': 'image/jpeg',
//     '.wav': 'audio/wav',
//     '.mp3': 'audio/mpeg',
//     '.svg': 'image/svg+xml',
//     '.pdf': 'application/pdf',
//     '.doc': 'application/msword'
//   };

//   fs.stat(pathname, function (exist) {
//     let relativePath = path.resolve(getResourceRootPath(), pathname);

//     if (!exist) {
//       // if the file is not found, return 404
//       res.statusCode = 404;
//       res.end(`File ${relativePath} not found!`);
//       return;
//     }

//     // if is a directory search for index file matching the extention
//     if (fs.statSync(relativePath).isDirectory()) {
//       relativePath += '/index' + ext
//     };

//     // read file from file system
//     fs.readFile(relativePath, function (err, data) {
//       if (err) {
//         res.statusCode = 500;
//         res.end(`Error getting the file: ${err}.`);
//       } else {
//         // if the file is found, set Content-type and send data
//         res.setHeader('Content-type', map[ext] || 'text/plain');
//         res.end(data);
//       }
//     });
//   });


// }).listen(parseInt(port));

// console.log(`📃 Resource Server listening on port ${port}`);


import serve from 'koa-static';
import Koa from 'koa';

const port = process.argv[2] || "9001";
const app = new Koa();
app.use(logger())

app.use(serve(getResourceRootPath()));


app.listen(port);

console.log(`⚙ static-server listening on port http://localhost:${port}`);