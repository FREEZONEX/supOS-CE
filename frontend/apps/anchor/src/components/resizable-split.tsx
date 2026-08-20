// 左右可拖拽分栏（对齐源端 three-panel-layout：百分比宽度 + 4px col-resize 拖拽柄）
import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react';

export function ResizableSplit({
  left,
  right,
  defaultLeftWidth = 38,
  minLeftWidth = 24,
  maxLeftWidth = 48,
  className,
}: {
  left: ReactNode;
  right: ReactNode;
  defaultLeftWidth?: number; // 百分比 (0-100)
  minLeftWidth?: number;
  maxLeftWidth?: number;
  className?: string;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [leftWidth, setLeftWidth] = useState(defaultLeftWidth);
  const [dragging, setDragging] = useState(false);

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    setDragging(true);
  }, []);

  useEffect(() => {
    if (!dragging) return;
    const handleMouseMove = (e: MouseEvent) => {
      const container = containerRef.current;
      if (!container) return;
      const rect = container.getBoundingClientRect();
      const next = ((e.clientX - rect.left) / rect.width) * 100;
      setLeftWidth(Math.min(Math.max(next, minLeftWidth), maxLeftWidth));
    };
    const handleMouseUp = () => setDragging(false);
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    };
  }, [dragging, minLeftWidth, maxLeftWidth]);

  return (
    <div ref={containerRef} className={`flex h-full min-h-0 w-full ${className ?? ''}`}>
      <div className="h-full min-w-0 overflow-hidden" style={{ width: `${leftWidth}%` }}>
        {left}
      </div>
      <div
        className="h-full w-1 shrink-0 cursor-col-resize"
        style={{ background: dragging ? 'var(--ui-theme-color)' : 'var(--ui-header-splitter-color, #e0e0e0)' }}
        onMouseDown={handleMouseDown}
        onMouseEnter={(e) => {
          if (!dragging) (e.currentTarget as HTMLDivElement).style.background = 'var(--ui-theme-color)';
        }}
        onMouseLeave={(e) => {
          if (!dragging)
            (e.currentTarget as HTMLDivElement).style.background = 'var(--ui-header-splitter-color, #e0e0e0)';
        }}
      />
      <div className="h-full min-w-0 flex-1 overflow-hidden">{right}</div>
    </div>
  );
}
