import { exportAllDeclaration } from '@babel/types';
import {Dictionary} from '../../src/infra/Dictionary';

describe('Test lookup2', () => {
    it('王力古漢語字典-論', () => {
       const dict = new Dictionary('testwangli', 
       '王力古汉语','王力古汉语词典',
       'testdict/testdict1/王力古漢語字典 (2000)/王力古漢語字典 (2000).mdx',
       ['testdict/testdict1/王力古漢語字典 (2000)/王力古漢語字典 (2000).mdd'],'test');
       const result = dict.lookup('論');
       expect(result).not.toBe(null);
       expect(result.keyText).toBe('論');
    });
    it('王力古漢語字典-滋', () => {
       const dict = new Dictionary('testwangli', 
       '王力古汉语','王力古汉语词典',
       'testdict/testdict1/王力古漢語字典 (2000)/王力古漢語字典 (2000).mdx',
       ['testdict/testdict1/王力古漢語字典 (2000)/王力古漢語字典 (2000).mdd'],'test');
       const result = dict.lookup('滋');
       console.log(result)
       console.log('滋', unicodeLiteral(result.keyText));

       const result2 = dict.lookup('滋');
       console.log(result2)
       console.log('@@@LINK=滋\r\n\x00', unicodeLiteral('滋'))

       expect(result).not.toBe(null);
       expect(result.keyText).toBe('滋');
    });
});


/* Creates a unicode literal based on the string */    
function unicodeLiteral(str: string){
    let result = "";
    for(let i = 0; i < str.length; ++i){
        /* You should probably replace this by an isASCII test */
        if(str.charCodeAt(i) > 126 || str.charCodeAt(i) < 32)
            result += "\\u" + fixedHex(str.charCodeAt(i),4);
        else
            result += str[i];
    }
    return result;
}

function fixedHex(number: number, length: number){
    var str = number.toString(16).toUpperCase();
    while(str.length < length){
        str = "0" + str;
    }
    return str;
}
