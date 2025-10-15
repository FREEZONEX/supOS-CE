import { forwardRef, type ReactNode } from 'react';

export const OverlayItem = forwardRef<HTMLDivElement, { overlayChildren?: ReactNode }>(({ overlayChildren }, ref) => {
  return (
    <div ref={ref} style={{ cursor: 'pointer', opacity: 1 }}>
      {overlayChildren}
    </div>
  );
});
