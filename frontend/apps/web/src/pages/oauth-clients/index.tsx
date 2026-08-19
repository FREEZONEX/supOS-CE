import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { App, Button, Form, Input, Space, Tag, Tooltip } from 'antd';
import { Copy } from '@carbon/icons-react';
import {
  createOAuthClient,
  deleteOAuthClient,
  queryOAuthClientList,
  updateOAuthClient,
  type OAuthClientPayload,
} from '@/apis/core-api/open-data';
import { ComContent, ComLayout, ProTable, ComEmpty } from '@/components';
import { CheckmarkFilled, Edit, TrashCan } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import { MAX_LENGTHS } from '@/utils/limits';
import { formatTimestamp } from '@/utils';
import { copyToClipboard } from '@/utils/common';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import styles from './index.module.scss';

type OAuthClient = {
  id: number;
  clientId: string;
  clientName: string;
  redirectUris?: string[];
  allowedOrigins?: string[];
  enabled?: boolean;
  createdTime?: number;
  updatedTime?: number;
};

type OAuthClientForm = {
  clientName: string;
  redirectUris: string;
  allowedOrigins?: string;
};

const redirectUriMaxLength = 512;
const fixedColumnWidth = {
  creationTime: 166,
  operation: 86,
};
const flexibleColumnConfig = {
  name: { min: 96, max: 240, padding: 48 },
  clientId: { min: 150, max: 280, padding: 46 },
  redirectUri: { min: 170, max: 420, padding: 40 },
};

const clampWidth = (width: number, min: number, max: number) => Math.min(max, Math.max(min, width));

const estimateTextWidth = (value: string, config: { min: number; max: number; padding: number }) =>
  clampWidth(config.padding + value.length * 7, config.min, config.max);

const getElementContentWidth = (element: HTMLElement) => {
  const style = window.getComputedStyle(element);
  const horizontalPadding = (Number.parseFloat(style.paddingLeft) || 0) + (Number.parseFloat(style.paddingRight) || 0);
  return Math.max(0, element.getBoundingClientRect().width - horizontalPadding);
};

const resolveFlexibleColumnWidths = (records: OAuthClient[], containerWidth: number) => {
  const widths = {
    name: flexibleColumnConfig.name.min,
    clientId: flexibleColumnConfig.clientId.min,
    redirectUri: flexibleColumnConfig.redirectUri.min,
  };

  records.forEach((record) => {
    widths.name = Math.max(widths.name, estimateTextWidth(record.clientName || '', flexibleColumnConfig.name));
    widths.clientId = Math.max(
      widths.clientId,
      estimateTextWidth(record.clientId || '', flexibleColumnConfig.clientId)
    );
    widths.redirectUri = Math.max(
      widths.redirectUri,
      estimateTextWidth((record.redirectUris || [])[0] || '', flexibleColumnConfig.redirectUri)
    );
  });

  const fixedWidth = fixedColumnWidth.creationTime + fixedColumnWidth.operation;
  const availableWidth = Math.max(
    containerWidth ? containerWidth - fixedWidth : 0,
    flexibleColumnConfig.name.min + flexibleColumnConfig.clientId.min + flexibleColumnConfig.redirectUri.min
  );
  const flexibleWidth = widths.name + widths.clientId + widths.redirectUri;

  if (flexibleWidth > availableWidth) {
    const shrinkableWidth =
      widths.name -
      flexibleColumnConfig.name.min +
      widths.clientId -
      flexibleColumnConfig.clientId.min +
      widths.redirectUri -
      flexibleColumnConfig.redirectUri.min;
    const overflow = flexibleWidth - availableWidth;
    if (shrinkableWidth > 0) {
      const shrinkRatio = Math.min(1, overflow / shrinkableWidth);
      widths.name -= (widths.name - flexibleColumnConfig.name.min) * shrinkRatio;
      widths.clientId -= (widths.clientId - flexibleColumnConfig.clientId.min) * shrinkRatio;
      widths.redirectUri -= (widths.redirectUri - flexibleColumnConfig.redirectUri.min) * shrinkRatio;
    }
  }

  return {
    name: Math.round(widths.name),
    clientId: Math.round(widths.clientId),
    redirectUri: Math.round(widths.redirectUri),
  };
};

const splitLines = (value?: string) =>
  String(value || '')
    .split(/\r?\n|,/)
    .map((item) => item.trim())
    .filter(Boolean);

const joinLines = (value?: string[] | string) => {
  if (Array.isArray(value)) {
    return value.join('\n');
  }
  return String(value || '');
};

const normalizeClient = (item: any): OAuthClient => ({
  ...item,
  redirectUris: Array.isArray(item?.redirectUris) ? item.redirectUris : splitLines(item?.redirectUris),
  allowedOrigins: Array.isArray(item?.allowedOrigins) ? item.allowedOrigins : splitLines(item?.allowedOrigins),
});

const OAuthClients = () => {
  const [dataSource, setDataSource] = useState<OAuthClient[]>([]);
  const [loading, setLoading] = useState(false);
  const [tableShellWidth, setTableShellWidth] = useState(0);
  const tableShellRef = useRef<HTMLDivElement>(null);
  const [form] = Form.useForm<OAuthClientForm>();
  const { modal, message } = App.useApp();
  const formatMessage = useTranslate();

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const resp = await queryOAuthClientList();
      setDataSource((resp?.list || []).map(normalizeClient));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    const element = tableShellRef.current;
    if (!element) return undefined;
    const updateWidth = () => setTableShellWidth(Math.floor(getElementContentWidth(element)));
    updateWidth();
    const resizeObserver = new ResizeObserver(updateWidth);
    resizeObserver.observe(element);
    return () => resizeObserver.disconnect();
  }, []);

  const toPayload = useCallback(
    (values: OAuthClientForm): OAuthClientPayload => ({
      clientName: values.clientName.trim(),
      redirectUris: splitLines(values.redirectUris),
      allowedOrigins: splitLines(values.allowedOrigins),
    }),
    []
  );

  const renderForm = useCallback(
    () => (
      <Form form={form} layout="vertical" style={{ marginTop: 16 }} preserve={false}>
        <Form.Item
          name="clientName"
          label={formatMessage('common.name')}
          rules={[{ required: true, message: formatMessage('oauth.nameRequired') }]}
        >
          <Input maxLength={MAX_LENGTHS.connectionName} />
        </Form.Item>
        <Form.Item
          name="redirectUris"
          label={formatMessage('oauth.redirectUri')}
          rules={[
            { required: true, message: formatMessage('oauth.redirectUriRequired') },
            {
              max: redirectUriMaxLength,
              message: formatMessage(
                'oauth.redirectUriMaxLength',
                { length: redirectUriMaxLength },
                `Redirect URI cannot exceed ${redirectUriMaxLength} characters`
              ),
            },
          ]}
        >
          <Input.TextArea
            rows={4}
            maxLength={redirectUriMaxLength}
            showCount
            placeholder={formatMessage('oauth.redirectUriPlaceholder')}
          />
        </Form.Item>
        <Form.Item name="allowedOrigins" label={formatMessage('oauth.allowedOrigin')}>
          <Input.TextArea rows={3} placeholder={formatMessage('oauth.allowedOriginPlaceholder')} />
        </Form.Item>
      </Form>
    ),
    [form, formatMessage]
  );

  const copyClientId = useCallback(
    (clientId: string) => {
      copyToClipboard(clientId, (success) => {
        if (success) {
          message.success(formatMessage('common.copySuccess'));
        } else {
          message.error(formatMessage('common.copyFail'));
        }
      });
    },
    [formatMessage, message]
  );

  const openCreate = () => {
    form.setFieldsValue({ clientName: '', redirectUris: '', allowedOrigins: '' });
    modal.confirm({
      title: formatMessage('oauth.createTitle'),
      icon: null,
      closable: true,
      width: 520,
      className: 'oauth-client-form-modal',
      content: renderForm(),
      okText: formatMessage('common.confirm'),
      cancelText: formatMessage('common.cancel'),
      onOk: async () => {
        const values = await form.validateFields();
        const result = await createOAuthClient(toPayload(values));
        message.success(formatMessage('common.optsuccess'));
        await load();
        const clientId = result?.client?.clientId;
        if (clientId) {
          modal.success({
            title: null,
            icon: null,
            width: 420,
            className: 'oauth-client-success-modal',
            content: (
              <div className="oauth-client-success-content">
                <CheckmarkFilled size={20} className="oauth-client-success-icon" aria-hidden />
                <h3 className="oauth-client-success-title">{formatMessage('common.optsuccess')}</h3>
                <div className="oauth-client-success-label">{formatMessage('oauth.clientId')}</div>
                <Space.Compact style={{ width: '100%' }}>
                  <Input value={clientId} readOnly />
                  <Button icon={<Copy size={16} />} onClick={() => copyClientId(clientId)} />
                </Space.Compact>
              </div>
            ),
            okText: formatMessage('common.ok'),
          });
        }
      },
    });
  };

  const openEdit = useCallback(
    (record: OAuthClient) => {
      form.setFieldsValue({
        clientName: record.clientName,
        redirectUris: joinLines(record.redirectUris),
        allowedOrigins: joinLines(record.allowedOrigins),
      });
      modal.confirm({
        title: formatMessage('oauth.editTitle'),
        icon: null,
        closable: true,
        width: 520,
        className: 'oauth-client-form-modal',
        content: renderForm(),
        okText: formatMessage('common.confirm'),
        cancelText: formatMessage('common.cancel'),
        onOk: async () => {
          const values = await form.validateFields();
          await updateOAuthClient(record.id, toPayload(values));
          message.success(formatMessage('common.optsuccess'));
          await load();
        },
      });
    },
    [form, formatMessage, load, message, modal, renderForm, toPayload]
  );

  const remove = useCallback(
    (record: OAuthClient) => {
      modal.confirm(
        mergeDeleteConfirmProps(
          {
            ...createDeleteConfirmOptions({
              title: formatMessage('oauth.deleteTitle'),
              name: record.clientName || record.clientId,
              formatMessage,
            }),
            onOk: async () => {
              await deleteOAuthClient(record.id);
              message.success(formatMessage('common.deleteSuccessfully'));
              await load();
            },
          },
          formatMessage
        )
      );
    },
    [formatMessage, load, message, modal]
  );

  const columns: any = useMemo(() => {
    const flexibleColumnWidth = resolveFlexibleColumnWidths(dataSource, tableShellWidth);
    return [
      {
        title: formatMessage('common.name'),
        dataIndex: 'clientName',
        width: flexibleColumnWidth.name,
        minWidth: flexibleColumnConfig.name.min,
        ellipsis: true,
        render: (clientName: string) => (
          <span className={styles.clientName} title={clientName || undefined}>
            {clientName || '-'}
          </span>
        ),
      },
      {
        title: formatMessage('oauth.clientId'),
        dataIndex: 'clientId',
        width: flexibleColumnWidth.clientId,
        minWidth: flexibleColumnConfig.clientId.min,
        render: (clientId: string) => (
          <div className={styles.copyTextCell}>
            <span className={styles.clientId} title={clientId}>
              {clientId}
            </span>
            <Button
              className={styles.copyButton}
              type="text"
              size="small"
              title={formatMessage('common.copy')}
              icon={<Copy size={14} />}
              onClick={() => copyClientId(clientId)}
            />
          </div>
        ),
      },
      {
        title: formatMessage('oauth.redirectUri'),
        dataIndex: 'redirectUris',
        width: flexibleColumnWidth.redirectUri,
        minWidth: flexibleColumnConfig.redirectUri.min,
        render: (items: string[]) => {
          const redirectUris = items || [];
          const firstUri = redirectUris[0];
          if (!firstUri) return '-';
          return (
            <div className={styles.redirectUriCell}>
              <Tooltip title={firstUri}>
                <Tag className={styles.redirectUriTag}>{firstUri}</Tag>
              </Tooltip>
              {redirectUris.length > 1 ? (
                <Tag className={styles.redirectUriMore}>+{redirectUris.length - 1}</Tag>
              ) : null}
            </div>
          );
        },
      },
      {
        title: formatMessage('common.creationTime'),
        dataIndex: 'createdTime',
        width: fixedColumnWidth.creationTime,
        minWidth: fixedColumnWidth.creationTime,
        fixed: 'right',
        render: (value: number) => (
          <span className={styles.timeCell} title={value ? formatTimestamp(value) : undefined}>
            {value ? formatTimestamp(value) : '-'}
          </span>
        ),
      },
      {
        title: formatMessage('common.operation'),
        dataIndex: 'operation',
        width: fixedColumnWidth.operation,
        minWidth: fixedColumnWidth.operation,
        fixed: 'right',
        render: (_: unknown, record: OAuthClient) => (
          <span className={styles.operationCell}>
            <Tooltip title={formatMessage('common.edit')}>
              <span className={styles.operationIcon}>
                <Edit
                  className="custom-operation"
                  size={16}
                  style={{ cursor: 'pointer' }}
                  onClick={() => openEdit(record)}
                />
              </span>
            </Tooltip>
            <Tooltip title={formatMessage('common.delete')}>
              <span className={styles.operationIcon}>
                <TrashCan
                  className="custom-operation"
                  size={16}
                  style={{ cursor: 'pointer' }}
                  onClick={() => remove(record)}
                />
              </span>
            </Tooltip>
          </span>
        ),
      },
    ];
  }, [copyClientId, dataSource, formatMessage, openEdit, remove, tableShellWidth]);

  return (
    <ComLayout loading={loading}>
      <ComContent
        title={formatMessage('oauth.title')}
        extra={
          <Button className={styles.addButton} type="primary" onClick={openCreate}>
            {formatMessage('oauth.addClient')}
          </Button>
        }
        mustHasBack={false}
        style={{ display: 'flex', flexDirection: 'column' }}
      >
        <div className={styles.tableShell} ref={tableShellRef}>
          <ProTable
            className={styles.oauthTable}
            columns={columns}
            dataSource={dataSource}
            rowKey="id"
            resizeable
            scroll={{ x: tableShellWidth || '100%' }}
            pagination={false}
            locale={{ emptyText: <ComEmpty /> }}
          />
        </div>
      </ComContent>
    </ComLayout>
  );
};

export default OAuthClients;
