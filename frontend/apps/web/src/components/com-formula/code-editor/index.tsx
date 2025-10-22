import { forwardRef, useEffect, useImperativeHandle, useRef } from 'react';
import { forEach, map } from 'lodash-es';
import { FUNCTIONS_DATA } from '../constants';
import codemirror from 'codemirror';
import 'codemirror/addon/hint/show-hint';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/hint/show-hint.css';
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material-darker.css';
import './index.scss';
import classNames from 'classnames';

export interface CodeEditorProp {
  formulaName?: string;
  hiddenHead?: boolean;
  readonly?: boolean;
}

export interface CodeEditorRef {
  getEditor: () => any;
}

const getKeywords = () => {
  let keywords: Array<string> = [];
  for (const i in FUNCTIONS_DATA) {
    // eslint-disable-next-line no-prototype-builtins
    if (FUNCTIONS_DATA.hasOwnProperty(i)) {
      keywords = keywords.concat(map(FUNCTIONS_DATA[i], (item) => item.name.replace(/[(|)]/g, '')));
    }
  }
  return keywords;
};

function nextUntilUnescaped(stream: any, end: any) {
  let escaped = false;
  let next = stream.next();
  while (next !== null && next !== undefined) {
    if (next === end && !escaped) {
      return false;
    }
    escaped = !escaped && next === '\\';
    next = stream.next();
  }
  return escaped;
}
const CodeEditor = forwardRef<CodeEditorRef | undefined, CodeEditorProp>(function CodeEditor(
  { formulaName, hiddenHead, readonly },
  ref
) {
  const editorContainer = useRef<HTMLDivElement>(null);
  const editor = useRef<any>(null);
  useImperativeHandle(ref, () => ({ getEditor: () => editor.current }));
  useEffect(() => {
    if (!editorContainer.current) return;
    const keywords: any = getKeywords();
    codemirror.defineMode('formula', () => {
      function zipObject(arr: Array<string>) {
        const map: any = {};
        forEach(arr, (item: string) => {
          map[item] = true;
        });
        return map;
      }
      const atomMap = zipObject(['false', 'true']);
      const keywordMap = zipObject(keywords);
      const deprecateMap = zipObject(['MAP']);

      return {
        startState: function () {
          return {
            tokens: [],
          };
        },
        token: function (stream: any) {
          if (stream.eatSpace()) {
            return null;
          }
          const ch = stream.next();

          if (ch === '"' || ch === "'") {
            nextUntilUnescaped(stream, ch);
            return 'string';
          }

          if (/[[\],()]/.test(ch)) {
            return 'bracket';
          }
          if (/[+\-*/=<>!&|]/.test(ch)) {
            return 'operator';
          }
          if (/\d/.test(ch)) {
            stream.eatWhile(/[\d.]/);
            return 'number';
          }
          stream.eatWhile(/[\w]/);
          const val = stream.current();
          // eslint-disable-next-line no-prototype-builtins
          if (atomMap.hasOwnProperty(val)) {
            return 'atom';
          }
          // eslint-disable-next-line no-prototype-builtins
          if (keywordMap.hasOwnProperty(val)) {
            return 'keyword';
          }
          // eslint-disable-next-line no-prototype-builtins
          if (deprecateMap.hasOwnProperty(val)) {
            return 'deprecate';
          }
          return 'negative';
        },
      };
    });
    codemirror.defineMIME('text/fx-formula', 'formula');
    // 挂载提示
    codemirror.registerHelper('hint', 'formula', (cm: any) => {
      const cur = cm.getCursor();
      const token = cm.getTokenAt(cur);
      if (token.end > cur.ch) {
        token.end = cur.ch;
        token.string = token.string.slice(0, cur.ch - token.start);
      }
      const list: Array<string> = [];
      const tmp = token.string.toUpperCase();
      if (tmp) {
        forEach(keywords, (item) => {
          if (item.lastIndexOf(tmp, 0) === 0 && list.indexOf(item) === -1) {
            list.push(item);
          }
        });
      }
      const data = {
        list: list.map((item) => {
          return {
            displayText: item,
            text: `${item}()`,
          };
        }),
        from: {
          line: cur.line,
          ch: token.start,
        },
        to: {
          line: cur.line,
          ch: token.end,
        },
      };
      codemirror.on(data, 'pick', () => {
        const pos = cm.getCursor();
        pos.ch -= 1;
        cm.setCursor(pos);
      });
      return data;
    });
    editor.current = codemirror(editorContainer.current, {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-expect-error
      keywords,
      theme: 'material-darker',
      textWrapping: true,
      lineWrapping: true,
      lineNumbers: false,
      matchBrackets: {
        maxScanLines: 2000,
        maxHighlightLineLength: 1000,
      },
      readOnly: readonly,
      // eslint-disable-next-line no-control-regex
      specialChars: /[\u0000-\u001f\u007f\u00ad\u200c-\u200f\u2028\u2029\ufeff]/,
      mode: 'formula',
    });
    editor.current.setSize('auto', '100%');
    const hintChange = (t: any) => {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-expect-error
      codemirror.showHint(t, codemirror.hint.formula!, {
        completeSingle: false,
        completeOnSingleClick: false,
        // closeOnPick: false,
        // closeOnUnfocus: false,
      });
    };
    editor.current.on('change', hintChange);
    return () => {
      editor.current.off('change', hintChange);
      editor.current = null;
    };
  }, []);
  return (
    <div className={classNames('formulaCodeEditor', { formulaCodeEditorReadonly: readonly })}>
      {!hiddenHead && (
        <div className="formulaHead">
          <span className="formulaName" title={formulaName}>
            {formulaName}
          </span>
          <span className="formulaEqual">=</span>
        </div>
      )}
      <div className="formulaCodeMirrorWrap" ref={editorContainer} />
    </div>
  );
});

export default CodeEditor;
