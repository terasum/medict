import { FileOpenService } from '../mainsvc/FileOpenService';

const fileOpenService = new FileOpenService();
export class FileOpenAPI {
  syncShowOpenDialog(arg: {fileExtensions: string[] | undefined, multiFile: boolean}) {
    return fileOpenService.showOpenDialog(arg);
  }

  syncShowOpenDirDialog() {
    return fileOpenService.showOpenDirDialog();
  }
}

