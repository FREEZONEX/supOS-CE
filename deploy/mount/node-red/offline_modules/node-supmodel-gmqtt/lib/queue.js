class Queue {
    constructor() {
      this.items = {};
      this.headIndex = 0; // 队列头部指针
      this.tailIndex = 0; // 队列尾部指针
    }
  
    // 入队到队尾
    offer(item) {
      this.items[this.tailIndex] = item;
      this.tailIndex++;
      if (this.tailIndex > 1000000) {
        console.log(`<queue>Warn: 队列长度超出100万, 需要对队列索引进行重置, 当前最后一个对象值为：${item}</queue>`);
        this.reset();
      }
    }
  
    // 出队
    poll() {
      if (this.isEmpty()) {
        this.reset();
        return null;
      }
      const item = this.items[this.headIndex];
      delete this.items[this.headIndex];
      this.headIndex++;
      
      return item;
    }
  
    isEmpty() {
      return this.tailIndex === this.headIndex;
    }
  
    size() {
      return this.tailIndex - this.headIndex;
    }

    reset() {
        if (this.headIndex > 0) {
            this.headIndex = 0;
        }
        if (this.tailIndex > 0) {
            this.tailIndex = 0;
        }
    }
}

function newQueue() {
    return new Queue();
}

module.exports = { newQueue }; 