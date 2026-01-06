import { type FC, useEffect, useState, useRef } from 'react';
import { Form, Select, Button, Divider, Flex } from 'antd';
import { SubtractAlt, AddAlt } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect';

import type { FieldItem } from '@/pages/uns/types';
import ComFormula from '@/components/com-formula';
import HelpTooltip from '@/components/help-tooltip';

type ReferType = { label: string; value: string; option?: { dataType?: number } };
type ReferItemType = {
  refer: ReferType;
  option?: { fields: FieldItem };
};

interface ExpressionFormProps {
  variable?: string;
  refersName?: string;
  expressionName?: string | string[];
  timeReferenceName?: string;
  showTimeReference?: boolean;
  showTitle?: boolean;
  apiParams?: any;
}

const ExpressionForm: FC<ExpressionFormProps> = ({
  variable = 'a',
  refersName = 'refers',
  expressionName = 'expression',
  timeReferenceName = 'timeReference',
  showTimeReference = false,
  showTitle = true,
  apiParams,
}) => {
  const [formulaList, setFormulaList] = useState<any>([]);
  const [timeReferenceOptions, setTimeReferenceOptions] = useState<ReferType[]>([]);
  const form = Form.useFormInstance();
  const formatMessage = useTranslate();
  const refers = Form.useWatch(refersName);
  const expression = Form.useWatch(expressionName, form);
  const timeReference = Form.useWatch(timeReferenceName, form);
  const formulaRef: any = useRef(null);

  const onChange = (e: ReferItemType, index: number) => {
    form.setFieldValue([refersName, index, 'fields'], e?.option?.fields || []);
    form.setFieldValue([refersName, index, 'field'], undefined);
  };

  useEffect(() => {
    const _formulaList = refers?.map((_: ReferItemType, index: number) => {
      return {
        label: `${formatMessage('uns.variable')}${index + 1}`,
        value: `${variable}${index + 1}`,
      };
    });
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setFormulaList(_formulaList);
  }, [refers]);

  useEffect(() => {
    if (!formulaRef?.current) return;
    if (formulaList?.length > 0) {
      formulaRef?.current?.setValue(form.getFieldValue(expressionName), formulaList);
    }
  }, [formulaList]);

  useEffect(() => {
    if (!showTimeReference) return;
    const _timeReferenceOptions: ReferType[] = [];

    refers?.forEach((item: ReferItemType, index: number) => {
      if (
        item &&
        item?.refer?.value &&
        expression?.includes(`$${variable}${index + 1}#`) &&
        item?.refer?.option?.dataType !== 2 &&
        !_timeReferenceOptions.find((option) => option.value === item?.refer?.value)
      ) {
        _timeReferenceOptions.push(item.refer);
      }
    });
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setTimeReferenceOptions(_timeReferenceOptions);

    //选中的时间引用不在表达式或refers中则清空
    if (timeReference && !_timeReferenceOptions.some((item) => timeReference === item?.value)) {
      form.setFieldValue(timeReferenceName, undefined);
    }
  }, [refers, expression, timeReference, showTimeReference]);

  return (
    <>
      {showTitle && (
        <Flex align="center" gap={8} style={{ paddingBottom: '10px' }}>
          <div>{formatMessage('uns.key')}</div>
          <HelpTooltip title={formatMessage('uns.variablePickerTooltip')} />
        </Flex>
      )}
      <Form.List name={refersName}>
        {(fields, { add, remove }) => {
          return (
            <>
              {fields.map(({ key, name, ...restField }, index) => (
                <Flex key={key} align="flex-start" gap={8}>
                  <div style={{ lineHeight: '32px', width: 'calc((100% - 70px) / 5)' }}>
                    {formatMessage('uns.variable')}
                    {index + 1}
                  </div>
                  <Form.Item
                    {...restField}
                    name={[name, 'refer']}
                    rules={[
                      {
                        required: true,
                        message: formatMessage('uns.pleaseInputNamespace'),
                      },
                    ]}
                    wrapperCol={{ span: 24 }}
                    style={{ width: 'calc((100% - 70px) * 3 / 5)' }}
                  >
                    <SearchSelect
                      placeholder={formatMessage('uns.namespace')}
                      onChange={(e) => onChange(e, index)}
                      popupMatchSelectWidth={490}
                      labelInValue
                      apiParams={apiParams}
                    />
                  </Form.Item>
                  <span style={{ lineHeight: '32px' }}>.</span>
                  <Form.Item
                    {...restField}
                    name={[name, 'field']}
                    rules={[
                      {
                        required: true,
                        message: formatMessage('uns.pleaseSelectKeyType'),
                      },
                    ]}
                    wrapperCol={{ span: 24 }}
                    style={{ width: 'calc((100% - 70px) / 5)' }}
                  >
                    <Select
                      placeholder={formatMessage('uns.key')}
                      options={refers?.[index]?.fields || []}
                      fieldNames={{ label: 'name', value: 'name' }}
                    />
                  </Form.Item>

                  <Button
                    color="default"
                    variant="filled"
                    icon={<SubtractAlt />}
                    onClick={() => {
                      remove(name);
                      formulaRef?.current?.restValue();
                    }}
                    style={{
                      border: '1px solid #CBD5E1',
                      color: 'var(--supos-text-color)',
                      backgroundColor: 'var(--supos-uns-button-color)',
                    }}
                    disabled={fields.length === 1}
                  />
                </Flex>
              ))}

              <Button
                color="default"
                variant="filled"
                onClick={() => add()}
                block
                style={{ color: 'var(--supos-text-color)', backgroundColor: 'var(--supos-uns-button-color)' }}
                icon={<AddAlt size={20} />}
              />
            </>
          );
        }}
      </Form.List>
      <Divider style={{ borderColor: '#c6c6c6' }} />
      <Form.Item
        name={expressionName}
        rules={[
          { required: true, message: formatMessage('rule.pleaseInput', { label: formatMessage('common.expression') }) },
        ]}
        label=""
        wrapperCol={{ span: 24 }}
        validateTrigger={['onChange', 'onBlur']}
      >
        <div>
          <ComFormula
            required
            formulaRef={formulaRef}
            fieldList={formulaList}
            defaultOpenCalculator={false}
            onChange={(value) => {
              form.setFieldValue(expressionName, value);
            }}
            showTooltip
          />
        </div>
      </Form.Item>

      <Divider style={{ borderColor: '#c6c6c6' }} />
      {showTimeReference && (
        <Form.Item
          name={timeReferenceName}
          label={formatMessage('uns.reference')}
          tooltip={{ title: formatMessage('uns.timeReferenceTooltip') }}
        >
          <Select options={timeReferenceOptions} allowClear />
        </Form.Item>
      )}
    </>
  );
};
export default ExpressionForm;
