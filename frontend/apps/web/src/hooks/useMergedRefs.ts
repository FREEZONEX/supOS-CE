import { useCallback, type Ref, type ForwardedRef } from 'react';

export const useMergedRefs = <T>(refs: ForwardedRef<T>[]): Ref<T> => {
  return useCallback((node: T | null) => {
    refs.forEach((ref) => {
      if (typeof ref === 'function') {
        ref(node);
      } else if (ref !== null && ref !== undefined) {
        ref.current = node;
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, refs);
};
