import { type CSSProperties } from 'react';
import CodeMirror, { type ReactCodeMirrorProps, type ReactCodeMirrorRef } from '@uiw/react-codemirror';
import codeStyles from '@/theme/codemirror.module.scss';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { Copy, ChevronDown, ChevronUp } from '@carbon/icons-react';
import { useClipboard, useTranslate } from '@/hooks';
import { useRef, useState } from 'react';
import { Button } from 'antd';

interface ProCodemirrorProps extends ReactCodeMirrorProps {
  wrapperStyle?: CSSProperties;
  showExpanded?: boolean;
  showHint?: boolean;
}

const ProCodemirror = (props: ProCodemirrorProps) => {
  const { wrapperStyle, showExpanded, showHint = true, ...restProps } = props;
  const formatMessage = useTranslate();
  const { copy } = useClipboard();
  const editorRef = useRef<ReactCodeMirrorRef>(null);
  const [isExpanded, setIsExpanded] = useState(false);

  return (
    <div
      style={{
        borderRadius: 4,
        border: '1px solid rgb(198, 198, 198)',
        padding: 16,
        position: 'relative',
        overflow: 'hidden',
        ...wrapperStyle,
      }}
      className={codeStyles['custom-theme']}
    >
      <div
        style={{
          position: 'absolute',
          right: 4,
          top: 0,
          color: 'var(--ui-text-color)',
          zIndex: 1,
        }}
      >
        {restProps?.value ? (
          <Copy
            style={{
              cursor: 'pointer',
              marginTop: 4,
            }}
            onClick={() => {
              copy(restProps?.value);
            }}
          />
        ) : (
          showHint && (
            <span
              style={{
                marginRight: 14,
                fontSize: '12px',
                pointerEvents: 'none',
                zIndex: 10,
                color: '#c6c6c6',
              }}
            >
              {formatMessage('uns.ctrlPQuickApplyExample')}
            </span>
          )
        )}
      </div>
      <CodeMirror
        ref={editorRef}
        theme={codemirrorTheme}
        {...restProps}
        onKeyDownCapture={(e) => {
          if (e.ctrlKey && e.key === 'Enter') {
            e.preventDefault();
            e.stopPropagation();
            if (editorRef.current?.view && !restProps?.value) {
              editorRef.current.view.dispatch({
                changes: {
                  from: 0,
                  to: editorRef.current.view.state.doc.length,
                  insert: String(restProps?.placeholder || ''),
                },
              });
            }
          }
        }}
        height={isExpanded ? undefined : restProps.height}
        onKeyDown={(e) => {
          if (e.ctrlKey && e.key === 'Enter') {
            e.preventDefault();
          }
        }}
      />
      {restProps.height && showExpanded && (
        <div
          style={{
            cursor: 'pointer',
            userSelect: 'none',
            position: 'absolute',
            right: 0,
            bottom: 0,
            color: 'var(--ui-text-color)',
            zIndex: 1,
          }}
          onClick={() => setIsExpanded(!isExpanded)}
        >
          {isExpanded ? (
            <Button color="default" variant="text">
              <ChevronUp size={16} />
              <span style={{ marginLeft: 4 }}>{formatMessage('uns.showLess')}</span>
            </Button>
          ) : (
            <Button color="default" variant="text">
              <ChevronDown size={16} />
              <span style={{ marginLeft: 4 }}>{formatMessage('uns.showMore')}</span>
            </Button>
          )}
        </div>
      )}
    </div>
  );
};

export default ProCodemirror;
