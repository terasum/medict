var crypto = require('crypto');
var md5 = crypto.createHash('md5');

export function md5hash(name :string) {
    let result = md5.update(name).digest('hex');
    return result;
}