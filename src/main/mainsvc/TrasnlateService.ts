import { logger } from '../../utils/logger';
import { translate as baiduApi } from '../../apis/baidu_translate';
import { translate as googleApi } from '../../apis/google_translate';
import { translate as youdaoApi } from '../../apis/youdao_translate';

import { ConfigAccessService } from './ConfigAccessorService';

export class TranslateService {
    config: ConfigAccessService

    constructor() {
        this.config = new ConfigAccessService();
    }


    async BaiduTranslate(query: string, from: string, to: string) {
        const config = this.config.loadTranslateApiConfig();
        if (!config.hasOwnProperty('baidu')) {
            throw new Error('无法读取百度翻译配置');
        }
        const appid = config.baidu.appid;
        const appkey = config.baidu.appkey;
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
        const config = this.config.loadTranslateApiConfig();
        if (!config.hasOwnProperty('youdao')) {
            throw new Error('无法读取有道翻译配置');
        }
        const appid = config.youdao.appid;
        const appkey = config.youdao.appkey;
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
