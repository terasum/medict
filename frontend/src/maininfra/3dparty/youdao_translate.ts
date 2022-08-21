const API_URL = 'http://openapi.youdao.com/api';
import axios from 'axios'
import sha256 from 'crypto-js/sha256';;
import Hex from 'crypto-js/enc-hex';
import querystring from 'querystring';

function truncate(q: string) {
  var len = q.length;
  if (len <= 20) return q;
  return q.substring(0, 10) + len + q.substring(len - 10, len);
}

export async function translate(
  appid: string,
  appkey: string,
  from: string,
  to: string,
  query: string
) {
  if (from === 'zh') {
    from = 'zh-CHS';
  }
  if (from === 'jp') {
    from = 'ja';
  }
  if (to === 'zh') {
    to = 'zh-CHS';
  }
  if (to === 'jp') {
    to = 'ja';
  }

  const salt = (new Date).getTime();
  let curtime = Math.round(new Date().getTime() / 1000);
  let str1 = appid + truncate(query) + salt + curtime + appkey;
  let vocabId = '';

  let sign = sha256(str1).toString(Hex);

  let data = {
    q: query,
    appKey: appid,
    salt: salt,
    from: from,
    to: to,
    sign: sign,
    signType: "v3",
    curtime: curtime,
    vocabId: vocabId,
  }

  var queryData = querystring.stringify(data);
  const newUrl = `${API_URL}?${queryData}`;
  return axios.post(newUrl)
}
