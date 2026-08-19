import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router';
import { Alert, Button, Form, Input, Select, Spin, message } from 'antd';
import { Copy, CheckmarkFilled } from '@carbon/icons-react';
import { cliAuthBind, getCliAuthStatus, type CliAuthBindResp } from '@/apis/core-api/cli-auth';
import { getUserInfo } from '@/apis/core-api/auth';
import { LOGIN_URL } from '@/common-types/constans';
import {
  ApiKeyPermissionBadge,
  getApiKeyPermissionLabel,
  getApiKeyPermissionOptions,
} from '@/components/api-key-permission';
import { MAX_LENGTHS } from '@/utils/limits';
import type { UserInfoProps } from '@/stores/types';
import { copyToClipboard } from '@/utils/common';
import './index.scss';

type Step = 'validating' | 'configuring_key' | 'creating_key' | 'show_key' | 'done' | 'error';

interface FormValues {
  name: string;
  permission: string;
}

const permissionOptions = getApiKeyPermissionOptions();

const buildRedirectPath = (setupCode: string) => `/cli-auth?setup=${encodeURIComponent(setupCode)}`;

const displayUserName = (user?: UserInfoProps | null) => {
  for (const value of [user?.firstName, user?.preferredUsername, user?.email, user?.sub]) {
    if (typeof value === 'string' && value.trim()) {
      return value.trim();
    }
    if (typeof value === 'number') {
      return String(value);
    }
  }
  return 'Current user';
};

const CliAuthPage = () => {
  const [searchParams] = useSearchParams();
  const setupCode = useMemo(() => (searchParams.get('setup') || '').trim(), [searchParams]);
  const [form] = Form.useForm<FormValues>();

  const [step, setStep] = useState<Step>('validating');
  const [error, setError] = useState('');
  const [user, setUser] = useState<UserInfoProps | null>(null);
  const [createdKey, setCreatedKey] = useState<CliAuthBindResp | null>(null);

  useEffect(() => {
    let canceled = false;

    const redirectToLogin = () => {
      window.location.replace(`${LOGIN_URL}?redirectUri=${encodeURIComponent(buildRedirectPath(setupCode))}`);
    };

    const init = async () => {
      if (!setupCode) {
        setError('Setup code is required.');
        setStep('error');
        return;
      }

      setStep('validating');
      try {
        const status = await getCliAuthStatus({ setupCode });
        if (canceled) {
          return;
        }
        if (status.status === 'completed') {
          setStep('done');
          return;
        }
        if (status.status === 'expired') {
          setError('This CLI authorization request has expired. Please run tier0 login again.');
          setStep('error');
          return;
        }
        if (status.status === 'denied') {
          setError('This CLI authorization request was denied.');
          setStep('error');
          return;
        }

        const currentUser = await getUserInfo();
        if (canceled) {
          return;
        }
        setUser(currentUser);
        form.setFieldsValue({
          name: `cli-${setupCode.toLowerCase()}`,
          permission: 'full',
        });
        setStep('configuring_key');
      } catch (err: any) {
        if (canceled) {
          return;
        }
        if (err?.response?.status === 401 || err?.code === 401) {
          redirectToLogin();
          return;
        }
        setError(err?.msg || err?.message || 'Failed to prepare CLI authorization.');
        setStep('error');
      }
    };

    init();
    return () => {
      canceled = true;
    };
  }, [form, setupCode]);

  const handleGenerateKey = async () => {
    const values = await form.validateFields();
    setStep('creating_key');
    setError('');
    try {
      const result = await cliAuthBind({
        setupCode,
        name: values.name.trim(),
        permission: values.permission,
      });
      setCreatedKey(result);
      setStep('show_key');
    } catch (err: any) {
      if (err?.response?.status === 401 || err?.code === 401) {
        window.location.replace(`${LOGIN_URL}?redirectUri=${encodeURIComponent(buildRedirectPath(setupCode))}`);
        return;
      }
      setError(err?.msg || err?.message || 'Failed to generate API key.');
      setStep('configuring_key');
    }
  };

  const handleCopyKey = () => {
    if (!createdKey?.apiKey) {
      return;
    }
    copyToClipboard(createdKey.apiKey, (success) => {
      if (success) {
        message.success('Copied');
        return;
      }
      message.error('Copy failed');
    });
  };

  const renderContent = () => {
    if (step === 'validating') {
      return (
        <div className="cli-auth-state">
          <Spin />
          <span>Preparing CLI authorization...</span>
        </div>
      );
    }

    if (step === 'error') {
      return (
        <div className="cli-auth-state">
          <Alert type="error" showIcon message={error || 'CLI authorization failed.'} />
          <Button type="primary" onClick={() => window.location.replace('/uns')}>
            Back to Console
          </Button>
        </div>
      );
    }

    if (step === 'done') {
      return (
        <div className="cli-auth-state">
          <CheckmarkFilled size={48} className="cli-auth-success-icon" />
          <h1>CLI authorization completed</h1>
          <p>You can return to the terminal now.</p>
        </div>
      );
    }

    if (step === 'show_key' && createdKey) {
      return (
        <>
          <h1>Your API key generated!</h1>
          <Alert
            type="warning"
            showIcon
            message={
              <span>
                Copy it now! <strong>This key is only shown once.</strong>
              </span>
            }
          />
          <div className="cli-auth-key-row">
            <Input value={createdKey.apiKey} readOnly />
            <Button icon={<Copy size={20} />} onClick={handleCopyKey} aria-label="Copy API key" />
          </div>
          <div className="cli-auth-result">
            <span>Key Name</span>
            <strong>{createdKey.name}</strong>
            <span>Permission</span>
            <ApiKeyPermissionBadge permission={createdKey.permission} />
          </div>
          <div className="cli-auth-actions">
            <Button type="primary" onClick={() => setStep('done')}>
              Done
            </Button>
          </div>
        </>
      );
    }

    return (
      <>
        <h1>Create an API Key for CLI</h1>
        <Form form={form} layout="vertical" className="cli-auth-form" requiredMark={false}>
          <Form.Item
            name="name"
            label="Key Name"
            rules={[
              { required: true, message: 'Key name is required.' },
              { max: MAX_LENGTHS.apiKeyName, message: 'Key name is too long.' },
            ]}
          >
            <Input size="large" autoFocus disabled={step === 'creating_key'} />
          </Form.Item>
          <Form.Item name="permission" label="Permission" rules={[{ required: true }]}>
            <Select
              className="api-key-permission-select"
              popupClassName="api-key-permission-select-dropdown"
              size="large"
              options={permissionOptions}
              disabled={step === 'creating_key'}
              optionRender={(option) => (
                <ApiKeyPermissionBadge permission={String(option.value)} label={option.label} />
              )}
              labelRender={(item) => (
                <ApiKeyPermissionBadge
                  permission={String(item.value)}
                  label={getApiKeyPermissionLabel(String(item.value), item.label)}
                />
              )}
            />
          </Form.Item>
        </Form>
        {error && <Alert type="error" showIcon message={error} />}
        <div className="cli-auth-actions">
          <Button type="primary" loading={step === 'creating_key'} onClick={handleGenerateKey}>
            Generate Key
          </Button>
        </div>
      </>
    );
  };

  return (
    <main className="cli-auth-page">
      <header className="cli-auth-brand" aria-label="Tier0 CLI authorization">
        <div className="cli-auth-logo" aria-label="TIER 0">
          <span className="cli-auth-logo-text">TIER</span>
          <span className="cli-auth-logo-zero" />
        </div>
        <span className="cli-auth-cli-badge">CLI</span>
      </header>
      <section className="cli-auth-card">
        {user && (
          <>
            <div className="cli-auth-user">
              <div className="cli-auth-avatar">{displayUserName(user).slice(0, 1).toUpperCase()}</div>
              <span>{displayUserName(user)}</span>
            </div>
            <div className="cli-auth-divider" />
          </>
        )}
        {renderContent()}
      </section>
    </main>
  );
};

export default CliAuthPage;
