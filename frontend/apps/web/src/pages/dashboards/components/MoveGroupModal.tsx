import { forwardRef, useImperativeHandle, useState } from 'react';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import OperationForm from '@/components/operation-form';
import { App, Flex, Form } from 'antd';
import ComButton from '@/components/com-button';
import { getGroupList, optGroup } from '@/apis/inter-api/group.ts';

export interface MoveGroupModalRef {
  onOpen: (type: number, props: any) => void;
  onClose: () => void;
}

export interface MoveGroupModalProps {
  [key: string]: any;
}

const MoveGroupModal = forwardRef<MoveGroupModalRef, MoveGroupModalProps>(({ refreshRequest }, ref) => {
  const [visible, setVisible] = useState(false);
  const [type, setType] = useState(1);
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message } = App.useApp();

  const onOpen = (type: number, props: any) => {
    form.setFieldsValue(props);
    setType(type);
    setVisible(true);
  };

  const onClose = () => {
    form.resetFields();
    setVisible(false);
  };

  const onSave = async () => {
    const value = await form.validateFields();
    return optGroup({
      bizId: value?.bizId,
      id: value?.id,
      status: true,
    }).then(() => {
      onClose?.();
      refreshRequest?.();
      message.success(formatMessage('common.optsuccess'));
    });
  };

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  return (
    <ProModal
      open={visible}
      onCancel={onClose}
      title={formatMessage('uns.moveToGroup')}
      width={500}
      styles={{
        body: {
          paddingBlockStart: 0,
        },
      }}
    >
      {() => {
        return (
          <OperationForm
            form={form}
            onCancel={onClose}
            onSave={onSave}
            formConfig={{
              layout: 'vertical',
              labelCol: { span: 24 },
              wrapperCol: { span: 124 },
            }}
            formItemOptions={[
              {
                name: 'type',
                hidden: true,
              },
              {
                name: 'bizId',
                hidden: true,
              },
              {
                name: 'id',
                type: 'Select',
                label: formatMessage('common.group'),
                properties: {
                  api: (key?: string) =>
                    getGroupList({ page: 1, pageSize: 1000, key, type }).then((res) => {
                      return (
                        res.map((item: any) => ({
                          ...item,
                          label: item.name,
                          value: item.id + '',
                        })) || []
                      );
                    }),
                  showSearch: true,
                },
              },
            ]}
            style={{ padding: 0 }}
            footer={
              <Flex gap="10px" justify="end">
                <ComButton
                  style={{
                    backgroundColor: 'var(--supos-uns-button-color)',
                    color: 'var(--supos-text-color)',
                  }}
                  color="default"
                  variant="filled"
                  onClick={onClose}
                  title={formatMessage('common.cancel')}
                >
                  {formatMessage('common.cancel')}
                </ComButton>
                <ComButton type="primary" variant="solid" onClick={onSave} title={formatMessage('common.save')}>
                  {formatMessage('common.save')}
                </ComButton>
              </Flex>
            }
          />
        );
      }}
    </ProModal>
  );
});

export default MoveGroupModal;
