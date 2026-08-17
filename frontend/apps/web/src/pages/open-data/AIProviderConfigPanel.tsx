import { useCallback, useEffect, useMemo, useState } from 'react';
import { App, Form, Input, Space, Spin, Tag } from 'antd';
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
/** 仅用于探测 Base URL 是否提供 /embeddings；不代表用户最终选用的模型。 */
const EMBED_PROBE_MODEL = 'text-embedding-3-small';

type AIConfigFormValues = Pick<AIProviderConfigPayload, 'baseUrl' | 'model' | 'apiKey' | 'embeddingModel'>;

type EmbedSupport = 'unknown' | 'supported' | 'unsupported';

/** 影响探测/实测结果的连接字段；一变即失效上次测试态。 */
const CONNECTION_TEST_FIELDS: (keyof AIConfigFormValues)[] = ['baseUrl', 'apiKey', 'model', 'embeddingModel'];

const AIProviderConfigPanel = () => {
  const formatMessage = useTranslate('OpenData');
  const formatCommon = useTranslate();
  const { message, modal } = App.useApp();
  const [form] = Form.useForm<AIConfigFormValues>();
  const [config, setConfig] = useState<AIProviderConfig>();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [embeddingTest, setEmbeddingTest] = useState<AIEmbeddingTestResult>();
  const [embedSupport, setEmbedSupport] = useState<EmbedSupport>('unknown');
  const [chatTestError, setChatTestError] = useState<string>();
  const [embeddingFieldError, setEmbeddingFieldError] = useState<string>();
  const embeddingModel = Form.useWatch('embeddingModel', form);

  const embeddingEnabled = embedSupport === 'supported';

  // 本页不再提供「单独配置」编辑：已有 custom 配置，或改了 URL/Key 尚未重测时，
  // 保存 payload 不带 embedding 字段，避免后端把 embedding 列一并改写/清空。
  const shouldWriteEmbeddingFields =
    config?.embeddingMode !== 'custom' && (embedSupport === 'supported' || embedSupport === 'unsupported');

  const embeddingHelpText = useMemo(() => {
    if (embeddingFieldError) return embeddingFieldError;
    if (embedSupport === 'unsupported') return formatMessage('aiEmbedUnsupportedHint');
    if (!embeddingEnabled) return formatMessage('aiEmbedDisabledHint');
    if (!embeddingModel?.trim()) return formatMessage('aiEmbedEmptyWarn');
    return formatMessage('aiEmbedUnsOnlyHint');
  }, [embeddingEnabled, embeddingFieldError, embeddingModel, embedSupport, formatMessage]);

  const embeddingHelpIsError = Boolean(embeddingFieldError || embedSupport === 'unsupported');

  const applyConfig = useCallback(
    (current?: AIProviderConfig) => {
      setConfig(current);
      setEmbeddingTest(undefined);
      setChatTestError(undefined);
      setEmbeddingFieldError(undefined);
      const savedEmbed = Boolean(current?.embeddingModel?.trim());
      // custom 模式仍展示已存模型名，但不视为本页可改的 inherit 探测态
      setEmbedSupport(current?.embeddingMode === 'custom' ? 'unknown' : savedEmbed ? 'supported' : 'unknown');
      form.setFieldsValue({
        baseUrl: current?.baseUrl || 'https://api.openai.com/v1',
        model: current?.model || '',
        apiKey: '',
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
      workspaceId: WORKSPACE_ID,
      name: config?.name || CONFIG_NAME,
      provider: 'openai_compatible',
      baseUrl: values.baseUrl,
      model: values.model,
      apiKey: values.apiKey,
    };
    if (!payload.apiKey?.trim()) {
      delete payload.apiKey;
    }
    if (shouldWriteEmbeddingFields) {
      payload.embeddingMode = 'inherit';
      // unsupported：明确清空 inherit 下的模型；supported：写入表单值（可为空表示透传）
      payload.embeddingModel = embedSupport === 'supported' ? values.embeddingModel || '' : '';
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
    // 未打算改写 embedding 列时直接保存（保留 custom / 未重测前的已存模型）
    if (!shouldWriteEmbeddingFields) {
      await doSave(values);
      return;
    }
    const embedName = embedSupport === 'supported' ? values.embeddingModel?.trim() : '';
    const testPassed = embeddingTest?.status === 'ok' && embeddingTest.dimensionsMatch;
    if (!embedName || testPassed) {
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

  const mapEmbeddingFieldError = (result: AIEmbeddingTestResult, userFilledModel: boolean): string | undefined => {
    if (!userFilledModel) return undefined;
    const params = {
      model: result.model,
      dims: result.dimensions,
      expected: result.expectedDimensions,
      code: result.httpStatus,
    };
    if (result.status === 'ok' && result.dimensionsMatch) return undefined;
    if (result.status === 'ok') return formatMessage('aiEmbedTestDim', params);
    if (result.status === 'model_not_found') return formatMessage('aiEmbedUnavailable');
    if (result.status === 'unauthorized') return formatMessage('aiEmbedTestAuth', params);
    if (result.status === 'unsupported') return formatMessage('aiEmbedUnsupportedHint');
    return formatMessage('aiEmbedUnavailable');
  };

  const testConnection = async () => {
    const values = await form.validateFields(['baseUrl', 'apiKey']);
    const allValues = { ...form.getFieldsValue(), ...values };
    setTesting(true);
    setChatTestError(undefined);
    setEmbeddingFieldError(undefined);
    setEmbeddingTest(undefined);

    const payload = buildPayload(allValues);
    let chatOk = false;
    try {
      await testAIProviderConnection(payload);
      chatOk = true;
    } catch {
      setChatTestError(formatMessage('aiChatUnavailable'));
    }

    const userFilledModel = Boolean(allValues.embeddingModel?.trim());
    const probeModel = userFilledModel ? allValues.embeddingModel!.trim() : EMBED_PROBE_MODEL;
    let nextSupport: EmbedSupport = 'unknown';
    let embedOkForUser = !userFilledModel;

    try {
      const result = await testAIEmbeddingConnection({
        ...payload,
        embeddingMode: 'inherit',
        embeddingModel: probeModel,
      });
      setEmbeddingTest(userFilledModel ? result : undefined);

      if (result.status === 'unsupported') {
        nextSupport = 'unsupported';
        embedOkForUser = true;
        form.setFieldValue('embeddingModel', '');
        setEmbeddingFieldError(undefined);
      } else {
        // API 路径存在（含 model_not_found / dim 不符 / unauthorized），开放填写。
        nextSupport = 'supported';
        const fieldError = mapEmbeddingFieldError(result, userFilledModel);
        setEmbeddingFieldError(fieldError);
        embedOkForUser = !fieldError;
      }
    } catch {
      if (userFilledModel) {
        setEmbeddingFieldError(formatMessage('aiEmbedUnavailable'));
        embedOkForUser = false;
      }
    }

    setEmbedSupport(nextSupport);

    if (chatOk && (nextSupport === 'unsupported' || embedOkForUser)) {
      message.success(formatMessage('aiTestSuccess'));
    } else if (chatOk && nextSupport === 'supported' && userFilledModel && !embedOkForUser) {
      message.warning(formatMessage('aiTestChatOkEmbedFail'));
    }

    setTesting(false);
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
            if (CONNECTION_TEST_FIELDS.some((field) => field in changed)) {
              setEmbeddingTest(undefined);
              setChatTestError(undefined);
              if ('baseUrl' in changed || 'apiKey' in changed) {
                setEmbedSupport('unknown');
                setEmbeddingFieldError(undefined);
              } else if ('embeddingModel' in changed) {
                setEmbeddingFieldError(undefined);
              } else if ('model' in changed) {
                setChatTestError(undefined);
              }
            }
          }}
        >
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
          <Form.Item
            name="model"
            label={formatMessage('aiModel')}
            validateStatus={chatTestError ? 'error' : undefined}
            help={chatTestError}
          >
            <Input placeholder="e.g. gpt-4o" maxLength={128} />
          </Form.Item>
          <Form.Item
            name="embeddingModel"
            label={
              <Space size={4}>
                {formatMessage('aiEmbedModel')}
                <HelpTooltip title={formatMessage('aiEmbedModelHint', { expected: 1536 })} />
              </Space>
            }
            validateStatus={embeddingHelpIsError ? 'error' : undefined}
            help={
              embeddingHelpText ? (
                <span className={embeddingHelpIsError ? styles.aiFieldError : styles.aiFieldHint}>
                  {embeddingHelpText}
                </span>
              ) : undefined
            }
          >
            <Input
              placeholder="text-embedding-3-small"
              maxLength={128}
              disabled={!embeddingEnabled}
            />
          </Form.Item>

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
