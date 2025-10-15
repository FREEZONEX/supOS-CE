import type { MenuProps } from './types.ts';

export function listToTree(items: any[]): MenuProps[] {
  // 创建 id 到节点的映射
  const idMap = new Map<string, MenuProps>();
  items.forEach((item) => {
    if (item.type !== 5) {
      // 资源类型1-目录 2-菜单 3-按钮（操作权限）4-Tab  5- 子菜单  子菜单不显示
      idMap.set(item.id, {
        ...item,
        label: item.showName,
        children: [],
        operationChildren: [],
        tabChildren: [],
      });
    }
  });

  // 构建树结构并实时排序
  const tree: MenuProps[] = [];

  // 排序函数
  const sortChildren = (nodes: MenuProps[]) => {
    return nodes.sort((a, b) => (a.sort || 0) - (b.sort || 0));
  };

  idMap.forEach((node) => {
    if (!node.parentId) {
      tree.push(node);
    } else {
      const parent = idMap.get(node.parentId);
      if (parent) {
        // 根据资源类型放入不同的children属性并立即排序
        // type 1-目录 2-菜单 3-按钮（操作权限）4-Tab
        if (node.type === 1 || node.type === 2) {
          parent.children?.push(node);
          parent.children = sortChildren(parent.children || []);
        } else if (node.type === 3) {
          parent.operationChildren?.push(node);
          parent.operationChildren = sortChildren(parent.operationChildren || []);
        } else if (node.type === 4) {
          // 为tab节点创建专门的容器节点
          let tabContainer: any = parent.children?.find((child) => child.showTabs);
          if (!tabContainer) {
            tabContainer = {
              id: `tab_container_${parent.id}`,
              showTabs: true,
              tabChildren: [],
              children: [],
              parentId: parent.id,
              type: 2,
            };
            parent.children?.push(tabContainer);
            parent.children = sortChildren(parent.children || []);
          }
          tabContainer.tabChildren?.push(node);
          tabContainer.tabChildren = sortChildren(tabContainer.tabChildren || []);
        }
      }
    }
  });

  // 对根节点进行排序
  return sortChildren(tree);
}
