import { translate } from '../../src/apis/baidu_translate';
import { assert } from 'chai';

import { appid, appkey } from './baidu_translate.api.key';

describe('百度翻译测试', () => {
  it('英文->中文', () => {
    return translate(appid, appkey, 'en', 'zh', 'apple').then(resp => {
      assert.deepEqual(resp.status, 200);
      assert.deepEqual(resp.data, {
        from: 'en',
        to: 'zh',
        trans_result: [{ src: 'apple', dst: '苹果' }],
      });
    });
  });
  it('中文->英文', () => {
    return translate(appid, appkey, 'zh', 'en', '这是一本词典').then(resp => {
      assert.deepEqual(resp.status, 200);
      assert.deepEqual(resp.data, {
        from: 'zh',
        to: 'en',
        trans_result: [{ src: '这是一本词典', dst: 'This is a dictionary' }],
      });
    });
  });
});
