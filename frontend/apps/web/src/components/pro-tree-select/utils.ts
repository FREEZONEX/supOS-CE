import type { Key } from 'react';
import type { TreeDataNodeProps, PageParamsType } from './interface.ts';
import type { TreeSelectProps } from 'antd';

export const RootId = '0';
export const SelectAllId = '_SELECT_ALL';

// 判断是否有更多数据
export const hasMoreData = (restResponse: PageParamsType): boolean => {
  const { total, pageNo, pageSize } = restResponse;
  return total > pageNo * pageSize;
};

// 先建更多节点
export const createLoadMoreNode = (
  parentKey: string,
  currentPage: number,
  loadMoreText?: string,
  fieldNames?: TreeSelectProps['fieldNames']
): TreeDataNodeProps => {
  return {
    title: loadMoreText || '加载更多...',
    key: `${parentKey}-loadmore`,
    value: `${parentKey}-loadmore`,
    [fieldNames?.value ? fieldNames.value : 'value']: `${parentKey}-loadmore`,
    [fieldNames?.label ? fieldNames.label : 'title']: loadMoreText || '加载更多...',
    parentKey: parentKey,
    isLeaf: true,
    isLoadMoreNode: true,
    currentPage,
    disabled: true,
  };
};

// 追加分页数据方法
export const appendTreeData = (list: TreeDataNodeProps[], key: Key, childrenToAppend: TreeDataNodeProps[]) => {
  // 使用深度优先搜索找到目标节点并直接修改
  const findAndAppend = (nodes: TreeDataNodeProps[]): boolean => {
    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i];
      if (node.key === key) {
        // 移出更多，然后加上新值
        const existingChildren = node.children?.filter((child) => !(child as any).isLoadMoreNode) || [];
        node.children = [...existingChildren, ...childrenToAppend];
        return true;
      }

      if (node.children && node.children.length > 0) {
        if (findAndAppend(node.children)) {
          return true;
        }
      }
    }
    return false; // 未找到目标节点
  };
  findAndAppend(list);
};

// 2个数组查找不同
export function findMissingValuesOptimized(arr1: Key[] = [], arr2: Key[] = []) {
  const set2 = new Set(arr2);
  return arr1.filter((item) => !set2.has(item));
}

// 将树平铺
export function findTreeNode(tree: any, id: Key, idKey = 'id', childrenKey = 'children') {
  // 处理数组形式的树
  if (Array.isArray(tree)) {
    for (const node of tree) {
      const result: any = findTreeNode(node, id, idKey, childrenKey);
      if (result) return result;
    }
    return null;
  }

  // 检查当前节点
  if (tree[idKey] === id) {
    return tree;
  }

  // 递归检查子节点
  if (tree[childrenKey] && Array.isArray(tree[childrenKey])) {
    for (const child of tree[childrenKey]) {
      const result: any = findTreeNode(child, id, idKey, childrenKey);
      if (result) return result;
    }
  }

  return null;
}
