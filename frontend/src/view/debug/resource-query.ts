export function composeResourceRequestURL(serverUrl: string, inputValue: string, dictID: string) {
  const input = inputValue.trim();
  if (!serverUrl || !input) return '';
  if (/^https?:\/\//i.test(input)) return input;

  const server = new URL(serverUrl);
  const url = input.startsWith('/')
    ? new URL(input, server.origin)
    : new URL(input.replace(/^\.\//, ''), `${serverUrl.replace(/\/$/, '')}/`);

  if (dictID) url.searchParams.set('dict_id', dictID);
  url.searchParams.set('d', '0');
  return url.toString();
}
