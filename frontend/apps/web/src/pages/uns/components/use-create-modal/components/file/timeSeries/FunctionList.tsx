import { type FC, useEffect, useState, useRef } from 'react';
import { Form, Flex, Select, Divider, Input } from 'antd';
import { ArrowRight } from '@carbon/icons-react';
import { FUNCTION_TYPES } from '@/pages/uns/components/use-create-modal/components/file/timeSeries/FunctionConfig';
import { useTranslate } from '@/hooks';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect';
import ComFormula from '@/components/com-formula';

const FunctionList: FC<any> = () => {
  const form = Form.useFormInstance();
  const formatMessage = useTranslate();
  const [havingFieldList, setHavingFieldList] = useState([]);
  const fields = Form.useWatch('fields', form) || form.getFieldValue('fields') || [];
  const functions = Form.useWatch('functions', form) || [];
  const whereFieldList = Form.useWatch('whereFieldList', form) || [];
  const whereFormulaRef: any = useRef(null);
  const havingFormulaRef: any = useRef(null);

  const onChange = (e: any) => {
    const { fields } = e.option || {};
    const _whereFieldList = fields
      ? fields?.map(({ name, type }: any) => {
          return { label: name, value: name, type };
        })
      : [];
    form.setFieldValue('whereFieldList', _whereFieldList);
    form.setFieldValue(['streamOptions', 'window', 'options', 'field'], undefined);
    fillFunctions(true);
    resetFormula();
  };

  const resetFormula = () => {
    form.setFieldValue('whereCondition', '');
    form.setFieldValue('havingCondition', '');
    if (whereFormulaRef?.current?.restValue) {
      whereFormulaRef.current.restValue();
    }
    if (havingFormulaRef?.current?.restValue) {
      havingFormulaRef.current.restValue();
    }
  };

  const fillFunctions = (clear?: any) => {
    let _functions: any = [];
    if (clear) {
      _functions = Array(fields.length).fill({});
    } else {
      _functions = form.getFieldValue('functions') || Array(fields.length).fill({});
    }
    form.setFieldValue('functions', _functions);
  };

  useEffect(() => {
    fillFunctions();
  }, []);

  useEffect(() => {
    if (functions?.length > 0) {
      const _havingFieldList = functions
        .filter((e: any) => e.functionType && e.key)
        .map(({ functionType, key }: any) => {
          return {
            label: `${functionType}(${key})`,
            value: `${functionType}(${key})`,
          };
        });
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setHavingFieldList(_havingFieldList);
    } else {
      setHavingFieldList([]);
    }
  }, [functions]);

  useEffect(() => {
    if (!whereFormulaRef?.current) return;
    if (whereFieldList?.length > 0) {
      whereFormulaRef?.current?.setValue(form.getFieldValue('whereCondition'), whereFieldList);
    }
  }, [whereFieldList]);

  useEffect(() => {
    if (!havingFormulaRef?.current) return;
    if (havingFieldList?.length > 0) {
      havingFormulaRef?.current?.setValue(form.getFieldValue('havingCondition'), havingFieldList);
    }
  }, [havingFieldList]);

  return (
    <>
      <Form.Item hidden name="whereFieldList">
        <Input />
      </Form.Item>
      <Form.Item
        label={formatMessage('streams.dataSource')}
        name="DataSource"
        rules={[
          {
            required: true,
          },
        ]}
      >
        <SearchSelect onChange={onChange} labelInValue />
      </Form.Item>

      <Divider style={{ borderColor: '#c6c6c6' }} />
      <Form.List name="functions">
        {(items) => (
          <>
            {items.map(({ key, name }) => (
              <Flex key={key} gap={40} justify="space-between" align="center" className="functionBox">
                <div className="leftSection" style={{ flex: 1 }}>
                  <Form.Item
                    label={formatMessage('streams.functionType')}
                    name={[name, 'functionType']}
                    rules={[{ required: true }]}
                    labelCol={{ span: 12 }}
                  >
                    <Select options={FUNCTION_TYPES} onChange={resetFormula} />
                  </Form.Item>
                  <Form.Item
                    label={formatMessage('uns.key')}
                    name={[name, 'key']}
                    rules={[{ required: true }]}
                    labelCol={{ span: 12 }}
                  >
                    <Select options={whereFieldList} onChange={resetFormula} />
                  </Form.Item>
                </div>
                <ArrowRight style={{ marginTop: '-24px' }} />
                <div className="rightSection">
                  <div>
                    {formatMessage('common.name')}: {fields[name]?.name}
                  </div>
                  <div>
                    {formatMessage('uns.type')}: {fields[name]?.type}
                  </div>
                </div>
              </Flex>
            ))}
          </>
        )}
      </Form.List>
      <div className="keyBox">
        <div className="keyLabel">{formatMessage('streams.whereCondition')}</div>
        <div style={{ width: 600, marginTop: '10px' }}>
          <ComFormula
            fieldList={whereFieldList}
            formulaRef={whereFormulaRef}
            defaultOpenCalculator={false}
            // value={form.getFieldValue('whereCondition')}
            onChange={(value) => {
              form.setFieldValue('whereCondition', value);
            }}
          />
        </div>
      </div>
      {/* <Divider style={{ borderColor: '#c6c6c6' }} />
      <div className="keyBox">
        <div className="keyLabel">{formatMessage('streams.havingCondition')}</div>
        <div style={{ width: 600, marginTop: '10px' }}>
          <ComFormula
            fieldList={havingFieldList}
            formulaRef={havingFormulaRef}
            defaultOpenCalculator={false}
            // value={form.getFieldValue('havingCondition')}
            onChange={(value) => {
              form.setFieldValue('havingCondition', value);
            }}
          />
        </div>
      </div> */}
      <Divider style={{ borderColor: '#c6c6c6' }} />
    </>
  );
};
export default FunctionList;
