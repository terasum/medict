const API_URL = 'http://api.fanyi.baidu.com/api/trans/vip/translate';
import axios from 'axios';
import { sign } from './baidu_translate_utils/utils';
import random from '@shuaninfo/random';

var querystring = require('querystring');

export async function translate(
  appid: string,
  appkey: string,
  from: string,
  to: string,
  query: string
) {
  const salt = random({ length: 10, type: 'numeric' });
  const signData = sign(appid, query, salt, appkey);

  var queryData = querystring.stringify({
    q: query,
    from,
    to,
    appid,
    salt,
    sign: signData,
  });

  const newUrl = `${API_URL}?${queryData}`;
  return await axios.get(newUrl);
}
