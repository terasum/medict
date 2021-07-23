import { getConfigJsonPath } from '../../config/config';
import { StorageService } from './svc/StorageServcice';
const storageService = new StorageService(getConfigJsonPath());

export class StubConfigAccessor {
  loadTranslateApiConfig() {
    const rawApis = storageService.getDataByKey('translateApis');
    if (!rawApis) {
      return {
        baidu: {
          appkey: '',
          appid: '',
        },
      };
    }
    return rawApis as {
      baidu: {
        appkey: string;
        appid: string;
      };
    };
  }
  saveTranslateBaiduApiConfig(args: { appid: string; appkey: string }) {
    let rawApis = storageService.getDataByKey('translateApis');
    if (!rawApis) {
      rawApis = {
        baidu: {
          appkey: args.appkey,
          appid: args.appid,
        },
      };
    } else {
      rawApis['baidu'] = {
        appkey: args.appkey,
        appid: args.appid,
      };
    }
    return storageService.setDataByKey('translateApis', rawApis);
  }
}
