import { Form, Input, Space, Button } from 'antd';
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons';
import type { FC } from 'react';
import type { FieldItem } from '@/pages/uns/types.tsx';
import { useTranslate } from '@/hooks';

export interface McpServerConfig {
  id: string;
  name: string;
  endpoint: string;
  transportType: 'sse' | 'streamable-http' | 'stdio';
  description?: string;
  enabled: boolean;
  status?: 'connected' | 'disconnected' | 'error';
  lastConnected?: string;
  config?: {
    // SSE 配置
    url?: string;
    headers?: Record<string, string>;
    // Stdio 配置
    command?: string;
    args?: string[];
    env?: Record<string, string>;
    // 通用配置
    timeout?: number;
    retryCount?: number;
  };
}

export const McpTransportForm: FC = () => {
  const formatMessage = useTranslate();
  const transportType = Form.useWatch('transportType');
  const renderSseConfig = () => (
    <>
      <Form.Item name={['config', 'url']} label="URL" rules={[{ required: true }]}>
        <Input placeholder={`${formatMessage('copilotkit.example')}: http://localhost:3000/sse`} />
      </Form.Item>
      {/*<Form.Item label="请求头" style={{ marginBottom: 0 }}>*/}
      {/*  <Form.List name={['config', 'headers']}>*/}
      {/*    {(fields, { add, remove }) => (*/}
      {/*      <>*/}
      {/*        {fields.map(({ key, name, ...restField }) => (*/}
      {/*          <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">*/}
      {/*            <Form.Item {...restField} name={[name, 'key']} rules={[{ required: true, message: '请输入键名' }]}>*/}
      {/*              <Input placeholder="键名" />*/}
      {/*            </Form.Item>*/}
      {/*            <Form.Item {...restField} name={[name, 'value']} rules={[{ required: true, message: '请输入值' }]}>*/}
      {/*              <Input placeholder="值" />*/}
      {/*            </Form.Item>*/}
      {/*            <MinusCircleOutlined onClick={() => remove(name)} />*/}
      {/*          </Space>*/}
      {/*        ))}*/}
      {/*        <Form.Item>*/}
      {/*          <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>*/}
      {/*            添加请求头*/}
      {/*          </Button>*/}
      {/*        </Form.Item>*/}
      {/*      </>*/}
      {/*    )}*/}
      {/*  </Form.List>*/}
      {/*</Form.Item>*/}
    </>
  );

  const renderStreamableHttpConfig = () => (
    <>
      <Form.Item name={['config', 'url']} label="URL" rules={[{ required: true }]}>
        <Input placeholder={`${formatMessage('copilotkit.example')}: http://localhost:3001/mcp`} />
      </Form.Item>
    </>
  );

  const validateFieldsRequired = (_: any, value: FieldItem[]) => {
    if (value?.length === 0) {
      return Promise.reject(new Error('请输入参数'));
    } else {
      return Promise.resolve();
    }
  };

  const renderStdioConfig = () => (
    <>
      <Form.Item name={['config', 'command']} label={formatMessage('common.name')} rules={[{ required: true }]}>
        <Input placeholder={`${formatMessage('copilotkit.example')}: npx`} />
      </Form.Item>

      <Form.Item label={formatMessage('copilotkit.params')} style={{ marginBottom: 0 }}>
        <Form.List name={['config', 'args']} rules={[{ validator: validateFieldsRequired }]}>
          {(fields, { add, remove }, { errors }) => (
            <>
              {fields.map(({ key, name, ...restField }) => (
                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                  <Form.Item {...restField} name={name} rules={[{ required: true }]}>
                    <Input />
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} />
                </Space>
              ))}
              <Form.Item>
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  {formatMessage('copilotkit.add') + formatMessage('copilotkit.params')}
                </Button>
              </Form.Item>
              <Form.ErrorList errors={errors} />
            </>
          )}
        </Form.List>
      </Form.Item>

      <Form.Item label={formatMessage('copilotkit.env')} style={{ marginBottom: 0 }}>
        <Form.List name={['config', 'env']}>
          {(fields, { add, remove }) => (
            <>
              {fields.map(({ key, name, ...restField }) => (
                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                  <Form.Item {...restField} name={[name, 'key']} rules={[{ required: true, message: '请输入变量名' }]}>
                    <Input />
                  </Form.Item>
                  <Form.Item {...restField} name={[name, 'value']} rules={[{ required: true }]}>
                    <Input />
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} />
                </Space>
              ))}
              <Form.Item>
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  {formatMessage('copilotkit.add') + formatMessage('copilotkit.env')}
                </Button>
              </Form.Item>
            </>
          )}
        </Form.List>
      </Form.Item>
    </>
  );

  return (
    <>
      {transportType === 'sse' && renderSseConfig()}
      {transportType === 'streamable-http' && renderStreamableHttpConfig()}
      {transportType === 'stdio' && renderStdioConfig()}
    </>
  );
};
