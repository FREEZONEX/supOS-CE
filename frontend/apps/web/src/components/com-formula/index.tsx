import { type CSSProperties, type FC, useEffect, useImperativeHandle, useRef, useState } from 'react';
import CodeEditor, { type CodeEditorRef } from './code-editor';
import ControlPanel from './control-panel';
import FieldPanel from '../com-formula/field-panel';
import classNames from 'classnames';
import { Flex } from 'antd';
import { usePropsValue, useTranslate } from '@/hooks';
import { forEach, get } from 'lodash-es';
import codemirror from 'codemirror';
import './index.scss';
import { Calculator } from '@carbon/icons-react';
import { useUpdateEffect } from 'ahooks';
import HelpTooltip from '../help-tooltip';

export interface ComFormulaProps {
  style?: CSSProperties;
  className?: string;
  openCalculator?: boolean;
  defaultOpenCalculator?: boolean;
  onOpenCalculator?: (openCalculator: boolean) => void;
  // 公式
  value?: string;
  onChange?: (v: any, isError: boolean) => void;
  fieldList?: { label: string; value: string }[];
  formulaRef?: any;
  showDetail?: boolean;
  readonly?: boolean;
  required?: boolean;
  showTooltip?: boolean;
}

const ComFormula: FC<ComFormulaProps> = ({
  style,
  className,
  openCalculator,
  onOpenCalculator,
  defaultOpenCalculator,
  value,
  onChange,
  fieldList,
  formulaRef,
  showDetail,
  readonly,
  required,
  showTooltip,
}) => {
  const formatMessage = useTranslate();
  const codeEditorRef = useRef<CodeEditorRef>(null);
  const [code, setCodeV] = useState<string>();
  const [openCalc, setOpenCalc] = usePropsValue({
    value: openCalculator,
    onChange: onOpenCalculator,
    defaultValue: defaultOpenCalculator,
  });
  const getFormula = (e?: any) => {
    const editor: any = e ?? codeEditorRef.current?.getEditor();
    const result: string[] = [];
    let isError = false;
    const lines = editor.display.lineDiv.children;
    forEach(lines, (pre: HTMLElement, index: number) => {
      if (index > 0) {
        result.push('\n');
      }
      const list = pre.children[0].querySelectorAll('span');
      forEach(list, (item: HTMLElement) => {
        const className = item.className;
        if (!className.includes('CodeMirror-widget')) {
          if (item.getAttribute('data-id')) {
            result.push(`$${item.getAttribute('data-id')}#`);
            if (className.includes('cm-field-invalid')) {
              isError = true;
            }
          } else {
            // 过滤非法输入
            if (!className.includes('cm-negative')) {
              // 过滤编辑器产生的特殊字符
              result.push(get(item, ['innerText'], '').replace(/\u200b/g, ''));
            }
          }
        }
      });
    });
    return onChange?.(result?.join(''), isError);
  };
  const markField = ({
    className,
    label,
    dataId,
    from,
    to,
  }: {
    className: string;
    label: string;
    dataId?: string;
    from: object;
    to: object;
  }) => {
    const editor: any = codeEditorRef.current?.getEditor();
    const span = document.createElement('span');
    span.className = className;
    span.innerHTML = label;
    if (dataId) {
      span.setAttribute('data-id', dataId);
    }
    editor.markText(from, to, {
      handleMouseEvents: true,
      atomic: true,
      className,
      replacedWith: span,
    });
  };
  const setFormulaValue = (value: any, fieldList: any[] = []) => {
    const editor: any = codeEditorRef.current?.getEditor();
    const fields: Array<object> = [];
    if (value) {
      const result: Array<string> = [];
      const arr = value.split('\n');
      forEach(arr, (str, line) => {
        const fieldReg = /(\$[0-9a-zA-Z._@()]+#)/g;
        let tempStr = '';
        const arr = str.split(fieldReg);
        forEach(arr, (item) => {
          if (fieldReg.test(item)) {
            let className = 'cm-field';
            let label = '';
            const fieldItem = fieldList?.find((f) => f.value === item.replace('$', '').replace('#', ''));
            if (fieldItem) {
              label = fieldItem.label;
            } else {
              label = '已删除字段';
              className += ' cm-field-invalid';
            }
            const temp = item.replace('$', '').split('#');
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            const from = codemirror.Pos(line, tempStr.length);
            tempStr += label;
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            const to = codemirror.Pos(line, tempStr.length);
            fields.push({
              className,
              label,
              dataId: temp[0],
              from,
              to,
            });
          } else {
            tempStr += item;
          }
        });
        result.push(tempStr);
      });
      editor.setValue(result.join('\n'));
      forEach(fields, (field: any) => {
        markField(field);
      });
    } else {
      editor.setValue('');
    }
    // 清空历史记录, 避免ctrl+z回退到其他公式的编辑历史
    editor.setHistory({
      done: [],
      undone: [],
    });
  };

  useImperativeHandle(formulaRef, () => ({
    getFormula: () => getFormula?.(),
    restValue: () => codeEditorRef?.current?.getEditor()?.setValue(''),
    setValue: setFormulaValue,
  }));
  useEffect(() => {
    setFormulaValue(value, fieldList);
  }, [value, fieldList]);

  useEffect(() => {
    const editor: any = codeEditorRef.current?.getEditor();
    const fn = (instance: any, changeObj: any) => {
      if (changeObj && changeObj?.origin !== 'setValue') {
        setCodeV(instance.getValue());
      }
    };
    editor.on('change', fn);
    return () => {
      editor.off('change', fn);
    };
  }, []);

  useUpdateEffect(() => {
    getFormula();
  }, [code]);

  const markText = (label: string, className: string, dataId?: string) => {
    const editor: any = codeEditorRef.current?.getEditor();
    editor.focus();
    let from = editor.getCursor();
    const selection = editor.getSelection();
    if (selection.length) {
      from = { ...from, ch: from.ch - selection.length };
    }
    editor.replaceSelection(label);
    const to = editor.getCursor();
    markField({
      className,
      label,
      dataId,
      from,
      to,
    });
  };

  const onFieldClick = (item: any) => {
    markText(item.label, 'cm-field', item.value);
  };

  const fastInput = (operatorValue: string, offset = 0) => {
    const editor: any = codeEditorRef.current?.getEditor();
    editor.focus();
    const pos = editor.getCursor();
    editor.replaceSelection(operatorValue);
    editor.setCursor({
      line: pos.line,
      ch: pos.ch + (operatorValue.length - offset),
    });
  };

  const onOperatorClick = (item: any) => {
    fastInput(item.operatorValue);
  };

  const onFunctionNameClick = (item: any) => {
    fastInput(item.name, 1);
  };

  return (
    <div className={classNames('comFormula', className)} style={style}>
      <FieldPanel
        onClick={onFieldClick}
        fieldList={fieldList}
        style={{ display: readonly ? 'none' : 'flex' }}
        tooltip={showTooltip ? <HelpTooltip title={formatMessage('uns.variableChipsTooltip')} /> : false}
      />
      <Flex
        justify="space-between"
        className={classNames('expression', {
          'expression-readonly': readonly,
        })}
      >
        <Flex align="center" gap={8}>
          <div className={required ? 'require' : undefined}>{formatMessage('common.expression')}</div>
          {showTooltip && <HelpTooltip title={formatMessage('uns.expressionTooltip')} />}
        </Flex>
        {!readonly && (
          <Flex
            onClick={() => {
              setOpenCalc(!openCalc);
            }}
            className="calc"
            gap={4}
          >
            <Calculator size={16} />
            {formatMessage('common.calculator')}
            {showTooltip && <HelpTooltip title={formatMessage('uns.functionCalculatorTooltip')} />}
          </Flex>
        )}
      </Flex>
      <CodeEditor hiddenHead ref={codeEditorRef} readonly={readonly} />
      {!readonly && openCalc && (
        <ControlPanel
          showDetail={showDetail}
          onOperatorClick={onOperatorClick}
          onFunctionNameClick={onFunctionNameClick}
        />
      )}
    </div>
  );
};

export default ComFormula;
