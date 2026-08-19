import {
  forwardRef,
  useState,
  useRef,
  type ChangeEvent,
  type InputHTMLAttributes,
  type KeyboardEvent,
  useEffect,
} from 'react';
import { Search, Close } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useMergedRefs } from '@/hooks/useMergedRefs';
import cx from 'classnames';
import './index.scss';

type InputPropsBase = Omit<InputHTMLAttributes<HTMLInputElement>, 'size'>;
export interface ASearchProps extends InputPropsBase {
  closeButtonLabelText?: string;
  onClear?: () => void;
  onSearch?: (value: string) => void;
  size?: 'sm' | 'md' | 'lg';
}

const ProSearch = forwardRef<HTMLInputElement, ASearchProps>(
  (
    {
      autoComplete = 'off',
      className,
      value,
      onChange,
      onClear,
      onSearch,
      closeButtonLabelText,
      style,
      size = 'md',
      onKeyDown,
      ...restProps
    },
    searchRef
  ) => {
    const [val, setVal] = useState(value);
    const inputRef = useRef<HTMLInputElement>(null);
    const ref = useMergedRefs<HTMLInputElement>([searchRef, inputRef]);

    const searchClasses = cx(
      {
        'custom-search': true,
        'custom-search-sm': size === 'sm',
        'custom-search-md': size === 'md',
        'custom-search-lg': size === 'lg',
      },
      className
    );
    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
      setVal(e.target.value || '');
      onChange?.(e);
    };

    const handleClear = () => {
      const inputTarget = Object.assign({}, inputRef.current, { value: '' });
      handleChange({ target: inputTarget, type: 'change' } as ChangeEvent<HTMLInputElement>);
      onClear?.();
      onSearch?.('');
      inputRef.current?.focus();
    };

    const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
      if (event.key === 'Enter') {
        event.preventDefault();
        onSearch?.(String(val ?? ''));
      }
      onKeyDown?.(event);
    };

    useEffect(() => {
      setVal(value);
    }, [value]);

    return (
      <div className={searchClasses}>
        <Search className="custom-search-icon" {...toolbarIconProps} />
        <input
          {...restProps}
          autoComplete={autoComplete}
          ref={ref}
          className="custom-search-input"
          value={val}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          style={{
            ...style,
            ...(val ? { paddingRight: size === 'sm' ? 32 : size === 'lg' ? 48 : 40 } : {}),
          }}
        />
        {val && (
          <button className="custom-search-clear" onClick={handleClear} title={closeButtonLabelText}>
            <Close {...toolbarIconProps} />
          </button>
        )}
      </div>
    );
  }
);

export default ProSearch;
