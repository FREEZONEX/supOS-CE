import { forwardRef } from 'react';
import cx from 'classnames';
import type { TreeItemProps } from '../type.ts';
import styles from './TreeItem.module.scss';
import { Draggable } from '@carbon/icons-react';
import { useThemeStore } from '@/stores/theme-store.ts';

const colorObj: any = {
  blue: {
    light: '#E8F1FF',
    dark: '#061833',
  },
  chartreuse: {
    light: '#F0FBD2',
    dark: '#242F06',
  },
};

export const TreeItem = forwardRef<HTMLDivElement, TreeItemProps>(
  (
    {
      wrapperRef,
      style,
      handleProps,
      depth,
      indentationWidth,
      label,
      leftExtra,
      rightExtra,
      wrapperStyle,
      fixed,
      disableSelection,
      disableInteraction,
      ghost,
      clone,
      indicator,
      allowDrop,
      node,
      onSelect,
      selected,
      disabledSelect,
      ...restProps
    },
    ref
  ) => {
    const { selectBgColor } = useThemeStore((state) => ({
      selectBgColor: colorObj?.[state.primaryColor]?.[state.theme],
    }));

    const className = cx(
      styles.TreeItemWrapper,
      clone && styles.clone,
      ghost && styles.ghost,
      fixed && styles.fixed,
      disabledSelect && styles.disabledSelect,
      {
        [styles.notAllow]: allowDrop === false,
      },
      selected && styles.selected,
      indicator && styles.indicator,
      disableSelection && styles.disableSelection,
      disableInteraction && styles.disableInteraction
    );
    return (
      <div
        ref={wrapperRef}
        {...restProps}
        className={className}
        style={{
          '--supos-select-bg-color': selectBgColor,
          '--spacing': `${indentationWidth * depth}px`,
          ...wrapperStyle,
        }}
        {...handleProps}
        onClick={() => {
          if (selected) {
            onSelect?.();
          } else {
            onSelect?.(node?.id, node);
          }
        }}
      >
        <div ref={ref} style={style} className={styles.TreeItem}>
          {!clone && !fixed && <Draggable style={{ cursor: 'grab', flexShrink: 0 }} {...handleProps} />}
          {leftExtra && <div className={styles.leftExtra}>{leftExtra}</div>}
          <div className={styles.Text} title={typeof label === 'string' ? label : ''}>
            {label}
          </div>
          {rightExtra && <div className={styles.rightExtra}>{rightExtra}</div>}
        </div>
      </div>
    );
  }
);
