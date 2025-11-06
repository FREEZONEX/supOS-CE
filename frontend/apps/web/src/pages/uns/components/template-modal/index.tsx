import { App, Button, Divider, Drawer, Flex, Form, Input, Select, TreeSelect, InputNumber } from 'antd';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useTranslate } from '@/hooks';
import { AddAlt, Close, ConnectSource, SubtractAlt } from '@carbon/icons-react';
import { addTemplate, getTemplateDetail, getTreeData, getTypes } from '@/apis/inter-api/uns.ts';
import styles from './index.module.scss';
const { SHOW_ALL } = TreeSelect;

import type { UnsTreeNode, InitTreeDataFnType, FieldItem, SelectTreeNode } from '@/pages/uns/types';
import type { TreeStoreActions } from '../../store/types';

type SelectNodeType = { value: string };
export interface TemplateModalProps {
  successCallBack: InitTreeDataFnType;
  changeCurrentPath: (node?: UnsTreeNode) => void;
  scrollTreeNode: (id: string) => void;
  setTreeMap: TreeStoreActions['setTreeMap'];
}

const getSelectedNodes = (values: SelectNodeType[], data: SelectTreeNode[]) => {
  const _values = values?.map((item: SelectNodeType) => item.value);
  const result: FieldItem[] = [];
  const loop = (data: SelectTreeNode[]) => {
    data.forEach((item: SelectTreeNode) => {
      if (_values.includes(item.id)) {
        result.push(...(item?.fields || []));
      }
      if (item.children) {
        loop(item.children);
      }
    });
  };
  loop(data);
  return result;
};
const useTemplateModal = ({ successCallBack, changeCurrentPath, scrollTreeNode, setTreeMap }: TemplateModalProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [showSource, setSource] = useState(false);
  const [form] = Form.useForm();
  const [types, setTypes] = useState([]);
  const [loading, setLoading] = useState(false);
  const [type, setType] = useState('addTemplate');
  const [treeData, setTreeData] = useState<SelectTreeNode[]>([]);
  const treeValueList = useRef<any>([]);

  const fieldList = Form.useWatch('fields', form);

  const openTemplateModal = useCallback(
    (type: string, id?: string) => {
      setType(type);
      setOpen(true);
      if (type === 'copyTemplate' && id) {
        getTemplateDetail({ id }).then((data: any) => {
          form.setFieldsValue({
            path: data?.path + '_Copy',
            fields: data?.fields,
            description: data?.description,
          });
        });
      } else {
        form.setFieldsValue({
          fields: [{}],
        });
      }
    },
    [setType, setOpen, form]
  );

  const onClose = () => {
    setSource(false);
    setOpen(false);
    form.resetFields();
  };

  useEffect(() => {
    getTypes().then((res: any) => {
      setTypes(res?.map?.((r: string) => ({ label: r, value: r })) || []);
    });
  }, [open]);
  const onSave = async () => {
    const values = await form.validateFields();

    setLoading(true);
    addTemplate(values)
      .then((data: any) => {
        onClose();
        successCallBack?.({ queryType: 'addTemplate', newNodeKey: data }, () => {
          changeCurrentPath({ key: data, id: data, pathType: 1 });
          setTreeMap(false);
          scrollTreeNode(data);
          message.success(formatMessage('common.optsuccess'));
        });
      })
      .finally(() => {
        setLoading(false);
      });
  };
  const onSourceHandle = () => {
    setSource(true);
    getTreeData({}).then((res: any) => {
      setTreeData(res);
    });
  };

  //重复键名校验
  const validateUnique = (_: any, value: string) => {
    const values = form.getFieldValue('fields') || []; // 获取所有表单项的值
    console.log(values);
    const isDuplicate = value && values.filter((item: FieldItem) => item?.name === value).length > 1; // 检查是否有重复值

    if (isDuplicate) {
      return Promise.reject(new Error(formatMessage('uns.duplicateKeyNameTip')));
    } else {
      return Promise.resolve();
    }
  };

  const triggerNameFieldValidation = () => {
    const currentNames = form.getFieldValue('fields');
    if (!Array.isArray(currentNames)) return;

    const fieldsToValidate = currentNames.map((_, i) => ['fields', i, 'name']);
    setTimeout(() => {
      form.validateFields(fieldsToValidate).catch(() => {});
    }, 0);
  };

  // 自定义校验空格
  const validateTrim = (_: any, value: string) => {
    if (value && value.trim() === '') {
      return Promise.reject(new Error(formatMessage('common.prohibitSpacesTip')));
    }
    return Promise.resolve();
  };

  const Dom = (
    <Drawer
      rootClassName={styles['template-modal']}
      title={formatMessage(`uns.${type}`)}
      open={open}
      closable={false}
      extra={<Close size={20} onClick={onClose} style={{ cursor: 'pointer' }} />}
      style={{
        backgroundColor: 'var(--supos-header-bg-color)',
        color: 'var(--supos-text-color)',
      }}
      maskClosable={false}
      destroyOnHidden={false}
      width={680}
    >
      <div>
        <Form
          form={form}
          colon={false}
          style={{ color: 'var(--supos-text-color)' }}
          labelCol={{ span: 6 }}
          wrapperCol={{ span: 18 }}
          labelAlign="left"
          labelWrap
        >
          <Form.Item
            label={formatMessage('common.name')}
            name="name"
            rules={[
              {
                required: true,
                message: formatMessage('rule.required'),
              },
              { max: 63 },
              { validator: validateTrim },
              { pattern: /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/, message: formatMessage('uns.nameFormat') },
            ]}
          >
            <Input placeholder={formatMessage('common.name')} />
          </Form.Item>
          <Form.Item
            label={formatMessage('uns.templateDescription')}
            name="description"
            rules={[
              {
                max: 255,
                message: formatMessage('uns.labelMaxLength', {
                  label: formatMessage('uns.templateDescription'),
                  length: 255,
                }),
              },
            ]}
          >
            <Input.TextArea rows={2} placeholder={formatMessage('uns.templateDescription')} />
          </Form.Item>
          <Divider style={{ borderColor: '#c6c6c6' }} />
          <Flex className={styles['key-title']} justify="space-between" align="center">
            {formatMessage('uns.key')}

            <Button
              color="default"
              variant="filled"
              size="small"
              disabled={showSource}
              icon={<ConnectSource size={14} />}
              onClick={onSourceHandle}
              style={{
                color: !showSource ? 'var(--supos-text-color)' : 'var(--supos-select-d-color)',
                backgroundColor: !showSource ? '#C6C6C6' : 'var(--supos-uns-button-color)',
              }}
            >
              {formatMessage('uns.source')}
            </Button>
          </Flex>
          {!showSource ? (
            <Form.List name="fields">
              {(fields, { add, remove }) => (
                <>
                  {fields.map(({ key, name, ...restField }, index) => (
                    <Flex key={key} gap="8px">
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
                        style={{ width: '110px' }}
                      >
                        <Select
                          placeholder={formatMessage('uns.type')}
                          title={fieldList?.[index]?.type || formatMessage('uns.type')}
                          options={types}
                          onChange={(type) => {
                            if (type.toLowerCase() !== 'string') {
                              form.setFieldValue(['fields', key, 'maxLen'], undefined);
                            }
                          }}
                        />
                      </Form.Item>
                      <Form.Item {...restField} name={[name, 'maxLen']} wrapperCol={{ span: 24 }} style={{ flex: 1 }}>
                        <InputNumber
                          disabled={fieldList?.[key]?.type?.toLowerCase() !== 'string'}
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
                      <Form.Item {...restField} name={[name, 'remark']} wrapperCol={{ span: 24 }} style={{ flex: 1 }}>
                        <Input
                          placeholder={`${formatMessage('uns.remark')}(${formatMessage('uns.optional')})`}
                          title={
                            fieldList?.[index]?.remark ||
                            `${formatMessage('uns.remark')}(${formatMessage('uns.optional')})`
                          }
                        />
                      </Form.Item>
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
                        }}
                        disabled={fields.length === 1}
                      />
                    </Flex>
                  ))}
                  <Button
                    color="default"
                    variant="filled"
                    onClick={() => {
                      add();
                      form.setFieldValue('functions', undefined);
                    }}
                    block
                    style={{ color: 'var(--supos-text-color)', backgroundColor: 'var(--supos-uns-button-color)' }}
                    icon={<AddAlt size={20} />}
                  />
                </>
              )}
            </Form.List>
          ) : (
            <>
              <TreeSelect
                popupClassName={styles['tree-select']}
                fieldNames={{ label: 'pathName', value: 'id' }}
                rootClassName={styles['tree-select-popup']}
                showSearch
                placeholder={formatMessage('common.select')}
                style={{ width: '100%' }}
                dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
                allowClear
                multiple
                treeDefaultExpandAll
                treeCheckable
                treeCheckStrictly
                showCheckedStrategy={SHOW_ALL}
                treeNodeFilterProp="name"
                onChange={(value) => {
                  treeValueList.current = value;
                }}
                treeData={treeData}
              />

              <Flex justify="flex-end" align="center" gap={10} style={{ marginTop: 10 }}>
                <Button
                  color="default"
                  variant="filled"
                  size="small"
                  onClick={() => {
                    if (treeValueList.current?.length > 0) {
                      const newValue = getSelectedNodes(treeValueList.current, treeData);
                      form.setFieldValue('fields', [...form.getFieldValue('fields'), ...newValue]);
                      setSource(false);
                    } else {
                      message.error(formatMessage('common.select') + formatMessage('uns.source'));
                    }
                  }}
                  style={{ color: 'var(--supos-text-color)', backgroundColor: '#C6C6C6' }}
                >
                  {formatMessage('common.confirm')}
                </Button>
                <Button
                  color="default"
                  variant="filled"
                  size="small"
                  onClick={() => {
                    setSource(false);
                  }}
                  style={{ color: 'var(--supos-text-color)', backgroundColor: '#C6C6C6' }}
                >
                  {formatMessage('common.cancel')}
                </Button>
              </Flex>
            </>
          )}

          <Divider style={{ borderColor: '#c6c6c6' }} />
          <Flex justify="flex-end" align="center">
            <Button color="primary" variant="solid" size="small" onClick={onSave} loading={loading}>
              {formatMessage('common.save')}
            </Button>
          </Flex>
        </Form>
      </div>
    </Drawer>
  );
  return {
    TemplateModal: Dom,
    openTemplateModal,
  };
};

export default useTemplateModal;
