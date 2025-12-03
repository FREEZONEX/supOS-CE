import { Form, Input, Select } from 'antd';
import CodeMirror from '@uiw/react-codemirror';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { json } from '@codemirror/lang-json';
import { McpTransportForm } from './McpTransportForm.tsx';
const placeholder = JSON.stringify(
  {
    mcpServers: {
      'demo-streamable': {
        url: 'http://localhost:3000/mcp',
        transportType: 'streamable-http',
      },
      'demo-sse': {
        url: 'http://localhost:3001/sse',
        transportType: 'sse',
      },
      'demo-stdio': {
        transportType: 'stdio',
        command: 'npx',
        args: ['-y', '@supos-os-edge/demo-mcp-server'],
      },
    },
  },
  null,
  2
);
const McpTypeForm = () => {
  const type = Form.useWatch('type');
  const form = Form.useFormInstance();
  return type === 'json' ? (
    <Form.Item name="json" rules={[{ required: true, message: '请输入Json配置' }]}>
      <CodeMirror
        style={{
          border: '1px solid rgb(198, 198, 198)',
          borderRadius: 4,
          padding: 16,
        }}
        placeholder={placeholder}
        onKeyDown={(e) => {
          // 监听Ctrl+P快捷键
          if (e.ctrlKey && e.key === 'p') {
            e.preventDefault();
            form.setFieldsValue({ json: placeholder });
          }
        }}
        theme={codemirrorTheme}
        height="200px"
        extensions={[json()]}
      />
    </Form.Item>
  ) : (
    <>
      <Form.Item name="name" label="服务器名称" rules={[{ required: true, message: '请输入服务器名称' }]}>
        <Input placeholder="请输入服务器名称" />
      </Form.Item>
      <Form.Item name="transportType" label="传输模式" initialValue="stdio">
        <Select
          options={[
            {
              label: 'sse',
              value: 'sse',
            },
            {
              label: 'streamable-http',
              value: 'streamable-http',
            },
            {
              label: 'stdio',
              value: 'stdio',
            },
          ]}
        />
      </Form.Item>
      <McpTransportForm />
    </>
  );
};

export default McpTypeForm;
