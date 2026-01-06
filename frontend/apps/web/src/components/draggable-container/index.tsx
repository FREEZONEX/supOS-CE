import { useState, useEffect, useRef, type FC, type ReactNode } from 'react';
import { Draggable } from '@carbon/icons-react';
import './index.scss';

const DraggableContainer: FC<{ children: ReactNode }> = ({ children }) => {
  const [position, setPosition] = useState({ x: 6, y: 8 });
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let isDragging = false;
    let offsetX: number, offsetY: number; // 相对于容器左上角的偏移量

    const handleMouseDown = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      if (
        ['navTop', 'navContent', 'navBottom', 'navHandle', 'navHandleIcon'].includes(target.className) ||
        target.closest('.navHandleIcon')
      ) {
        isDragging = true;
        document.body.style.cursor = 'move';
        if (containerRef.current) {
          offsetX = event.clientX - containerRef.current.offsetLeft;
          offsetY = event.clientY - containerRef.current.offsetTop;
        }
        if (document.getElementById('iframeMask')) {
          document.getElementById('iframeMask')!.style.display = 'block';
        }
        document.addEventListener('mousemove', handleMouseMove);
        document.addEventListener('mouseup', handleMouseUp);
      }
    };

    const handleMouseMove = (event: MouseEvent) => {
      if (!isDragging) return;
      const x = event.clientX - offsetX > 0 ? event.clientX - offsetX : 0;
      const y = event.clientY - offsetY > 0 ? event.clientY - offsetY : 0;
      setPosition({
        x,
        y,
      });
    };

    const handleMouseUp = () => {
      isDragging = false;
      document.body.style.cursor = 'default';
      if (document.getElementById('iframeMask')) {
        document.getElementById('iframeMask')!.style.display = 'none';
      }
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };

    // eslint-disable-next-line @typescript-eslint/no-unused-expressions
    containerRef.current && containerRef.current.addEventListener('mousedown', handleMouseDown);

    return () => {
      // eslint-disable-next-line @typescript-eslint/no-unused-expressions
      containerRef.current && containerRef.current.removeEventListener('mousedown', handleMouseDown);
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [containerRef, setPosition]); // 注意依赖项

  return (
    <div
      className="custom-draggable"
      ref={containerRef}
      style={{
        position: 'fixed',
        left: `${position.x}px`,
        top: `${position.y}px`,
        zIndex: 9999,
        // Add any other styles you need here
      }}
    >
      {children}
      <div className="navHandle">
        <Draggable className="navHandleIcon" />
      </div>
    </div>
  );
};

export default DraggableContainer;
