import { Form, Input, type InputProps } from 'antd';

const CodeInput = ({ disabled, ...rest }: InputProps) => {
  return (
    <Form.Item
      shouldUpdate={(prevValues, currentValues) => {
        return prevValues?.source?.routeSource !== currentValues?.source?.routeSource;
      }}
      noStyle
    >
      {({ getFieldValue }) => {
        const sourceDiabled = getFieldValue('source')?.routeSource === 2;
        return <Input disabled={disabled || sourceDiabled} {...rest} />;
      }}
    </Form.Item>
  );
};

export default CodeInput;
