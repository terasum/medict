import { ConfigAccessService } from '../mainsvc/ConfigAccessorService';
import { TranslateService } from '../mainsvc/TrasnlateService';

export class TranslateApi {
  config: ConfigAccessService
  translate: TranslateService

  constructor() {
    this.config = new ConfigAccessService();
    this.translate = new TranslateService();
  }

  __sendReturn(event: any, engine: string, data: any, error?: Error, code?: Number) {
    if (error) {
      event.sender.send('onAsyncTranslate', {
        engine: engine,
        data: undefined,
        code: code || -1,
        message: 'Failed:' + error,
      });
    } else {
      event.sender.send('onAsyncTranslate', {
        engine: engine,
        data: data,
        code: 0,
        message: 'Success',
      });
    }

  }


  asyncBaiduTranslate(event: any, arg: { query: string; from: string; to: string; }) {
    if (!arg || !arg.from || !arg.to || !arg.query) {
      this.__sendReturn(event, 'baidu', undefined, new Error('invalid args'), -1);
    }
    this.translate.BaiduTranslate(arg.from, arg.to, arg.query)
      .then(text => {
        this.__sendReturn(event, 'baidu', text);
      }).catch(error => {
        this.__sendReturn(event, 'baidu', undefined, error, -2);
      })
  }

  asyncGoogleTranslate(event: any, arg: { query: string; from: string; to: string; }) {
    if (!arg || !arg.from || !arg.to || !arg.query) {
      this.__sendReturn(event, 'google', undefined, new Error('invalid args'), -1);
    }
    this.translate.GoogleTranslate(arg.from, arg.to, arg.query)
      .then(text => {
        this.__sendReturn(event, 'google', text);
      }).catch(error => {
        this.__sendReturn(event, 'google', undefined, error, -2);
      })
  }

  asyncYoudaoTranslate(event: any, arg: { query: string; from: string; to: string; }) {
    if (!arg || !arg.from || !arg.to || !arg.query) {
      this.__sendReturn(event, 'youdao', undefined, new Error('invalid args'), -1);
    }
    this.translate.YoudaoTranslate(arg.from, arg.to, arg.query)
      .then(text => {
        this.__sendReturn(event, 'youdao', text);
      }).catch(error => {
        this.__sendReturn(event, 'youdao', undefined, error, -2);
      })
  }
}

