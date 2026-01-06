import type { CSSProperties, FC, ReactNode } from 'react';
import { useSortable } from '@dnd-kit/sortable';
import type { ProTreeProps } from './types.ts';
import './SortableTreeNode.scss';

const SortableTreeNode: FC<{
  node: any;
  isActive?: boolean;
  children: ReactNode;
  isInset: boolean;
  draggable?: boolean;
  drapOverChanges: ProTreeProps['drapOverChanges'];
}> = ({ node, isActive, children, isInset, drapOverChanges, draggable }) => {
  const { attributes, listeners, setNodeRef, isDragging, over } = useSortable({
    id: node.id as string,
    data: {
      type: 'tree-node',
      node,
    },
    animateLayoutChanges: () => false, // 禁用布局变化动画
  });
  if (!draggable) {
    return children;
  }
  const isIndicator = !isActive && node.id === over?.id;
  const style: CSSProperties = {
    opacity: isDragging ? 0.5 : 1,
    cursor: 'grab',
  };
  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={
        isIndicator
          ? drapOverChanges
            ? drapOverChanges({
                node,
                isInset,
                classNames: {
                  in: 'sortableTreeNodeIn',
                  out: 'sortableTreeNodeOut',
                },
              })
            : isInset
              ? 'sortableTreeNodeIn'
              : 'sortableTreeNodeOut'
          : ''
      }
    >
      {children}
    </div>
  );
};

export default SortableTreeNode;
