import type { FC } from 'react';
import { Button, Flex, Form, Input } from 'antd';
import type { FormItemProps } from 'antd';
import { AddAlt, SubtractAlt } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';

export interface ExpandedKeyFormListProps {
  formProps?: FormItemProps;
}

const ExpandedKeyFormList: FC<ExpandedKeyFormListProps> = ({ formProps }) => {
  const formatMessage = useTranslate();
  return (
    <Form.Item style={{ marginBottom: 20 }} label={formatMessage('uns.expandedInformation')} {...formProps}>
      <Form.List name="extend">
        {(fields, { add, remove }) => (
          <>
            {fields.map(({ key, name, ...restField }) => (
              <Flex key={key} gap="8px">
                <Form.Item
                  {...restField}
                  name={[name, 'key']}
                  wrapperCol={{ span: 24 }}
                  style={{ flex: 1 }}
                  rules={[
                    { required: true, message: formatMessage('uns.pleaseInputKeyName') },
                    {
                      max: 32,
                      message: formatMessage('uns.labelMaxLength', { label: formatMessage('uns.key'), length: 32 }),
                    },
                  ]}
                >
                  <Input placeholder={formatMessage('uns.key')} />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'value']}
                  wrapperCol={{ span: 24 }}
                  style={{ flex: 1 }}
                  rules={[
                    { required: true, message: formatMessage('uns.pleaseInputValue') },
                    {
                      max: 128,
                      message: formatMessage('uns.labelMaxLength', { label: formatMessage('uns.value'), length: 128 }),
                    },
                  ]}
                >
                  <Input placeholder={formatMessage('uns.value')} />
                </Form.Item>
                <Button
                  color="default"
                  variant="filled"
                  icon={<SubtractAlt />}
                  onClick={() => {
                    remove(name);
                  }}
                  style={{
                    border: '1px solid #CBD5E1',
                    color: 'var(--supos-text-color)',
                    backgroundColor: 'var(--supos-uns-button-color)',
                  }}
                />
              </Flex>
            ))}
            {fields?.length < 3 && (
              <Button
                color="default"
                variant="filled"
                onClick={() => {
                  add();
                }}
                block
                style={{
                  color: 'var(--supos-text-color)',
                  backgroundColor: 'var(--supos-uns-button-color)',
                }}
                icon={<AddAlt size={20} />}
              />
            )}
          </>
        )}
      </Form.List>
    </Form.Item>
  );
};

export default ExpandedKeyFormList;
