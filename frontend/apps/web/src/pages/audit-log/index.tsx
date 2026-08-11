import { getAuditLogDetail, getAuditLogPage, type AuditLogItem } from '@/apis/core-api/audit-log';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ProTable from '@/components/pro-table';
import ProSearch from '@/components/pro-search';
import { useTranslate } from '@/hooks';
import {
  getAuditActionLabel,
  getAuditModuleLabel,
  getAuditOperatorLabel,
  getAuditResultLabel,
  getAuditResourceLabel,
  getAuditScopeLabel,
} from './presentation';
import { formatTimestamp } from '@/utils/format';
import { App, Space, Typography } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { AuditLogDetailDrawer } from './AuditLogDetailDrawer';
import styles from './index.module.scss';

const { Text } = Typography;

const fixedColumnWidth = {
  time: 166,
  operation: 86,
};
const tableScrollbarReserve = 12;
const auditMinimumTableWidth = 1120;
const flexibleColumnConfig = {
  scope: { min: 88, max: 150, padding: 46 },
  module: { min: 112, max: 180, padding: 46 },
  action: { min: 96, max: 150, padding: 46 },
  resource: { min: 110, max: 180, padding: 46 },
  operator: { min: 120, max: 210, padding: 46 },
  result: { min: 96, max: 140, padding: 46 },
};

const clampWidth = (width: number, min: number, max: number) => Math.min(max, Math.max(min, width));

const estimateTextWidth = (value: string, config: { min: number; max: number; padding: number }) =>
  clampWidth(config.padding + value.length * 7, config.min, config.max);

const getElementContentWidth = (element: HTMLElement) => {
  const style = window.getComputedStyle(element);
  const horizontalPadding = (Number.parseFloat(style.paddingLeft) || 0) + (Number.parseFloat(style.paddingRight) || 0);
  return Math.max(0, element.getBoundingClientRect().width - horizontalPadding);
};

const renderEllipsisText = (value: string) => (
  <Text className={styles.ellipsisCell} ellipsis={{ tooltip: value }}>
    {value || '-'}
  </Text>
);

const renderTimeText = (value: string) => {
  const formatted = formatTimestamp(value) || '-';
  return (
    <Text className={styles.timeCell} ellipsis={{ tooltip: formatted }}>
      {formatted}
    </Text>
  );
};

const flexibleColumnKeys = ['scope', 'module', 'action', 'resource', 'operator', 'result'] as const;
type FlexibleColumnKey = (typeof flexibleColumnKeys)[number];

const distributeExtraWidth = (widths: Record<FlexibleColumnKey, number>, extraWidth: number) => {
  let restWidth = extraWidth;
  while (restWidth > 0.5) {
    const expandableColumns = flexibleColumnKeys.filter((key) => widths[key] < flexibleColumnConfig[key].max);
    if (!expandableColumns.length) break;
    const widthPerColumn = restWidth / expandableColumns.length;
    let usedWidth = 0;
    expandableColumns.forEach((key) => {
      const addedWidth = Math.min(widthPerColumn, flexibleColumnConfig[key].max - widths[key]);
      widths[key] += addedWidth;
      usedWidth += addedWidth;
    });
    if (usedWidth <= 0) break;
    restWidth -= usedWidth;
  }
};

const resolveFlexibleColumnWidths = (
  records: AuditLogItem[],
  containerWidth: number,
  translate: (id: string, opt?: any, defaultMessage?: string) => string
) => {
  const widths = {
    scope: flexibleColumnConfig.scope.min,
    module: flexibleColumnConfig.module.min,
    action: flexibleColumnConfig.action.min,
    resource: flexibleColumnConfig.resource.min,
    operator: flexibleColumnConfig.operator.min,
    result: flexibleColumnConfig.result.min,
  };

  records.forEach((record) => {
    widths.scope = Math.max(
      widths.scope,
      estimateTextWidth(getAuditScopeLabel(record, translate), flexibleColumnConfig.scope)
    );
    widths.module = Math.max(
      widths.module,
      estimateTextWidth(
        getAuditModuleLabel(record.scopeType, record.scopeName, record.resType, translate),
        flexibleColumnConfig.module
      )
    );
    widths.action = Math.max(
      widths.action,
      estimateTextWidth(getAuditActionLabel(record.businessType, translate), flexibleColumnConfig.action)
    );
    widths.resource = Math.max(
      widths.resource,
      estimateTextWidth(getAuditResourceLabel(record, translate), flexibleColumnConfig.resource)
    );
    widths.operator = Math.max(
      widths.operator,
      estimateTextWidth(getAuditOperatorLabel(record), flexibleColumnConfig.operator)
    );
    widths.result = Math.max(
      widths.result,
      estimateTextWidth(getAuditResultLabel(record, translate), flexibleColumnConfig.result)
    );
  });

  const fixedWidth = fixedColumnWidth.time + fixedColumnWidth.operation;
  const availableWidth = Math.max(
    containerWidth ? containerWidth - fixedWidth : 0,
    flexibleColumnConfig.scope.min +
      flexibleColumnConfig.module.min +
      flexibleColumnConfig.action.min +
      flexibleColumnConfig.resource.min +
      flexibleColumnConfig.operator.min +
      flexibleColumnConfig.result.min
  );
  const flexibleWidth =
    widths.scope + widths.module + widths.action + widths.resource + widths.operator + widths.result;

  if (flexibleWidth > availableWidth) {
    const shrinkableWidth =
      widths.scope -
      flexibleColumnConfig.scope.min +
      widths.module -
      flexibleColumnConfig.module.min +
      widths.action -
      flexibleColumnConfig.action.min +
      widths.resource -
      flexibleColumnConfig.resource.min +
      widths.operator -
      flexibleColumnConfig.operator.min +
      widths.result -
      flexibleColumnConfig.result.min;
    const overflow = flexibleWidth - availableWidth;
    if (shrinkableWidth > 0) {
      const shrinkRatio = Math.min(1, overflow / shrinkableWidth);
      widths.scope -= (widths.scope - flexibleColumnConfig.scope.min) * shrinkRatio;
      widths.module -= (widths.module - flexibleColumnConfig.module.min) * shrinkRatio;
      widths.action -= (widths.action - flexibleColumnConfig.action.min) * shrinkRatio;
      widths.resource -= (widths.resource - flexibleColumnConfig.resource.min) * shrinkRatio;
      widths.operator -= (widths.operator - flexibleColumnConfig.operator.min) * shrinkRatio;
      widths.result -= (widths.result - flexibleColumnConfig.result.min) * shrinkRatio;
    }
  } else if (containerWidth && flexibleWidth < availableWidth) {
    distributeExtraWidth(widths, availableWidth - flexibleWidth);
  }

  return {
    scope: Math.round(widths.scope),
    module: Math.round(widths.module),
    action: Math.round(widths.action),
    resource: Math.round(widths.resource),
    operator: Math.round(widths.operator),
    result: Math.round(widths.result),
  };
};

const AuditLogPage = () => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const translate = (id: string, opt?: any, defaultMessage?: string) => formatMessage(id, opt, defaultMessage);
  const t = (id: string, fallback: string, opt?: Record<string, unknown>) => translate(id, opt, fallback);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [data, setData] = useState<AuditLogItem[]>([]);
  const [total, setTotal] = useState(0);
  const [pageNo, setPageNo] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [operatorKeyword, setOperatorKeyword] = useState('');
  const [detail, setDetail] = useState<AuditLogItem | null>(null);
  const [detailOpen, setDetailOpen] = useState(false);
  const [tableShellWidth, setTableShellWidth] = useState(0);
  const tableShellRef = useRef<HTMLDivElement>(null);
  const tableAvailableWidth = tableShellWidth ? Math.max(0, tableShellWidth - tableScrollbarReserve) : 0;
  const auditTableWidth = Math.max(tableAvailableWidth, auditMinimumTableWidth);
  const flexibleColumnWidth = resolveFlexibleColumnWidths(data, auditTableWidth, translate);

  const fetchList = async (nextPageNo = pageNo, nextPageSize = pageSize, nextKeyword = operatorKeyword) => {
    setLoading(true);
    try {
      const resp = await getAuditLogPage({
        pageNo: nextPageNo,
        pageSize: nextPageSize,
        operatorKeyword: nextKeyword || undefined,
      });
      setData(resp?.data || []);
      setTotal(resp?.total || 0);
      setPageNo(resp?.pageNo || nextPageNo);
      setPageSize(resp?.pageSize || nextPageSize);
    } catch {
      message.error(t('auditLog.loadFailed', 'Failed to load audit logs'));
    } finally {
      setLoading(false);
    }
  };

  // Initial page load uses the default filter state once.

  useEffect(() => {
    void fetchList(1, 20, '');
  }, []);

  useEffect(() => {
    const element = tableShellRef.current;
    if (!element) return undefined;
    const updateWidth = () => setTableShellWidth(Math.floor(getElementContentWidth(element)));
    updateWidth();
    const resizeObserver = new ResizeObserver(updateWidth);
    resizeObserver.observe(element);
    return () => resizeObserver.disconnect();
  }, []);

  const openDetail = async (record: AuditLogItem) => {
    setDetailOpen(true);
    setDetailLoading(true);
    try {
      const resp = await getAuditLogDetail(record.id);
      setDetail(resp);
    } catch {
      message.error(t('auditLog.detailFailed', 'Failed to load audit log detail'));
    } finally {
      setDetailLoading(false);
    }
  };

  return (
    <ComLayout loading={loading}>
      <ComContent
        className={styles.auditContent}
        title={t('auditLog.title', 'Audit Log')}
        hasBack={false}
        style={{ display: 'flex', flexDirection: 'column' }}
        extra={
          <Space>
            <ProSearch
              className={styles.auditSearch}
              size="sm"
              value={operatorKeyword}
              placeholder={t('auditLog.operatorPlaceholder', 'Search operator')}
              onChange={(e) => setOperatorKeyword(e.target.value)}
              onSearch={(value) => {
                setOperatorKeyword(value);
                void fetchList(1, pageSize, value);
              }}
              style={{ width: 260 }}
            />
          </Space>
        }
      >
        <div ref={tableShellRef} style={{ width: '100%', minWidth: 0 }}>
          <ProTable
            key={auditTableWidth}
            className={styles.auditTable}
            rowKey="id"
            dataSource={data}
            tableLayout="fixed"
            columnFit={false}
            scroll={{ y: 'calc(100vh - 240px)', x: auditTableWidth }}
            columns={[
              {
                title: () => t('auditLog.time', 'Time'),
                dataIndex: 'createdAt',
                key: 'createdAt',
                ellipsis: true,
                width: fixedColumnWidth.time,
                render: (value: string) => renderTimeText(value),
              },
              {
                title: () => t('auditLog.scope', 'Scope'),
                dataIndex: 'scopeName',
                key: 'scopeName',
                ellipsis: true,
                width: flexibleColumnWidth.scope,
                render: (_: unknown, record: any) =>
                  renderEllipsisText(getAuditScopeLabel(record as AuditLogItem, translate)),
              },
              {
                title: () => t('auditLog.module', 'Module'),
                dataIndex: 'resType',
                key: 'resType',
                ellipsis: true,
                width: flexibleColumnWidth.module,
                render: (_: unknown, record: any) =>
                  renderEllipsisText(
                    getAuditModuleLabel(
                      (record as AuditLogItem).scopeType,
                      (record as AuditLogItem).scopeName,
                      (record as AuditLogItem).resType,
                      translate
                    )
                  ),
              },
              {
                title: () => t('auditLog.action', 'Action'),
                dataIndex: 'businessType',
                key: 'businessType',
                ellipsis: true,
                width: flexibleColumnWidth.action,
                render: (value: string) => renderEllipsisText(getAuditActionLabel(value, translate)),
              },
              {
                title: () => t('auditLog.resource', 'Resource'),
                dataIndex: 'resName',
                key: 'resName',
                ellipsis: true,
                width: flexibleColumnWidth.resource,
                render: (_: unknown, record: any) =>
                  renderEllipsisText(getAuditResourceLabel(record as AuditLogItem, translate)),
              },
              {
                title: () => t('auditLog.operator', 'Operator'),
                dataIndex: 'operatorName',
                key: 'operatorName',
                ellipsis: true,
                width: flexibleColumnWidth.operator,
                render: (_: unknown, record: any) => renderEllipsisText(getAuditOperatorLabel(record as AuditLogItem)),
              },
              {
                title: () => t('auditLog.result', 'Result'),
                dataIndex: 'code',
                key: 'code',
                ellipsis: true,
                width: flexibleColumnWidth.result,
                render: (_: unknown, record: any) =>
                  renderEllipsisText(getAuditResultLabel(record as AuditLogItem, translate)),
              },
              {
                title: () => t('common.operation', 'Operation'),
                key: 'operation',
                width: fixedColumnWidth.operation,
                fixed: 'right',
                render: (_: unknown, record: any) => (
                  <a className={styles.operationLink} onClick={() => void openDetail(record)}>
                    {t('auditLog.detail', 'Detail')}
                  </a>
                ),
              },
            ]}
            pagination={{
              total,
              current: pageNo,
              pageSize,
              onChange: (current, size) => {
                void fetchList(current, size, operatorKeyword);
              },
            }}
          />
        </div>
        <AuditLogDetailDrawer
          open={detailOpen}
          loading={detailLoading}
          detail={detail}
          onClose={() => setDetailOpen(false)}
        />
      </ComContent>
    </ComLayout>
  );
};

export default AuditLogPage;
