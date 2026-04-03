import { useEffect, useState } from 'react';
import { Button, Form, Input } from 'antd';
import { loginApi } from '@/apis/inter-api/auth';
import { getToken } from '@/utils/auth';
import { LOGIN_URL } from '@/common-types/constans.ts';
import loginLogo from '@/assets/custom-nav/logo-white.svg';
import loginBackground from '@/assets/login/login-background-reference.svg';
import './index.scss';

type LoginFormValues = {
  username: string;
  password: string;
};

const defaultPostLoginPath = '/uns';

const normalizeRedirectUri = (redirectUri?: string | null) => {
  const next = redirectUri?.trim();
  if (!next || !next.startsWith('/')) {
    return '';
  }
  if (next === '/' || next === '/?isLogin=true' || next === LOGIN_URL) {
    return '';
  }
  return next;
};

const LoginPage = () => {
  const [loading, setLoading] = useState(false);
  const redirectUri = normalizeRedirectUri(new URLSearchParams(window.location.search).get('redirectUri'));

  useEffect(() => {
    if (getToken()) {
      window.location.replace(redirectUri || defaultPostLoginPath);
    }
  }, [redirectUri]);

  const onFinish = async (values: LoginFormValues) => {
    setLoading(true);
    try {
      const userInfo = await loginApi(values);
      window.location.replace(redirectUri || normalizeRedirectUri(userInfo?.homePage) || defaultPostLoginPath);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="tier0-login">
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
              Sign in to your account
            </h1>
            <Form<LoginFormValues>
              layout="vertical"
              requiredMark={false}
              onFinish={onFinish}
              autoComplete="off"
              className="tier0-login__form"
            >
              <Form.Item label="Username" name="username" rules={[{ required: true, message: 'Please enter your username' }]}>
                <Input size="large" autoComplete="username" />
              </Form.Item>
              <Form.Item label="Password" name="password" rules={[{ required: true, message: 'Please enter your password' }]}>
                <Input.Password size="large" autoComplete="current-password" />
              </Form.Item>
              <Button type="primary" htmlType="submit" size="large" block loading={loading}>
                Sign In
              </Button>
            </Form>
          </main>
        </section>
      </div>
    </div>
  );
};

export default LoginPage;
