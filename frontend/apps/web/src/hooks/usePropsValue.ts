import { type SetStateAction, useRef } from 'react';
import { useMemoizedFn, useUpdate } from 'ahooks';

interface Options<T> {
  value?: T;
  defaultValue?: T;
  onChange?: (v: T) => void;
}

// 内部控制value和onChange
function usePropsValue<T>(options: Options<T>): any {
  const { value, defaultValue, onChange } = options;
  const update = useUpdate();
  const stateRef = useRef<T | undefined>(value ?? defaultValue);
  if (value !== undefined) {
    stateRef.current = value;
  }

  const setState = useMemoizedFn((v: SetStateAction<T>, forceTrigger: boolean = false) => {
    const nextValue = typeof v === 'function' ? (v as (prevState: T) => T)(stateRef.current as T) : v;
    if (!forceTrigger && nextValue === stateRef.current) return;
    stateRef.current = nextValue;
    update();
    return onChange?.(nextValue);
  });
  return [stateRef.current, setState] as const;
}

export default usePropsValue;
