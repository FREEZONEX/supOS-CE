import {
  cloneElement,
  type CSSProperties,
  type Dispatch,
  type FC,
  forwardRef,
  type HTMLAttributes,
  type MouseEvent,
  type ReactElement,
  type SetStateAction,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
  type WheelEvent,
} from 'react';
import { Tag, type TagProps as AntTagProps } from 'antd';
import type { DragEndEvent } from '@dnd-kit/core';
import { DndContext, PointerSensor, useSensor } from '@dnd-kit/core';
import { arrayMove, horizontalListSortingStrategy, SortableContext, useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { restrictToHorizontalAxis } from '@dnd-kit/modifiers';
import './index.scss';
import type { KeepAliveTab } from '@/layout/useTabs.ts';
import { ChevronLeft, ChevronRight } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import { useDeepCompareEffect, useSize } from 'ahooks';
import ControlledDropdown, { type ControlledDropdownRef } from '@/components/controlled-dropdown';

interface DraggableTabPaneProps extends HTMLAttributes<HTMLDivElement> {
  'data-node-key': string;
}

interface TagProps extends Omit<AntTagProps, 'onClick' | 'onClose'> {
  key: string;
  onClick?: (key: string) => void;
  onClose?: (key: string) => void;
}
interface ComTagsProps {
  style?: CSSProperties;
  activeTag?: string;
  options?: TagProps[];
  onClose?: (key: string) => void;
  onRefresh?: (key: string) => void;
  onCloseOther?: (key: string) => void;
  setTabs?: Dispatch<SetStateAction<KeepAliveTab[]>>;
}

const DraggableTagNode = forwardRef<HTMLDivElement, DraggableTabPaneProps>((props: DraggableTabPaneProps, ref) => {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({
    id: props['data-node-key'],
  });

  const style: CSSProperties = {
    ...props.style,
    transform: CSS.Transform.toString(transform && { ...transform, scaleX: 1 }),
    transition,
  };

  // eslint-disable-next-line react-hooks/refs
  return cloneElement(props.children as ReactElement, {
    ref: (node: HTMLDivElement) => {
      setNodeRef(node);
      if (typeof ref === 'function') {
        ref(node);
      } else if (ref) {
        ref.current = node;
      }
    },
    style,
    ...attributes,
    ...listeners,
  });
});

const ComTags: FC<ComTagsProps> = ({ options = [], activeTag, onClose, onCloseOther, onRefresh, setTabs }) => {
  const tabsContainerRef = useRef<HTMLDivElement>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const wrapperSize = useSize(wrapperRef);
  const dropdownRef = useRef<ControlledDropdownRef>(null);

  const [showPrev, setShowPrev] = useState(false);
  const [showNext, setShowNext] = useState(false);
  const formatMessage = useTranslate();

  const handleWheelScroll = (e: WheelEvent) => {
    if (tabsContainerRef.current) {
      tabsContainerRef.current.scrollLeft += e.deltaY; // 横向滚动
    }
  };

  const handleContextMenu = (e: MouseEvent<HTMLSpanElement>, key: string) => {
    e.preventDefault(); // 阻止默认右键菜单
    const menuItems = [
      {
        label: formatMessage('common.refresh'),
        key: 'REFRESH',
        onClick: () => onRefresh?.(key),
      },
      options?.length > 1
        ? {
            label: formatMessage('common.close'),
            key: 'CLOSE',
            onClick: () => onClose?.(key),
          }
        : null,
      options?.length > 1
        ? {
            label: formatMessage('common.closeOther'),
            key: 'CLOSEOTHER',
            onClick: () => onCloseOther?.(key),
          }
        : null,
    ].filter((o) => o !== null);
    dropdownRef?.current?.showDropdown(e, menuItems);
  };

  const onDragEnd = ({ active, over }: DragEndEvent) => {
    if (active.id !== over?.id) {
      setTabs?.((prev) => {
        const activeIndex = prev.findIndex((i) => i.routePath === active.id);
        const overIndex = prev.findIndex((i) => i.routePath === over?.id);
        return arrayMove(prev, activeIndex, overIndex);
      });
    }
  };
  const sensor = useSensor(PointerSensor, { activationConstraint: { distance: 10 } });

  const handleScroll = () => {
    const container = tabsContainerRef.current;
    if (container) {
      setTimeout(() => {
        const { scrollLeft, scrollWidth, clientWidth } = container;

        if (scrollWidth > clientWidth) {
          setShowPrev(scrollLeft > 0);
          setShowNext(scrollLeft < scrollWidth - clientWidth);
        } else {
          setShowPrev(false);
          setShowNext(false);
        }
      }, 50); // 延迟 50ms 检查
    }
  };

  const handlePrevClick = () => {
    if (tabsContainerRef.current) {
      tabsContainerRef.current.scrollLeft -= 100; // 每次滚动 100 像素
    }
  };

  const handleNextClick = () => {
    if (tabsContainerRef.current) {
      tabsContainerRef.current.scrollLeft += 100;
    }
  };

  useEffect(() => {
    const container = tabsContainerRef.current;
    if (container) {
      container.addEventListener('scroll', handleScroll);
    }
    window.addEventListener('resize', handleScroll); // 窗口尺寸变化时也更新状态
    return () => {
      if (container) {
        container.removeEventListener('scroll', handleScroll);
      }
      window.removeEventListener('resize', handleScroll);
    };
  }, []);

  useDeepCompareEffect(() => {
    handleScroll();
  }, [options, wrapperSize?.width]);

  useLayoutEffect(() => {
    if (wrapperRef.current) {
      const wrapper = wrapperRef.current;
      const activeTabDom = wrapper.querySelector(`[data-key="${activeTag}"]`);

      if (activeTabDom) {
        // 获取容器与激活 Tab 的位置信息
        setTimeout(() => {
          const wrapperRect = wrapper.getBoundingClientRect();
          const tabRect = activeTabDom.getBoundingClientRect();
          if (wrapperRect.right < tabRect.right || wrapperRect.left > tabRect.left) {
            activeTabDom.scrollIntoView({
              behavior: 'smooth',
              block: 'end',
              inline: 'start',
            });
          }
        }, 50);
      }
    }
  }, [activeTag]);

  return (
    <div
      ref={wrapperRef}
      style={{ width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'flex-start' }}
    >
      <DndContext sensors={[sensor]} onDragEnd={onDragEnd} modifiers={[restrictToHorizontalAxis]}>
        <SortableContext items={options.map((i) => i.key)} strategy={horizontalListSortingStrategy}>
          {showPrev && (
            <div className="scroll-button left" onClick={handlePrevClick}>
              <ChevronLeft />
            </div>
          )}
          <div className="com-tags" ref={tabsContainerRef} onWheel={handleWheelScroll}>
            {options?.map(({ children, key, onClick, onClose, ...restProps }) => (
              <DraggableTagNode data-node-key={key} key={key} style={{ opacity: activeTag === key ? 1 : 0.6 }}>
                <Tag
                  data-key={key}
                  className="com-tags-item"
                  bordered={false}
                  onClose={(e) => {
                    e.stopPropagation();
                    e.preventDefault();
                    onClose?.(key);
                  }}
                  onClick={() => onClick?.(key)}
                  {...restProps}
                  onContextMenu={(e) => handleContextMenu(e, key)}
                >
                  {children}
                </Tag>
              </DraggableTagNode>
            ))}
          </div>
          {showNext && (
            <div className="scroll-button right" onClick={handleNextClick}>
              <ChevronRight />
            </div>
          )}
        </SortableContext>
      </DndContext>
      <ControlledDropdown ref={dropdownRef} />
    </div>
  );
};

export default ComTags;
