import { type CSSProperties, type FC, type ReactNode, useState, useRef } from 'react';
import classNames from 'classnames';
import styles from './index.module.scss';
import { debounce } from 'lodash-es';
import { useMediaSize } from '@/hooks';
import { Close } from '@carbon/icons-react';
import { useUnsTreeMapContext } from '@/UnsTreeMapContext';
import { useThemeStore } from '@/stores/theme-store.ts';

export interface ComLeftProps {
  children?: ReactNode;
  style?: CSSProperties;
  className?: string;
  styleCallBack?: (b: boolean) => CSSProperties;
  resize?: boolean;
  title?: ReactNode;
  defaultWidth?: number;
}

const ComLeft: FC<ComLeftProps> = ({
  children,
  style,
  className,
  styleCallBack,
  resize,
  title,
  defaultWidth = 300,
}) => {
  const { isH5 } = useMediaSize();
  const { setTreeMapVisible } = useUnsTreeMapContext();
  const isTop = useThemeStore((state) => state.isTop);
  const [width, setWidth] = useState(defaultWidth); // 默认宽度为300px
  const [isResizing, setIsResizing] = useState(false); // 调整大小
  const [mouseEnter, setMouseEnter] = useState(false); // 鼠标移入
  const isResizingRef: any = useRef(false);
  const containerRef: any = useRef(null);
  const initialXRef: any = useRef(0);
  const initialWidthRef: any = useRef(defaultWidth);
  // const [isFloating, setIsFloating] = useState(isH5); // 控制是否显示为悬浮按钮
  // const formatMessage = useTranslate();
  // 防抖函数，延迟 200ms 执行
  // eslint-disable-next-line react-hooks/refs
  const handleResize = debounce(() => {
    if (containerRef.current) {
      const resizeEvent = new Event('resize');
      // 在window对象上触发这个事件
      window.dispatchEvent(resizeEvent);
    }
  }, 200); // 200ms 防抖延迟
  const handleMouseDown = (e: any) => {
    isResizingRef.current = true;
    setIsResizing(true);
    initialXRef.current = e.clientX;
    initialWidthRef.current = width;
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  };

  const handleMouseMove = (e: any) => {
    if (isResizingRef.current) {
      const diffX = e.clientX - initialXRef.current;
      let newWidth = initialWidthRef.current + diffX;
      if (newWidth > 650) newWidth = 650;
      if (newWidth < 250) newWidth = 250;
      setWidth(newWidth);
    }
  };

  const handleMouseUp = () => {
    isResizingRef.current = false;
    setIsResizing(false);
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
    handleResize();
  };
  return (
    <>
      <div
        className={classNames(
          styles['com-left'],
          className,
          !isH5 && (isResizing || mouseEnter) ? styles['resize-focus'] : ''
          // isFloating && isH5 ? styles['floating'] : ''
        )}
        style={{ width: isH5 ? '80%' : `${width}px` }}
        ref={containerRef}
      >
        {isH5 && (
          <Close
            size={16}
            className="title-close"
            onClick={() => {
              // setIsFloating(false);
              setTreeMapVisible(false);
            }}
          />
        )}
        {title && <div className="title">{title}</div>}
        <div
          style={{
            overflow: 'auto',
            flex: 1,
            ...style,
            ...(styleCallBack?.(!isTop) || {}),
          }}
        >
          {children}
        </div>
        {resize && (
          <div
            className="col-resize-wrap"
            onMouseDown={handleMouseDown}
            onMouseEnter={() => {
              setMouseEnter(true);
            }}
            onMouseLeave={() => {
              setMouseEnter(false);
            }}
          />
        )}
      </div>
      {/* {isH5 && isFloating && (
        <Flex onClick={() => setIsFloating(false)} className={styles['floating-btn']}>
          <TreeViewIcon size={20} style={{ color: 'var(--supos-theme-color)' }} />
          <span style={{ color: 'var(--supos-theme-color)', paddingLeft: '4px' }}>{formatMessage('uns.treeList')}</span>
        </Flex>
      )} */}
      {isResizing && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            zIndex: 9999,
            cursor: 'col-resize',
          }}
        />
      )}
    </>
  );
};

export default ComLeft;
