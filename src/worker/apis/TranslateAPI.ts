import { TranslateService } from '../worksvc/TranslateService';

function buildResp(engine: string, data: any, error?: Error, code?: Number) {
    if (error) {
        return {
            engine: engine,
            data: undefined,
            code: code || -1,
            message: 'Failed:' + error,
        };
    } else {
        return {
            engine: engine,
            data: data,
            code: 0,
            message: 'Success',
        };
    }

}

export class TranslateAPI {
    translate: TranslateService

    constructor() {
        this.translate = new TranslateService();
    }




    async asyncBaiduTranslate(arg: { query: string; from: string; to: string; }) {
        if (!arg || !arg.from || !arg.to || !arg.query) {
            return buildResp('baidu', undefined, new Error('invalid args'), -1);
        }
        try {
            const text = await this.translate.BaiduTranslate(arg.from, arg.to, arg.query)
            if (text) {
                return buildResp('baidu', text);
            }
            return buildResp('baidu', undefined, new Error('baidu translate result is null'), -2);
        } catch (err) {
            return buildResp('baidu', undefined, err as Error, -3);
        }
    }

    async asyncGoogleTranslate(arg: { query: string; from: string; to: string; }) {
        if (!arg || !arg.from || !arg.to || !arg.query) {
            return buildResp('google', undefined, new Error('invalid args'), -1);
        }

        try {
            const text = await this.translate.GoogleTranslate(arg.from, arg.to, arg.query)
            if (text) {
                return buildResp('google', text);
            }
            return buildResp('google', undefined, new Error('google translate result is null'), -2);
        } catch (err) {
            return buildResp('google', undefined, err as Error, -3);
        }
    }

    async asyncYoudaoTranslate(arg: { query: string; from: string; to: string; }) {
        if (!arg || !arg.from || !arg.to || !arg.query) {
            return buildResp('youdao', undefined, new Error('invalid args'), -1);
        }

        try {
            const text = await this.translate.YoudaoTranslate(arg.from, arg.to, arg.query)
            if (text) {

                return buildResp('youdao', text);
            }
            return buildResp('youdao', undefined, new Error('youdao translate result is null'), -2);
        } catch (err) {
            return buildResp('youdao', undefined, err as Error, -3);
        }

    }
}

