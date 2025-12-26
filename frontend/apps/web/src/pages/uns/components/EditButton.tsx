import { useState, useEffect } from 'react';
import { Button, Flex, App, Form, Input, Select, InputNumber, Divider, ConfigProvider } from 'antd';
import { AddAlt, SubtractAlt } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import { getTypes, detectModel, editModel } from '@/apis/inter-api/uns';
import Icon from '@ant-design/icons';

import type { FieldItem } from '@/pages/uns/types';
import { AuthButton } from '@/components/auth';
import ProModal from '@/components/pro-modal';
import FileEdit from '@/components/svg-components/FileEdit';
import { useBaseStore } from '@/stores/base';

const EditButton = ({ modelInfo, getModel, auth, editType }: any) => {
  const { alias, dataType, fields = [], moduleId, extendFieldUsed } = modelInfo || {};
  const showMoreBtn = editType === 'file' && dataType === 1 && !moduleId;
  const { message, modal } = App.useApp();
  const [form] = Form.useForm();
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [show, setShow] = useState(false);
  const [types, setTypes] = useState([]);
  const [fieldSelected, setFieldSelected] = useState<string[]>([]);

  const fieldList = Form.useWatch('fields', form) || [];

  const { qualityName = 'quality', timestampName = 'timeStamp' } = useBaseStore((state) => state.systemInfo);

  const onClose = () => {
    setShow(false);
  };

  const editRequest = (fields: FieldItem[]) => {
    setLoading(true);
    editModel({ alias, [dataType === 8 ? 'jsonFields' : 'fields']: fields, extendFieldUsed })
      .then(() => {
        message.success(formatMessage('uns.editSuccessful'));
        setLoading(false);
        onClose();
        getModel();
      })
      .catch(() => {
        setLoading(false);
      });
  };

  //重复键名校验
  const validateUnique = (_: any, value: string) => {
    const values = form.getFieldValue('fields') || []; // 获取所有表单项的值
    const isDuplicate = value && values.filter((item: FieldItem) => item?.name === value).length > 1; // 检查是否有重复值

    if (isDuplicate) {
      return Promise.reject(new Error(formatMessage('uns.duplicateKeyNameTip')));
    } else {
      return Promise.resolve();
    }
  };

  //校验系统字段
  // const validateSystemField = (_: any, value: string) => {
  //   if (value && [qualityName, timestampName].includes(value)) {
  //     return Promise.reject(new Error(formatMessage('uns.systemFieldTip')));
  //   }
  //   return Promise.resolve();
  // };

  const triggerNameFieldValidation = () => {
    const currentNames = form.getFieldValue('fields');
    if (!Array.isArray(currentNames)) return;

    const fieldsToValidate = currentNames.map((_, i) => ['fields', i, 'name']);
    setTimeout(() => {
      form.validateFields(fieldsToValidate).catch(() => {});
    }, 0);
  };

  const moveDefaultFieldsToEnd = (arr: FieldItem[], fieldNames: string[]) => {
    // 创建副本以避免修改原始数组
    const arrCopy = [...arr];

    // 过滤出有效的字段名（非空）
    const validFieldNames = new Set(fieldNames.filter((name) => name));

    // 分别收集要保留和要移动的字段
    const [filtered, toMove] = arrCopy.reduce(
      ([keep, move], item) => {
        if (validFieldNames.has(item.name)) {
          return [keep, [...move, item]];
        } else {
          return [[...keep, item], move];
        }
      },
      [[], []] as [FieldItem[], FieldItem[]]
    );

    // 返回：未被移动的字段 + 被移动的字段
    return [...filtered, ...toMove];
  };

  const onSave = () => {
    if (fieldList.length === 0 && editType === 'template') {
      return message.error(formatMessage('uns.pleaseEnterAtLeastOneAttribute'));
    }

    form
      .validateFields()
      .then((values) => {
        const _fields = values.fields.map(
          ({ name, type, displayName, remark, maxLen, unit, upperLimit, lowerLimit, decimal, unique }: FieldItem) => {
            return dataType === 1
              ? {
                  name,
                  type,
                  displayName,
                  remark,
                  maxLen,
                  unit: fieldSelected?.includes('unit') ? unit : undefined,
                  upperLimit: fieldSelected?.includes('upperLimit') ? upperLimit : undefined,
                  lowerLimit: fieldSelected?.includes('lowerLimit') ? lowerLimit : undefined,
                  decimal: fieldSelected?.includes('decimal') ? decimal : undefined,
                }
              : {
                  name,
                  type,
                  displayName,
                  remark,
                  maxLen,
                  ...(editType === 'file' && dataType === 2 ? { unique } : {}),
                };
          }
        );
        if (dataType === 8) {
          editRequest(_fields);
        } else {
          detectModel({
            alias,
            fields: _fields,
          }).then((res: any) => {
            if (res && res.referred) {
              modal.confirm({
                content: res.tips,
                zIndex: 9001,
                onOk() {
                  editRequest(_fields);
                },
                onCancel() {},
                okButtonProps: {
                  title: formatMessage('common.confirm'),
                },
                cancelButtonProps: {
                  title: formatMessage('common.cancel'),
                },
              });
            } else {
              editRequest(_fields);
            }
          });
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  useEffect(() => {
    if (show) {
      getTypes()
        .then((res: any) => {
          setTypes(res ? res.map((type: string) => ({ label: type, value: type })) : []);
        })
        .catch((err) => {
          console.log(err);
        });
      const _fields = fields?.map((field: FieldItem) => ({ ...field, readOnly: true }));
      form.setFieldsValue({
        fields:
          editType === 'file' && dataType === 1
            ? moveDefaultFieldsToEnd(_fields, [timestampName, qualityName])
            : _fields,
      });
    }
  }, [show]);

  useEffect(() => {
    setFieldSelected(extendFieldUsed || []);
  }, [extendFieldUsed]);

  const validateFieldsRequired = (_: any, value: FieldItem[]) => {
    if (editType === 'file' && [1, 2].includes(dataType) && value?.filter((e) => !e?.systemField).length === 0) {
      return Promise.reject(new Error(formatMessage('uns.fieldsRequiredTip')));
    } else {
      return Promise.resolve();
    }
  };

  // const handleTypes = (dataType: number, types: { label: string; value: string }[]) => {
  //   switch (dataType) {
  //     case 2:
  //       return types.slice(0, 7);
  //     default:
  //       return types;
  //   }
  // };

  // const fieldSelectorContent = (
  //   <Checkbox.Group
  //     style={{
  //       display: 'flex',
  //       flexDirection: 'column',
  //       gap: 8,
  //     }}
  //     options={[
  //       { label: formatMessage('uns.unit'), value: 'unit' },
  //       { label: formatMessage('uns.upperLimit'), value: 'upperLimit' },
  //       { label: formatMessage('uns.lowerLimit'), value: 'lowerLimit' },
  //       { label: formatMessage('uns.decimal'), value: 'decimal' },
  //     ]}
  //     value={fieldSelected}
  //     onChange={(values) => {
  //       setFieldSelected(values);
  //       editModel({ alias, fields, extendFieldUsed: values }).then(() => {
  //         getModel();
  //       });
  //     }}
  //   />
  // );

  const commonStyle = {
    width: 'calc((100% - 32px) / 5)',
    marginBottom: 0,
  };
  const renderMoreField = (fieldName: string, name: number, restField: any, disabled?: boolean) => {
    switch (fieldName) {
      case 'unit':
        return (
          <Form.Item
            {...restField}
            name={[name, 'unit']}
            style={commonStyle}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            key={fieldName + name}
          >
            <Input placeholder={formatMessage('uns.unit')} disabled={disabled} maxLength={5} />
          </Form.Item>
        );
      case 'upperLimit':
        return (
          <Form.Item
            {...restField}
            name={[name, 'upperLimit']}
            style={commonStyle}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            key={fieldName + name}
          >
            <InputNumber placeholder={formatMessage('uns.upperLimit')} disabled={disabled} style={{ width: '100%' }} />
          </Form.Item>
        );
      case 'lowerLimit':
        return (
          <Form.Item
            {...restField}
            name={[name, 'lowerLimit']}
            style={commonStyle}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            key={fieldName + name}
          >
            <InputNumber placeholder={formatMessage('uns.lowerLimit')} disabled={disabled} style={{ width: '100%' }} />
          </Form.Item>
        );
      case 'decimal':
        return (
          <Form.Item
            {...restField}
            name={[name, 'decimal']}
            style={commonStyle}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            key={fieldName + name}
          >
            <InputNumber
              placeholder={formatMessage('uns.decimal')}
              disabled={disabled}
              precision={0}
              style={{ width: '100%' }}
            />
          </Form.Item>
        );
      default:
        return null;
    }
  };

  return (
    <>
      <Flex gap={8}>
        {/* {dataType === 1 && (
          <Popover
            content={fieldSelectorContent}
            trigger="click"
            getPopupContainer={(triggerNode) => triggerNode}
            placement="bottomLeft"
            // onOpenChange={(open) => {
            //   if (!open)
            //     editModel({ alias, fields, extendFieldUsed: fieldSelected }).then(() => {
            //       getModel();
            //     });
            // }}
          >
            <Button icon={<Add />} color="default" variant="filled">
              {formatMessage('uns.fieldSelector')}
            </Button>
          </Popover>
        )} */}
        <AuthButton
          auth={auth}
          onClick={() => setShow(true)}
          style={{ border: '1px solid #C6C6C6', background: 'var(--supos-uns-button-color)' }}
          icon={
            <Icon
              component={FileEdit}
              style={{
                fontSize: 17,
                color: 'var(--supos-text-color)',
              }}
            />
          }
        />
      </Flex>
      <ProModal
        title={formatMessage('common.edit')}
        onCancel={onClose}
        open={show}
        className="editModalWrap"
        width={720}
        styles={{
          content: { padding: 0 },
          header: { padding: '20px 24px 10px', margin: 0 },
          body: { padding: '0 24px 30px', margin: 0, maxHeight: 'calc(100vh - 62px)', overflowY: 'auto' },
        }}
      >
        <Form form={form} name="editModelForm" colon={false} initialValues={{ fields: fields }} disabled={loading}>
          <ConfigProvider
            theme={{
              components: {
                Form: {
                  itemMarginBottom: fieldSelected.length > 0 && showMoreBtn ? 8 : 16,
                },
              },
            }}
          >
            <Form.List name="fields" rules={[{ validator: validateFieldsRequired }]}>
              {(fields, { add, remove }, { errors }) => (
                <>
                  {fields.map(({ key, name, ...restField }, index) => {
                    const moreField = fieldSelected.length > 0 && showMoreBtn && !fieldList?.[index]?.readOnly;
                    return (
                      <div key={key}>
                        <Flex align={moreField ? 'center' : 'flex-start'} gap={8}>
                          <Flex vertical style={{ flex: 1 }}>
                            <Flex gap={8}>
                              {fieldList[index]?.readOnly ? (
                                <div
                                  className="readOnlyField"
                                  style={{
                                    flex: 1,
                                    minHeight: 32,
                                    marginBottom: 24,
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 8,
                                    borderBottom: '1px solid var(--supos-table-tr-color)',
                                    wordBreak: 'break-all',
                                  }}
                                >
                                  <span style={{ flex: 1 }}>{fieldList[index]?.name}</span>
                                  <span style={{ width: '110px' }}>{fieldList[index]?.type}</span>
                                  <span style={{ flex: 1 }}>{fieldList[index]?.maxLen}</span>
                                  <span style={{ flex: 1 }}>{fieldList[index]?.displayName}</span>
                                  <span style={{ flex: 1 }}>{fieldList[index]?.remark}</span>
                                </div>
                              ) : (
                                <>
                                  <Form.Item
                                    {...restField}
                                    name={[name, 'name']}
                                    rules={[
                                      {
                                        required: true,
                                        message: formatMessage('uns.pleaseInputKeyName'),
                                      },
                                      {
                                        pattern: /^[A-Za-z][A-Za-z0-9_]*$/,
                                        message: formatMessage('uns.keyNameFormat'),
                                      },
                                      { validator: validateUnique }, // 添加自定义校验规则
                                      {
                                        max: 63,
                                        message: formatMessage('uns.labelMaxLength', {
                                          label: formatMessage('common.name'),
                                          length: 63,
                                        }),
                                      },
                                      // { validator: validateSystemField }, // 添加自定义校验规则
                                    ]}
                                    wrapperCol={{ span: 24 }}
                                    style={{ flex: 1 }}
                                  >
                                    <Input
                                      placeholder={formatMessage('common.name')}
                                      title={fieldList?.[index]?.name || formatMessage('common.name')}
                                      onChange={triggerNameFieldValidation}
                                    />
                                  </Form.Item>
                                  <Form.Item
                                    {...restField}
                                    name={[name, 'type']}
                                    rules={[
                                      {
                                        required: true,
                                        message: formatMessage('uns.pleaseSelectKeyType'),
                                      },
                                    ]}
                                    wrapperCol={{ span: 24 }}
                                    style={{ flex: 1 }}
                                  >
                                    <Select
                                      placeholder={formatMessage('uns.type')}
                                      title={fieldList?.[index]?.type || formatMessage('uns.type')}
                                      options={types}
                                      onChange={(type) => {
                                        if (type.toLowerCase() !== 'string') {
                                          form.setFieldValue(['fields', index, 'maxLen'], undefined);
                                        }
                                      }}
                                    />
                                  </Form.Item>
                                  <Form.Item
                                    {...restField}
                                    name={[name, 'maxLen']}
                                    wrapperCol={{ span: 24 }}
                                    style={{ flex: 1 }}
                                  >
                                    <InputNumber
                                      disabled={fieldList?.[index]?.type?.toLowerCase() !== 'string'}
                                      style={{ width: '100%' }}
                                      min={1}
                                      max={10485760}
                                      step={1}
                                      precision={0}
                                      placeholder={formatMessage('common.length')}
                                      title={fieldList?.[index]?.maxLen || formatMessage('common.length')}
                                    />
                                  </Form.Item>
                                  <Form.Item
                                    {...restField}
                                    name={[name, 'displayName']}
                                    wrapperCol={{ span: 24 }}
                                    style={{ flex: 1 }}
                                  >
                                    <Input
                                      placeholder={`${formatMessage('uns.displayName')}(${formatMessage('uns.optional')})`}
                                      title={
                                        fieldList?.[index]?.displayName ||
                                        `${formatMessage('uns.displayName')}(${formatMessage('uns.optional')})`
                                      }
                                    />
                                  </Form.Item>
                                  <Form.Item
                                    {...restField}
                                    name={[name, 'remark']}
                                    wrapperCol={{ span: 24 }}
                                    style={{ flex: 1 }}
                                  >
                                    <Input
                                      placeholder={`${formatMessage('uns.remark')}(${formatMessage('uns.optional')})`}
                                      title={
                                        fieldList?.[index]?.remark ||
                                        `${formatMessage('uns.remark')}(${formatMessage('uns.optional')})`
                                      }
                                    />
                                  </Form.Item>
                                </>
                              )}
                            </Flex>
                            {moreField && (
                              <Flex gap={8}>
                                {[...fieldSelected].map((item: string) =>
                                  renderMoreField(item, name, restField, fieldList?.[index]?.systemField)
                                )}
                              </Flex>
                            )}
                          </Flex>
                          {/* 删除按钮 */}
                          {editType === 'file' &&
                          dataType === 1 &&
                          [qualityName, timestampName].includes(fieldList[index]?.name) ? (
                            <div style={{ width: 32 }} />
                          ) : (
                            <Button
                              color="default"
                              variant="filled"
                              icon={<SubtractAlt />}
                              onClick={() => {
                                remove(name);
                                triggerNameFieldValidation();
                              }}
                              style={{
                                border: '1px solid #CBD5E1',
                                color: 'var(--supos-text-color)',
                                backgroundColor: 'var(--supos-uns-button-color)',
                                height: moreField ? '70px' : '32px',
                              }}
                            />
                          )}
                        </Flex>
                        {moreField && <Divider dashed style={{ borderColor: '#c6c6c6', margin: '12px 0' }} />}
                      </div>
                    );
                  })}
                  <Button
                    color="default"
                    variant="filled"
                    onClick={() => {
                      // if (dataType === 1) {
                      //   const insertIndex = fields.length - 2 > 0 ? fields.length - 2 : 0;
                      //   add({}, insertIndex);
                      // } else {
                      //   add();
                      // }
                      add();
                    }}
                    block
                    style={{
                      color: 'var(--supos-text-color)',
                      backgroundColor: 'var(--supos-uns-button-color)',
                    }}
                    icon={<AddAlt size={20} />}
                  />
                  <Form.ErrorList errors={errors} />
                </>
              )}
            </Form.List>
          </ConfigProvider>
        </Form>
        <Button loading={loading} color="primary" variant="solid" block onClick={onSave} style={{ marginTop: 20 }}>
          {formatMessage('common.save')}
        </Button>
      </ProModal>
    </>
  );
};

export default EditButton;
