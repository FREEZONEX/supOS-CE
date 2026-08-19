import { useCallback, useEffect, useState } from 'react';
import { Alert, Button, Form, Input, Select, message } from 'antd';
import { useTranslate } from '@/hooks';
import { getUserInfo, loginApi } from '@/apis/core-api/auth';
import { updatePersonConfigApi } from '@/apis/core-api/uns';
import { getToken, removeToken } from '@/utils/auth';
import { LOGIN_URL } from '@/common-types/constans.ts';
import { replaceWithLicenseActivation, resolveLicenseGate } from '@/utils/license-auth';
import { isLaunchpadStandaloneAllowedPath, isLaunchpadStandalonePort } from '@/utils/launchpad-site';
import { I18nEnum, initI18n, useI18nStore } from '@/stores/i18n-store.ts';
import localEnUS from '@/locale/en-US.json';
import localZhCN from '@/locale/zh-CN.json';
import './index.scss';

const loginLogo = '/login/logo-dark.svg';
const loginBackground = '/login/login-background-reference.svg';

type LoginFormValues = {
  username: string;
  password: string;
};

const defaultPostLoginPath = '/uns';
const launchpadPostLoginPath = '/launchpad';
const getDefaultPostLoginPath = () => (isLaunchpadStandalonePort() ? launchpadPostLoginPath : defaultPostLoginPath);
const loginFallbackMessages: Record<string, Record<string, string>> = {
  [I18nEnum.EnUS]: localEnUS as Record<string, string>,
  [I18nEnum.ZhCN]: localZhCN as Record<string, string>,
};

const normalizeRedirectUri = (redirectUri?: string | null) => {
  const next = redirectUri?.trim();
  if (!next || !next.startsWith('/') || next.startsWith('//')) {
    return '';
  }
  const [pathname] = next.split('?');
  if (
    next === '/' ||
    next === '/?isLogin=true' ||
    next === '/login' ||
    next === LOGIN_URL ||
    (isLaunchpadStandalonePort() && !isLaunchpadStandaloneAllowedPath(pathname))
  ) {
    return '';
  }
  return next;
};

const LoginPage = () => {
  const [loading, setLoading] = useState(false);
  const [licenseReady, setLicenseReady] = useState(false);
  const [licenseWarningDays, setLicenseWarningDays] = useState<number | null>(null);
  const formatMessage = useTranslate();
  const lang = useI18nStore((state) => state.lang);
  const [selectedLanguage, setSelectedLanguage] = useState(lang);
  const redirectUri = normalizeRedirectUri(new URLSearchParams(window.location.search).get('redirectUri'));
  const loginMessage = useCallback(
    (id: string) => {
      const translated = formatMessage(id);
      if (translated && translated !== id) {
        return translated;
      }
      const fallbackLang = selectedLanguage === I18nEnum.ZhCN ? I18nEnum.ZhCN : I18nEnum.EnUS;
      return loginFallbackMessages[fallbackLang]?.[id] || loginFallbackMessages[I18nEnum.EnUS]?.[id] || id;
    },
    [formatMessage, selectedLanguage]
  );

  const languageOptions = [
    { label: loginMessage('login.zhCN'), value: I18nEnum.ZhCN },
    { label: loginMessage('login.enUS'), value: I18nEnum.EnUS },
  ];

  const getLoginErrorMessage = useCallback(
    (error: unknown) => {
      const rawMsg = String(
        (error as { response?: { data?: { msg?: string } }; msg?: string })?.response?.data?.msg ||
          (error as { msg?: string })?.msg ||
          ''
      ).trim();

      if (!rawMsg || rawMsg === 'invalid username or password') {
        return loginMessage('login.invalidCredentials');
      }
      return rawMsg;
    },
    [loginMessage]
  );

  useEffect(() => {
    setSelectedLanguage(lang);
  }, [lang]);

  useEffect(() => {
    let active = true;
    const redirectIfTokenValid = async () => {
      if (!getToken()) {
        return false;
      }
      try {
        const userInfo = await getUserInfo();
        if (active) {
          window.location.replace(redirectUri || normalizeRedirectUri(userInfo?.homePage) || getDefaultPostLoginPath());
        }
        return true;
      } catch (error) {
        console.error('Stored token is invalid, clearing token before login:', error);
        removeToken();
        return false;
      }
    };
    const init = async () => {
      try {
        const licenseGate = await resolveLicenseGate();
        if (!active) {
          return;
        }
        if (licenseGate.status === 'not_activated' || licenseGate.status === 'expired') {
          replaceWithLicenseActivation(licenseGate.status);
          return;
        }
        if (
          licenseGate.phase === 'warning' &&
          typeof licenseGate.daysLeft === 'number' &&
          licenseGate.daysLeft >= 0 &&
          licenseGate.daysLeft <= 30
        ) {
          setLicenseWarningDays(licenseGate.daysLeft);
        } else {
          setLicenseWarningDays(null);
        }
        if (await redirectIfTokenValid()) {
          return;
        }
        setLicenseReady(true);
      } catch (error) {
        console.error('Failed to check license status on login page:', error);
        if (active) {
          if (await redirectIfTokenValid()) {
            return;
          }
          setLicenseReady(true);
        }
      }
    };
    init();
    return () => {
      active = false;
    };
  }, [redirectUri]);

  const onFinish = async (values: LoginFormValues) => {
    setLoading(true);
    try {
      const userInfo = await loginApi(values);
      if (userInfo?.sub) {
        try {
          await updatePersonConfigApi({ userId: userInfo.sub, mainLanguage: selectedLanguage });
        } catch (error) {
          console.error('登录后更新用户语言失败:', error);
        }
      }
      window.location.replace(redirectUri || normalizeRedirectUri(userInfo?.homePage) || getDefaultPostLoginPath());
    } catch (error) {
      message.error(getLoginErrorMessage(error));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="tier0-login">
      <Select
        className="tier0-login__language"
        popupClassName="tier0-login__language-popup"
        aria-label={loginMessage('login.language')}
        value={selectedLanguage}
        options={languageOptions}
        onChange={(value) => {
          setSelectedLanguage(value);
          void initI18n(value, {}, { skipRemoteMessages: true });
        }}
      />
      <div className="tier0-login__layout">
        <aside
          className="tier0-login__visual"
          style={{
            backgroundImage: `url(${loginBackground})`,
          }}
        >
          <img className="tier0-login__logo" src={loginLogo} alt="Tier0 Edge" />
        </aside>

        <section className="tier0-login__content">
          <main className="tier0-login__card" id="main">
            <h1 className="tier0-login__title" id="kc-page-title">
              {loginMessage('login.title')}
            </h1>
            {licenseWarningDays !== null ? (
              <Alert
                className="tier0-login__warning"
                type="warning"
                showIcon
                message={loginMessage('login.licenseExpiring')}
              />
            ) : null}
            <Form<LoginFormValues>
              layout="vertical"
              requiredMark={false}
              onFinish={onFinish}
              autoComplete="off"
              className="tier0-login__form"
            >
              <Form.Item
                label={loginMessage('login.username')}
                name="username"
                rules={[{ required: true, message: loginMessage('login.usernameRequired') }]}
              >
                <Input size="large" autoComplete="username" disabled={!licenseReady} />
              </Form.Item>
              <Form.Item
                label={loginMessage('login.password')}
                name="password"
                rules={[{ required: true, message: loginMessage('login.passwordRequired') }]}
              >
                <Input.Password size="large" autoComplete="current-password" disabled={!licenseReady} />
              </Form.Item>
              <Button type="primary" htmlType="submit" size="large" block loading={loading} disabled={!licenseReady}>
                {loginMessage('login.signIn')}
              </Button>
            </Form>
          </main>
        </section>
      </div>
    </div>
  );
};

export default LoginPage;
