import { ConfigAccessService } from '../mainsvc/ConfigAccessorService';

export class ConfigAccessorApi {
  service: ConfigAccessService;
  constructor() {
    this.service = new ConfigAccessService();
  }

  loadTranslateApiConfig() {
    return this.service.loadTranslateApiConfig();
  }

  saveTranslateBaiduApiConfig(args: { appid: string; appkey: string }) {
    return this.service.saveTranslateBaiduApiConfig(args);
  }
  saveTranslateYoudaoApiConfig(args: { appid: string; appkey: string }) {
    return this.service.saveTranslateBaiduApiConfig(args);
  }
}

