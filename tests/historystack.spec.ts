import { HistoryStack } from '../src/renderer/utils/history_stack';
var chai = require('chai');
var assert = chai.assert; // Using Assert style

describe('history stack test', () => {
  it('stack push pop', () => {
    const stack = new HistoryStack(5, 'word');
    stack.stackPush('s1');
    assert.deepEqual(['word', 's1'], stack.linkedList());
    stack.stackPush('s2');
    assert.deepEqual(['word', 's1', 's2'], stack.linkedList());
    stack.stackPush('s3');
    assert.deepEqual(['word', 's1', 's2', 's3'], stack.linkedList());
    stack.stackPush('s4');
    assert.deepEqual(['word', 's1', 's2', 's3', 's4'], stack.linkedList());
    stack.stackPush('s5');
    assert.deepEqual(['s1', 's2', 's3', 's4', 's5'], stack.linkedList());
    stack.stackPush('s6');
    assert.deepEqual(['s2', 's3', 's4', 's5', 's6'], stack.linkedList());
    stack.stackPush('s7');
    assert.deepEqual(['s3', 's4', 's5', 's6', 's7'], stack.linkedList());
    stack.stackPop();
    assert.deepEqual(['s3', 's4', 's5', 's6'], stack.linkedList());
    stack.stackPop();
    assert.deepEqual(['s3', 's4', 's5'], stack.linkedList());
    stack.stackPop();
    assert.deepEqual(['s3', 's4'], stack.linkedList());
    stack.stackPop();
    assert.deepEqual(['s3'], stack.linkedList());
    // will keep last one
    stack.stackPop();
    assert.deepEqual(['s3'], stack.linkedList());
    assert.strictEqual(stack.getSize(), 1);
    stack.stackPop();
    assert.deepEqual(['s3'], stack.linkedList());
    assert.strictEqual(stack.getSize(), 1);
  });
});
