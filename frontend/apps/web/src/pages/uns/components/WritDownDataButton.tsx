import { useState, useEffect } from 'react';
import { Button, Flex, App, Form, Input, Select, InputNumber, DatePicker, Tooltip } from 'antd';
import { useTranslate } from '@/hooks';
import { batchWriteFileValue } from '@/apis/inter-api/uns';
import { DocumentDownload } from '@carbon/icons-react';

import type { FieldItem } from '@/pages/uns/types';
import { AuthWrapper } from '@/components/auth';
import ProModal from '@/components/pro-modal';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';

dayjs.extend(utc);

interface CustomFieldItem extends FieldItem {
  value?: string | number | boolean;
  systemField?: boolean;
}

const WritDownData = ({ fileInfo = {}, websocketData = {}, auth }: any) => {
  const { alias, fields = [], dataType, refers } = fileInfo;
  const { data = {} } = websocketData;
  const { message } = App.useApp();
  const [form] = Form.useForm();
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [show, setShow] = useState(false);
  const [checkValue, setCheckValue] = useState(false);

  const fieldList = Form.useWatch('fields', form) || [];

  const onClose = () => {
    setShow(false);
    setLoading(false);
    setCheckValue(false);
  };

  const onSave = () => {
    form
      .validateFields()
      .then((values) => {
        const fieldsObj: any = {};
        values.fields.forEach(({ name, value, type }: CustomFieldItem) => {
          if (!['BLOB', 'LBLOB'].includes(type)) {
            fieldsObj[name] = [null, undefined].includes(value as any) ? undefined : `${value}`;
          }
        });
        const params = [
          {
            alias: dataType === 7 ? refers?.[0]?.alias : alias,
            data: fieldsObj,
          },
        ];
        setLoading(true);
        batchWriteFileValue(params)
          .then((res) => {
            if (res?.code === 200) {
              message.success(formatMessage('appGui.saveSuccess'));
              onClose();
            }
            if (res?.code === 206) {
              message.error(Object.values(res?.data?.errorFields || {}).join(', '));
            }
            setLoading(false);
          })
          .catch((err) => {
            console.error(err);
            setLoading(false);
          });
      })
      .catch((err) => {
        setLoading(false);
        console.error(err);
      });
  };

  useEffect(() => {
    if (show) {
      form.setFieldsValue({
        fields: fields
          ?.filter((e: CustomFieldItem) => !e.systemField)
          ?.map((field: CustomFieldItem) => ({
            ...field,
            value:
              field.type === 'DATETIME'
                ? data[field.name]
                  ? dayjs(data[field.name]).utc().format()
                  : undefined
                : field.type === 'STRING'
                  ? (data[field.name] || '').replace(/\0/g, '') //去除空字符
                  : data[field.name],
          })),
      });
    }
  }, [show]);

  const getFormItemRestProps = (type: string, index: number) => {
    switch (type) {
      case 'DATETIME':
        return {
          getValueProps: (value: any) => ({ value: value && dayjs(value) }),
          normalize: (value: any) => value && dayjs(value).utc().format(),
        };
      case 'STRING':
        return {
          rules: fieldList[index]?.maxLen
            ? [
                { required: checkValue, message: formatMessage('uns.pleaseInputValue') },
                {
                  max: fieldList[index]?.maxLen,
                  message: formatMessage('uns.labelMaxLength', {
                    label: formatMessage('uns.value'),
                    length: fieldList[index]?.maxLen,
                  }),
                },
              ]
            : [{ required: checkValue, message: formatMessage('uns.pleaseInputValue') }],
        };
      default:
        return null;
    }
  };

  const renderFormItemChild = (type: string) => {
    switch (type) {
      case 'INTEGER':
      case 'LONG':
        return <InputNumber style={{ width: '100%' }} precision={0} stringMode />;
      case 'FLOAT':
      case 'DOUBLE':
        return <InputNumber style={{ width: '100%' }} stringMode />;
      case 'BOOLEAN':
        return (
          <Select
            options={[
              { label: formatMessage('uns.true'), value: true },
              { label: formatMessage('uns.false'), value: false },
            ]}
          />
        );
      case 'DATETIME':
        return <DatePicker style={{ width: '100%' }} showTime />;
      case 'STRING':
        return <Input />;
      case 'BLOB':
      case 'LBLOB':
        return <Input disabled />;
      default:
        return <span />;
    }
  };

  useEffect(() => {
    setCheckValue(fieldList.every((e: CustomFieldItem) => [undefined, '', null].includes(e.value as any)));
  }, [fieldList]);

  useEffect(() => {
    const valuesNameList = fieldList.map((_: CustomFieldItem, index: number) => ['fields', index, 'value']);
    form.validateFields(valuesNameList);
  }, [checkValue, form]);

  return (
    <>
      <AuthWrapper auth={auth}>
        <Flex align="center" onClick={() => setShow(true)} style={{ cursor: 'pointer' }}>
          <DocumentDownload />
        </Flex>
      </AuthWrapper>
      <ProModal
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>{formatMessage('uns.writDownData')}</span>
          </div>
        }
        onCancel={onClose}
        open={show}
        className="writDownDataWrap"
        width={720}
        styles={{
          content: { padding: 0 },
          header: { padding: '20px 24px 10px', margin: 0 },
          body: { padding: '0 24px 30px', margin: 0, maxHeight: 'calc(100vh - 62px)', overflowY: 'auto' },
        }}
        afterClose={() => form.resetFields()}
      >
        <Form form={form} colon={false} disabled={loading}>
          <Flex
            align="center"
            gap="8px"
            style={{
              height: 40,
              background: 'var(--supos-table-head-color)',
              marginBottom: 24,
              paddingLeft: 5,
            }}
          >
            <span style={{ flex: 1 }}>{formatMessage('uns.key')}</span>
            <span style={{ width: '110px' }}>{formatMessage('uns.type')}</span>
            <span style={{ flex: 1 }}>{formatMessage('common.length')}</span>
            <span style={{ flex: 1.2 }}>{formatMessage('uns.value')}</span>
          </Flex>
          <Form.List name="fields">
            {(fields) => (
              <>
                {fields.map(({ key, name, ...restField }, index) => (
                  <Flex
                    key={key}
                    align="baseline"
                    gap="8px"
                    style={{
                      minHeight: 32,
                      marginBottom: 24,
                      borderBottom: '1px solid var(--supos-table-tr-color)',
                      wordBreak: 'break-all',
                      paddingLeft: 5,
                    }}
                  >
                    <span style={{ flex: 1 }}>{fieldList[index]?.name}</span>
                    <span style={{ width: '110px' }}>{fieldList[index]?.type}</span>
                    <span style={{ flex: 1 }}>{fieldList[index]?.maxLen}</span>
                    {['BLOB', 'LBLOB'].includes(fieldList[index]?.type) ? (
                      <Tooltip title={formatMessage('uns.writDownDataTip')}>
                        <Form.Item
                          {...restField}
                          name={[name, 'value']}
                          wrapperCol={{ span: 24 }}
                          style={{ flex: 1.2 }}
                          {...getFormItemRestProps(fieldList[index]?.type, index)}
                        >
                          {renderFormItemChild(fieldList[index]?.type)}
                        </Form.Item>
                      </Tooltip>
                    ) : (
                      <Form.Item
                        {...restField}
                        name={[name, 'value']}
                        wrapperCol={{ span: 24 }}
                        style={{ flex: 1.2 }}
                        rules={[
                          {
                            required: checkValue,
                            message: formatMessage('uns.pleaseInputValue'),
                          },
                        ]}
                        {...getFormItemRestProps(fieldList[index]?.type, index)}
                      >
                        {renderFormItemChild(fieldList[index]?.type)}
                      </Form.Item>
                    )}
                  </Flex>
                ))}
              </>
            )}
          </Form.List>
        </Form>
        <Button loading={loading} color="primary" variant="solid" block onClick={onSave} style={{ marginTop: 20 }}>
          {formatMessage('common.save')}
        </Button>
      </ProModal>
    </>
  );
};

export default WritDownData;
