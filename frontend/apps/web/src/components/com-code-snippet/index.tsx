import { type ComponentProps, type CSSProperties, type FC, useEffect, useRef } from 'react';
import classNames from 'classnames';
import './index.scss';
import { useSize } from 'ahooks';
import CodeSnippet from '../code-snippet';

type CodeSnippetProps = ComponentProps<typeof CodeSnippet>;

interface ComCodeSnippetProps extends CodeSnippetProps {
  style?: CSSProperties;
  onSizeChange?: (size?: { height: number; width: number }) => void;
  copyPosition?: boolean;
}

const ComCodeSnippet: FC<ComCodeSnippetProps> = ({
  className,
  onSizeChange,
  style,
  copyPosition = true,
  ...restProps
}) => {
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);

  useEffect(() => {
    if (!ref.current) return;
    onSizeChange?.({ height: ref.current?.offsetHeight, width: ref.current?.offsetWidth });
  }, []);

  useEffect(() => {
    onSizeChange?.(size);
  }, [size]);
  return (
    <div className="com-code-snippet" ref={ref} style={style}>
      <CodeSnippet
        className={classNames('code-snippet-wrapper', className, { 'com-copy-code-snippet': copyPosition })}
        type="multi"
        minCollapsedNumberOfRows={26}
        maxCollapsedNumberOfRows={26}
        {...restProps}
      />
    </div>
  );
};

export default ComCodeSnippet;
