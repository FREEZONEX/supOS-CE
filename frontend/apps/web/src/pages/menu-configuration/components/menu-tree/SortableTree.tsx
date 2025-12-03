import { type FC, useMemo, useState } from 'react';
import { createPortal } from 'react-dom';
import {
  closestCenter,
  defaultDropAnimation,
  DndContext,
  type DragEndEvent,
  type DragMoveEvent,
  type DragOverEvent,
  DragOverlay,
  type DragStartEvent,
  type DropAnimation,
  KeyboardSensor,
  MeasuringStrategy,
  type Modifier,
  MouseSensor,
  PointerSensor,
  TouchSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import { arrayMove, SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import type { FlattenedItem, SortableTreeProps, UniqueIdentifier } from './type.ts';
import { Empty, Flex, Spin } from 'antd';
import { removeChildrenOf, flattenTree, getProjection, buildTree } from './utils.ts';
import { SortableTreeItem } from './TreeItem';
import usePropsValue from '@/hooks/usePropsValue.ts';
import styles from './index.module.scss';

const measuring = {
  droppable: {
    strategy: MeasuringStrategy.Always,
  },
};

const dropAnimationConfig: DropAnimation = {
  keyframes({ transform }) {
    return [
      { opacity: 1, transform: CSS.Transform.toString(transform.initial) },
      {
        opacity: 0,
        transform: CSS.Transform.toString({
          ...transform.final,
          x: transform.final.x + 5,
          y: transform.final.y + 5,
        }),
      },
    ];
  },
  easing: 'ease-out',
  sideEffects({ active }) {
    active.node.animate([{ opacity: 0 }, { opacity: 1 }], {
      duration: defaultDropAnimation.duration,
      easing: defaultDropAnimation.easing,
    });
  },
};

export const SortableTree: FC<SortableTreeProps> = ({
  indicator = true,
  renderLabel,
  indentationWidth = 50,
  rightExtra,
  leftExtra,
  allowDrop,
  selectedKey,
  onSelect,
  style,
  treeData = [],
  loading,
  disabledSelected,
  onHandleDragEnd,
  // ...restProps
}) => {
  const [data, setData] = usePropsValue({ value: treeData });
  const [selectValue, setSelect] = usePropsValue({
    value: selectedKey,
  });

  const [activeId, setActiveId] = useState<UniqueIdentifier | null>(null);
  const [overId, setOverId] = useState<UniqueIdentifier | null>(null);
  const [offsetLeft, setOffsetLeft] = useState(0);

  const flattenedItems = useMemo(() => {
    const flattenedTree = flattenTree(data);
    const collapsedItems = flattenedTree.reduce<any[]>(
      (acc, { children, collapsed, id }) => (collapsed && children.length ? [...acc, id] : acc),
      []
    );

    return removeChildrenOf(flattenedTree, activeId != null ? [activeId, ...collapsedItems] : collapsedItems);
  }, [activeId, data]);

  const activeItem = activeId ? flattenedItems.find(({ id }) => id === activeId) : null;

  const sensors = useSensors(
    useSensor(MouseSensor),
    useSensor(TouchSensor),
    useSensor(KeyboardSensor),
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    })
  );
  const resetState = () => {
    setOverId(null);
    setActiveId(null);
    setOffsetLeft(0);
    document.body.style.setProperty('cursor', '');
  };

  const handleDragStart = ({ active: { id: activeId } }: DragStartEvent) => {
    setActiveId(activeId);
    setOverId(activeId);
    document.body.style.setProperty('cursor', 'grabbing');
  };

  const handleDragMove = ({ delta }: DragMoveEvent) => {
    setOffsetLeft(delta.x);
  };
  const handleDragOver = ({ over }: DragOverEvent) => {
    setOverId(over?.id ?? null);
  };
  const handleDragEnd = ({ active, over }: DragEndEvent) => {
    resetState();

    if (projected && over) {
      if (!allowDrop?.({ drop: projected?.parent, drag: projected?.drag })) return;
      const { depth, parentId } = projected;
      const clonedItems: FlattenedItem[] = JSON.parse(JSON.stringify(flattenTree(data)));
      const overIndex = clonedItems.findIndex(({ id }) => id === over.id);
      const activeIndex = clonedItems.findIndex(({ id }) => id === active.id);
      const activeTreeItem = clonedItems[activeIndex];

      clonedItems[activeIndex] = { ...activeTreeItem, depth, parentId };

      const sortedItems = arrayMove(clonedItems, activeIndex, overIndex);
      const newItems = buildTree(sortedItems);
      if (onHandleDragEnd) {
        onHandleDragEnd?.(flattenTree(newItems), newItems);
      } else {
        setData(newItems);
      }
    }
  };

  const handleDragCancel = () => {
    resetState();
  };

  const projected =
    activeId && overId ? getProjection(flattenedItems, activeId, overId, offsetLeft, indentationWidth) : null;

  const _allowDrop = useMemo(() => {
    const __allowDrop = projected && allowDrop ? allowDrop({ drop: projected?.parent, drag: projected?.drag }) : true;
    if (activeId) {
      document.body.style.setProperty('cursor', __allowDrop ? 'grabbing' : 'no-drop');
    }
    return __allowDrop;
  }, [projected, allowDrop, activeId]);

  const sortedIds = useMemo(() => flattenedItems.map(({ id }) => id), [flattenedItems]);

  return (
    <Spin spinning={!!loading} wrapperClassName={styles['sortable-tree-loading']}>
      <div style={style} className={styles['sortable-tree-wrap']}>
        {!treeData?.length ? (
          <Flex justify="center" align="center" style={{ height: '100%' }}>
            {<Empty />}
          </Flex>
        ) : (
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            measuring={measuring}
            onDragStart={handleDragStart}
            onDragMove={handleDragMove}
            onDragOver={handleDragOver}
            onDragEnd={handleDragEnd}
            onDragCancel={handleDragCancel}
          >
            <SortableContext items={sortedIds} strategy={verticalListSortingStrategy}>
              <Flex gap={8} vertical>
                {flattenedItems?.map((item) => {
                  const disabledSelect = disabledSelected?.(item) || false;
                  if (item.showTabs) {
                    const tabChildren = item.tabChildren?.filter((f) => f.id !== '102');
                    return (
                      <Flex wrap style={{ paddingLeft: item.depth * indentationWidth }} gap={8} key={item.id}>
                        {tabChildren?.map((tab) => {
                          const length = tabChildren && tabChildren?.length < 3 ? tabChildren?.length : 3;
                          const disabledSelect = disabledSelected?.(tab) || false;
                          return (
                            <SortableTreeItem
                              node={tab}
                              onSelect={(v, node) => {
                                if (disabledSelect) return;
                                setSelect(v);
                                onSelect?.(v, node);
                              }}
                              disabledSelect={disabledSelect}
                              key={tab.id}
                              indicator={indicator}
                              fixed={tab.fixed}
                              wrapperStyle={{
                                width: `calc((100% - ${(length - 1) * 8}px)/ ${length})`,
                              }}
                              id={tab.id}
                              depth={0}
                              label={renderLabel ? renderLabel(tab) : tab.label}
                              rightExtra={typeof rightExtra === 'function' ? rightExtra?.(tab) : rightExtra}
                              leftExtra={typeof leftExtra === 'function' ? leftExtra?.(tab) : leftExtra}
                              indentationWidth={0}
                              selected={tab.id === selectValue}
                            />
                          );
                        })}
                      </Flex>
                    );
                  }
                  return (
                    <SortableTreeItem
                      node={item}
                      onSelect={(v, node) => {
                        if (disabledSelect) return;
                        setSelect(v);
                        onSelect?.(v, node);
                      }}
                      disabledSelect={disabledSelect}
                      key={item.id}
                      indicator={indicator}
                      fixed={item.fixed}
                      id={item.id}
                      depth={item.id === activeId && projected ? projected.depth : item.depth}
                      label={renderLabel ? renderLabel(item) : item.label}
                      rightExtra={typeof rightExtra === 'function' ? rightExtra?.(item) : rightExtra}
                      leftExtra={typeof leftExtra === 'function' ? leftExtra?.(item) : leftExtra}
                      indentationWidth={indentationWidth}
                      allowDrop={activeId === item.id ? _allowDrop : true}
                      selected={item.id === selectValue}
                    />
                  );
                })}
              </Flex>
              {createPortal(
                <DragOverlay dropAnimation={dropAnimationConfig} modifiers={indicator ? [adjustTranslate] : undefined}>
                  {activeId && activeItem ? (
                    <SortableTreeItem
                      id={activeId}
                      depth={activeItem.depth}
                      clone
                      // childCount={getChildCount(items, activeId) + 1}
                      label={renderLabel ? renderLabel(activeItem) : activeItem.label}
                      indentationWidth={indentationWidth}
                      indicator={indicator}
                      allowDrop={_allowDrop}
                    />
                  ) : null}
                </DragOverlay>,
                document.body
              )}
            </SortableContext>
          </DndContext>
        )}
      </div>
    </Spin>
  );
};

const adjustTranslate: Modifier = ({ transform }) => {
  return {
    ...transform,
    y: transform.y - 25,
  };
};
