import { createContext, type FC, type PropsWithChildren, useContext, useEffect, useRef, useState } from 'react';
import { Form, type FormInstance, Input, type InputRef } from 'antd';

const EditableContext = createContext<FormInstance<any> | null>(null);

export const EditableRow: FC = ({ ...props }) => {
  const [form] = Form.useForm();
  return (
    <Form form={form} component={false}>
      <EditableContext.Provider value={form}>
        <tr {...props} />
      </EditableContext.Provider>
    </Form>
  );
};

interface EditableCellProps {
  title: string;
  editable: boolean;
  dataIndex: any;
  record: any;
  handleSave: (record: any) => void;
}

export const EditableCell: FC<PropsWithChildren<EditableCellProps>> = ({
  editable,
  children,
  dataIndex,
  record,
  handleSave,
  ...restProps
}) => {
  const [editing, setEditing] = useState(false);
  const inputRef = useRef<InputRef>(null);
  const form = useContext(EditableContext)!;
  const _dataIndex = Array.isArray(dataIndex) ? dataIndex[1] : dataIndex;
  useEffect(() => {
    if (editing) {
      inputRef.current?.focus();
    }
  }, [editing]);

  const toggleEdit = () => {
    setEditing(!editing);
    const fieldValue = Array.isArray(dataIndex)
      ? dataIndex.reduce((obj, key) => obj?.[key], record)
      : record[dataIndex];
    form.setFieldsValue({ [_dataIndex]: fieldValue });
  };

  const save = async () => {
    try {
      const values = await form.validateFields();
      toggleEdit();
      handleSave({ ...record, values: { ...record.values, ...values } });
    } catch (errInfo) {
      console.log('Save failed:', errInfo);
    }
  };

  let childNode = children;

  if (editable) {
    childNode = editing ? (
      <Form.Item style={{ margin: 0 }} name={_dataIndex}>
        <Input ref={inputRef} onPressEnter={save} onBlur={save} />
      </Form.Item>
    ) : (
      <div className="editable-cell-value-wrap" style={{ paddingInlineEnd: 24 }} onClick={toggleEdit}>
        {children}
      </div>
    );
  }

  return (
    <td
      {...restProps}
      title={Array.isArray(dataIndex) ? dataIndex.reduce((obj, key) => obj?.[key], record) : record?.[dataIndex]}
    >
      {childNode}
    </td>
  );
};
