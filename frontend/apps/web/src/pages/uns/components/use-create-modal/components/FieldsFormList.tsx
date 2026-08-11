import { type CSSProperties, type FC, useEffect } from 'react';
import { Form, Flex, Input, Select, Button } from 'antd';
import { SubtractAlt, AddAlt } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import Icon from '@ant-design/icons';
import { getDefaultFields } from '@/pages/uns/components/CONST';
import './FieldsFormList.scss';

import type { FieldItem } from '@/pages/uns/types';
import ComPopupGuide from '@/components/com-popup-guide';
import HelpTooltip from '@/components/help-tooltip';
import MainKey from '@/components/svg-components/MainKey';
import { useBaseStore } from '@/stores/base';
import { MAX_LENGTHS } from '@/utils/limits';

const { Option } = Select;

export interface FieldsFormListProps {
  types?: string[];
  disabled?: boolean;
  isCreateFolder?: boolean;
  addNamespaceForAi?: { [key: string]: any };
  setAddNamespaceForAi?: (e: any) => void;
  showMainKey?: boolean;
  showWrap?: boolean;
  showTooltip?: boolean;
  dataTypeName?: string | (string | number)[];
  fieldsName?: string | (string | number)[];
  mainKeyName?: string | (string | number)[];
  hasDefaultVal?: boolean;
  showMoreBtn?: boolean;
  requiredFields?: boolean;
  style?: CSSProperties;
}

const FieldsFormList: FC<FieldsFormListProps> = ({
  types = [],
  disabled,
  isCreateFolder,
  addNamespaceForAi,
  setAddNamespaceForAi,
  showMainKey = true,
  showWrap = true,
  showTooltip = true,
  dataTypeName = 'dataType',
  fieldsName = 'fields',
  mainKeyName = 'mainKey',
  requiredFields = true,
  style,
}) => {
  const formatMessage = useTranslate();
  const form = Form.useFormInstance();
  const dataType = Form.useWatch(dataTypeName, form);
  const calculationType = Form.useWatch('calculationType');
  const fieldList = Form.useWatch(fieldsName, form) || [];
  const mainKey = Form.useWatch(mainKeyName, form);

  const { qualityName = '_quality', timestampName = '_timestamp' } = useBaseStore((state) => state.systemInfo);
  const defaultFields = getDefaultFields(qualityName, timestampName);

  const setMainKey = (index?: number) => {
    form.setFieldValue(mainKeyName, index);
  };

  //重复键名校验
  const validateUnique = (_: any, value: string) => {
    const values = form.getFieldValue(fieldsName) || []; // 获取所有表单项的值
    const isDuplicate = value && values.filter((item: FieldItem) => item?.name === value).length > 1; // 检查是否有重复值

    if (isDuplicate) {
      return Promise.reject(new Error(formatMessage('uns.duplicateKeyNameTip')));
    } else {
      return Promise.resolve();
    }
  };

  //校验系统字段
  // const validateSystemField = (value: string, systemField: boolean) => {
  //   if (systemField) {
  //     return Promise.resolve();
  //   } else {
  //     if (value && [qualityName, timestampName].includes(value)) {
  //       return Promise.reject(new Error(formatMessage('uns.systemFieldTip')));
  //     }
  //     return Promise.resolve();
  //   }
  // };

  //fields必填校验
  const validateFieldsRequired = (_: any, value: FieldItem[]) => {
    if ([1, 2].includes(dataType) && !isCreateFolder && value?.filter((e) => !e?.systemField).length === 0) {
      return Promise.reject(new Error(formatMessage('uns.fieldsRequiredTip')));
    } else {
      return Promise.resolve();
    }
  };

  const triggerNameFieldValidation = () => {
    const currentNames = form.getFieldValue(fieldsName) || [];
    if (!Array.isArray(currentNames)) return;

    const fieldsToValidate = currentNames
      .map((e: FieldItem, idx: number) => ({ ...e, idx }))
      .filter((m: FieldItem) => m.name)
      .map((n) => [...(Array.isArray(fieldsName) ? fieldsName : [fieldsName]), n.idx, 'name']);
    setTimeout(() => {
      form.validateFields(fieldsToValidate).catch(() => {});
    }, 0);
  };

  useEffect(() => {
    if (isCreateFolder) return;
    if (
      [1, 3].includes(dataType) &&
      Array.isArray(fieldList) &&
      JSON.stringify(fieldList.slice(-2)) !== JSON.stringify(defaultFields)
    ) {
      const removeDefaultFields = fieldList?.filter((e: FieldItem) => !e?.systemField);
      form.setFieldValue(fieldsName, [...removeDefaultFields, ...defaultFields]);
    }
    if (![1, 3].includes(dataType) && fieldList?.some((e: FieldItem) => e?.systemField)) {
      const removeDefaultFields = fieldList?.filter((e: FieldItem) => !e?.systemField);
      form.setFieldValue(fieldsName, removeDefaultFields?.length > 0 ? removeDefaultFields : [{}]);
      triggerNameFieldValidation();
    }
  }, [dataType, fieldList]);

  const defaultDisabled = (item: FieldItem) => {
    const { systemField } = item || {};
    return !isCreateFolder && [1, 3].includes(dataType) && systemField;
  };

  const handleChangeType = (type: string, index: number) => {
    if (index === mainKey && !['integer', 'long', 'string'].includes(type.toLowerCase())) {
      setMainKey(undefined);
    }
  };

  const getTypes = (dataType: number, types: string[]) => {
    switch (dataType) {
      case 3:
        return types.slice(0, 4);
      default:
        return types;
    }
  };

  const content = (
    <>
      <Flex align="center" justify="space-between" style={{ paddingBottom: '10px' }}>
        <Flex align="center" gap={8}>
          {setAddNamespaceForAi && addNamespaceForAi ? (
            <ComPopupGuide
              stepName={'fileFields'}
              steps={addNamespaceForAi?.steps}
              currentStep={addNamespaceForAi?.currentStep}
              placement="left"
              onBegin={(_, __, info) => form.setFieldsValue(info?.value)}
              onFinish={(_, nextStepName) => setAddNamespaceForAi({ ...addNamespaceForAi, currentStep: nextStepName })}
            >
              <div>{formatMessage('uns.attribute')}</div>
            </ComPopupGuide>
          ) : (
            <div>{formatMessage('uns.attribute')}</div>
          )}
          {showTooltip && <HelpTooltip title={formatMessage('uns.keyTooltip')} />}
        </Flex>
      </Flex>

      <Form.Item name={mainKeyName} hidden>
        <Input />
      </Form.Item>
      <Form.List name={fieldsName} rules={[{ validator: validateFieldsRequired }]}>
        {(fields, { add, remove }, { errors }) => (
          <>
            {fields.map(({ key, name, ...restField }, index) => (
              <div key={key}>
                <Flex align="flex-start" gap={8}>
                  <Flex gap={8} style={{ flex: 1, minWidth: 0 }}>
                    {dataType === 2 && showMainKey && (
                      <Button
                        className={mainKey === index ? 'activeKeyIndexBtn' : 'keyIndexBtn'}
                        color="default"
                        variant="filled"
                        icon={<Icon component={MainKey} />}
                        onClick={() => setMainKey(mainKey === index ? undefined : index)}
                        style={{
                          color: 'var(--ui-text-color)',
                          backgroundColor: 'var(--ui-uns-button-color)',
                        }}
                        disabled={
                          !(
                            fieldList[index]?.type &&
                            ['integer', 'long', 'string'].includes(fieldList[index]?.type?.toLowerCase())
                          )
                        }
                      />
                    )}

                    <Form.Item
                      {...restField}
                      name={[name, 'name']}
                      rules={[
                        { required: requiredFields, message: formatMessage('uns.pleaseInputKeyName') },
                        ...(!fieldList[index]?.systemField
                          ? [
                              {
                                pattern: /^[A-Za-z][A-Za-z0-9_]*$/,
                                message: formatMessage('uns.keyNameFormat'),
                              },
                            ]
                          : []),
                        { validator: validateUnique },
                        {
                          max: MAX_LENGTHS.name,
                          message: formatMessage('uns.labelMaxLength', {
                            label: formatMessage('common.name'),
                            length: MAX_LENGTHS.name,
                          }),
                        },
                      ]}
                      wrapperCol={{ span: 24 }}
                      style={{ flex: 1 }}
                    >
                      <Input
                        disabled={disabled || defaultDisabled(fieldList[index])}
                        placeholder={formatMessage('common.name')}
                        title={fieldList?.[index]?.name || formatMessage('common.name')}
                        onChange={triggerNameFieldValidation}
                      />
                    </Form.Item>

                    <Form.Item
                      {...restField}
                      name={[name, 'type']}
                      rules={[{ required: requiredFields, message: formatMessage('uns.pleaseSelectKeyType') }]}
                      wrapperCol={{ span: 24 }}
                      style={{ width: '97px' }}
                    >
                      <Select
                        disabled={disabled || defaultDisabled(fieldList[index])}
                        placeholder={formatMessage('uns.type')}
                        title={fieldList?.[index]?.type || formatMessage('uns.type')}
                        onChange={(type) => handleChangeType(type, index)}
                      >
                        {getTypes(dataType, types).map((e: string) => (
                          <Option key={e} value={e}>
                            {e}
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>

                    <Form.Item {...restField} name={[name, 'unit']} wrapperCol={{ span: 24 }} style={{ flex: 1 }}>
                      <Input
                        disabled={disabled || defaultDisabled(fieldList[index])}
                        placeholder={formatMessage('uns.unit')}
                        title={fieldList?.[index]?.unit || formatMessage('uns.unit')}
                      />
                    </Form.Item>
                  </Flex>
                  {!disabled && !(dataType === 3 && calculationType === 3) && !defaultDisabled(fieldList[index]) ? (
                    <Button
                      color="default"
                      variant="filled"
                      icon={<SubtractAlt />}
                      onClick={() => {
                        remove(name);
                        form.setFieldValue('functions', undefined);
                        if (mainKey === index) setMainKey(undefined);
                        triggerNameFieldValidation();
                      }}
                      style={{ border: '1px solid var(--ui-line-color)', flexShrink: 0, height: '32px' }}
                      disabled={fields.length === 1 && !isCreateFolder}
                    />
                  ) : dataType !== 3 && calculationType !== 3 && defaultDisabled(fieldList[index]) && !disabled ? (
                    <span style={{ width: '32px', flexShrink: 0 }} />
                  ) : null}
                </Flex>
              </div>
            ))}
            {!disabled && (dataType !== 3 || (dataType === 3 && calculationType === 4)) && (
              <Button
                color="default"
                variant="filled"
                onClick={() => {
                  if (!isCreateFolder && dataType === 1) {
                    const insertIndex = fields.length - 2 > 0 ? fields.length - 2 : 0;
                    add({}, insertIndex);
                  } else {
                    add();
                  }
                  form.setFieldValue('functions', undefined);
                }}
                block
                icon={<AddAlt size={20} />}
              />
            )}
            <Form.ErrorList errors={errors} />
          </>
        )}
      </Form.List>
    </>
  );

  return showWrap ? (
    <div className="dashedWrap" style={style}>
      {content}
    </div>
  ) : (
    content
  );
};
export default FieldsFormList;
