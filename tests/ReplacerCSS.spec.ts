import { extractKeys } from '../src/infra/ReplacerCSS';
import { assert } from 'chai';
// var assert = chai.assert; // Using Assert style
describe('extractKeys', () => {
  it('normal-case', () => {
    assert.isTrue(extractKeys('<link href="test.css"></link>').has('test.css'));
  });
});
