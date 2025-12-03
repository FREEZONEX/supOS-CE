import { forwardRef, useState, useRef, type ChangeEvent, type InputHTMLAttributes, useEffect } from 'react';
import { Search, Close } from '@carbon/icons-react';
import { useMergedRefs } from '@/hooks/useMergedRefs';
import cx from 'classnames';
import './index.scss';

type InputPropsBase = Omit<InputHTMLAttributes<HTMLInputElement>, 'size'>;
export interface ASearchProps extends InputPropsBase {
  closeButtonLabelText?: string;
  onClear?: () => void;
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
      closeButtonLabelText,
      style,
      size = 'md',
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
      inputRef.current?.focus();
    };

    useEffect(() => {
      setVal(value);
    }, [value]);

    return (
      <div className={searchClasses}>
        <Search className="custom-search-icon" />
        <input
          {...restProps}
          autoComplete={autoComplete}
          ref={ref}
          className="custom-search-input"
          value={val}
          onChange={handleChange}
          style={{ ...style, paddingRight: val ? '32px' : '10px' }}
        />
        {val && (
          <button className="custom-search-clear" onClick={handleClear} title={closeButtonLabelText}>
            <Close />
          </button>
        )}
      </div>
    );
  }
);

export default ProSearch;
