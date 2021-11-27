import { extractKeys } from '../src/infra/ReplacerJs';
import { assert } from 'chai';
// var assert = chai.assert; // Using Assert style
describe('extractKey', () => {
  it('normal-case', () => {
    assert.isTrue(
      extractKeys('<script src="test.js"></script>').has('test.js')
    );
  });
});
