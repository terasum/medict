import { HistoryStack } from '../src/renderer/utils/history_stack';

describe('history stack test', () => {
  it('stack push pop', () => {
    const stack = new HistoryStack(5, 'word');
    stack.stackPush('s1');
    console.log(stack.linkedList());
    stack.stackPush('s2');
    console.log(stack.linkedList());
    stack.stackPush('s3');
    console.log(stack.linkedList());
    stack.stackPush('s4');
    console.log(stack.linkedList());
    stack.stackPush('s5');
    console.log(stack.linkedList());
    stack.stackPush('s6');
    console.log(stack.linkedList());
    stack.stackPush('s7');
    console.log(stack.linkedList());
    console.log('-------');
    stack.stackPop();
    console.log(stack.linkedList());
    stack.stackPop();
    console.log(stack.linkedList());
    stack.stackPop();
    console.log(stack.linkedList());
    stack.stackPop();
    console.log(stack.linkedList());
  });
});
