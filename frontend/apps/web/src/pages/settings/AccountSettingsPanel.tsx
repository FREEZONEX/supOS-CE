import { useEffect, useMemo, useRef, useState } from 'react';
import { App, Button, Form, Input, Modal, Tabs, Alert } from 'antd';
import { Copy, Edit, Password, UserAvatar } from '@carbon/icons-react';
import ComSelect from '@/components/com-select';
import HomePageSelect from '@/components/home-page-select';
import { setHomePageApi, updateCurrentUserProfile, userResetPwd } from '@/apis/core-api/user-manage';
import { logoutApi } from '@/apis/core-api/auth';
import {
  queryIdentity,
  queryLicenseStatus,
  replaceLicenseBundle,
  replaceLicenseKey,
  type LicenseStatusResp,
} from '@/apis/core-api/license';
import { LOGIN_URL, OMC_MODEL, APP_USER_GUIDE_ROUTES, APP_USER_TIPS_ENABLE } from '@/common-types/constans';
import { useTranslate } from '@/hooks';
import { fetchBaseStore, fetchSystemInfo, updateForUserInfo, useBaseStore } from '@/stores/base';
import { initI18n, useI18nStore } from '@/stores/i18n-store.ts';
import { removeToken } from '@/utils/auth';
import { copyToClipboard } from '@/utils/common';
import { passwordRegex, validNameRegex } from '@/utils/pattern';
import { preloadPluginLang } from '@/utils/plugin.ts';
import { storageOpt } from '@/utils/storage';
import { updatePersonConfigApi } from '@/apis/core-api/uns';
import { isLaunchpadStandalonePort } from '@/utils/launchpad-site';
import Cookies from 'js-cookie';
import styles from './account-settings-panel.module.scss';

type SettingsTab = 'profile' | 'preferences' | 'security' | 'license';

type AccountSettingsPanelProps = {
  activeTab?: SettingsTab;
};

const logout = (path?: string) => {
  logoutApi().then(() => {
    removeToken();
    Cookies.remove(OMC_MODEL, { path: '/' });
    storageOpt.remove(APP_USER_GUIDE_ROUTES);
    storageOpt.remove(APP_USER_TIPS_ENABLE);
    storageOpt.remove('personInfo');
    location.href = path || LOGIN_URL;
  });
};

const maskLicenseKey = (licenseKey?: string) => {
  if (!licenseKey) return '----';
  if (licenseKey.includes('*')) return licenseKey;
  if (licenseKey.length <= 8) return licenseKey;
  return `${licenseKey.slice(0, 5)}****${licenseKey.slice(-4)}`;
};

const AccountSettingsPanel = ({ activeTab = 'profile' }: AccountSettingsPanelProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const { currentUserInfo, systemInfo, menuGroup, pluginList } = useBaseStore((state) => ({
    currentUserInfo: state.currentUserInfo,
    systemInfo: state.systemInfo,
    pluginList: state.pluginList,
    menuGroup: state.menuGroup,
  }));
  const { lang, langList } = useI18nStore((state) => ({
    lang: state.lang,
    langList: state.langList
      ?.filter((f) => f.hasUsed)
      ?.map((item: any) => ({ ...item, value: item?.value?.replace?.('_', '-') })),
  }));
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [homeForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [licenseLoading, setLicenseLoading] = useState(false);
  const [isEditingProfile, setIsEditingProfile] = useState(false);
  const [isEditingPassword, setIsEditingPassword] = useState(false);
  const [isEditingLicense, setIsEditingLicense] = useState(false);
  const [licenseDeviceToken, setLicenseDeviceToken] = useState('');
  const [licenseActiveTab, setLicenseActiveTab] = useState('online');
  const [licenseKeyInput, setLicenseKeyInput] = useState('');
  const [bundleJsonInput, setBundleJsonInput] = useState('');
  const [licenseError, setLicenseError] = useState('');
  const preferencesRoutesLoadedRef = useRef(false);
  const [licenseInfo, setLicenseInfo] = useState<{
    licenseKey?: string;
    status?: 'active' | 'expiring_soon' | 'expired' | 'unknown';
    validUntil?: string;
    daysRemaining?: number;
  }>({});
  const name = currentUserInfo.firstName || currentUserInfo.preferredUsername || '';
  const avatarInitial = name.trim().slice(0, 1).toLocaleUpperCase();
  const canEditAccount = currentUserInfo?.sub !== 'guest';
  const showHomePagePreference = !isLaunchpadStandalonePort();
  const menuGroupNoSub = useMemo(() => menuGroup?.filter((item: any) => !item.subMenu) || [], [menuGroup]);

  const resetProfileForm = () => {
    profileForm.setFieldsValue({
      firstName: currentUserInfo?.firstName,
      phone: currentUserInfo?.phone,
      email: currentUserInfo?.email,
    });
  };

  const fetchLicenseStatus = async () => {
    setLicenseLoading(true);
    try {
      const data: LicenseStatusResp = await queryLicenseStatus();
      const phase = data?.phase;
      const isActivated = data?.activated;
      let status: 'active' | 'expiring_soon' | 'expired' | 'unknown' = 'unknown';
      if (!isActivated || phase === 'none') {
        status = 'unknown';
      } else if (phase === 'warning') {
        status = 'expiring_soon';
      } else if (phase === 'grace' || phase === 'hard_expired') {
        status = 'expired';
      } else {
        status = 'active';
      }

      setLicenseInfo({
        licenseKey: data?.licenseKey,
        validUntil: data?.expireAt,
        daysRemaining: data?.daysLeft,
        status,
      });
    } finally {
      setLicenseLoading(false);
    }
  };

  const fetchDeviceToken = async () => {
    const resp = await queryIdentity();
    if (resp?.device_token) {
      setLicenseDeviceToken(resp.device_token);
    }
  };

  const resolveUpdateLicenseKeyErrorMessage = (rawMessage?: string) => {
    const messageText = (rawMessage || '').trim();
    if (!messageText) {
      return formatMessage('license.updateLicenseKeyFailed');
    }

    // license.error.* 开头的 key 直接走 i18n 翻译
    if (messageText.startsWith('license.error.')) {
      return formatMessage(messageText, {}, messageText);
    }

    const normalized = messageText.toLowerCase();
    if (
      messageText.includes('License激活被服务端拒绝') ||
      messageText.includes('License 激活被服务端拒绝') ||
      messageText.includes('许可证密钥未通过服务端校验') ||
      normalized.includes('license activation rejected by server') ||
      normalized.includes('license key was rejected by the server')
    ) {
      return formatMessage('license.updateLicenseKeyRejectedDetail');
    }
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

  const onSaveProfile = async () => {
    const info = await profileForm.validateFields();
    setLoading(true);
    updateCurrentUserProfile(info)
      .then(() => {
        message.success(formatMessage('common.settingSuccess'));
        updateForUserInfo({
          ...info,
        });
        setIsEditingProfile(false);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const onSavePassword = async () => {
    const info = await passwordForm.validateFields();
    setLoading(true);
    userResetPwd({
      newPassword: info.password,
      password: info.oldPassword,
      userId: currentUserInfo?.sub,
      username: currentUserInfo?.preferredUsername,
    })
      .then(() => {
        message.success(formatMessage('common.settingSuccess'));
        setIsEditingPassword(false);
        passwordForm.resetFields();
        logout(systemInfo.loginPath);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const onChangeLanguage = (value: string) => {
    if (!currentUserInfo?.sub) return;
    updatePersonConfigApi({ userId: currentUserInfo.sub, mainLanguage: value }).then(async () => {
      const pluginLang = await preloadPluginLang(
        pluginList
          ?.filter((item: any) => item.installStatus === 'installed')
          ?.filter((item: any) => item?.plugInfoYml?.route?.name)
          ?.map((item: any) => ({ name: `/${item?.plugInfoYml?.route?.name}`, backendName: item?.name })) || [],
        value
      );
      await initI18n(value, pluginLang);
      await fetchSystemInfo(true);
      message.success(formatMessage('common.settingSuccess'));
    });
  };

  const onChangeHomePage = (value?: string) => {
    if (!value) return;
    setHomePageApi({ homePage: value })?.then((config: any) => {
      const homePage = config?.homePage || value;
      homeForm.setFieldValue('homePage', homePage);
      updateForUserInfo({ homePage });
      message.success(formatMessage('common.settingSuccess'));
    });
  };

  const resetLicenseModal = () => {
    setIsEditingLicense(false);
    setLicenseActiveTab('online');
    setLicenseKeyInput('');
    setBundleJsonInput('');
    setLicenseError('');
  };

  const onSaveLicense = async () => {
    setLicenseError('');
    if (licenseActiveTab === 'online') {
      if (!licenseKeyInput.trim()) {
        setLicenseError(formatMessage('rule.required'));
        return;
      }
    } else {
      const trimmed = bundleJsonInput.trim();
      if (!trimmed) {
        setLicenseError(formatMessage('rule.required'));
        return;
      }
      try {
        JSON.parse(trimmed);
      } catch {
        setLicenseError(formatMessage('license.invalidJson'));
        return;
      }
    }

    setLicenseLoading(true);
    try {
      if (licenseActiveTab === 'online') {
        await replaceLicenseKey(licenseKeyInput.trim());
      } else {
        await replaceLicenseBundle(bundleJsonInput.trim());
      }
      message.success(formatMessage('license.updateLicenseKeySuccess'));
      resetLicenseModal();
      void fetchLicenseStatus;
    } catch (error: any) {
      message.error(resolveUpdateLicenseKeyErrorMessage(error?.msg));
    } finally {
      setLicenseLoading(false);
    }
  };

  const getLicenseStatusInfo = (status?: string) => {
    switch (status) {
      case 'active':
        return formatMessage('license.status.active');
      case 'expiring_soon':
        return formatMessage('license.status.expiring');
      case 'expired':
        return formatMessage('license.status.expired');
      default:
        return formatMessage('license.status.unknown');
    }
  };

  const getLicenseStatusClassName = (status?: string) => {
    switch (status) {
      case 'active':
        return styles['status-active'];
      case 'expiring_soon':
        return styles['status-warning'];
      case 'expired':
        return styles['status-expired'];
      default:
        return styles['status-unknown'];
    }
  };

  useEffect(() => {
    resetProfileForm();
  }, [currentUserInfo?.email, currentUserInfo?.firstName, currentUserInfo?.phone]);

  useEffect(() => {
    if (activeTab !== 'preferences') return;
    if (!showHomePagePreference) return;
    if (preferencesRoutesLoadedRef.current) return;
    preferencesRoutesLoadedRef.current = true;
    fetchBaseStore?.();
  }, [activeTab, showHomePagePreference]);

  useEffect(() => {
    if (activeTab !== 'preferences') return;
    if (!showHomePagePreference) return;
    homeForm.setFieldValue('homePage', currentUserInfo?.homePage);
  }, [activeTab, currentUserInfo?.homePage, homeForm, showHomePagePreference]);

  useEffect(() => {
    void fetchLicenseStatus;
    void fetchDeviceToken;
  }, []);

  if (activeTab === 'preferences') {
    return (
      <div className={styles['settings-panel']}>
        <div className={styles['panel-header']}>
          <h1 className={styles['panel-title']}>{formatMessage('settings.preferences')}</h1>
          <div className={styles['panel-subtitle']}>{formatMessage('settings.preferencesDesc')}</div>
        </div>
        <div className={styles.section}>
          <div className={styles.row}>
            <div className={styles['row-label']}>{formatMessage('common.language')}</div>
            <div className={styles['row-value']}>
              <ComSelect
                className={styles.select}
                disabled={!currentUserInfo?.sub}
                value={lang}
                options={langList}
                onChange={onChangeLanguage}
              />
            </div>
          </div>
          {showHomePagePreference ? (
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.homePage')}</div>
              <div className={styles['row-value']}>
                <Form form={homeForm} className={`${styles['select-form']} ${styles['home-page-select-form']}`}>
                  <Form.Item name="homePage" noStyle>
                    <HomePageSelect
                      className={`${styles.select} ${styles['home-page-select']}`}
                      enabled={activeTab === 'preferences'}
                      resources={menuGroupNoSub}
                      placeholder={formatMessage('common.searchPage')}
                      showSearch
                      onChange={onChangeHomePage}
                    />
                  </Form.Item>
                </Form>
              </div>
            </div>
          ) : null}
        </div>
      </div>
    );
  }

  if (activeTab === 'security') {
    return (
      <div className={styles['settings-panel']}>
        <div className={styles['panel-header']}>
          <h1 className={styles['panel-title']}>{formatMessage('settings.security')}</h1>
          <div className={styles['panel-subtitle']}>{formatMessage('settings.securityDesc')}</div>
        </div>
        <div className={styles.section}>
          <h2 className={styles['section-title']}>{formatMessage('settings.login')}</h2>
          {currentUserInfo?.forceChangePassword ? (
            <Alert
              type="warning"
              showIcon
              className={styles['default-password-alert']}
              message={formatMessage('account.defaultPasswordStillInUse')}
              action={
                <Button size="small" type="primary" onClick={() => setIsEditingPassword(true)}>
                  {formatMessage('account.updatePassword')}
                </Button>
              }
            />
          ) : null}
          <div className={styles.row}>
            <div className={styles['row-label']}>{formatMessage('common.password')}</div>
            <div className={styles['row-value']}>
              <Button
                icon={<Password size={14} />}
                disabled={!canEditAccount}
                onClick={() => setIsEditingPassword(true)}
              >
                {formatMessage('account.updatePassword')}
              </Button>
            </div>
          </div>
          <div className={styles.row}>
            <div className={styles['row-label']}>{formatMessage('common.logout')}</div>
            <div className={styles['row-value']}>
              <Button onClick={() => logout(systemInfo.loginPath)}>{formatMessage('common.logout')}</Button>
            </div>
          </div>
        </div>
        <Modal
          title={formatMessage('account.updatePassword')}
          open={isEditingPassword}
          onCancel={() => {
            setIsEditingPassword(false);
            passwordForm.resetFields();
          }}
          onOk={onSavePassword}
          confirmLoading={loading}
          okButtonProps={{ disabled: !canEditAccount }}
          destroyOnHidden
        >
          <Form form={passwordForm} layout="vertical">
            <Form.Item
              label={formatMessage('account.oldPassWord')}
              name="oldPassword"
              rules={[
                { required: true, message: '' },
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
              <Input.Password placeholder={formatMessage('appGui.password')} disabled={!canEditAccount} />
            </Form.Item>
            <Form.Item
              label={formatMessage('account.newpassWord')}
              name="password"
              dependencies={['oldPassword']}
              rules={[
                { required: true, message: '' },
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
                    if (!value || getFieldValue('oldPassword') !== value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error(formatMessage('account.passwordSame')));
                  },
                }),
              ]}
            >
              <Input.Password placeholder={formatMessage('appGui.password')} disabled={!canEditAccount} />
            </Form.Item>
            <Form.Item
              label={formatMessage('account.confirmpassWord')}
              name="confirm_password"
              dependencies={['password']}
              rules={[
                { required: true, message: '' },
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
              <Input.Password placeholder={formatMessage('appGui.password')} disabled={!canEditAccount} />
            </Form.Item>
          </Form>
        </Modal>
      </div>
    );
  }

  if (activeTab === 'license') {
    return (
      <div className={`${styles['settings-panel']} ${styles['license-panel']}`}>
        <div className={styles['panel-header']}>
          <h1 className={styles['panel-title']}>{formatMessage('common.license')}</h1>
          <div className={styles['panel-subtitle']}>{formatMessage('settings.licenseDesc')}</div>
        </div>
        <div className={styles['license-section']}>
          <div className={styles['license-row']}>
            <div className={styles['license-label']}>{formatMessage('license.licenseKey')}</div>
            <div className={styles['license-value']}>
              <span className={styles.mono}>{maskLicenseKey(licenseInfo.licenseKey)}</span>
              <Button icon={<Edit size={14} />} onClick={() => setIsEditingLicense(true)} />
            </div>
          </div>
          <div className={styles['license-row']}>
            <div className={styles['license-label']}>{formatMessage('common.status')}</div>
            <div className={styles['license-value']}>
              <span className={`${styles.status} ${getLicenseStatusClassName(licenseInfo.status)}`}>
                <span className={styles['status-dot']} />
                {getLicenseStatusInfo(licenseInfo.status)}
              </span>
            </div>
          </div>
          <div className={styles['license-row']}>
            <div className={styles['license-label']}>{formatMessage('license.validUntil')}</div>
            <div className={styles['license-value']}>{licenseInfo.validUntil || '----'}</div>
          </div>
          <div className={styles['license-row']}>
            <div className={styles['license-label']}>{formatMessage('license.deviceToken')}</div>
            <div className={styles['license-value']}>
              <span className={`${styles.mono} ${styles['token-value']}`}>{licenseDeviceToken || '----'}</span>
              {licenseDeviceToken ? (
                <Button
                  icon={<Copy size={14} />}
                  onClick={() => {
                    copyToClipboard(licenseDeviceToken, (success) => {
                      if (success) {
                        message.success(formatMessage('license.copySuccess'));
                      } else {
                        message.error(formatMessage('common.copyFail'));
                      }
                    });
                  }}
                />
              ) : null}
            </div>
          </div>
        </div>
        <Modal
          open={isEditingLicense}
          title={formatMessage('common.license')}
          footer={null}
          width={560}
          className={styles.licenseModal}
          onCancel={resetLicenseModal}
          destroyOnClose
        >
          <div className={styles['license-modal-content']}>
            <Tabs
              activeKey={licenseActiveTab}
              onChange={setLicenseActiveTab}
              items={[
                { key: 'online', label: formatMessage('license.onlineActivation') },
                { key: 'offline', label: formatMessage('license.offlineActivation') },
              ]}
            />
            {licenseActiveTab === 'online' && (
              <div className={styles['form-field']}>
                <label>{formatMessage('license.licenseKey')}</label>
                <Input
                  value={licenseKeyInput}
                  onChange={(e) => setLicenseKeyInput(e.target.value)}
                  placeholder={formatMessage('license.enterLicenseKey')}
                />
              </div>
            )}
            {licenseActiveTab === 'offline' && (
              <>
                <Alert type="info" showIcon message={formatMessage('license.offlineInfo')} />
                <div className={styles['form-field']}>
                  <label>{formatMessage('license.deviceTokenLabel')}</label>
                  <div className={styles['token-row']}>
                    <Input value={licenseDeviceToken} readOnly />
                    <Button
                      icon={<Copy size={14} />}
                      onClick={() => {
                        copyToClipboard(licenseDeviceToken, (success) => {
                          if (success) {
                            message.success(formatMessage('license.copySuccess'));
                          } else {
                            message.error(formatMessage('common.copyFail'));
                          }
                        });
                      }}
                    >
                      {formatMessage('license.copyDeviceToken')}
                    </Button>
                  </div>
                </div>
                <div className={styles['form-field']}>
                  <label>{formatMessage('license.bundleJson')}</label>
                  <Input.TextArea
                    value={bundleJsonInput}
                    onChange={(e) => setBundleJsonInput(e.target.value)}
                    rows={6}
                    placeholder={formatMessage('license.pasteBundle')}
                  />
                </div>
              </>
            )}
            {licenseError && <div className={styles['license-error']}>{licenseError}</div>}
            <div className={styles.actions}>
              <Button onClick={resetLicenseModal}>{formatMessage('common.cancel')}</Button>
              <Button type="primary" loading={licenseLoading} onClick={onSaveLicense}>
                {licenseActiveTab === 'online' ? formatMessage('common.save') : formatMessage('license.importBundle')}
              </Button>
            </div>
          </div>
        </Modal>
      </div>
    );
  }

  return (
    <div className={styles['settings-panel']}>
      <div className={styles['panel-header']}>
        <h1 className={styles['panel-title']}>{formatMessage('account.profile')}</h1>
        <div className={styles['panel-subtitle']}>{formatMessage('settings.profileDesc')}</div>
      </div>
      <div className={`${styles.section} ${isEditingProfile ? styles['profile-section-editing'] : ''}`}>
        <div className={styles['section-toolbar']}>
          <h2 className={styles['section-title']}>{formatMessage('settings.basicInfo')}</h2>
          {!isEditingProfile ? (
            <Button icon={<Edit size={14} />} disabled={!canEditAccount} onClick={() => setIsEditingProfile(true)}>
              {formatMessage('common.edit')}
            </Button>
          ) : null}
        </div>
        <div className={styles.row}>
          <div className={styles['row-label']}>{formatMessage('settings.avatar')}</div>
          <div className={styles['row-value']}>
            <div className={styles.avatar}>{avatarInitial || <UserAvatar size={20} />}</div>
          </div>
        </div>
        <div className={styles.row}>
          <div className={styles['row-label']}>{formatMessage('account.account')}</div>
          <div className={styles['row-value']}>{currentUserInfo?.preferredUsername || '-'}</div>
        </div>
        {isEditingProfile ? (
          <Form form={profileForm} layout="vertical" className={styles['profile-form']}>
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.displayName')}</div>
              <div className={styles['row-value']}>
                <Form.Item
                  name="firstName"
                  rules={[
                    { required: true, message: formatMessage('rule.required') },
                    {
                      type: 'string',
                      min: 1,
                      max: 200,
                      message: formatMessage('rule.characterLimit'),
                    },
                    { pattern: validNameRegex, message: formatMessage('rule.invalidChars') },
                  ]}
                >
                  <Input placeholder={formatMessage('account.displayName')} disabled={!canEditAccount} />
                </Form.Item>
              </div>
            </div>
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.email')}</div>
              <div className={styles['row-value']}>
                <Form.Item name="email" rules={[{ type: 'email' }]}>
                  <Input placeholder={formatMessage('account.email')} disabled={!canEditAccount} />
                </Form.Item>
              </div>
            </div>
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.phone')}</div>
              <div className={styles['row-value']}>
                <Form.Item name="phone">
                  <Input placeholder={formatMessage('account.phone')} disabled={!canEditAccount} />
                </Form.Item>
              </div>
            </div>
            <div className={`${styles.actions} ${styles['profile-actions']}`}>
              <Button
                onClick={() => {
                  resetProfileForm();
                  setIsEditingProfile(false);
                }}
              >
                {formatMessage('common.cancel')}
              </Button>
              <Button type="primary" loading={loading} disabled={!canEditAccount} onClick={onSaveProfile}>
                {formatMessage('common.save')}
              </Button>
            </div>
          </Form>
        ) : (
          <>
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.displayName')}</div>
              <div className={styles['row-value']}>{currentUserInfo?.firstName || '-'}</div>
            </div>
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.email')}</div>
              <div className={styles['row-value']}>{currentUserInfo?.email || '-'}</div>
            </div>
            <div className={styles.row}>
              <div className={styles['row-label']}>{formatMessage('account.phone')}</div>
              <div className={styles['row-value']}>{currentUserInfo?.phone || '-'}</div>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default AccountSettingsPanel;
