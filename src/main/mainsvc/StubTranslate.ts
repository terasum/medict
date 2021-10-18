import { logger } from '../../utils/logger';
import { StubConfigAccessor } from './StubConfigAccessor';
import { translate as baiduApi } from '../../apis/baidu_translate';
import { translate as googleApi } from '../../apis/google_translate';
import { translate as youdaoApi } from '../../apis/youdao_translate';

const configAccessor = new StubConfigAccessor();

export class StubTranslate {
  asyncBaiduTranslate(
    event: any,
    arg: {
      query: string;
      from: string;
      to: string;
    }
  ) {
    const config = configAccessor.loadTranslateApiConfig();
    if (config.hasOwnProperty('baidu')) {
      const appid = config.baidu.appid;
      const appkey = config.baidu.appkey;
      baiduApi(appid, appkey, arg.from, arg.to, arg.query).then(resp => {
        if (resp && resp.status === 200) {
          event.sender.send('onAsyncTranslate', {
            engine: "baidu",
            data: resp.data,
            code: 0,
            message: 'success',
          });
        } else {
          event.sender.send('onAsyncTranslate', {
            engine: "baidu",
            data: resp.data,
            code: -1,
            message: '翻译失败:' + resp.status,
          });
        }
      });
    }
  }

  asyncGoogleTranslate(
    event: any,
    arg: {
      query: string;
      from: string;
      to: string;
    }
  ) {
    // https://github.com/vitalets/google-translate-api/blob/master/languages.js
    if (arg.from === 'zh') {
      arg.from = 'zh-CN'
    }
    googleApi("appid", "appkey", arg.from, arg.to, arg.query).then(resp => {
      console.log(resp.text);
      event.sender.send('onAsyncTranslate', {
        engine: "google",
        data: resp.text,
        code: 0,
        message: 'success',
      });
    }).catch(err => {
      event.sender.send('onAsyncTranslate', {
        engine: "google",
        data: "",
        code: -1,
        message: err,
      });
    });
  }

  asyncYoudaoTranslate(
    event: any,
    arg: {
      query: string;
      from: string;
      to: string;
    }
  ) {
    const config = configAccessor.loadTranslateApiConfig();
    if (config.hasOwnProperty('youdao')) {
      const appid = config.youdao.appid;
      const appkey = config.youdao.appkey;
      youdaoApi(appid, appkey, arg.from, arg.to, arg.query).then(resp => {
        if (resp && resp.status === 200) {
          if (resp.data.errorCode === '0' && resp.data.translation && resp.data.translation.length > 0) {
            event.sender.send('onAsyncTranslate', {
              engine: "youdao",
              data: resp.data.translation[0],
              code: 0,
              message: 'success',
            });
          } else {
            event.sender.send('onAsyncTranslate', {
              engine: "youdao",
              data: resp.data,
              code: -1,
              message: 'failed, errorcode: ' + resp.data.errorCode,
            });
          }
        } else {
          event.sender.send('onAsyncTranslate', {
            engine: "youdao",
            data: resp.data,
            code: -1,
            message: '翻译失败:' + resp.status,
          });
        }
      });
    }
  }

}
