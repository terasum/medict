import { createByProc } from '@terasum/electron-call';
import { ClipboardApi } from './apis/ClipboardApi';
import { ConfigAccessorApi } from './apis/ConfigAccessorApi';
import { DictAccessorApi } from './apis/DictAccessorApi';
import { FileOpenApi } from './apis/FileOpenApi';
import { MessageApi } from './apis/MessageApi';
import { TranslateApi } from './apis/TranslateApi';
import { WindowApi } from './apis/WindowApi';
import {registerServices} from './rpc.sync.main.register';

const stubByMain = createByProc('main', 'error');

stubByMain.provide(['renderer','worker'],'ClipboardApi', new ClipboardApi());
stubByMain.provide(['renderer','worker'],'ConfigAccessorApi', new ConfigAccessorApi());
stubByMain.provide(['renderer','worker'],'DictAccessorApi', new DictAccessorApi());
stubByMain.provide(['renderer','worker'],'FileOpenApi', new FileOpenApi());
stubByMain.provide(['renderer','worker'],'MessageApi', new MessageApi());
stubByMain.provide(['renderer','worker'],'TranslateApi', new TranslateApi());
stubByMain.provide(['renderer','worker'],'WindowApi', new WindowApi());

registerServices();