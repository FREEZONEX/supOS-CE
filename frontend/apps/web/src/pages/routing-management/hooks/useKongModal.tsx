import { useState, useCallback } from 'react';
import { Modal, Form, message } from 'antd';

interface UseKongModalParams {
  title: string;
  createApi: (data: Record<string, unknown>) => Promise<any>;
  updateApi: (id: string, data: Record<string, unknown>) => Promise<any>;
  onSuccess: () => void;
  width?: number;
  renderForm: (form: any, editingRecord: any) => React.ReactNode;
  transformValues?: (values: any, isEdit: boolean) => Record<string, unknown>;
}

const useKongModal = ({
  title,
  createApi,
  updateApi,
  onSuccess,
  width = 640,
  renderForm,
  transformValues,
}: UseKongModalParams) => {
  const [visible, setVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<any>(null);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [form] = Form.useForm();

  const open = useCallback(
    (record?: any) => {
      setEditingRecord(record ?? null);
      if (record) {
        form.setFieldsValue(record);
      } else {
        form.resetFields();
      }
      setVisible(true);
    },
    [form]
  );

  const close = useCallback(() => {
    setVisible(false);
    setEditingRecord(null);
    form.resetFields();
  }, [form]);

  const handleOk = useCallback(async () => {
    try {
      const values = await form.validateFields();
      setConfirmLoading(true);
      const isEdit = !!editingRecord?.id;
      const payload = transformValues ? transformValues(values, isEdit) : values;
      if (isEdit) {
        await updateApi(editingRecord.id, payload);
      } else {
        await createApi(payload);
      }
      message.success(isEdit ? 'Updated successfully' : 'Created successfully');
      close();
      onSuccess();
    } catch (err: any) {
      if (err?.message) {
        message.error(err.message);
      }
    } finally {
      setConfirmLoading(false);
    }
  }, [form, editingRecord, createApi, updateApi, onSuccess, close, transformValues]);

  const ModalDom = (
    <Modal
      title={editingRecord ? `Edit ${title}` : `Add ${title}`}
      open={visible}
      onOk={handleOk}
      onCancel={close}
      confirmLoading={confirmLoading}
      width={width}
      destroyOnClose
      maskClosable={false}
    >
      <Form form={form} layout="vertical" autoComplete="off" preserve={false}>
        {renderForm(form, editingRecord)}
      </Form>
    </Modal>
  );

  return { ModalDom, open, close };
};

export default useKongModal;
