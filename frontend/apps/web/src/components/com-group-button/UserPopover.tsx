import { type FC, type ReactNode, useEffect, useState } from 'react';
import {
  Popover,
  Divider,
  Flex,
  type PopoverProps,
  type SelectProps,
  Button,
  Form,
  Input,
  App,
  ConfigProvider,
  Tabs,
} from 'antd';
import { useTranslate } from '@/hooks';
import ComSelect from '../com-select';
import ProModal from '../pro-modal';
import { setHomePageApi, updateCurrentUserProfile, userResetPwd } from '@/apis/core-api/user-manage';
import { LOGIN_URL, OMC_MODEL } from '@/common-types/constans';
import { APP_USER_GUIDE_ROUTES, APP_USER_TIPS_ENABLE } from '@/common-types/constans';
import { removeToken } from '@/utils/auth';
import { storageOpt } from '@/utils/storage';
import { passwordRegex, validNameRegex } from '@/utils/pattern';
import { fetchBaseStore, fetchSystemInfo, updateForUserInfo, useBaseStore } from '@/stores/base';
import { setTheme, ThemeType, useThemeStore } from '@/stores/theme-store.ts';
import Cookies from 'js-cookie';
import { updatePersonConfigApi } from '@/apis/core-api/uns';
import { logoutApi } from '@/apis/core-api/auth';
import { initI18n, useI18nStore } from '@/stores/i18n-store.ts';
import { preloadPluginLang } from '@/utils/plugin.ts';
import { copyToClipboard } from '@/utils/common';
import { queryLicenseStatus, replaceLicenseKey, queryIdentity, type LicenseStatusResp } from '@/apis/core-api/license';
import { Copy } from '@/components/lucide-icon/carbon';
import { Contrast, Globe, LogOut, Settings2 } from 'lucide-react';
import { useNavigate } from 'react-router';
import HomePageSelect from '@/components/home-page-select';
import { resolveProductVersion } from '@/utils/product-version';

const logout = (path?: string) => {
  logoutApi().then(() => {
    removeToken();
    Cookies.remove(OMC_MODEL, { path: '/' });
    // 退出时删除guide routes信息
    storageOpt.remove(APP_USER_GUIDE_ROUTES);
    // 退出时重置tips信息
    storageOpt.remove(APP_USER_TIPS_ENABLE);
    // 清空
    storageOpt.remove('personInfo');
    location.href = path || LOGIN_URL;
  });
};

const ComList: FC<{
  list: {
    icon?: ReactNode;
    label?: ReactNode;
    children?: ReactNode;
    key: string;
    onClick?: () => void;
    disabled?: boolean;
  }[];
}> = ({ list }) => {
  return (
    <>
      {list?.map((item) => {
        return (
          <Flex
            key={item.key}
            justify="space-between"
            align="center"
            style={{
              width: '100%',
              padding: '6px 8px',
              cursor: item?.disabled ? 'not-allowed' : 'pointer',
              opacity: item?.disabled ? 0.5 : undefined,
            }}
            onClick={!item?.disabled ? item?.onClick : undefined}
          >
            <Flex justify="flex-start" align="center" gap={8} style={{ flex: 1 }}>
              {item.icon}
              {item.label}
            </Flex>
            {item.children && (
              <div onClick={(event) => event.stopPropagation()} onMouseDown={(event) => event.stopPropagation()}>
                {item.children}
              </div>
            )}
          </Flex>
        );
      })}
    </>
  );
};
const popoverSelectProps: Pick<SelectProps, 'getPopupContainer' | 'styles'> = {
  getPopupContainer: () => document.body,
  styles: { popup: { root: { zIndex: 10050 } } },
};

const UserPopover: FC<PopoverProps> = ({ children, ...restProps }) => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const { currentUserInfo, systemInfo, pluginList, menuGroupNoSub } = useBaseStore((state) => ({
    currentUserInfo: state.currentUserInfo,
    systemInfo: state.systemInfo,
    pluginList: state.pluginList,
    menuGroupNoSub: state.menuGroup?.filter((f) => !f.subMenu),
  }));
  const { _theme } = useThemeStore((state) => ({
    _theme: state._theme,
  }));
  const { lang, langList } = useI18nStore((state) => ({
    lang: state.lang,
    langList: state.langList
      ?.filter((f) => f.hasUsed)
      ?.map((m: any) => ({ ...m, value: m?.value?.replace?.('_', '-') })),
  }));
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [form1] = Form.useForm();
  const [form2] = Form.useForm();
  const [form3] = Form.useForm();
  const [form4] = Form.useForm();
  // License 状态
  const [licenseInfo, setLicenseInfo] = useState<{
    licenseKey?: string;
    status?: 'active' | 'expiring_soon' | 'expired' | 'unknown';
    validUntil?: string;
    daysRemaining?: number;
  }>({});
  const [licenseLoading, setLicenseLoading] = useState(false);
  const [isEditingLicense, setIsEditingLicense] = useState(false);
  const [licenseDeviceToken, setLicenseDeviceToken] = useState('');
  const { message } = App.useApp();

  const toggleTheme = (v: string) => {
    setTheme(v as ThemeType);
  };
  const name = currentUserInfo.firstName || currentUserInfo.preferredUsername;
  const version = `v${resolveProductVersion(systemInfo)}`;
  const userContent = (
    <div className="userPopoverWrap">
      <div className="userAvatar">{name?.slice(0, 1)?.toLocaleUpperCase()}</div>
      <div className="userName">{name}</div>
      {currentUserInfo.roleString?.trim() ? (
        <Flex
          title={currentUserInfo.roleString}
          className="userRole"
          justify="center"
          align="center"
          style={{ width: '100%' }}
        >
          <div
            style={{
              overflow: 'hidden',
              whiteSpace: 'nowrap',
              textOverflow: 'ellipsis',
            }}
          >
            {currentUserInfo.roleString}
          </div>
        </Flex>
      ) : null}
      {currentUserInfo.email && (
        <div className="userEmail" title={currentUserInfo.email}>
          {currentUserInfo.email}
        </div>
      )}
      <Divider
        style={{
          background: '#c6c6c6',
          margin: '15px auto',
        }}
      />
      <ComList
        list={[
          {
            icon: <Contrast color="var(--ui-text-color)" size={18} strokeWidth={1.75} aria-hidden />,
            label: <div style={{ color: 'var(--ui-text-color)' }}>{formatMessage('common.theme')}</div>,
            key: 'theme',
            children: (
              <ComSelect
                value={_theme}
                style={{ height: 28, width: 94, backgroundColor: 'var(--ui-bg-color) !important' }}
                onChange={toggleTheme}
                {...popoverSelectProps}
                options={[
                  {
                    label: formatMessage('common.light'),
                    value: ThemeType.Light,
                  },
                  {
                    label: formatMessage('common.dark'),
                    value: ThemeType.Dark,
                  },
                  {
                    label: formatMessage('common.followSystem'),
                    value: ThemeType.System,
                  },
                ]}
              />
            ),
          },
          {
            icon: <Globe color="var(--ui-text-color)" size={18} strokeWidth={1.75} aria-hidden />,
            label: <div style={{ color: 'var(--ui-text-color)' }}>{formatMessage('common.language')}</div>,
            key: 'language',
            children: (
              <ComSelect
                disabled={!currentUserInfo?.sub}
                onChange={(v) => {
                  if (currentUserInfo?.sub) {
                    // 重新过滤插件国际化文件;
                    updatePersonConfigApi({ userId: currentUserInfo.sub!, mainLanguage: v }).then(async () => {
                      const pluginLang = await preloadPluginLang(
                        pluginList
                          ?.filter((f: any) => f.installStatus === 'installed')
                          ?.filter((f: any) => f?.plugInfoYml?.route?.name)
                          ?.map((m: any) => ({ name: `/${m?.plugInfoYml?.route?.name}`, backendName: m?.name })) || [],
                        v
                      );
                      // 先切换前端国际化，再刷新依赖后端国际化的菜单/路由名称
                      await initI18n(v, pluginLang);
                      return fetchSystemInfo(true);
                    });
                  }
                }}
                value={lang}
                style={{ height: 28, width: 94, backgroundColor: 'var(--ui-bg-color) !important' }}
                {...popoverSelectProps}
                options={langList}
              />
            ),
          },
        ]}
      />
      <Divider
        style={{
          background: '#c6c6c6',
          margin: '15px auto',
        }}
      />
      <ComList
        list={[
          // {
          //   icon: (
          //     <Badge count={100} size={'small'} styles={{ indicator: { fontSize: 10, padding: '0 2px' } }}>
          //       <Alarm color="var(--ui-text-color)" size={18} />
          //     </Badge>
          //   ),
          //   label: <div style={{ color: 'var(--ui-text-color)' }}>{formatMessage('common.information')}</div>,
          //   key: 'information',
          //   onClick: () => {
          //     setInformationOpen(true);
          //   },
          // },
          {
            icon: <Settings2 color="var(--ui-text-color)" size={18} strokeWidth={1.75} aria-hidden />,
            label: <div style={{ color: 'var(--ui-text-color)' }}>{formatMessage('common.settings')}</div>,
            key: 'setting',
            onClick: () => {
              setPopoverOpen(false);
              navigate('/settings/profile');
            },
          },
          {
            icon: <LogOut color="var(--ui-text-color)" size={18} strokeWidth={1.75} aria-hidden />,
            label: <div style={{ color: 'var(--ui-text-color)' }}>{formatMessage('common.logout')}</div>,
            key: 'layout',
            onClick: () => {
              setPopoverOpen(false);
              logout(systemInfo.loginPath);
            },
          },
        ]}
      />
      <span style={{ marginTop: 10 }} className="userEmail" title={version}>
        {version}
      </span>
    </div>
  );
  const rest = () => {
    form1.resetFields();
    form2.resetFields();
  };
  const onSave1 = async () => {
    const info = await form1.validateFields();
    setLoading(true);
    updateCurrentUserProfile(info)
      .then(() => {
        message.success(formatMessage('common.settingSuccess'));
        // 修改用户名，手动去更新
        updateForUserInfo({
          ...info,
        });
        setOpen(false);
        form1.resetFields();
      })
      .finally(() => {
        setLoading(false);
      });
  };
  const onSave2 = async () => {
    const info = await form2.validateFields();
    setLoading(true);
    userResetPwd({
      newPassword: info.password,
      password: info.oldPassword,
      userId: currentUserInfo?.sub,
      username: currentUserInfo?.preferredUsername,
    })
      .then(() => {
        message.success(formatMessage('common.settingSuccess'));
        logout(systemInfo.loginPath);
      })
      .finally(() => {
        setLoading(false);
      });
  };
  const onSave3 = async () => {
    const info = await form3.validateFields();
    if (info) {
      setHomePageApi({ homePage: info.homePage })?.then((config: any) => {
        const homePage = config?.homePage || info.homePage;
        updateForUserInfo({ homePage });
        message.success(formatMessage('common.settingSuccess'));
      });
    }
  };

  // 获取 License 状态
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
    } catch (error) {
      console.error('Failed to fetch license status:', error);
    } finally {
      setLicenseLoading(false);
    }
  };

  // 获取设备标识
  const fetchDeviceToken = async () => {
    try {
      const resp = await queryIdentity();
      if (resp?.device_token) {
        setLicenseDeviceToken(resp.device_token);
      }
    } catch (error) {
      console.error('Failed to fetch device token:', error);
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

  // 更新 License Key
  const onSave4 = async () => {
    const info = await form4.validateFields();
    setLicenseLoading(true);
    try {
      await replaceLicenseKey(info.licenseKey);
      message.success(formatMessage('license.updateLicenseKeySuccess'));
      setIsEditingLicense(false);
      fetchLicenseStatus();
    } catch (error: any) {
      console.error('Failed to update license key:', error);
      message.error(resolveUpdateLicenseKeyErrorMessage(error?.msg));
    } finally {
      setLicenseLoading(false);
    }
  };

  // 获取状态颜色和文本
  const getLicenseStatusInfo = (status?: string) => {
    switch (status) {
      case 'active':
        return {
          color: '#52c41a',
          bgColor: '#f6ffed',
          borderColor: '#b7eb8f',
          text: formatMessage('license.status.active'),
        };
      case 'expiring_soon':
        return {
          color: '#faad14',
          bgColor: '#fffbe6',
          borderColor: '#ffe58f',
          text: formatMessage('license.status.expiring'),
        };
      case 'expired':
        return {
          color: '#ff4d4f',
          bgColor: '#fff2f0',
          borderColor: '#ffccc7',
          text: formatMessage('license.status.expired'),
        };
      default:
        return {
          color: '#d9d9d9',
          bgColor: '#fafafa',
          borderColor: 'var(--ui-line-color)',
          text: formatMessage('license.status.unknown'),
        };
    }
  };

  useEffect(() => {
    if (open) {
      // 更新路由
      fetchBaseStore?.().then(() => {
        form3.setFieldValue('homePage', currentUserInfo?.homePage);
      });
      // 获取 License 状态和设备标识
      // Keep the hidden callbacks type-checked without contacting Enterprise-only endpoints.
      void fetchLicenseStatus;
      void fetchDeviceToken;
    }
  }, [open]);

  const items: any[] = [
    {
      label: formatMessage('account.profile'),
      key: 1,
      children: (
        <Flex style={{ height: 300 }} vertical>
          <Form layout="vertical" form={form1} style={{ flex: 1 }}>
            <Form.Item
              label={formatMessage('account.updateDisplayName')}
              name="firstName"
              rules={[
                {
                  required: true,
                  message: formatMessage('rule.required'),
                },
                {
                  type: 'string',
                  min: 1,
                  max: 200,
                  message: formatMessage('rule.characterLimit'),
                },
                {
                  pattern: validNameRegex,
                  message: formatMessage('rule.invalidChars'),
                },
              ]}
            >
              <Input className={'input'} placeholder={formatMessage('account.displayName')} />
            </Form.Item>
            <Form.Item label={formatMessage('common.updateEmail')} name="email" rules={[{ type: 'email' }]}>
              <Input placeholder={formatMessage('account.email')} />
            </Form.Item>
            <Form.Item label={formatMessage('common.updatePhone')} name="phone">
              <Input placeholder={formatMessage('account.phone')} />
            </Form.Item>
          </Form>
          <Button onClick={onSave1} style={{ height: 32 }} block type="primary" loading={loading}>
            {formatMessage('common.save')}
          </Button>
        </Flex>
      ),
    },
    {
      label: formatMessage('common.password'),
      key: 2,
      children: (
        <Flex style={{ height: 300 }} vertical>
          <Form layout="vertical" form={form2} style={{ flex: 1 }}>
            <Form.Item
              label={formatMessage('account.oldPassWord')}
              name="oldPassword"
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
              label={formatMessage('account.newpassWord')}
              name="password"
              dependencies={['oldPassword']}
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
                    if (!value || getFieldValue('oldPassword') !== value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error(formatMessage('account.passwordSame')));
                  },
                }),
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
              <Input.Password placeholder={formatMessage('appGui.password')} />
            </Form.Item>
          </Form>
          <Button onClick={onSave2} style={{ height: 32 }} block type="primary" loading={loading}>
            {formatMessage('common.save')}
          </Button>
        </Flex>
      ),
    },
    {
      label: formatMessage('account.homePage'),
      key: 3,
      children: (
        <Form
          layout="vertical"
          form={form3}
          style={{ height: 300, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}
        >
          <Form.Item label={formatMessage('account.homePage')} name="homePage">
            <HomePageSelect
              enabled={open}
              resources={menuGroupNoSub}
              placeholder={formatMessage('common.searchPage')}
              allowClear
              showSearch
            />
          </Form.Item>
          <Form.Item shouldUpdate={(pre, cur) => pre.homePage !== cur.homePage} noStyle>
            {({ getFieldValue }) => {
              return (
                <Button
                  disabled={!getFieldValue('homePage')}
                  onClick={onSave3}
                  style={{ height: 32 }}
                  block
                  type="primary"
                  loading={loading}
                >
                  {formatMessage('common.save')}
                </Button>
              );
            }}
          </Form.Item>
        </Form>
      ),
    },
    {
      label: formatMessage('common.license'),
      key: 4,
      children: (
        <Flex style={{ minHeight: 300 }} vertical>
          <Form form={form4} layout="vertical" style={{ flex: 1 }}>
            {/* License Key */}
            <Form.Item label={formatMessage('license.licenseKey')}>
              {isEditingLicense ? (
                <Form.Item
                  name="licenseKey"
                  rules={[{ required: true, message: formatMessage('rule.required') }]}
                  style={{ marginBottom: 0 }}
                >
                  <Input placeholder={formatMessage('license.enterLicenseKey')} />
                </Form.Item>
              ) : (
                <span style={{ fontFamily: 'monospace', letterSpacing: 1 }}>
                  {(() => {
                    if (!licenseInfo.licenseKey) return '----';
                    if (licenseInfo.licenseKey.includes('*')) return licenseInfo.licenseKey;
                    if (licenseInfo.licenseKey.length <= 8) return licenseInfo.licenseKey;
                    return `${licenseInfo.licenseKey.slice(0, 5)}****${licenseInfo.licenseKey.slice(-4)}`;
                  })()}
                </span>
              )}
            </Form.Item>

            {/* Status */}
            <Form.Item label={formatMessage('common.status')}>
              {(() => {
                const statusInfo = getLicenseStatusInfo(licenseInfo.status);
                return (
                  <span
                    style={{
                      display: 'inline-flex',
                      alignItems: 'center',
                      gap: 4,
                      padding: '2px 8px',
                      borderRadius: 4,
                      backgroundColor: statusInfo.bgColor,
                      border: `1px solid ${statusInfo.borderColor}`,
                      color: statusInfo.color,
                      fontSize: 12,
                    }}
                  >
                    <span
                      style={{
                        width: 6,
                        height: 6,
                        borderRadius: '50%',
                        backgroundColor: statusInfo.color,
                      }}
                    />
                    {statusInfo.text}
                  </span>
                );
              })()}
            </Form.Item>

            {/* Valid Until */}
            <Form.Item label={formatMessage('license.validUntil')}>
              <span>{licenseInfo.validUntil || '----'}</span>
            </Form.Item>

            {/* Device Token */}
            <Form.Item label={formatMessage('license.deviceToken')}>
              <Flex gap={8} align="center">
                <span style={{ fontFamily: 'monospace', letterSpacing: 1 }}>{licenseDeviceToken || '----'}</span>
                {licenseDeviceToken && (
                  <Copy
                    size={14}
                    style={{ cursor: 'pointer', color: 'var(--ui-theme-color)', flexShrink: 0 }}
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
                )}
              </Flex>
            </Form.Item>
          </Form>

          {/* Action Button */}
          {isEditingLicense ? (
            <Flex gap={8} wrap="wrap" style={{ marginTop: 8 }}>
              <Button
                onClick={() => {
                  setIsEditingLicense(false);
                  form4.resetFields();
                }}
                style={{ height: 32, flex: '1 1 96px', minWidth: 0 }}
              >
                {formatMessage('common.cancel')}
              </Button>
              <Button
                onClick={onSave4}
                style={{ height: 32, flex: '1 1 96px', minWidth: 0 }}
                type="primary"
                loading={licenseLoading}
              >
                {formatMessage('common.save')}
              </Button>
            </Flex>
          ) : (
            <Button
              onClick={() => setIsEditingLicense(true)}
              style={{ height: 32, color: 'var(--ui-theme-color)', border: '1px solid var(--ui-theme-color)' }}
            >
              {formatMessage('license.changeLicenseKey')}
            </Button>
          )}
        </Flex>
      ),
    },
  ];
  return (
    <>
      <Popover
        rootClassName="userPopover"
        placement="bottomRight"
        trigger="click"
        {...restProps}
        open={popoverOpen}
        onOpenChange={setPopoverOpen}
        content={userContent}
      >
        {children}
      </Popover>
      <ProModal
        size="sm"
        forceRender
        onCancel={() => {
          setOpen(false);
          rest();
        }}
        title={formatMessage('account.settings')}
        open={open}
        // open
        maskClosable={false}
      >
        <ConfigProvider
          theme={{
            components: {
              Form: {
                itemMarginBottom: 12,
              },
            },
          }}
        >
          <Tabs tabPosition="left" items={items.filter((item) => item.key !== 4).map((item) => ({ ...item, forceRender: true }))} />

          {/*<Divider*/}
          {/*  style={{*/}
          {/*    background: '#c6c6c6',*/}
          {/*    margin: '16px auto',*/}
          {/*  }}*/}
          {/*/>*/}

          {/*<Divider*/}
          {/*  style={{*/}
          {/*    background: '#c6c6c6',*/}
          {/*    margin: '16px auto',*/}
          {/*  }}*/}
          {/*/>*/}
        </ConfigProvider>
      </ProModal>
    </>
  );
};

export default UserPopover;
