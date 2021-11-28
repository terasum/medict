import { logger } from '../../utils/logger';
import { translate as baiduApi } from '../../apis/baidu_translate';
import { translate as googleApi } from '../../apis/google_translate';
import { translate as youdaoApi } from '../../apis/youdao_translate';

import Configuration from './Configuration.svc';

export class TranslateService {
    config: Configuration;

    constructor() {
        this.config = Configuration.newInstance();
    }


    async BaiduTranslate(query: string, from: string, to: string) {
        const appid = this.config.getBaiduAppID();
        const appkey = this.config.getBaiduAppKey();
        let resp = await baiduApi(appid, appkey, from, to, query);
        if (resp && resp.status === 200) {
            return resp.data;
        } else {
            throw new Error('翻译失败:' + resp.status)
        }
    }

    async GoogleTranslate(query: string, from: string, to: string) {
        // https://github.com/vitalets/google-translate-api/blob/master/languages.js
        if (from === 'zh') {
            from = 'zh-CN'
        }
        let resp = await googleApi("appid", "appkey", from, to, query);
        return resp.text;
    }

    async YoudaoTranslate(query: string, from: string, to: string) {
        const appid = this.config.getYoudaoAppID();
        const appkey = this.config.getYoudaoAppKey();
        let resp = await youdaoApi(appid, appkey, from, to, query);
        if (!resp || resp.status !== 200) {
            throw new Error('翻译失败:' + resp.status);
        }
        if (resp.data.errorCode === '0' && resp.data.translation && resp.data.translation.length > 0) {
            return resp.data.translation[0];
        } else {
            throw new Error('翻译失败, ERROR: ' + resp.data.errorCode);
        }
    }

}
