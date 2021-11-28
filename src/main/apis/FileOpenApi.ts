import { FileOpenService } from '../mainsvc/FileOpenService';
import { getUserResourceRootPath, getResourceRootPath, getLoggerFilePath, webviewPreloadFilePath } from '../../config/config';

const fileOpenService = new FileOpenService();
export class FileOpenApi {
  syncShowOpenDialog(arg: {fileExtensions: string[] | undefined, multiFile: boolean}) {
    return fileOpenService.showOpenDialog(arg);
  }
  syncShowMainLoggerPath(arg?: any) {
    return getLoggerFilePath();
  }

  syncGetResourceRootPath(arg?: any) {
    return getResourceRootPath();
  }
  
  syncGetUserResourceRootPath(arg?: any) {
    return getUserResourceRootPath();
  }

  syncGetWebviewPreliadFilePath(arg?:any) {
    return webviewPreloadFilePath();
  }
}

