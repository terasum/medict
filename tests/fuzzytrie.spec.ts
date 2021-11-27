import FuzzyTrie from '../src/infra/indexing/FuzzyTrie';
import Mdict from 'js-mdict';
var chai = require('chai');
var assert = chai.assert; // Using Assert style

describe('test FuzzyTrie', () => {
    it('FuzzyTrie add', () => {
        const ft = new FuzzyTrie();
        ft.add('hell');
        ft.add('hello');
        ft.add('hellworld');
        ft.add('fight');
        ft.add('future');
        ft.add('finger');
        expect(ft.has('hell')).not.toBe(null);
        expect(ft.has('hello')).not.toBe(null);
        expect(ft.has('hellworld')).not.toBe(null);
        expect(ft.has('fight')).not.toBe(null);
        expect(ft.has('future')).not.toBe(null);
        expect(ft.has('finger')).not.toBe(null);
        expect(ft.has('finger11123')).toBe(null);
        expect(ft.all()).toEqual(['hell', 'hello', 'hellworld', 'fight', 'finger', 'future']);
        expect(ft._size).toBe(24);
        expect(ft._level).toBe(9);
        expect(ft._bytesSize).toBe(24);
        expect(ft.prefix('fi')).toEqual([
            { key: 'fight', data: undefined },
            { key: 'finger', data: undefined }
        ])
        expect(ft.prefix('he')).toEqual([
            { key: 'hell', data: undefined },
            { key: 'hello', data: undefined },
            { key: 'hellworld', data: undefined }
        ])
    });
    
    it('Clowded FuzzySearch', () => {
        const mdict = new Mdict('./testdict/testdict1/新世纪汉英大词典/新世纪汉英大词典.mdx');
        console.time('rangeKey');
        const keyWords = mdict.rangeKeyWords();
        console.timeEnd('rangeKey');
        console.log(`key word size ${keyWords.length}`);

        const ft = new FuzzyTrie();
        console.time('build-trie');
        for (let word of keyWords) {
            ft.add(word.keyText, word);
        }
        console.timeEnd('build-trie');
        console.log(`fuzzyTrie level: ${ft._level}`)

        console.time('prefix');
        console.log(ft.prefix('一丝'));
        console.timeEnd('prefix');
        console.log(ft.has('一丝不苟'));
    })
});
