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
  const [oldId, setOldId] = useState('5');
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message } = App.useApp();

  const onOpen = (type: number, props: any) => {
    form.setFieldsValue({
      ...props,
      id: props?.id ? props.id + '' : '-9999',
    });
    if (props.id) {
      setOldId(props.id);
    }
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
      id: value?.id === '-9999' ? oldId : value?.id,
      status: value?.id !== '-9999',
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
      destroyOnHidden
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
                rules: [{ required: true, message: formatMessage('common.select') }],
                properties: {
                  isRequest: visible,
                  api: (name?: string) =>
                    getGroupList({ page: 1, pageSize: 1000, name, type }).then((res) => {
                      return [
                        {
                          label: formatMessage('uns.rootDirectory'),
                          value: '-9999',
                        },
                        ...(res?.map((item: any) => ({
                          ...item,
                          label: item.name,
                          value: item.id + '',
                        })) || []),
                      ];
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
