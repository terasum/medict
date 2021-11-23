import { Config } from '../../model/Config';
import { getConfigJsonPath } from '../../config/config';
import { StorageService } from '../mainsvc/StorageServcice';

const storageService = new StorageService(getConfigJsonPath());

export class ConfigAccessService {

loadTranslateApiConfig() {
    const defaultConfig = new Config();
    const rawApis = storageService.getDataByKey('translateApis');
    if (!rawApis) {
        return defaultConfig.translateApis;
      }

    if (!rawApis.youdao) {
        return defaultConfig.translateApis;
    }

    if (!rawApis.baidu) {
        return defaultConfig.translateApis;
    }

    return rawApis as {
      baidu: {
        appkey: string;
        appid: string;
      };
      youdao: {
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
  saveTranslateYoudaoApiConfig(args: { appid: string; appkey: string }) {
    let rawApis = storageService.getDataByKey('translateApis');
    if (!rawApis) {
      rawApis = {
        youdao: {
          appkey: args.appkey,
          appid: args.appid,
        }
      };
    } else {
      rawApis['youdao'] = {
        appkey: args.appkey,
        appid: args.appid,
      };
    }
    return storageService.setDataByKey('translateApis', rawApis);
  }
}
