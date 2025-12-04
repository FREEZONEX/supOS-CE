import { Form, Input, Select } from 'antd';
import CodeMirror from '@uiw/react-codemirror';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { json } from '@codemirror/lang-json';
import { McpTransportForm } from './McpTransportForm.tsx';
import { useTranslate } from '@/hooks';
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
  const formatMessage = useTranslate();
  return type === 'json' ? (
    <Form.Item name="json" rules={[{ required: true }]}>
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
      <Form.Item name="name" label={formatMessage('common.name')} rules={[{ required: true }]}>
        <Input />
      </Form.Item>
      <Form.Item name="transportType" label={formatMessage('copilotkit.transportMode')} initialValue="stdio">
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
