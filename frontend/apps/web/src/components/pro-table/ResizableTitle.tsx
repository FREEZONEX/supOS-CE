import { type FC, useEffect, useRef, useState, type HTMLAttributes } from 'react';
import type { TitlePropsType } from './types.ts';

// 查找所有可滚动祖先元素（包括window）
const findScrollParents = (element: HTMLElement | null): (HTMLElement | Window)[] => {
  const parents: (HTMLElement | Window)[] = [];
  let current = element;

  while (current) {
    parents.push(current);
    const { overflow, overflowX, overflowY } = getComputedStyle(current);
    if (/(auto|scroll)/.test(overflow + overflowX + overflowY)) {
      break;
    }
    current = current.parentElement;
  }
  // 添加window作为可能的滚动容器
  parents.push(window);
  return parents;
};

export const ResizableTitle: FC<Readonly<HTMLAttributes<any> & TitlePropsType>> = (props) => {
  const { changeWidth, width, minWidth = 100, children, ...restProps } = props;
  const [tableRect, setTableRect] = useState<DOMRect | null>(null);
  const isResizingRef: any = useRef(false);
  const thRef = useRef<HTMLTableHeaderCellElement>(null);
  const thStartXRef = useRef(0);
  const thWidthRef = useRef(width || 0);
  const diffXRef = useRef(0);
  const highLineRef = useRef<any>(null);

  const changeBodyUserSelect = (value: string) => {
    const bodyStyle = document.body.style;
    bodyStyle.userSelect = value; // Non-prefixed version, currently supported by all major browsers
  };

  const handleMouseDown = (e: any) => {
    if (!thRef.current) return;
    isResizingRef.current = true;
    // 拖拽开始时记录当前th最右侧的x坐标
    const currentX = thRef.current?.clientWidth + thRef.current?.getBoundingClientRect()?.x;
    // 拖拽开始时记录当前鼠标的x坐标
    thStartXRef.current = e.clientX;
    // 拖拽开始时记录当前鼠标的x坐标与th最右侧的x坐标的差值
    diffXRef.current = e.clientX - currentX + 1;
    highLineRef.current.style.left = currentX - 1 + 'px';
    document.body.style.cursor = 'col-resize';
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    // 拖拽开始时禁止文本选择功能
    changeBodyUserSelect('none');
    //添加高亮线
    document.body.appendChild(highLineRef.current);
  };

  const handleMouseMove = (e: any) => {
    if (!isResizingRef.current || !thRef.current) return;
    //实际宽度 = th的宽度 + 鼠标移动的距离 - 鼠标按下时的x坐标
    const realWidth = thRef.current?.clientWidth + e.clientX - thStartXRef.current;
    //记录当前th的最左侧的x坐标
    const thRectX = thRef.current?.getBoundingClientRect()?.x;
    // if (minWidth && realWidth < minWidth) {
    thWidthRef.current = Math.max(minWidth, realWidth);
    highLineRef.current.style.left = Math.max(thRectX + minWidth - 1, e.clientX - diffXRef.current) + 'px';
    // } else {
    //   thWidthRef.current = Math.max(80, realWidth);
    //   highLineRef.current.style.left = Math.max(thRectX + 80 - 1, e.clientX - diffXRef.current) + 'px';
    // }
  };

  const handleMouseUp = () => {
    if (highLineRef.current) {
      //移除高亮线
      document.body?.removeChild(highLineRef.current);
    }
    isResizingRef.current = false;
    document.body.style.cursor = '';
    //鼠标抬起的时候再改变宽度
    changeWidth(thWidthRef.current);
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
    // 拖拽结束时恢复文本选择功能
    changeBodyUserSelect('');
  };

  useEffect(() => {
    highLineRef.current = document.createElement('div');
    highLineRef.current.style.cssText = `
      position: fixed;
      z-index: 10000;
      width: 1px;
      height: ${tableRect?.height}px;
      top: ${tableRect?.top}px;
      background: var(--supos-theme-color);
    `;
  }, [tableRect]);

  useEffect(() => {
    if (!thRef.current) return;
    const tableContainer: any = thRef.current.closest('.ant-table-container');
    if (!tableContainer) return;
    const updateSize = () => {
      const rect = tableContainer.getBoundingClientRect();
      setTableRect(rect);
    };
    // 初始化尺寸
    updateSize();
    // 监听容器尺寸变化
    const resizeObserver = new ResizeObserver(updateSize);
    resizeObserver.observe(tableContainer);
    // 监听滚动事件
    const scrollParents = findScrollParents(tableContainer);
    const handleScroll = () => {
      updateSize(); // 滚动时更新位置
    };

    scrollParents.forEach((parent: any) => {
      if (parent instanceof Window) {
        window.addEventListener('scroll', handleScroll, { passive: true });
      } else {
        parent.addEventListener('scroll', handleScroll, { passive: true });
      }
    });
    return () => {
      resizeObserver.disconnect();
      scrollParents.forEach((parent: any) => {
        if (parent instanceof Window) {
          window.removeEventListener('scroll', handleScroll);
        } else {
          (parent as HTMLElement).removeEventListener('scroll', handleScroll);
        }
      });
    };
  }, []);

  useEffect(() => {
    if (!highLineRef.current || !tableRect) return;
    // 动态更新高亮线位置
    highLineRef.current.style.height = `${tableRect.height}px`;
    highLineRef.current.style.top = `${tableRect.top}px`;
  }, [tableRect]);

  return (
    <th {...restProps} ref={thRef} title={Array.isArray(children) && children.length > 1 ? children[1] : undefined}>
      {children}

      <div className="react-resizable-line" onMouseDown={handleMouseDown} />
    </th>
  );
};
