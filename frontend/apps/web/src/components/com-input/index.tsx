import { Input, type InputProps, type InputRef } from 'antd';
import { type ChangeEvent, forwardRef, useEffect, useRef } from 'react';
import { useMergedRefs } from '@/hooks/useMergedRefs.ts';
import usePropsValue from '@/hooks/usePropsValue.ts';

export interface ComInputProps extends Omit<InputProps, 'onChange'> {
  /**
   * @description 多返回有参数，来判断是否中文输入完成
   * */
  onChange?: (e: ChangeEvent<HTMLInputElement>, isComposing?: boolean) => void;
}

const ComInput = forwardRef<InputRef, ComInputProps>(({ value, onChange, ...restProps }, searchRef) => {
  const inputRef = useRef<InputRef>(null);
  const ref = useMergedRefs<InputRef>([searchRef, inputRef]);
  const isComposingRef = useRef(false); // 中文输入是否完成..
  const [val, setVal] = usePropsValue({
    value,
  });
  useEffect(() => {
    // 判断中文输入
    const inputElement = inputRef.current?.input;
    if (inputElement) {
      const handleCompositionStart = () => {
        isComposingRef.current = true;
      };

      const handleCompositionEnd = (e: any) => {
        isComposingRef.current = false;
        onChange?.({ target: e.target } as ChangeEvent<HTMLInputElement>, false);
      };

      inputElement.addEventListener('compositionstart', handleCompositionStart);
      inputElement.addEventListener('compositionend', handleCompositionEnd);

      return () => {
        inputElement.removeEventListener('compositionstart', handleCompositionStart);
        inputElement.removeEventListener('compositionend', handleCompositionEnd);
      };
    }
  }, []);

  return (
    <Input
      value={val}
      onChange={(e) => {
        setVal(e.target.value);
        onChange?.(e, isComposingRef.current);
      }}
      {...restProps}
      ref={ref}
    />
  );
});

export default ComInput;
