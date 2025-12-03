import { Form, Input, Space, Button } from 'antd';
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons';
import type { FC } from 'react';
import type { FieldItem } from '@/pages/uns/types.tsx';

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
  const transportType = Form.useWatch('transportType');
  const renderSseConfig = () => (
    <>
      <Form.Item name={['config', 'url']} label="URL" rules={[{ required: true, message: '请输入SSE服务器URL' }]}>
        <Input placeholder="例如: http://localhost:3000/mcp" />
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
      <Form.Item name={['config', 'url']} label="URL" rules={[{ required: true, message: '请输入URL' }]}>
        <Input placeholder="例如: http://localhost:3001" />
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
      <Form.Item name={['config', 'command']} label="命令" rules={[{ required: true, message: '请输入命令' }]}>
        <Input placeholder="例如: npx" />
      </Form.Item>

      <Form.Item label="参数" style={{ marginBottom: 0 }}>
        <Form.List name={['config', 'args']} rules={[{ validator: validateFieldsRequired }]}>
          {(fields, { add, remove }, { errors }) => (
            <>
              {fields.map(({ key, name, ...restField }) => (
                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                  <Form.Item {...restField} name={name} rules={[{ required: true, message: '请输入参数' }]}>
                    <Input placeholder="参数值" />
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} />
                </Space>
              ))}
              <Form.Item>
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  添加参数
                </Button>
              </Form.Item>
              <Form.ErrorList errors={errors} />
            </>
          )}
        </Form.List>
      </Form.Item>

      <Form.Item label="环境变量" style={{ marginBottom: 0 }}>
        <Form.List name={['config', 'env']}>
          {(fields, { add, remove }) => (
            <>
              {fields.map(({ key, name, ...restField }) => (
                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                  <Form.Item {...restField} name={[name, 'key']} rules={[{ required: true, message: '请输入变量名' }]}>
                    <Input placeholder="变量名" />
                  </Form.Item>
                  <Form.Item
                    {...restField}
                    name={[name, 'value']}
                    rules={[{ required: true, message: '请输入变量值' }]}
                  >
                    <Input placeholder="变量值" />
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} />
                </Space>
              ))}
              <Form.Item>
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  添加环境变量
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
