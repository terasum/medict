import MD5 from './md5.js';

export function sign(
  appid: string,
  q: string,
  salt: string,
  appkey: string
): string {
  const line = appid + q + salt + appkey;
  //   return line;
  return MD5(line);
}
