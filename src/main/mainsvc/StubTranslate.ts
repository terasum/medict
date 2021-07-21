import { logger } from '../../utils/logger';
import { StubConfigAccessor } from './StubConfigAccessor';
import { translate as baiduApi } from '../../apis/baidu_translate';

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
          event.sender.send('onAsyncBaiduTranslate', {
            data: resp.data,
            code: 0,
            message: 'success',
          });
        } else {
          event.sender.send('onAsyncBaiduTranslate', {
            data: resp.data,
            code: -1,
            message: '翻译失败:' + resp.status,
          });
        }
      });
    }
  }
}
