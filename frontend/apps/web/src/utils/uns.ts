import { v4 as uuidv4 } from 'uuid';
import { pinyin } from 'pinyin-pro';
import type { UnsTreeNode } from '@/pages/uns/types';

export const noDuplicates = (arr: any) => {
  return new Set(arr).size === arr.length;
};

export const isPathOdd = (path: string) => {
  const trimmedPath = path.replace(/^\/|\/$/g, '');
  const segments = trimmedPath.split('/');
  return segments.length % 2 !== 0;
};

export const removeLastPathSegment = (path: string): string => {
  // 移除路径两端的斜杠并查找最后一个斜杠的位置
  const trimmedPath = path.replace(/^\/+|\/+$/g, '');
  const lastSlashIndex = trimmedPath.lastIndexOf('/');

  // 如果没有斜杠或路径为空，则返回空字符串
  if (lastSlashIndex <= 0) return '';

  // 返回去除最后一个路径片段后的路径
  return `${trimmedPath.substring(0, lastSlashIndex)}/`;
};

export const generatePaths = (input: string): string[] => {
  if (typeof input !== 'string') return [];
  const startsWithSlash = input.startsWith('/'); // 检查是否以斜杠开始
  const parts = input.split('/').filter((part) => part.length); // 分割并过滤掉空字符串
  let currentPath = startsWithSlash ? '/' : ''; // 根据原输入决定currentPath初始值

  return parts.slice(0, -1).map((part) => {
    // 忽略最后一个元素，并构造路径
    currentPath += `${part}/`;
    return currentPath;
  });
};

export const loadIbmFont = () => {
  //为了提前触发字体文件的加载，防止闪屏
  const div = document.createElement('div');
  div.style.fontWeight = 'bold';
  div.style.visibility = 'hidden';
  div.style.position = 'absolute';
  div.style.zIndex = '-9999';
  div.style.top = '-9999px';
  div.style.left = '-9999px';
  div.innerHTML = '中文';
  document.body.appendChild(div);
  setTimeout(() => {
    document.body.removeChild(div);
  }, 3000);
};

//递归查询模糊命中节点的所有父级节点
export const findParentIds = (searchString: string, tree: UnsTreeNode[]): string[] => {
  // 使用 const 定义结果数组
  const result: string[] = [];

  // 使用 const 定义辅助函数：递归搜索树
  const search = (node: UnsTreeNode, parentIds: string[]): void => {
    if (node?.name?.toLowerCase()?.includes(searchString.toLowerCase())) {
      // 如果匹配，则将当前父级path加入结果
      result.push(...parentIds);
    }
    if (Array.isArray(node.children)) {
      // 对每个子节点，更新其父级path链并递归搜索
      const updatedParentIds = [...parentIds, node.id] as string[];
      node.children.forEach((child: UnsTreeNode) => search(child, updatedParentIds));
    }
  };

  // 使用 const 定义根节点遍历逻辑，从根节点开始搜索，初始时没有父级ID
  tree.forEach((rootNode: any) => search(rootNode, []));

  // 去重后返回父级path数组
  return [...new Set(result)];
};

//获取目标节点下所有带children的节点path
export const collectChildrenIds = (tree: UnsTreeNode[], targetId: string): string[] => {
  const result: string[] = [];

  const recurseForChildren = (node: UnsTreeNode): void => {
    // 如果当前节点有children并且children数组不为空，则收集其id
    if (node.children && node.children.length > 0) {
      result.push(node.id as string);
      // 对每个子节点递归调用recurseForChildren
      node.children.forEach((child) => recurseForChildren(child));
    }
  };

  const searchTree = (nodes: UnsTreeNode[]): boolean => {
    for (let i = 0; i < nodes.length; i++) {
      const node: UnsTreeNode = nodes[i];
      // 如果targetPath为空字符串，则直接开始收集所有带children的节点
      if (targetId === '') {
        recurseForChildren(node);
      } else if (node.id === targetId) {
        // 如果找到了目标节点，则从该节点开始递归收集其下所有带children的节点id
        recurseForChildren(node);
        return true;
      } else if (node.children) {
        // 继续在子节点中搜索
        if (searchTree(node.children)) return true;
      }
    }
    return false; // 如果没有找到目标节点且targetPath不是空字符串
  };

  // 开始搜索
  if (targetId === '') {
    tree.forEach((node) => recurseForChildren(node)); // 直接遍历树的所有节点
  } else {
    searchTree(tree);
  }

  return result;
};

type ResultItem = {
  key: string;
  functionType: string;
};

export const parseArrayToObjects = (arr: string[]): ResultItem[] => {
  return arr.map((item) => {
    const _item = item.slice(0, -1);
    return {
      key: _item.split('(')[1],
      functionType: _item.split('(')[0],
    };
  });
};

export const parseTime = (input: number | string): [number] | [number, string] => {
  let timeString: string;

  // 如果输入是数字类型，则直接转换为字符串
  if (typeof input === 'number') {
    timeString = input.toString();
  } else {
    timeString = input;
  }

  // 检查是否为纯数字
  if (/^\d+\.?\d*$/.test(timeString)) {
    return [parseFloat(timeString)];
  }

  // 使用正则表达式匹配数字和单位
  const match = timeString.match(/^(\d+\.?\d*)([a-zA-Z]+)$/);

  // 如果匹配失败，抛出错误或者根据需要处理
  if (!match) {
    throw new Error("输入格式不正确：必须是数字紧跟单位的形式，例如 '5seconds' 或者纯数字");
  }

  // 提取数字和单位
  const number = parseFloat(match[1]); // 将匹配到的第一部分(数字)转换为浮点数
  const unit = match[2]; // 单位

  return [number, unit];
};

interface ExpressionItem {
  id: string;
  field: string;
}

export const getExpression = (items: ExpressionItem[], inputStr: string, history?: boolean): string => {
  if (!inputStr?.trim()) return '';

  const patterns: [RegExp, string][] = [];
  let cleanedStr = inputStr.replace(/"/g, '');

  // 预编译正则表达式
  items.forEach((item, index) => {
    const placeholder = `\\$${item.id}\\.${item.field}#`;
    patterns.push([new RegExp(placeholder, 'g'), history ? `$${item.field}#` : `$a${index + 1}#`]);
  });

  // 批量执行替换
  patterns.forEach(([regex, replacement]) => {
    cleanedStr = cleanedStr.replace(regex, replacement);
  });

  return cleanedStr;
};

export const generateAlias = (name: string) => {
  return `_${pinyin(name || '', { toneType: 'none', v: true })
    ?.replace(/\s+/g, '')
    ?.replace(/-/g, '_')
    .slice(0, 38)}_${uuidv4().replace(/-/g, '').slice(0, 20)}`;
};

//根据dataPath获取目标树节点信息
export const getTargetNode = (treeData: UnsTreeNode[], targetId: string): UnsTreeNode | null => {
  for (const node of treeData) {
    if (node.id === targetId) {
      return { ...node };
    }
    // 如果存在子节点并且还没有找到目标节点，则递归搜索子节点
    if (node.children && node.children.length > 0) {
      const foundNode = getTargetNode(node.children, targetId);
      if (foundNode) {
        return { ...foundNode }; // 找到目标节点后立即返回
      }
    }
  }
  return null; // 如果遍历完整个树都没有找到，则返回null
};

export const SelectAllId = '_SELECT_ALL';
export const getParamsForArray = (list: { value: string; label: string }[] = [], keyOrFn?: any, options?: any) => {
  if (list.length === 0) return {};
  if (list.some((s) => s.value === SelectAllId)) {
    return { exportType: 'ALL' };
  }
  if (!keyOrFn) {
    return { ids: list?.map((i) => i.value) || [] };
  } else {
    const { groups, extract } = options;

    return list.reduce((acc: any, item: any) => {
      let key;
      if (typeof keyOrFn === 'function') {
        key = keyOrFn(item);
      } else {
        key = item[keyOrFn];
      }
      if (groups && groups[key] !== undefined) {
        key = groups[key];
      }
      if (!acc[key]) {
        acc[key] = [];
      }
      const valueToPush = extract ? item[extract] : item;
      acc[key].push(valueToPush);

      return acc;
    }, {});
  }
};
