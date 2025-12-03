import { useRef, useState } from 'react';
import { useTranslate } from '@/hooks';
import { App, Button, Form, Input } from 'antd';
import { resetPwd } from '@/apis/inter-api/user-manage';
import ProModal from '@/components/pro-modal';
import { passwordRegex } from '@/utils';

const useResetPassword = ({ onSaveBack }: any) => {
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const payload = useRef<any>(null);

  const onOpen = (data: any) => {
    payload.current = data;
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
    form.resetFields();
  };
  const onSave = async () => {
    const info = await form.validateFields();
    setLoading(true);
    resetPwd({ password: info.password, userId: payload.current?.id })
      .then(() => {
        message.success(formatMessage('common.optsuccess'));
        onClose();
        onSaveBack?.();
      })
      .finally(() => {
        setLoading(false);
      });
  };
  const Dom = (
    <ProModal
      size="xxs"
      open={open}
      onCancel={onClose}
      maskClosable={false}
      title={formatMessage('account.resetpassword')}
    >
      <Form layout="vertical" form={form}>
        <Form.Item
          label={formatMessage('account.newpassWord')}
          name="password"
          rules={[
            {
              required: true,
              message: '',
            },
            {
              max: 10,
              message: formatMessage('uns.labelMaxLength', {
                label: formatMessage('appGui.password'),
                length: 10,
              }),
            },
            { pattern: passwordRegex, message: formatMessage('rule.password') },
          ]}
        >
          <Input.Password placeholder={formatMessage('appGui.password')} />
        </Form.Item>
        <Form.Item
          label={formatMessage('account.confirmpassWord')}
          name="confirm_password"
          dependencies={['password']}
          rules={[
            {
              required: true,
              message: '',
            },
            {
              max: 10,
              message: formatMessage('uns.labelMaxLength', {
                label: formatMessage('appGui.password'),
                length: 10,
              }),
            },
            { pattern: passwordRegex, message: formatMessage('rule.password') },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error(formatMessage('account.passwordMatch')));
              },
            }),
          ]}
        >
          <Input.Password className="password" placeholder={formatMessage('appGui.password')} />
        </Form.Item>
        <Button
          onClick={onSave}
          style={{ height: 32 }}
          block
          type="primary"
          loading={loading}
          title={formatMessage('common.save')}
        >
          {formatMessage('common.save')}
        </Button>
      </Form>
    </ProModal>
  );
  return {
    ModalDom: Dom,
    onOpen,
  };
};

export default useResetPassword;
