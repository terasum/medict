import { FileOpenApi } from './apis/FileOpenApi';
const fileOpenApi = new FileOpenApi();

export const asyncfn = {
};

export const syncfn = {
  syncShowMainLoggerPath: fileOpenApi.syncShowMainLoggerPath,
  syncGetResourceRootPath: fileOpenApi.syncGetResourceRootPath,
  syncGetWebviewPreliadFilePath: fileOpenApi.syncGetWebviewPreliadFilePath,
};
