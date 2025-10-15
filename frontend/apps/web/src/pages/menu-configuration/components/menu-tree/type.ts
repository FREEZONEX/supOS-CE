import type { CSSProperties, HTMLAttributes, ReactNode } from 'react';

export type UniqueIdentifier = string | number;

export interface SortableTreeProps {
  indicator?: boolean;
  renderLabel?: (node: TreeDataProps) => ReactNode;
  indentationWidth?: number;
  rightExtra?: ((node: TreeDataProps) => ReactNode) | ReactNode;
  leftExtra?: ((node: TreeDataProps) => ReactNode) | ReactNode;
  allowDrop?: (info: { drop?: TreeDataProps; drag: TreeDataProps }) => boolean;
  selectedKey?: UniqueIdentifier | null;
  onSelect?: (key?: UniqueIdentifier, node?: TreeDataProps) => void;
  disabledSelected?: (node: TreeDataProps) => boolean;
  style?: CSSProperties;
  loading?: boolean;
  treeData?: TreeDataProps[];
  onHandleDragEnd?: (list: TreeDataProps[], tree: TreeDataProps[]) => void;
}

export interface TreeItemProps extends Omit<HTMLAttributes<HTMLDivElement>, 'id' | 'title' | 'onSelect'> {
  wrapperRef?(node: HTMLDivElement): void;
  // 深度
  depth: number;
  // 被拖拽
  ghost?: boolean;
  // 禁止选中
  disableSelection?: boolean;
  // 禁止交互
  disableInteraction?: boolean;
  // 拖拽手柄交互
  handleProps?: any;
  // 缩进宽度
  indentationWidth: number;
  // 展示内容
  label?: ReactNode;
  // 右侧展示内容
  rightExtra?: ReactNode;
  leftExtra?: ReactNode;
  wrapperStyle?: CSSProperties;
  // 是否拖拽
  fixed?: boolean;
  // 是否使用指示器
  indicator?: boolean;
  clone?: boolean;
  // 是否运行放置
  allowDrop?: boolean;
  // 是否选中
  selected?: boolean;
  // 点击事件
  onSelect?: SortableTreeProps['onSelect'];
  node?: TreeDataProps;
  // 禁止点击
  disabledSelect?: boolean;
}

export interface SortableTreeItemProps extends TreeItemProps {
  id: UniqueIdentifier;
}

export interface TreeDataProps {
  label?: string;
  showTabs?: boolean;
  id: UniqueIdentifier;
  parentId?: UniqueIdentifier | null;
  collapsed?: boolean;
  // 是否拖拽
  fixed?: boolean;
  isLeaf?: boolean;
  children: TreeDataProps[];
  tabChildren?: TreeDataProps[];
}

export interface FlattenedItem extends TreeDataProps {
  depth: number;
  index: number;
}
