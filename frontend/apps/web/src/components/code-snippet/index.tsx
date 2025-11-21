import {
  type FC,
  type ReactNode,
  useRef,
  useEffect,
  useState,
  useCallback,
  type MouseEvent,
  useMemo,
  Children,
} from 'react';
import { Tooltip } from 'antd';
import { ChevronDown, Copy } from '@carbon/icons-react';
import classNames from 'classnames';
import { useTranslate } from '@/hooks';
import { useSize } from 'ahooks';
import './index.scss';
import useClipboard from '../../hooks/useClipboard.ts';

export type CodeSnippetType = 'single' | 'multi' | 'inline';
export type DeprecatedCodeSnippetAlignment =
  | 'top-left'
  | 'top-right'
  | 'bottom-left'
  | 'bottom-right'
  | 'left-bottom'
  | 'left-top'
  | 'right-bottom'
  | 'right-top';

export type NewCodeSnippetAlignment =
  | 'top'
  | 'bottom'
  | 'left'
  | 'right'
  | 'top-start'
  | 'top-end'
  | 'bottom-start'
  | 'bottom-end'
  | 'left-end'
  | 'left-start'
  | 'right-end'
  | 'right-start';

export type CodeSnippetAlignment = DeprecatedCodeSnippetAlignment | NewCodeSnippetAlignment;
export interface CodeSnippetProps {
  /**
   * 指定触发器与工具提示的对齐方式
   */
  align?: CodeSnippetAlignment;

  /**
   * 实验性功能：尝试自动对齐工具提示
   */
  autoAlign?: boolean;

  /**
   * 为文本框容器节点指定屏幕阅读器可读的标签
   */
  ['aria-label']?: string;

  /**
   * 为容器节点指定可选的类名
   */
  className?: string;

  /**
   * 指定复制按钮的描述文本
   */
  copyButtonDescription?: string;

  /**
   * 可选的复制文本。如果未指定，将使用 `children` 节点的 `innerText` 作为复制值
   */
  copyText?: string;

  /**
   * 指定代码片段是否应禁用
   */
  disabled?: boolean;

  /**
   * 指定代码片段复制后显示的提示文本
   */
  feedback?: string;

  /**
   * 指定提示消息的超时时间（毫秒）
   */
  feedbackTimeout?: number;

  /**
   * 指定是否隐藏/不渲染复制按钮
   */
  hideCopyButton?: boolean;

  /**
   * 指定是否使用代码片段的浅色变体，通常用于内联片段以显示替代颜色
   */
  light?: boolean;

  /**
   * 指定折叠视图下显示的最大行数
   */
  maxCollapsedNumberOfRows?: number;

  /**
   * 指定展开视图下显示的最大行数
   */
  maxExpandedNumberOfRows?: number;

  /**
   * 指定折叠视图下显示的最小行数
   */
  minCollapsedNumberOfRows?: number;

  /**
   * 指定展开视图下显示的最小行数
   */
  minExpandedNumberOfRows?: number;

  /**
   * 可选的点击事件监听器，用于响应复制按钮的 `onClick` 事件
   */
  onClick?: (e: MouseEvent) => void;

  /**
   * 指定代码片段交互后显示"显示更少"的文本
   */
  showLessText?: string;

  /**
   * 当代码片段超过15行时显示的"显示更多"文本
   */
  showMoreText?: string;

  /**
   * 指定代码片段的类型
   */
  type?: 'single' | 'inline' | 'multi';

  /**
   * 指定是否显示行号
   */
  showLineNumbers?: boolean;

  /**
   * 指定是否自动换行
   */
  wrapText?: boolean;
  children?: ReactNode;
}

const rowHeightInPixels = 16;
const defaultMaxCollapsedNumberOfRows = 15;
const defaultMaxExpandedNumberOfRows = 0;
const defaultMinCollapsedNumberOfRows = 3;
const defaultMinExpandedNumberOfRows = 16;

/**
 * 代码片段组件
 *
 * 用于展示格式化的代码片段，支持复制功能
 */
const CodeSnippet: FC<CodeSnippetProps> = ({
  align = 'bottom',
  // autoAlign = false,
  className,
  type = 'single',
  children,
  disabled,
  feedback,
  // feedbackTimeout,
  onClick,
  ['aria-label']: ariaLabel = 'Copy to clipboard',
  copyText,
  copyButtonDescription,
  light,
  showMoreText,
  showLessText,
  hideCopyButton,
  wrapText = false,
  showLineNumbers = false,
  maxCollapsedNumberOfRows = defaultMaxCollapsedNumberOfRows,
  maxExpandedNumberOfRows = defaultMaxExpandedNumberOfRows,
  minCollapsedNumberOfRows = defaultMinCollapsedNumberOfRows,
  minExpandedNumberOfRows = defaultMinExpandedNumberOfRows,
  ...rest
}) => {
  const [expandedCode, setExpandedCode] = useState(false);
  const [shouldShowMoreLessBtn, setShouldShowMoreLessBtn] = useState(false);
  const codeContentRef = useRef<HTMLPreElement>(null);
  const codeContainerRef = useRef<HTMLDivElement>(null);
  const innerCodeRef = useRef<HTMLElement>(null);
  const [hasLeftOverflow, setHasLeftOverflow] = useState(false);
  const [hasRightOverflow, setHasRightOverflow] = useState(false);
  const copyButtonRef = useRef<HTMLButtonElement>(null);
  const formatMessage = useTranslate();
  const { copy } = useClipboard();

  const getCodeRef = useCallback(() => {
    if (type === 'single') {
      return codeContainerRef;
    }
    if (type === 'multi') {
      return codeContentRef;
    } else {
      return innerCodeRef;
    }
  }, [type]);

  const getCodeRefDimensions = useCallback(() => {
    const {
      clientWidth: codeClientWidth = 0,
      scrollLeft: codeScrollLeft = 0,
      scrollWidth: codeScrollWidth = 0,
    } = getCodeRef().current || {};

    return {
      horizontalOverflow: codeScrollWidth > codeClientWidth,
      codeClientWidth,
      codeScrollWidth,
      codeScrollLeft,
    };
  }, [getCodeRef]);

  const handleScroll = useCallback(() => {
    if (
      type === 'inline' ||
      (type === 'single' && !codeContainerRef?.current) ||
      (type === 'multi' && !codeContentRef?.current)
    ) {
      return;
    }

    const { horizontalOverflow, codeClientWidth, codeScrollWidth, codeScrollLeft } = getCodeRefDimensions();

    setHasLeftOverflow(horizontalOverflow && !!codeScrollLeft);
    setHasRightOverflow(horizontalOverflow && codeScrollLeft + codeClientWidth !== codeScrollWidth);
  }, [type, getCodeRefDimensions]);

  useEffect(() => {
    handleScroll();
  }, [handleScroll]);

  // 使用useSize钩子监听代码片段大小变化
  const size = useSize(getCodeRef());

  useEffect(() => {
    if (codeContentRef?.current && type === 'multi') {
      const { height } = codeContentRef.current.getBoundingClientRect();

      if (
        maxCollapsedNumberOfRows > 0 &&
        (maxExpandedNumberOfRows <= 0 || maxExpandedNumberOfRows > maxCollapsedNumberOfRows) &&
        height > maxCollapsedNumberOfRows * rowHeightInPixels
      ) {
        setShouldShowMoreLessBtn(true);
      } else {
        setShouldShowMoreLessBtn(false);
      }
      if (expandedCode && minExpandedNumberOfRows > 0 && height <= minExpandedNumberOfRows * rowHeightInPixels) {
        setExpandedCode(false);
      }
    }
    if ((codeContentRef?.current && type === 'multi') || (codeContainerRef?.current && type === 'single')) {
      handleScroll();
    }
  }, [
    size,
    type,
    expandedCode,
    maxCollapsedNumberOfRows,
    maxExpandedNumberOfRows,
    minExpandedNumberOfRows,
    minCollapsedNumberOfRows,
    handleScroll,
  ]);

  const handleCopyClick = (evt: MouseEvent) => {
    if (copyText || innerCodeRef?.current) {
      const textToCopy = copyText ?? innerCodeRef?.current?.innerText ?? '';
      copy(textToCopy);
    }

    if (onClick) {
      onClick(evt as unknown as MouseEvent);
    }
  };

  const getChildrenString = useCallback((node: ReactNode): string => {
    let text = '';
    Children.forEach(node, (child) => {
      if (typeof child === 'string' || typeof child === 'number') {
        text += child;
      } else if (child && typeof child === 'object' && 'props' in child && child.props.children) {
        text += getChildrenString(child.props.children);
      }
    });
    return text;
  }, []);

  const processedChildren = useMemo(() => {
    if (showLineNumbers && type === 'multi') {
      const textContent = getChildrenString(children);
      const lines = textContent.split('\n');

      if (lines.length > 0 && lines[lines.length - 1] === '') {
        lines.pop();
      }

      return lines.map((line, index) => (
        <span className="code-snippet__line" key={index}>
          {line === '' ? '\u00A0' : line}
        </span>
      ));
    }
    return children;
  }, [children, showLineNumbers, type, getChildrenString]);

  const codeSnippetClasses = classNames(className, 'code-snippet', {
    [`code-snippet--${type}`]: type,
    'code-snippet--disabled': type !== 'inline' && disabled,
    'code-snippet--expand': expandedCode,
    'code-snippet--light': light,
    'code-snippet--no-copy': hideCopyButton,
    'code-snippet--wraptext': wrapText,
    'code-snippet--has-right-overflow': type === 'multi' && hasRightOverflow,
    'code-snippet--copy-button': !hideCopyButton,
    'code-snippet--with-line-numbers': showLineNumbers && type === 'multi',
  });

  const expandCodeBtnText = expandedCode
    ? (showLessText ?? formatMessage('uns.showLess'))
    : (showMoreText ?? formatMessage('uns.showMore'));

  if (type === 'inline') {
    if (hideCopyButton) {
      return (
        <span className={codeSnippetClasses}>
          <code ref={innerCodeRef}>{children}</code>
        </span>
      );
    }
    return (
      <Tooltip title={feedback || formatMessage('common.copy')} trigger="hover" placement={align as any}>
        <button className={codeSnippetClasses} onClick={handleCopyClick} aria-label={ariaLabel} {...rest}>
          <code ref={innerCodeRef}>{children}</code>
        </button>
      </Tooltip>
    );
  }

  type stylesType = { maxHeight?: number; minHeight?: number };
  type containerStyleType = { style?: stylesType };
  const containerStyle: containerStyleType = {};
  if (type === 'multi') {
    const styles: stylesType = {};

    if (expandedCode) {
      if (maxExpandedNumberOfRows > 0) {
        styles.maxHeight = maxExpandedNumberOfRows * rowHeightInPixels;
      }
      if (minExpandedNumberOfRows > 0) {
        styles.minHeight = minExpandedNumberOfRows * rowHeightInPixels;
      }
    } else {
      if (maxCollapsedNumberOfRows > 0) {
        styles.maxHeight = maxCollapsedNumberOfRows * rowHeightInPixels;
      }
      if (minCollapsedNumberOfRows > 0) {
        styles.minHeight = minCollapsedNumberOfRows * rowHeightInPixels;
      }
    }

    if (Object.keys(styles).length) {
      containerStyle.style = styles;
    }
  }
  return (
    <div {...rest} className={codeSnippetClasses}>
      <div
        ref={codeContainerRef}
        role={type === 'single' || type === 'multi' ? 'textbox' : undefined}
        tabIndex={(type === 'single' || type === 'multi') && !disabled ? 0 : undefined}
        className="code-snippet-container"
        aria-label={ariaLabel || 'code-snippet'}
        aria-readonly={type === 'single' || type === 'multi' ? true : undefined}
        aria-multiline={type === 'multi' ? true : undefined}
        onScroll={(type === 'single' && handleScroll) || undefined}
        {...containerStyle}
      >
        <pre ref={codeContentRef} onScroll={(type === 'multi' && handleScroll) || undefined}>
          <code ref={innerCodeRef}>{processedChildren}</code>
        </pre>
      </div>
      {/* 左侧溢出指示器 */}
      {hasLeftOverflow && <div className="code-snippet-overflow-indicator--left" />}
      {/* 右侧溢出指示器 */}
      {hasRightOverflow && type !== 'multi' && <div className="code-snippet-overflow-indicator--right" />}
      {!hideCopyButton && (
        <Tooltip title={feedback || formatMessage('common.copy')} trigger="hover" placement={align as any}>
          <button
            className="code-snippet-copy-btn"
            ref={copyButtonRef}
            disabled={disabled}
            onClick={handleCopyClick}
            aria-label={copyButtonDescription || formatMessage('common.copy')}
          >
            <Copy color="var(--supos-text-color)" style={{ cursor: 'pointer' }} size={16} />
          </button>
        </Tooltip>
      )}
      {shouldShowMoreLessBtn && (
        <button className="code-snippet-btn--expand" disabled={disabled} onClick={() => setExpandedCode(!expandedCode)}>
          <span className="code-snippet-btn--text">{expandCodeBtnText}</span>
          <ChevronDown className="code-snippet-icon-chevron--down code-snippet-icon" role="img" />
        </button>
      )}
    </div>
  );
};

export default CodeSnippet;
