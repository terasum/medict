import { createByProc } from '@terasum/electron-call';
import {ClipboardApi} from './apis/ClipboardApi';
import {ConfigAccessorApi} from './apis/ConfigAccessorApi';
import {DictAccessorApi} from './apis/DictAccessorApi';
import {FileOpenApi} from './apis/FileOpenApi';
import {MessageApi} from './apis/MessageApi';
import {TranslateApi} from './apis/TranslateApi';

const stubByMain = createByProc('main', 'error');


const clipboardApi =new ClipboardApi();
const configAccessorApi =new ConfigAccessorApi();
const dictAccessorApi =new DictAccessorApi();
const fileOpenApi =new FileOpenApi();
const messageApi =new MessageApi();
const translateApi =new TranslateApi();

stubByMain.provide(['renderer','worker'],'ClipboardApi', clipboardApi);
stubByMain.provide(['renderer','worker'],'ConfigAccessorApi', configAccessorApi);
stubByMain.provide(['renderer','worker'],'DictAccessorApi', dictAccessorApi);
stubByMain.provide(['renderer','worker'],'FileOpenApi', fileOpenApi);
stubByMain.provide(['renderer','worker'],'MessageApi', messageApi);
stubByMain.provide(['renderer','worker'],'TranslateApi', translateApi);