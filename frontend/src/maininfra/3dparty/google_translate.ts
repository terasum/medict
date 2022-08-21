import google_translate from '@vitalets/google-translate-api';
export async function translate(
    appid: string,
    appkey: string,
    from: string,
    to: string,
    query: string

) {
    if (from == "zh") {
        from = 'zh-CN'
    }
    if (from == "jp") {
        from = 'ja'
    }

    if (to == "zh") {
        to = 'zh-CN'
    }
    if (to == "jp") {
        to = 'ja'
    }

    // @ts-ignore
    let request = google_translate(query, { from: from, to: to, tld: 'cn' })
    return request
}
