if (process.env.NODE_ENV !== 'production') {
  require('dotenv').config();
}
console.log('Now the value for FOO is:', process.env.BAIDU_APPID);
console.log('Now the value for FOO is:', process.env.BAIDU_APP_KEY);
