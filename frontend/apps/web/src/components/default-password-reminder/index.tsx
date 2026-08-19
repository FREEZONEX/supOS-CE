import { useEffect, useMemo, useState } from 'react';
import { App, Form, Input, Modal, Typography } from 'antd';
import { useTranslate } from '@/hooks';
import { userResetPwd } from '@/apis/core-api/user-manage';
import { updateForUserInfo, useBaseStore } from '@/stores/base';
import { getToken } from '@/utils/auth';
import { passwordRegex } from '@/utils/pattern';

const defaultPassword = 'tier0';

const tokenFingerprint = (token?: string) => {
  let hash = 5381;
  const value = token || '';
  for (let index = 0; index < value.length; index += 1) {
    hash = (hash * 33) ^ value.charCodeAt(index);
  }
  return (hash >>> 0).toString(36);
};

const suppressKey = (userId?: string, token?: string) =>
  userId && token ? `tier0.defaultPassword.suppress:${userId}:${tokenFingerprint(token)}` : '';

const isSuppressed = (key: string) => {
  if (!key) return false;
  try {
    return window.sessionStorage.getItem(key) === '1';
  } catch {
    return false;
  }
};

const suppressForCurrentSession = (key: string) => {
  if (!key) return;
  try {
    window.sessionStorage.setItem(key, '1');
  } catch {
    // Ignore storage failures in strict privacy modes.
  }
};

const DefaultPasswordReminder = () => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [form] = Form.useForm();
  const currentUserInfo = useBaseStore((state) => state.currentUserInfo);
  const [open, setOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const token = getToken();
  const storageKey = useMemo(() => suppressKey(currentUserInfo?.sub, token), [currentUserInfo?.sub, token]);

  useEffect(() => {
    if (currentUserInfo?.forceChangePassword && storageKey && !isSuppressed(storageKey)) {
      setOpen(true);
      return;
    }
    setOpen(false);
  }, [currentUserInfo?.forceChangePassword, storageKey]);

  const handleLater = () => {
    suppressForCurrentSession(storageKey);
    setOpen(false);
    form.resetFields();
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setSubmitting(true);
    try {
      await userResetPwd({
        oldPassword: defaultPassword,
        password: defaultPassword,
        newPassword: values.password,
      });
      message.success(formatMessage('common.settingSuccess'));
      updateForUserInfo({ forceChangePassword: false });
      setOpen(false);
      form.resetFields();
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal
      open={open}
      title={formatMessage('account.defaultPasswordTitle')}
      onCancel={handleLater}
      onOk={handleSubmit}
      okText={formatMessage('account.changePasswordNow')}
      cancelText={formatMessage('common.later')}
      confirmLoading={submitting}
      destroyOnHidden
    >
      <Typography.Paragraph>{formatMessage('account.defaultPasswordDesc')}</Typography.Paragraph>
      <Form form={form} layout="vertical">
        <Form.Item
          label={formatMessage('account.newpassWord')}
          name="password"
          rules={[
            { required: true, message: formatMessage('rule.required') },
            {
              max: 10,
              message: formatMessage('uns.labelMaxLength', {
                label: formatMessage('appGui.password'),
                length: 10,
              }),
            },
            { pattern: passwordRegex, message: formatMessage('rule.password') },
            {
              validator(_, value) {
                if (!value || value !== defaultPassword) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error(formatMessage('account.defaultPasswordSame')));
              },
            },
          ]}
        >
          <Input.Password placeholder={formatMessage('appGui.password')} autoComplete="new-password" />
        </Form.Item>
        <Form.Item
          label={formatMessage('account.confirmpassWord')}
          name="confirm_password"
          dependencies={['password']}
          rules={[
            { required: true, message: formatMessage('rule.required') },
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
          <Input.Password placeholder={formatMessage('appGui.password')} autoComplete="new-password" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default DefaultPasswordReminder;
