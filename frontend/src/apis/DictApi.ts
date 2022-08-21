export default class DictAPI {
    async dicts() {
        return window.go.apis.DictsAPI.Dicts();
    }
    async suggest(word: string) {
        return window.go.apis.DictsAPI.Suggest(word);
    }
    async lookupDefinition(dict_id: string, raw_key_word: string, record_start: number) {
        return window.go.apis.DictsAPI.LookupDefinition(dict_id, raw_key_word, record_start);
    }
}