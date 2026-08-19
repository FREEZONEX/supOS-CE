import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router';
import { message, Input, Tabs } from 'antd';
import { Copy } from '@carbon/icons-react';
import { activateLicense, queryIdentity, importLicenseBundle, queryQuotas } from '@/apis/core-api/license';
import { LOGIN_URL } from '@/common-types/constans';
import { useTranslate } from '@/hooks';
import { copyToClipboard } from '@/utils/common';
import loginBackground from '@/assets/login/login-background.svg';
import logoDark from '@/assets/login/logo-dark.svg';
import './index.scss';

const { TextArea } = Input;

export interface LicenseActivationProps {
  licenseStatus?: 'not_activated' | 'expired';
  onActivationSuccess?: () => void;
}

const LicenseActivation: React.FC<LicenseActivationProps> = ({ licenseStatus: propStatus, onActivationSuccess }) => {
  const formatMessage = useTranslate();
  const [searchParams] = useSearchParams();

  const urlStatus = searchParams.get('status') as 'not_activated' | 'expired' | null;
  const licenseStatus = propStatus || urlStatus || 'not_activated';

  // 在线激活状态
  const [licenseKey, setLicenseKey] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  // 离线激活状态
  const [deviceToken, setDeviceToken] = useState('');
  const [bundleJson, setBundleJson] = useState('');
  const [offlineLoading, setOfflineLoading] = useState(false);
  const [offlineError, setOfflineError] = useState('');

  // 配额信息
  const [quotas, setQuotas] = useState<{
    maxClientConnections: number;
    maxUsers: number;
    maxApps: number;
    usedUsers: number;
    usedApps: number;
  } | null>(null);

  const isExpired = licenseStatus === 'expired';

  // 挂载时获取设备标识和配额信息
  useEffect(() => {
    queryIdentity()
      .then((resp) => {
        if (resp?.device_token) {
          setDeviceToken(resp.device_token);
        }
      })
      .catch((err) => {
        console.error('Failed to query identity:', err);
      });

    queryQuotas()
      .then((resp) => {
        if (resp) {
          setQuotas(resp);
        }
      })
      .catch((err) => {
        console.error('Failed to query quotas:', err);
      });
  }, []);

  const resolveActivationErrorMessage = (rawMessage?: string) => {
    const messageText = (rawMessage || '').trim();
    if (!messageText) {
      return formatMessage('common.settingFailed');
    }

    // license.error.* 开头的 key 直接走 i18n 翻译
    if (messageText.startsWith('license.error.')) {
      return formatMessage(messageText, {}, messageText);
    }

    const normalized = messageText.toLowerCase();
    if (
      messageText.includes('许可证密钥无效') ||
      normalized.includes('this license key is invalid') ||
      normalized.includes('license key not found on cloud server')
    ) {
      return formatMessage('license.invalidLicenseKey');
    }
    if (
      messageText.includes('许可证已被禁用') ||
      messageText.includes('该许可证已被禁用') ||
      normalized.includes('license has been disabled') ||
      normalized.includes('this license has been disabled')
    ) {
      return formatMessage('license.disabledLicenseKey');
    }

    return messageText;
  };

  const handleActivateSuccess = () => {
    setSuccess(true);
    message.success(formatMessage('license.activationSuccess'));
    setTimeout(() => {
      if (onActivationSuccess) {
        onActivationSuccess();
      } else {
        window.location.href = LOGIN_URL;
      }
    }, 2000);
  };

  // 在线激活
  const handleActivate = async () => {
    if (!licenseKey.trim()) {
      setError(formatMessage('rule.required'));
      return;
    }
    setError('');
    setLoading(true);
    try {
      const result = await activateLicense(licenseKey.trim());
      if (result?.success) {
        handleActivateSuccess();
      } else {
        setError(result.message || formatMessage('common.settingFailed'));
      }
    } catch (err) {
      const errInfo = err as any;
      setError(resolveActivationErrorMessage(errInfo?.msg || errInfo?.message));
    } finally {
      setLoading(false);
    }
  };

  // 离线导入
  const handleOfflineImport = async () => {
    if (!bundleJson.trim()) {
      setOfflineError(formatMessage('rule.required'));
      return;
    }
    try {
      JSON.parse(bundleJson.trim());
    } catch {
      setOfflineError(formatMessage('license.invalidJson'));
      return;
    }
    setOfflineError('');
    setOfflineLoading(true);
    try {
      const result = await importLicenseBundle(bundleJson.trim());
      if (result?.success) {
        handleActivateSuccess();
      } else {
        setOfflineError(result.message || formatMessage('common.settingFailed'));
      }
    } catch (err) {
      const errInfo = err as any;
      setOfflineError(resolveActivationErrorMessage(errInfo?.msg || errInfo?.message));
    } finally {
      setOfflineLoading(false);
    }
  };

  // 复制设备标识
  const handleCopyToken = () => {
    if (!deviceToken) return;
    copyToClipboard(deviceToken, (success) => {
      if (success) {
        message.success(formatMessage('license.copySuccess'));
        return;
      }
      message.error(formatMessage('common.settingFailed'));
    });
  };

  // 在线激活 Tab 内容
  const onlineContent = (
    <>
      <div className="form-field">
        <label className="field-label">{formatMessage('license.licenseKey')}</label>
        <Input
          size="large"
          placeholder={formatMessage('license.enterLicenseKey')}
          value={licenseKey}
          onChange={(e) => {
            setLicenseKey(e.target.value);
            setError('');
          }}
          onPressEnter={() => !loading && handleActivate()}
          disabled={loading}
          status={error ? 'error' : undefined}
        />
        {error && <div className="error-message">{error}</div>}
      </div>
      <div className="action-area">
        <button
          className={`activate-button ${loading ? 'loading' : ''}`}
          onClick={handleActivate}
          disabled={loading || !licenseKey.trim()}
        >
          {formatMessage('license.activateNow')}
        </button>
      </div>
    </>
  );

  // 离线激活 Tab 内容
  const offlineContent = (
    <>
      {/* 离线激活说明 */}
      <div className="info-banner info">
        <div className="icon">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z" />
          </svg>
        </div>
        <span className="text">{formatMessage('license.offlineInfo')}</span>
      </div>

      {/* 设备标识 */}
      <div className="form-field">
        <label className="field-label">{formatMessage('license.deviceTokenLabel')}</label>
        <div className="token-input-wrapper">
          <Input size="large" value={deviceToken} readOnly placeholder="--" />
          <button
            className="copy-button"
            onClick={handleCopyToken}
            disabled={!deviceToken}
            title={formatMessage('license.copyDeviceToken')}
          >
            <Copy size={16} />
          </button>
        </div>
      </div>

      <div className="divider" />

      {/* Bundle JSON 输入 */}
      <div className="form-field">
        <label className="field-label">{formatMessage('license.bundleJson')}</label>
        <TextArea
          rows={8}
          placeholder={formatMessage('license.pasteBundle')}
          value={bundleJson}
          onChange={(e) => {
            setBundleJson(e.target.value);
            setOfflineError('');
          }}
          disabled={offlineLoading}
          status={offlineError ? 'error' : undefined}
        />
        {offlineError && <div className="error-message">{offlineError}</div>}
      </div>

      <div className="action-area">
        <button
          className={`activate-button ${offlineLoading ? 'loading' : ''}`}
          onClick={handleOfflineImport}
          disabled={offlineLoading || !bundleJson.trim()}
        >
          {formatMessage('license.importBundle')}
        </button>
      </div>
    </>
  );

  return (
    <div className="license-activation-page">
      {/* 左侧背景区域 */}
      <div className="background-section">
        <div className="logo-container">
          <img src={logoDark} alt="TIER 0 Edge" className="logo" />
        </div>
        <div className="background-image" style={{ backgroundImage: `url(${loginBackground})` }} />
      </div>

      {/* 右侧表单区域 */}
      <div className="form-section">
        <div className="form-container">
          {success ? (
            <div className="success-message">
              <div className="success-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M20 6L9 17l-5-5" />
                </svg>
              </div>
              <h2>{formatMessage('license.activationSuccess')}</h2>
              <p>{formatMessage('license.redirecting') || 'Redirecting to login page...'}</p>
            </div>
          ) : (
            <>
              <h2 className="form-title">{formatMessage('license.activation')}</h2>

              {/* 信息/警告 Banner */}
              {isExpired && (
                <div className="info-banner warning">
                  <div className="icon">
                    <svg viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
                    </svg>
                  </div>
                  <span className="text">{formatMessage('license.expiredWarning')}</span>
                </div>
              )}

              <Tabs
                defaultActiveKey="online"
                items={[
                  {
                    key: 'online',
                    label: formatMessage('license.onlineActivation'),
                    children: onlineContent,
                  },
                  {
                    key: 'offline',
                    label: formatMessage('license.offlineActivation'),
                    children: offlineContent,
                  },
                ]}
              />

              <div className="divider" />

              {/* 配额信息展示——仅许可证过期时显示，未激活状态不展示 */}
              {licenseStatus !== 'not_activated' &&
                quotas &&
                (quotas.maxUsers > 0 || quotas.maxApps > 0 || quotas.maxClientConnections > 0) && (
                  <div className="quota-section">
                    <h3>{formatMessage('license.quotas')}</h3>
                    <div className="quota-grid">
                      <div className="quota-item">
                        <span className="quota-label">{formatMessage('license.users')}</span>
                        <span className="quota-value">
                          {quotas.usedUsers} / {quotas.maxUsers > 0 ? quotas.maxUsers : formatMessage('license.total')}
                        </span>
                      </div>
                      <div className="quota-item">
                        <span className="quota-label">{formatMessage('license.apps')}</span>
                        <span className="quota-value">
                          {quotas.usedApps} / {quotas.maxApps > 0 ? quotas.maxApps : formatMessage('license.total')}
                        </span>
                      </div>
                      <div className="quota-item">
                        <span className="quota-label">{formatMessage('license.connections')}</span>
                        <span className="quota-value">
                          {quotas.maxClientConnections > 0
                            ? quotas.maxClientConnections
                            : formatMessage('license.total')}
                        </span>
                      </div>
                    </div>
                  </div>
                )}

              {licenseStatus !== 'not_activated' &&
                quotas &&
                (quotas.maxUsers > 0 || quotas.maxApps > 0 || quotas.maxClientConnections > 0) && (
                  <div className="divider" />
                )}

              {/* 帮助区域 */}
              <div className="help-section">
                <h3>{formatMessage('license.needHelp')}</h3>
                <ul>
                  <li>
                    {formatMessage('license.contactSupport')}{' '}
                    <a href="mailto:productsupport@tier0.dev">productsupport@tier0.dev</a>
                  </li>
                  <li>
                    {formatMessage('license.visitDocumentation')}{' '}
                    <a href="https://tier0.app/" target="_blank" rel="noopener noreferrer">
                      {formatMessage('license.documentation')}
                    </a>{' '}
                    {formatMessage('license.activationGuides')}
                  </li>
                </ul>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default LicenseActivation;
