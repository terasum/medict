import { BIconArrowReturnLeft } from 'bootstrap-vue';

export class HistoryStack {
  first: HistoryStackNode;
  last: HistoryStackNode;
  size: number;
  limit: number;

  constructor(limit: number, data: any) {
    const node = new HistoryStackNode(null, null, data);
    this.first = node;
    this.last = node;
    this.size = 1;
    this.limit = limit;
  }

  linkedList(): HistoryStackNode[] {
    let node: HistoryStackNode | null = this.first;
    const data = [];
    while (node != null) {
      data.push(node.data);
      node = node.next;
    }
    return data;
  }

  getSize(): number {
    return this.size;
  }

  top(): HistoryStackNode {
    return this.last;
  }

  stackPop(): any {
    if (this.size === 1) {
      return this.last;
    }

    const last = this.last;
    // unreachable case
    if (last.prev === null) {
      // the last one, do nothing
      return last;
    }

    this.last = last.prev;
    this.last.next = null;
    this.size--;
    return last;
  }

  stackPush(data: any) {
    const node = new HistoryStackNode(this.last, null, data);
    if (this.size === this.limit) {
      this.first = this.first.next === null ? node : this.first.next;
      this.last.next = node;
      this.last = node;
    } else {
      this.last.next = node;
      this.last = node;
      this.size++;
    }
  }
}

class HistoryStackNode {
  next: HistoryStackNode | null;
  prev: HistoryStackNode | null;
  data: any;

  constructor(
    prev: HistoryStackNode | null,
    next: HistoryStackNode | null,
    data: any
  ) {
    this.data = data;
    this.next = next;
    this.prev = prev;
  }
}
