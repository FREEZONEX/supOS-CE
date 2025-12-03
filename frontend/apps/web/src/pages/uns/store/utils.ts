import type { Key } from 'react';
import { cloneDeep } from 'lodash-es';
import type { UnsTreeNode } from '@/pages/uns/types';

// 根据当前节点id获取链路
export const getParentNodes = (tree: UnsTreeNode[], targetId: Key): UnsTreeNode[] => {
  const result: UnsTreeNode[] = [];

  // 定义递归查找函数
  const findNode = (nodes: UnsTreeNode[], targetId: Key): boolean => {
    for (const node of nodes) {
      // 如果当前节点匹配目标节点
      if (node.id === targetId) {
        result.push(node);
        return true;
      }

      // 如果有子节点，递归查找
      if (node.children && node.children.length > 0) {
        if (findNode(node.children, targetId)) {
          result.push(node); // 找到目标节点后，将当前节点加入结果
          return true;
        }
      }
    }
    return false;
  };

  // 开始查找
  const found = findNode(tree, targetId);

  // 如果找到目标节点，返回结果数组（反转以确保顺序从父到子）
  return found ? result.reverse() : [];
};

// 数据重组出pathName
export const handlerTreeData = (treeData: UnsTreeNode[]) => {
  const modifyTreeData = (data: UnsTreeNode[], path?: string) => {
    data.forEach((node: UnsTreeNode) => {
      node.parentPath = path || '';
      // 真实展示的name
      const pathName = node?.path?.split('/').slice(-1)[0];
      node.pathName = pathName;
      node.title = pathName;
      node.key = node.id as string;
      // 不需要设置叶子节点,因为是整棵树，会根据children来判断
      // node.isLeaf = !node.hasChildren;
      if (node.children && node.children.length) {
        modifyTreeData(node.children, node.path);
      }
    });
  };
  modifyTreeData(treeData);
  return treeData;
};

// 递归获取所有后代 key
export function getDescendantKeys(nodes: UnsTreeNode[], parentKey: any) {
  let result: any = [];
  for (const node of nodes) {
    if (node.key === parentKey && node.children) {
      for (const child of node.children) {
        result.push(child.key);
        result = result.concat(getDescendantKeys(nodes, child.key));
      }
    } else if (node.children) {
      result = result.concat(getDescendantKeys(node.children, parentKey));
    }
  }
  return result;
}

// 递归清空子chidren
export function setTreeDescendantChildren(draft: UnsTreeNode[], nodeKey: string) {
  // 递归查找并更新节点
  const updateNode = (nodes: any) => {
    for (const node of nodes) {
      if (node.key === nodeKey && node.children?.length > 0) {
        node.children = [];
        return true;
      } else if (node.children) {
        if (updateNode(node.children)) {
          return true;
        }
      }
    }
    return false;
  };
  updateNode(draft);
}

// 递归查找信息
export function findNodeInfoById(draft: any, nodeKey: string) {
  const findNodeInTree = (nodes: UnsTreeNode[]): UnsTreeNode | undefined => {
    for (const node of nodes) {
      if (node.key === nodeKey || node.id === nodeKey) {
        return node;
      }

      if (node.children && node.children.length > 0) {
        const foundInChildren = findNodeInTree(node.children);
        if (foundInChildren) {
          return foundInChildren;
        }
      }
    }

    return undefined;
  };

  return findNodeInTree(draft);
}

// 追加分页数据方法
export const appendTreeData = (
  list: UnsTreeNode[],
  key: Key,
  childrenToAppend: UnsTreeNode[],
  type?: string,
  nodeDetail?: UnsTreeNode,
  cb?: () => void
): UnsTreeNode[] => {
  // 使用深度优先搜索找到目标节点并直接修改
  const findAndAppend = (nodes: UnsTreeNode[]): boolean => {
    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i];
      if (node.key === key) {
        // 移出更多，然后加上新值
        const existingChildren = node.children?.filter((child) => !(child as any).isLoadMoreNode) || [];
        if (childrenToAppend?.[0] && !childrenToAppend?.[0]?.preId) {
          childrenToAppend[0].preId = existingChildren?.[existingChildren.length - 1]?.id;
        }
        node.children = [...existingChildren, ...childrenToAppend];
        // 特殊处理
        if (type === 'addFile') {
          node.hasChildren = true;
          node.isLeaf = false;
          // 递归数量
          const updateAncestors = (key: string) => {
            let parent = findNodeInfoById(list, key);
            while (parent) {
              parent.countChildren = (parent.countChildren ?? 0) + 1;
              parent = findNodeInfoById(list, parent.parentId as string);
            }
          };
          updateAncestors(key as string);
        } else if (type === 'addFolder') {
          node.hasChildren = true;
          node.isLeaf = false;
          const updateAncestors = (key: string) => {
            let parent = findNodeInfoById(list, key);
            while (parent) {
              parent.countChildren =
                (parent.countChildren ?? 0) + (nodeDetail?.countChildren ? nodeDetail?.countChildren : 0);
              parent = findNodeInfoById(list, parent.parentId as string);
            }
          };
          updateAncestors(key as string);
        } else if (type === 'deleteFile' || type === 'deleteFolder') {
          const countChild = type === 'deleteFolder' ? nodeDetail?.countChildren || 0 : 1;
          // 递归处理数量减掉
          const updateAncestors = (key: string) => {
            let parent = findNodeInfoById(list, key);
            while (parent) {
              parent.countChildren = (parent.countChildren ?? 0) - countChild || 0;
              parent = findNodeInfoById(list, parent.parentId as string);
            }
          };
          updateAncestors(key as string);
        }
        if (node.children?.length === 0) {
          node.hasChildren = false;
          node.isLeaf = true;
          cb?.();
        }
        return true; // 找到并修改成功
      }

      if (node.children && node.children.length > 0) {
        if (findAndAppend(node.children)) {
          return true; // 在子节点中找到并修改成功
        }
      }
    }
    return false; // 未找到目标节点
  };

  // 判断是否在immer环境中
  if (Object.isFrozen(list)) {
    // 非immer环境，创建副本
    const result = cloneDeep(list);
    findAndAppend(result);
    return result;
  } else {
    // immer环境，直接修改
    findAndAppend(list);
    return list;
  }
};

/**
 * 将API返回的数据转换为树节点格式
 * @param data API返回的数据
 * @param parentPath 父节点路径
 * @param preId 上个节点id
 * @returns 转换后的树节点数据
 */
export const formatNodeData = (data: any[], parentPath: string = '', preId?: any): UnsTreeNode[] => {
  return data.map((item: any, index: number) => ({
    ...item,
    title: item.pathName,
    parentPath: parentPath,
    key: item.id,
    isLeaf: !item.hasChildren,
    preId: index === 0 && preId ? preId : data?.[index - 1]?.id,
    nextId: data?.[index + 1]?.id,
  }));
};

export const formatNodeDataForTemplate = (data: any[], preId?: any): UnsTreeNode[] => {
  return data.map((item: any, index: number) => ({
    ...item,
    pathType: 1,
    value: 0,
    title: item.name,
    isLeaf: true,
    key: item.id,
    preId: index === 0 && preId ? preId : data?.[index - 1]?.id,
    nextId: data?.[index + 1]?.id,
  }));
};

/**
 * 创建"加载更多"节点
 * @param parentId 父节点ID
 * @param parentInfo 父节点信息
 * @returns 加载更多节点
 */
export const createLoadMoreNode = (
  parentId: string | number,
  currentPage: number,
  parentInfo?: UnsTreeNode,
  loadMoreText?: string
): UnsTreeNode => {
  return {
    parentInfo,
    title: loadMoreText || '加载更多...',
    key: `${parentId}-loadmore`,
    parentKey: parentId,
    isLeaf: true,
    isLoadMoreNode: true,
    currentPage,
  };
};

/**
 * 检查是否有更多数据
 * @param total 总数
 * @param pageNo 当前页码
 * @param pageSize 每页大小
 * @returns 是否有更多数据
 */
export const hasMoreData = (restResponse: { total: number; pageNo: number; pageSize: number }): boolean => {
  const { total, pageNo, pageSize } = restResponse;
  return total > pageNo * pageSize;
};

export const uniqueArr = (arr: any[]) => {
  return arr.reduce((acc, current) => {
    if (!acc.find((item: any) => item.id === current.id)) {
      acc.push(current);
    }
    return acc;
  }, []);
};

// 平铺树
export const flattenTreeData = (nodes: UnsTreeNode[]): UnsTreeNode[] => {
  const result: UnsTreeNode[] = [];

  const traverse = (nodeList: UnsTreeNode[]) => {
    nodeList.forEach((node) => {
      result.push(node);
      if (node.children && node.children.length > 0) {
        traverse(node.children);
      }
    });
  };

  traverse(nodes);
  return result;
};

// SHOW_CHILD 只显示子; SHOW_ALL 显示所有; SHOW_PARENT 显示父
export const processedCheckedKeys = ({
  checkedKeys,
  strategy = 'SHOW_PARENT',
  treeData = [],
}: {
  // 字符串 的 id集合
  checkedKeys: Key[];
  strategy?: 'SHOW_CHILD' | 'SHOW_ALL' | 'SHOW_PARENT';
  // 树数据，非平铺
  treeData?: UnsTreeNode[];
}) => {
  if (!treeData || treeData.length === 0 || !checkedKeys || checkedKeys.length === 0) {
    return [];
  }

  // 平铺树数据，方便查找
  const flatTreeData = flattenTreeData(treeData);

  // 创建节点映射表（使用id作为key）
  const nodeMap = new Map<Key, UnsTreeNode>();
  flatTreeData.forEach((node) => {
    nodeMap.set(node.id!, node);
  });

  const checkedKeySet = new Set(checkedKeys);
  let resultKeys: Key[] = [];

  if (strategy === 'SHOW_ALL') {
    // SHOW_ALL: 返回所有选中的节点
    resultKeys = Array.from(checkedKeySet);
  } else if (strategy === 'SHOW_CHILD') {
    // SHOW_CHILD: 只返回选中的子节点 (其父节点未被选中的已选节点)
    resultKeys = checkedKeys.filter((key) => {
      const node = nodeMap.get(key);
      // 如果节点没有父节点，或者父节点未被选中，则保留该节点
      if (!node || !node.parentId) {
        return true;
      }
      return !checkedKeySet.has(node.parentId);
    });
  } else if (strategy === 'SHOW_PARENT') {
    // SHOW_PARENT: antd TreeSelect的SHOW_PARENT策略逻辑：
    // 当一个父节点的所有直接子节点都被选中时，只显示该父节点，其子节点将被隐藏。
    // 这个逻辑会从下至上应用，如果子节点成组后其父节点又满足条件，会继续向上合并。
    const resultKeySet = new Set(checkedKeySet);

    // 从扁平数据中筛选出所有父节点
    const parentNodes = flatTreeData.filter((node) => node.children && node.children.length > 0);

    // 逆序遍历父节点，实现从下至上的检查，确保子节点先被处理
    for (const parentNode of parentNodes.reverse()) {
      // 确保 parentNode.children 存在
      if (!parentNode.children) continue;

      // 检查该父节点的所有直接子节点是否都在 resultKeySet 中
      const allChildrenInSet = parentNode.children.every((child) => resultKeySet.has(child.id!));

      if (allChildrenInSet) {
        // 如果所有子节点都在，则从结果集中移除它们
        parentNode.children.forEach((child) => {
          resultKeySet.delete(child.id!);
        });
        // 并将父节点添加到结果集中
        resultKeySet.add(parentNode.id!);
      }
    }
    resultKeys = Array.from(resultKeySet);
  } else {
    // 默认返回原始checkedKeys
    resultKeys = Array.from(checkedKeySet);
  }

  // 根据最终的keys返回对应的节点信息
  return resultKeys.map((id) => nodeMap.get(id)).filter((node): node is UnsTreeNode => node !== undefined);
};
