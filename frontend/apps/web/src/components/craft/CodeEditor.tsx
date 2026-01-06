import { type FC, useMemo, useRef } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { javascript } from '@codemirror/lang-javascript';
import { xml } from '@codemirror/lang-xml';
import { html } from '@codemirror/lang-html';
import { sql } from '@codemirror/lang-sql';
import { formatStreamSQL } from './sqlFormatter';
import { Copy } from '@carbon/icons-react';
import { useClipboard } from '@/hooks';
import './index.scss';

const CodeEditor: FC<any> = ({ height, code, setCode, width, readOnly = false, theme, language = 'javascript' }) => {
  const buttonRef = useRef<any>(null);
  useClipboard(buttonRef, code, 'copy success！');

  const extensions = useMemo(() => {
    if (language === 'sql') {
      try {
        const formattedCode = formatStreamSQL(code);
        if (code !== formattedCode) {
          setCode(formattedCode);
        }
      } catch (error) {
        console.warn('SQL formatting failed:', error);
      }
      return [sql()];
    }
    return [javascript({ jsx: true }), xml(), html(), sql()];
  }, [language, code]);

  return (
    <div
      style={{
        height: '100%',
        width: '100%',
        position: 'relative',

        '--cdm-c-width': width ? width + 'px' : '100%',
        '--cdm-c-height': height ? height + 'px' : '100%',
      }}
    >
      <div
        ref={buttonRef}
        style={{
          position: 'absolute',
          right: 4,
          top: 4,
          color: theme === 'light' ? '#000' : '#fff',
          zIndex: 1,
          cursor: 'pointer',
        }}
      >
        <Copy />
      </div>
      <CodeMirror
        className="cm-code-editor"
        value={code}
        height={'100%'}
        width={'100%'}
        extensions={extensions}
        onChange={(value) => {
          setCode(value);
        }}
        theme={theme || 'dark'}
        basicSetup={{ foldGutter: true }}
        readOnly={readOnly}
      />
    </div>
  );
};

export default CodeEditor;
