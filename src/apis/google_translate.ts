import google_translate from '@vitalets/google-translate-api';
export async function translate(
    appid: string,
    appkey: string,
    from: string,
    to: string,
    query: string

) {
    // @ts-ignore
    let request = google_translate(query, { from: from, to: to, tld: 'cn' })
    return request
}
