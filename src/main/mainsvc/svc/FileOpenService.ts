import { dialog } from 'electron';

export class FileOpenService {
  showOpenDialog(fileExtensions: string[] | undefined) {
    console.log(fileExtensions);
    if (!fileExtensions) {
      fileExtensions = ['mdd', 'mdx'];
    }
    return dialog.showOpenDialogSync({
      properties: ['openFile', 'noResolveAliases', 'dontAddToRecent'],
      message: '选择文件',
      filters: [{ name: 'Custom File Type', extensions: fileExtensions }],
    });
  }
}
