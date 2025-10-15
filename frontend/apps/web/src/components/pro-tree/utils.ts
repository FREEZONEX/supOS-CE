import type { DataNodeProps, FlatNode } from '@/components/pro-tree';

function flatten(items: DataNodeProps[], parentId: string | number | null = null, depth: number = 0): FlatNode[] {
  return items.reduce<FlatNode[]>((acc, item, index) => {
    // 浅拷贝节点避免污染原数据
    const flatNode: FlatNode = {
      ...item,
      _parentId: parentId,
      _depth: depth,
      _index: index,
    };

    const children = item.children || [];
    return acc.concat(flatNode, flatten(children, item.id, depth + 1));
  }, []);
}

export function flattenTree(items: DataNodeProps[]): FlatNode[] {
  return flatten(items);
}
