import type { FlattenedItem, TreeDataProps, UniqueIdentifier } from './type.ts';
import { arrayMove } from '@dnd-kit/sortable';

export const iOS = /iPad|iPhone|iPod/.test(navigator.platform);

/**
 * 树 =》 list 方法
 */
function flatten(items: TreeDataProps[], parentId: UniqueIdentifier | null = null, depth = 0): FlattenedItem[] {
  return items.reduce<FlattenedItem[]>((acc, item, index) => {
    return [...acc, { ...item, parentId, depth, index }, ...flatten(item.children, item.id, depth + 1)];
  }, []);
}

/**
 * 树 =》 list
 */
export function flattenTree(items: TreeDataProps[]): FlattenedItem[] {
  return flatten(items);
}

/**
 * 移除 flattenItems里面 折叠的 children
 */
export function removeChildrenOf(items: FlattenedItem[], ids: UniqueIdentifier[]) {
  const excludeParentIds = [...ids];

  return items.filter((item) => {
    if (item.parentId && excludeParentIds.includes(item.parentId)) {
      if (item.children.length) {
        excludeParentIds.push(item.id);
      }
      return false;
    }

    return true;
  });
}

/**
 * 获取当前映射关系 =》 deepth
 * getDragDepth
 * getMaxDepth
 * getMinDepth
 * getProjection
 */

function getDragDepth(offset: number, indentationWidth: number) {
  return Math.round(offset / indentationWidth);
}

function getMaxDepth({ previousItem }: { previousItem: FlattenedItem }) {
  if (previousItem) {
    return previousItem.depth + 1;
  }

  return 0;
}

function getMinDepth({ nextItem }: { nextItem: FlattenedItem }) {
  if (nextItem) {
    return nextItem.depth;
  }

  return 0;
}
/**
 * 获取当前映射关系 =》 deepth
 */
export function getProjection(
  items: FlattenedItem[],
  activeId: UniqueIdentifier,
  overId: UniqueIdentifier,
  dragOffset: number,
  indentationWidth: number
) {
  const overItemIndex = items.findIndex(({ id }) => id === overId);
  const activeItemIndex = items.findIndex(({ id }) => id === activeId);
  const activeItem = items[activeItemIndex];
  const newItems = arrayMove(items, activeItemIndex, overItemIndex);
  const previousItem = newItems[overItemIndex - 1];
  const nextItem = newItems[overItemIndex + 1];
  const dragDepth = getDragDepth(dragOffset, indentationWidth);
  const projectedDepth = activeItem.depth + dragDepth;
  const maxDepth = getMaxDepth({
    previousItem,
  });
  const minDepth = getMinDepth({ nextItem });
  // console.log(
  //   overItemIndex,
  //   activeItemIndex,
  //   activeItem,
  //   newItems,
  //   previousItem,
  //   nextItem,
  //   dragDepth,
  //   projectedDepth,
  //   maxDepth,
  //   minDepth
  // );
  let depth = projectedDepth;

  if (projectedDepth >= maxDepth) {
    depth = maxDepth;
  } else if (projectedDepth < minDepth) {
    depth = minDepth;
  }
  const parentId = getParentId();
  const parent = items?.find((f) => f.id === parentId);

  return { depth, parentId, parent: parent, minDepth, maxDepth, drag: activeItem };
  function getParentId() {
    if (depth === 0 || !previousItem) {
      return null;
    }

    if (depth === previousItem.depth) {
      return previousItem.parentId;
    }

    if (depth > previousItem.depth) {
      return previousItem.id;
    }

    const newParent = newItems
      .slice(0, overItemIndex)
      .reverse()
      .find((item) => item.depth === depth)?.parentId;

    return newParent ?? null;
  }
}

export function findItem(items: TreeDataProps[], itemId: UniqueIdentifier) {
  return items.find(({ id }) => id === itemId);
}

export function buildTree(flattenedItems: FlattenedItem[]): TreeDataProps[] {
  const root: TreeDataProps = { id: 'root', children: [] };
  const nodes: Record<string, TreeDataProps> = { [root.id]: root };
  const items = flattenedItems.map((item) => ({ ...item, children: [] }));

  for (const item of items) {
    const { id, children } = item;
    const parentId = item.parentId ?? root.id;
    const parent = nodes[parentId] ?? findItem(items, parentId);

    nodes[id] = { id, children };
    parent.children.push(item);
  }

  return root.children;
}
