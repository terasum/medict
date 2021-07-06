import { convertToObject } from 'typescript';
import { FileOpenService } from './svc/FileOpenService';
const fileOpenService = new FileOpenService();

export class StubFileOpen {
  syncShowOpenDialog(arg: string[]) {
    console.log('syncShowOpenDialog - arg');
    console.log(arg);
    return fileOpenService.showOpenDialog(arg);
  }
}
