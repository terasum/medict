var crypto = require('crypto');

export function md5hash(name :string) {
    var md5 = crypto.createHash('md5');
    let result = md5.update(name).digest('hex');
    return result;
}