import { type FC, useState } from 'react';
import {
  BASIC_OPERATORS,
  CONDITION_OPERATORS,
  type FunctionConfig,
  FUNCTIONS_DATA,
  FUNCTIONS_TYPE,
} from '../../com-formula/constants';
import { Divider, Flex, Input } from 'antd';
import { Search } from '@carbon/icons-react';
import ComSelect from '../../com-select';
import { filter, flatten, get, trim, values } from 'lodash-es';
import classNames from 'classnames';
import './index.scss';
import { useTranslate } from '@/hooks';

interface operatorProps {
  operatorName: string;
  operatorValue: string;
}
export interface ControlPanel {
  basicOperators?: operatorProps[];
  conditionOperators?: operatorProps[];
  functions?: { [index: string]: FunctionConfig[] };
  onOperatorClick?: (item: any) => void;
  onFunctionNameClick?: (item: any) => void;
  defaultFuncType?: string;
  // 是否显示function详情
  showDetail?: boolean;
}

const createIntroPane = (item: any, { usage, example }: any, formatMessage: any) => {
  if (!item) return '';
  const nameReg = new RegExp('(' + item.name + ')', 'g');
  const paramReg = /\{(.+?)\}/g;
  const nameTemplate = `<span class="formulaName">$1</span>`;
  const paramTemplate = `<span class="formulaField">$1</span>`;
  return `<p>${formatMessage(item.intro).replace(nameReg, nameTemplate)}</p>
      <p>${usage}：${formatMessage(item.usage).replace(nameReg, nameTemplate)}</p>
      <p>${example}：${formatMessage(item.example).replace(nameReg, nameTemplate).replace(paramReg, paramTemplate)}</p>`;
};

const ControlPanel: FC<ControlPanel> = ({
  basicOperators = BASIC_OPERATORS,
  conditionOperators = CONDITION_OPERATORS,
  functions = FUNCTIONS_DATA,
  onOperatorClick,
  defaultFuncType = 'math',
  onFunctionNameClick,
  showDetail = true,
}) => {
  const formatMessage = useTranslate();
  const finalFunctions = flatten(values(functions));
  const [functionList, setFunctionList] = useState<any[]>(functions[defaultFuncType] || []);
  const [functionSearchKey, setFunctionSearchKey] = useState<string | undefined>();
  const [activeFunc, setActiveFunc] = useState<any>(functions[defaultFuncType][0]);

  const displayFunctions = functionSearchKey
    ? filter(finalFunctions, (item) => {
        return item.name.toLowerCase().indexOf(functionSearchKey.toLowerCase()) > -1;
      })
    : functionList;
  const onFuncTypeChange = (type: any) => {
    setFunctionList(functions[type] ?? []);
  };
  // 高亮显示name中匹配的关键词
  const highlight = (name: string, searchKey?: string) => {
    if (!searchKey) return name;
    // 转义输入的括号，避免正则表达式语法报错
    const keyword = searchKey.replace(/\(/g, '\\(').replace(/\)/g, '\\)');
    const reg = new RegExp(keyword, 'gi');
    return name.replace(reg, (matchStr) => `<span class="highlight">${matchStr}</span>`);
  };
  return (
    <div className="controlPanel">
      <div className="operationBox">
        <div className="operationTitle">{formatMessage('common.basic')}</div>
        <Flex wrap className="operationFlex">
          {basicOperators?.map((b) => {
            return (
              <div key={b.operatorValue} className="operationName" onClick={() => onOperatorClick?.(b)}>
                {b.operatorName}
              </div>
            );
          })}
        </Flex>
        <div className="operationTitle">{formatMessage('common.conditional')}</div>
        <Flex wrap className="operationFlex">
          {conditionOperators?.map((b) => {
            return (
              <div key={b.operatorValue} className="operationName" onClick={() => onOperatorClick?.(b)}>
                {b.operatorName}
              </div>
            );
          })}
        </Flex>
      </div>
      <div className="functionBox">
        <div className="operationTitle">{formatMessage('common.function')}</div>
        <div className="functionContent">
          <div className="functionSelect">
            <Input
              onKeyUp={(e: any) => setFunctionSearchKey(trim(e.target.value))}
              prefix={<Search />}
              variant="filled"
              placeholder={formatMessage('common.searchFunction')}
              allowClear
              onClear={() => setFunctionSearchKey(undefined)}
            />
            <Divider
              style={{
                margin: '8px 0',
                backgroundColor: '#C6C6C6',
              }}
            />
            <ComSelect
              variant="borderless"
              style={{ width: 140 }}
              defaultValue={defaultFuncType}
              options={FUNCTIONS_TYPE?.map?.((item: any) => {
                return {
                  label: formatMessage(item.cateName),
                  value: item.cateValue,
                };
              })}
              onChange={onFuncTypeChange}
            />
            <ul className="fieldList">
              {displayFunctions?.map?.((item) => {
                return (
                  <li
                    key={item.name}
                    onMouseEnter={() => {
                      setActiveFunc(item);
                    }}
                    className={classNames('fieldLi', { hover: item.name === activeFunc?.name })}
                    onClick={() => onFunctionNameClick?.(item)}
                  >
                    <span
                      className="fieldName"
                      dangerouslySetInnerHTML={{ __html: highlight(get(item, ['name'], ''), functionSearchKey) }}
                    />
                    <span className="fieldType" title={formatMessage(get(item, ['label'], ''))}>
                      {formatMessage(get(item, ['label'], ''))}
                    </span>
                  </li>
                );
              })}
            </ul>
          </div>
          {showDetail && (
            <div className="functionDetail">
              <div className="name">{activeFunc?.name}</div>
              <div
                className="description"
                dangerouslySetInnerHTML={{
                  __html: createIntroPane(
                    activeFunc,
                    {
                      usage: formatMessage('common.usage'),
                      example: formatMessage('common.example'),
                    },
                    formatMessage
                  ),
                }}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ControlPanel;
