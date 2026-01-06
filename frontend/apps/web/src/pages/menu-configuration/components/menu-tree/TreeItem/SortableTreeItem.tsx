import { type AnimateLayoutChanges, useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import type { CSSProperties, FC } from 'react';
import type { SortableTreeItemProps } from '../type.ts';
import { TreeItem } from './TreeItem';
import { iOS } from '../utils.ts';

const animateLayoutChanges: AnimateLayoutChanges = ({ isSorting, wasDragging }) => !(isSorting || wasDragging);

export const SortableTreeItem: FC<SortableTreeItemProps> = ({ id, depth, ...restProps }) => {
  const {
    attributes,
    isDragging,
    isSorting,
    listeners,
    setDraggableNodeRef,
    setDroppableNodeRef,
    transform,
    transition,
  } = useSortable({
    id,
    animateLayoutChanges,
    disabled: restProps.fixed,
  });
  const style: CSSProperties = {
    transform: CSS.Translate.toString(transform),
    transition,
  };

  return (
    <TreeItem
      ref={restProps.fixed ? undefined : setDraggableNodeRef}
      wrapperRef={restProps.fixed ? undefined : setDroppableNodeRef}
      style={style}
      depth={depth}
      ghost={isDragging}
      disableSelection={iOS}
      disableInteraction={isSorting}
      handleProps={{
        ...attributes,
        ...listeners,
      }}
      {...restProps}
    />
  );
};
