/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

export class FuzzyTrieNode {
    _char: string = '';
    _data: any;
    _level: number = 0;
    _children: Map<string, FuzzyTrieNode> = new Map();
    // _children: any = {};

    _endWord: boolean = false;

    constructor(char: string, level?: number, data?: any) {
        this._char = char;
        this._endWord = false;
        this._level = level || 0;
        this._data = data;
    }

    getChar() {
        return this._char;
    }

    getData() {
        return this._data;
    }

    isChar(char: string): boolean {
        return this._char === char;
    }

    isEnd(): boolean {
        return this._endWord;
    }

    hasChild(char: string) {
        // return this._children[char] !== undefined;
        return this._children.has(char);
    }

    addChild(char: string, data?: any) {
        // this._children[char] = new FuzzyTrieNode(char, isEnd);
        this._children.set(char, new FuzzyTrieNode(char, this._level + 1, data));
    }

    setEnd(isEnd: boolean) {
        this._endWord = isEnd;
    }

    getChild(char: string): FuzzyTrieNode {
        // return this._children[char]!;
        return this._children.get(char)!;
    }

}

export default class FuzzyTrie {
    _root: FuzzyTrieNode = new FuzzyTrieNode('');
    _size: number = 0;
    _level: number = 0;
    _bytesSize: number = 0;
    constructor() {
    }

    add(word: string, data?: any) {
        let node = this._root;
        for (let i = 0; i < word.length; ++i) {
            let c = word.charAt(i);
            if (!node.hasChild(c)) {
                // set last word char's end word flag = true
                node.addChild(c, data);
                this._size = this._size + 1;
                this._bytesSize = this._bytesSize + c.length;
            }
            node = node.getChild(c);
            this._level = this._level < node._level ? node._level : this._level;
            if (i == word.length - 1) {
                node.setEnd(true);
            }
        }
    }

    has(word: string) {
        var node = this._root;
        for (var i = 0; i < word.length; ++i) {
            var c = word.charAt(i);
            if (!node.hasChild(c)) {
                return null;
            }
            node = node.getChild(c);
        }
        return node.isEnd() ? node : null;
    }

    json() {
        const stack: FuzzyTrieNode[] = [];
        stack.push(this._root);
        while (stack.length > 0) {
            let node = stack.pop();
            console.log(node?._char)
            for (let c in node!._children.keys()) {
                stack.push(node!._children.get(c)!)
            }
        }
    }

    all() {
        let result: string[] = [];
        this._all(result, this._root, '');
        return result;
    }

    _all(result: string[], node: FuzzyTrieNode, stack: string) {
        if (!node) {
            return
        }

        for (let k of node._children.keys()) {
            let knode = node._children.get(k);
            if (knode?.isEnd()) {
                result.push(stack + k);
            }
            this._all(result, node._children.get(k)!, stack + k);
        }
    }

    // delete(word) {
    //     var node = this._root;
    //     var stack = [node];
    //     for (var i = 0; i < word.length; ++i) {
    //         var c = word.charAt(i);
    //         if (node[c] === undefined)
    //             return false;
    //         node = node[c];
    //         stack.push(node);
    //     }
    //     if (!node[stop])
    //         return false;
    //     delete node[stop];
    //     for (var i = stack.length - 2; i >= 0; --i) {
    //         if (!this._empty(stack[i + 1]))
    //             break;
    //         delete stack[i][word.charAt(i)];
    //     }
    //     return true;
    // }

    empty() {
        return this._empty(this._root);
    }

    _empty(node: FuzzyTrieNode) {
        if (node._char === '' && node._children.size == 0) {
            return true
        }
        return false;
    }

    tail(node: FuzzyTrieNode) {
        let result: string[] = [];
        this._all(result, node, '');
        return result;
    }

    _locate(result: { key: string, data?: any }[], node: FuzzyTrieNode, stack: string) {
        if (!node) {
            return
        }

        for (let k of node._children.keys()) {
            let knode = node._children.get(k);
            if (knode?.isEnd()) {
                result.push({ key: stack + k, data: knode.getData() });
            }
            this._locate(result, node._children.get(k)!, stack + k);
        }
    }


    prefix(pfx: string): { key: string, data?: any }[] {
        if (pfx.length < 1) {
            return [];
        }
        let node = this._root;
        for (var i = 0; i < pfx.length; ++i) {
            var c = pfx.charAt(i);
            if (!node.hasChild(c)) {
                break
            }
            node = node.getChild(c);
        }
        // still root
        if (node._char == '') {
            return []
        }
        let result: { key: string, data?: any }[] = [];
        this._locate(result, node, '');
        result.forEach(ele => {
            ele.key = pfx + ele.key;
        })
        return result;

    }

    // _find(word, result, errors, replaces, pos, node, stack, inserted, tail) {
    //     var c = (pos === word.length) ? undefined : word.charAt(pos);
    //     if (tail === undefined) {
    //         if (c === undefined) {
    //             if (node[stop] && (result[stack] === undefined || result[stack] < errors))
    //                 result[stack] = errors;
    //         } else if (node[c] !== undefined) {
    //             this._find(word, result, errors, 0, pos + 1, node[c], stack + c);
    //         }
    //         if (errors) {
    //             // delete
    //             if (c !== undefined && !inserted)
    //                 this._find(word, result, errors - 1, replaces + 1, pos + 1, node, stack);
    //             // transposition
    //             var maxpos = Math.min(pos + errors, word.length - 1);
    //             for (var i = pos + 1; i <= maxpos; ++i) {
    //                 var c2 = word.charAt(i);
    //                 if (c2 === c)
    //                     continue;
    //                 var node2 = node[c2];
    //                 if (node2 !== undefined) {
    //                     var node3 = node2[c];
    //                     // with delete of middle part
    //                     if (node3 !== undefined)
    //                         this._find(word, result, errors - (i - pos), 0, i + 1, node3, stack + c2 + c);
    //                     if (i === pos + 1) {
    //                         // with insert into the middle part
    //                         this._find(word, result, errors - 1, 0, pos + 2, node2, stack + c2, undefined, c);
    //                     }
    //                 }
    //             }
    //         }
    //     } else if (inserted && node[tail]) {
    //         this._find(word, result, errors, 0, pos, node[tail], stack + tail);
    //     }
    //     if (errors || replaces) {
    //         // insert / replace
    //         for (var k in node) {
    //             if (k === stop || (k === c && !replaces))
    //                 continue;
    //             this._find(word, result, replaces ? errors : errors - 1, replaces ? replaces - 1 : 0, pos, node[k], stack + k, true, tail)
    //         }
    //     }
    // }

}
