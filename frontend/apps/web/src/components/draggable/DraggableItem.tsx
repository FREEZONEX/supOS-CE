import { useDraggable } from '@dnd-kit/core';
import { Axis } from './index.tsx';
import { type CSSProperties, forwardRef, type ReactNode } from 'react';
import styles from './index.module.scss';
import classNames from 'classnames';
import IframeMask from '../iframe-mask';

interface DraggableItemProps {
  handle?: boolean;
  style?: CSSProperties;
  buttonStyle?: CSSProperties;
  axis?: Axis;
  top?: number;
  left?: number;
  children?: ReactNode;
  draggingId?: string | number;
  disableDrag?: boolean; // 是否禁用拖拽
  onMouseEnter?: (isDragging: boolean) => void;
  onMouseLeave?: (isDragging: boolean) => void;
}

const DraggableItem = forwardRef<HTMLDivElement, DraggableItemProps>(
  ({ style, top, left, children, draggingId, disableDrag, onMouseEnter, onMouseLeave }, ref) => {
    const { attributes, isDragging, listeners, setNodeRef, transform } = useDraggable({
      id: draggingId || 'draggable',
      disabled: disableDrag || false,
    });
    const handleRef = (node: HTMLDivElement | null) => {
      setNodeRef(node);
      if (typeof ref === 'function') {
        ref(node);
      } else if (ref) {
        ref.current = node;
      }
    };
    return (
      <>
        <div
          ref={handleRef}
          {...attributes}
          {...listeners}
          onMouseEnter={() => onMouseEnter?.(isDragging)}
          onMouseLeave={() => onMouseLeave?.(isDragging)}
          className={classNames(styles['draggable-item'])}
          style={
            {
              ...style,
              top,
              left,
              position: 'fixed',
              '--translate-x': `${transform?.x ?? 0}px`,
              '--translate-y': `${transform?.y ?? 0}px`,
              cursor: isDragging ? 'move' : 'pointer',
              zIndex: 10,
            } as CSSProperties
          }
        >
          {children}
        </div>
        <IframeMask style={{ display: isDragging ? 'block' : 'none' }} />
      </>
    );
  }
);

export default DraggableItem;
