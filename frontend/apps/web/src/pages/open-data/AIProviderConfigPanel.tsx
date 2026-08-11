import { useCallback, useEffect, useState } from 'react';
import { Alert, App, Form, Input, Radio, Space, Spin, Tag } from 'antd';
import {
  createAIProviderConfig,
  queryAIProviderConfigs,
  testAIEmbeddingConnection,
  testAIProviderConnection,
  updateAIProviderConfig,
  type AIEmbeddingTestResult,
  type AIProviderConfig,
  type AIProviderConfigPayload,
} from '@/apis/core-api/open-data';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton, HelpTooltip } from '@/components';
import { openConfirmModal } from '@/components/confirm-modal';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

const WORKSPACE_ID = 1;
const CONFIG_NAME = 'Enterprise AI';

type AIConfigFormValues = Pick<
  AIProviderConfigPayload,
  'baseUrl' | 'model' | 'apiKey' | 'embeddingMode' | 'embeddingBaseUrl' | 'embeddingApiKey' | 'embeddingModel'
>;

// 影响向量连接测试结果的字段；这些字段一变，上次的实测结果即失效。
const EMBEDDING_TEST_FIELDS: (keyof AIConfigFormValues)[] = [
  'baseUrl',
  'apiKey',
  'embeddingMode',
  'embeddingBaseUrl',
  'embeddingApiKey',
  'embeddingModel',
];

const AIProviderConfigPanel = () => {
  const formatMessage = useTranslate('OpenData');
  const formatCommon = useTranslate();
  const { message, modal } = App.useApp();
  const [form] = Form.useForm<AIConfigFormValues>();
  const [config, setConfig] = useState<AIProviderConfig>();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testingEmbedding, setTestingEmbedding] = useState(false);
  const [embeddingTest, setEmbeddingTest] = useState<AIEmbeddingTestResult>();
  const embeddingMode = Form.useWatch('embeddingMode', form) ?? 'inherit';
  const embeddingModel = Form.useWatch('embeddingModel', form);

  const applyConfig = useCallback(
    (current?: AIProviderConfig) => {
      setConfig(current);
      setEmbeddingTest(undefined);
      form.setFieldsValue({
        baseUrl: current?.baseUrl || 'https://api.openai.com/v1',
        model: current?.model || '',
        apiKey: '',
        embeddingMode: current?.embeddingMode || 'inherit',
        embeddingBaseUrl: current?.embeddingBaseUrl || '',
        embeddingApiKey: '',
        embeddingModel: current?.embeddingModel || '',
      });
    },
    [form]
  );

  const loadConfig = useCallback(async () => {
    setLoading(true);
    try {
      const result = await queryAIProviderConfigs(WORKSPACE_ID);
      applyConfig(result?.list?.[0]);
    } catch {
      message.error(formatMessage('aiLoadFailed'));
    } finally {
      setLoading(false);
    }
  }, [applyConfig, formatMessage, message]);

  useEffect(() => {
    const timer = window.setTimeout(() => void loadConfig(), 0);
    return () => window.clearTimeout(timer);
  }, [loadConfig]);

  const buildPayload = (values: AIConfigFormValues): AIProviderConfigPayload => {
    const payload: AIProviderConfigPayload = {
      ...values,
      workspaceId: WORKSPACE_ID,
      name: config?.name || CONFIG_NAME,
      provider: 'openai_compatible',
    };
    if (!payload.apiKey?.trim()) {
      delete payload.apiKey;
    }
    if (!payload.embeddingApiKey?.trim()) {
      delete payload.embeddingApiKey;
    }
    // 复用模式下不提交未启用的单独连接字段，避免表单里残留的值被后端校验拦下。
    if (payload.embeddingMode !== 'custom') {
      delete payload.embeddingBaseUrl;
      delete payload.embeddingApiKey;
    }
    return payload;
  };

  const doSave = async (values: AIConfigFormValues) => {
    const payload: AIProviderConfigPayload = {
      ...buildPayload(values),
      isDefault: true,
      enabled: true,
    };
    setSaving(true);
    try {
      const saved = config ? await updateAIProviderConfig(config.id, payload) : await createAIProviderConfig(payload);
      applyConfig(saved);
      message.success(formatMessage('aiSaveSuccess'));
    } catch {
      // The shared request interceptor displays the localized API error.
    } finally {
      setSaving(false);
    }
  };

  const saveConfig = async () => {
    const values = await form.validateFields();
    // 向量模型非必填：没填视为未启用（网关透传），页面常驻警告已提示影响，直接保存。
    // 填了但未实测通过（或维度不符）时给二次确认，避免静默保存后语义检索悄悄失效。
    const testPassed = embeddingTest?.status === 'ok' && embeddingTest.dimensionsMatch;
    if (!values.embeddingModel?.trim() || testPassed) {
      await doSave(values);
      return;
    }
    openConfirmModal(modal, {
      title: formatMessage('aiEmbedConfirmTitle'),
      content: formatMessage(
        embeddingTest?.status === 'ok' && !embeddingTest.dimensionsMatch ? 'aiEmbedConfirmDim' : 'aiEmbedConfirmTest',
        {
          dims: embeddingTest?.dimensions,
          expected: embeddingTest?.expectedDimensions ?? 1536,
        }
      ),
      okText: formatMessage('aiEmbedConfirmSave'),
      cancelText: formatCommon('common.cancel'),
      onOk: () => doSave(values),
    });
  };

  const testConnection = async () => {
    const values = await form.validateFields();
    setTesting(true);
    try {
      await testAIProviderConnection(buildPayload(values));
      message.success(formatMessage('aiTestSuccess'));
    } catch {
      // The shared request interceptor displays the localized API error.
    } finally {
      setTesting(false);
    }
  };

  const testEmbedding = async () => {
    if (!form.getFieldValue('embeddingModel')?.trim()) {
      message.warning(formatMessage('aiEmbedTestNeedModel'));
      return;
    }
    const fields: (keyof AIConfigFormValues)[] =
      embeddingMode === 'custom' ? ['embeddingBaseUrl', 'embeddingApiKey'] : ['baseUrl', 'apiKey'];
    await form.validateFields(fields);
    const values = form.getFieldsValue();
    setTestingEmbedding(true);
    try {
      const result = await testAIEmbeddingConnection(buildPayload(values));
      setEmbeddingTest(result);
      const params = {
        model: result.model,
        dims: result.dimensions,
        expected: result.expectedDimensions,
        code: result.httpStatus,
      };
      if (result.status === 'ok' && result.dimensionsMatch) {
        message.success(formatMessage('aiEmbedTestOk', params));
      } else if (result.status === 'ok') {
        message.warning(formatMessage('aiEmbedTestDim', params), 6);
      } else if (result.status === 'unsupported') {
        message.error(formatMessage('aiEmbedTestNoApi', params), 6);
      } else if (result.status === 'model_not_found') {
        message.error(formatMessage('aiEmbedTestNoModel', params), 6);
      } else if (result.status === 'unauthorized') {
        message.error(formatMessage('aiEmbedTestAuth', params));
      } else {
        message.error(formatMessage('aiEmbedTestFailed'));
      }
    } catch {
      setEmbeddingTest(undefined);
      // The shared request interceptor displays the localized API error.
    } finally {
      setTestingEmbedding(false);
    }
  };

  const editAuth = ButtonPermission[config ? 'OpenData.aiConfigEdit' : 'OpenData.aiConfigAdd'];

  return (
    <section className={styles.aiConfigSection}>
      <div className={styles.aiConfigHeader}>
        <h2>{formatMessage('aiConfigTitle')}</h2>
        <Tag color={config?.apiKeySet ? 'success' : 'default'}>
          {formatMessage(config?.apiKeySet ? 'aiConfigured' : 'aiUnconfigured')}
        </Tag>
      </div>

      <Spin spinning={loading}>
        <Form
          form={form}
          layout="vertical"
          requiredMark={false}
          className={styles.aiSingletonForm}
          onValuesChange={(changed) => {
            if (EMBEDDING_TEST_FIELDS.some((field) => field in changed)) {
              setEmbeddingTest(undefined);
            }
          }}
        >
          <Form.Item
            name="apiKey"
            label={
              <Space size={4}>
                {formatMessage('aiKey')}
                <HelpTooltip title={formatMessage(config ? 'aiKeyUnchanged' : 'aiKeyHint')} />
              </Space>
            }
            rules={[{ required: !config }]}
          >
            <Input.Password
              placeholder={config?.apiKeySet ? '••••••••••' : 'sk-********'}
              autoComplete="new-password"
            />
          </Form.Item>
          <Form.Item name="model" label={formatMessage('aiModel')}>
            <Input placeholder="e.g. gpt-4o" maxLength={128} />
          </Form.Item>
          <Form.Item
            name="baseUrl"
            label={
              <Space size={4}>
                {formatMessage('aiBaseUrl')}
                <HelpTooltip title={formatMessage('aiBaseUrlHint')} />
              </Space>
            }
            rules={[{ required: true, type: 'url' }]}
          >
            <Input />
          </Form.Item>

          <div className={styles.aiEmbeddingSection}>
            <div className={styles.aiEmbeddingTitle}>{formatMessage('aiEmbedTitle')}</div>
            <Alert type="info" showIcon message={formatMessage('aiEmbedIntro')} className={styles.aiEmbeddingIntro} />
            <Form.Item name="embeddingMode" initialValue="inherit">
              <Radio.Group>
                <Radio value="inherit">{formatMessage('aiEmbedInherit')}</Radio>
                <Radio value="custom">{formatMessage('aiEmbedCustom')}</Radio>
              </Radio.Group>
            </Form.Item>
            {embeddingMode === 'inherit' && (
              <p className={styles.aiEmbeddingNote}>{formatMessage('aiEmbedInheritNote')}</p>
            )}
            {embeddingMode === 'custom' && (
              <>
                <Form.Item
                  name="embeddingBaseUrl"
                  label={
                    <Space size={4}>
                      {formatMessage('aiBaseUrl')}
                      <HelpTooltip title={formatMessage('aiBaseUrlHint')} />
                    </Space>
                  }
                  rules={[{ required: true, type: 'url' }]}
                >
                  <Input />
                </Form.Item>
                <Form.Item
                  name="embeddingApiKey"
                  label={
                    <Space size={4}>
                      {formatMessage('aiKey')}
                      <HelpTooltip title={formatMessage(config?.embeddingApiKeySet ? 'aiKeyUnchanged' : 'aiKeyHint')} />
                    </Space>
                  }
                  rules={[{ required: !config?.embeddingApiKeySet }]}
                >
                  <Input.Password
                    placeholder={config?.embeddingApiKeySet ? '••••••••••' : 'sk-********'}
                    autoComplete="new-password"
                  />
                </Form.Item>
              </>
            )}
            <Form.Item
              name="embeddingModel"
              label={
                <Space size={4}>
                  {formatMessage('aiEmbedModel')}
                  <HelpTooltip title={formatMessage('aiEmbedModelHint', { expected: 1536 })} />
                </Space>
              }
            >
              <Input placeholder="text-embedding-3-small" maxLength={128} />
            </Form.Item>
            {!embeddingModel?.trim() && (
              <Alert
                type="warning"
                showIcon
                message={formatMessage('aiEmbedEmptyWarn')}
                className={styles.aiEmbeddingIntro}
              />
            )}
            <AuthButton auth={editAuth} loading={testingEmbedding} onClick={() => void testEmbedding()}>
              {formatMessage('aiEmbedTest')}
            </AuthButton>
          </div>

          <div className={styles.aiConfigActions}>
            <AuthButton auth={editAuth} loading={testing} onClick={() => void testConnection()}>
              {formatMessage('aiTestConnection')}
            </AuthButton>
            <AuthButton auth={editAuth} type="primary" loading={saving} onClick={() => void saveConfig()}>
              {formatMessage('aiSaveSettings')}
            </AuthButton>
          </div>
        </Form>
      </Spin>
    </section>
  );
};

export default AIProviderConfigPanel;
