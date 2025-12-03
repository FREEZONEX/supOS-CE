import React, { useState, useCallback, useRef, useEffect } from 'react';
import classNames from 'classnames';
import './index.scss';

export interface ClickTriggerProps {
  /** 点击达到指定次数后触发的回调函数 */
  onTrigger?: () => void;
  /** 触发所需的点击次数，默认为5 */
  triggerCount?: number;
  /** 子元素 */
  children?: React.ReactNode;
  /** 自定义样式类名 */
  className?: string;
  /** 自定义样式 */
  style?: React.CSSProperties;
  /** 是否显示点击计数器 */
  showCounter?: boolean;
  /** 点击后是否重置计数器 */
  resetAfterTrigger?: boolean;
  /** 触发的时间窗口（毫秒），在此时间内达到点击次数才触发，默认为2000ms */
  timeWindow?: number;
}

/**
 * 点击触发组件 - 在达到指定点击次数后触发回调函数
 */
const ClickTrigger: React.FC<ClickTriggerProps> = (props) => {
  const {
    onTrigger,
    triggerCount = 5,
    children,
    className,
    style,
    showCounter = false,
    resetAfterTrigger = true,
    timeWindow = 2000, // 默认时间窗口为2秒
  } = props;

  const [clickCount, setClickCount] = useState<number>(0);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const firstClickTimeRef = useRef<number | null>(null);

  // 组件卸载时清理定时器
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  const resetCounter = useCallback(() => {
    setClickCount(0);
    firstClickTimeRef.current = null;
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const handleClick = useCallback(() => {
    const currentTime = Date.now();

    // 如果是第一次点击，记录时间并设置定时器
    if (clickCount === 0) {
      firstClickTimeRef.current = currentTime;
      timerRef.current = setTimeout(() => {
        resetCounter();
      }, timeWindow);
    }

    // 检查是否在时间窗口内
    if (firstClickTimeRef.current && currentTime - firstClickTimeRef.current <= timeWindow) {
      const newCount = clickCount + 1;
      setClickCount(newCount);

      if (newCount >= triggerCount) {
        onTrigger?.();
        if (resetAfterTrigger) {
          resetCounter();
        }
      }
    } else {
      // 超出时间窗口，重置计数器并重新开始
      resetCounter();
      firstClickTimeRef.current = currentTime;
      timerRef.current = setTimeout(() => {
        resetCounter();
      }, timeWindow);
      setClickCount(1); // 将当前点击计为第一次
    }
  }, [clickCount, triggerCount, onTrigger, resetAfterTrigger, timeWindow, resetCounter]);

  return (
    <div className={classNames('com-click-trigger', className)} style={style} onClick={handleClick}>
      {children}
      {showCounter && (
        <div className="com-click-trigger-counter">
          {clickCount}/{triggerCount}
        </div>
      )}
    </div>
  );
};

export default ClickTrigger;
