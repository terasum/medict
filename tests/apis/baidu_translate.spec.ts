import { translate } from '../../src/apis/baidu_translate';
import './baidu_translate.api.key';
import mockAxios from 'jest-mock-axios';
import { assert } from 'chai';

const appid = process.env.BAIDU_APPID || '';
const appkey = process.env.BAIDU_APP_KEY || '';

afterEach(() => {
  // cleaning up the mess left behind the previous test
  mockAxios.reset();
});


describe('百度翻译测试', () => {
  it('英文->中文', () => {
    let testdata = {data: {
      from: 'en',
      to: 'zh',
      trans_result: [{ src: 'apple', dst: '苹果' }],
    }};
    const promise = translate(appid, appkey, 'en', 'zh', 'apple');
    mockAxios.mockResponse(testdata);

    return promise.then(resp => {
      assert.deepEqual(resp.status, 200);
      assert.deepEqual(resp.data, {
        from: 'en',
        to: 'zh',
        trans_result: [{ src: 'apple', dst: '苹果' }],
      });
    });
  });

  it('中文->英文', () => {
    
    let testdata = {data: {
        from: 'zh',
        to: 'en',
        trans_result: [{ src: '这是一本词典', dst: 'This is a dictionary' }],
      }};

    const promise = translate(appid, appkey, 'zh', 'en', '这是一本词典');
    mockAxios.mockResponse(testdata);

    return promise.then(resp => {
      assert.deepEqual(resp.status, 200);
      assert.deepEqual(resp.data, {
        from: 'zh',
        to: 'en',
        trans_result: [{ src: '这是一本词典', dst: 'This is a dictionary' }],
      });
    });
  });
});
