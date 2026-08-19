import { type FC, useCallback, useEffect, useMemo, useState } from 'react';
import { Line } from '@ant-design/charts';
import { Alert, Button, DatePicker, Flex, Spin, Table, Typography, type TimeRangePickerProps } from 'antd';
import { Renew } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { ComEmptyState } from '@/components';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import { useTranslate } from '@/hooks';
import { getUnsDashboardData } from '@/apis/core-api/uns';
import { formatTimestamp, simpleFormat } from '@/utils/format';
import { ThemeType, useThemeStore } from '@/stores/theme-store';

const { RangePicker } = DatePicker;

interface DetailDashboardProps {
  instanceInfo: { [key: string]: any };
}

interface DashboardField {
  name: string;
  type?: string;
}

interface DashboardPoint {
  timestamp: number;
  payload?: unknown;
}

const numericTypes = new Set(['integer', 'int', 'long', 'float', 'double', 'number', 'decimal']);

const isNumericField = (field: DashboardField) => numericTypes.has(String(field.type || '').toLowerCase());

const normalizeValue = (value: unknown) => {
  const text = simpleFormat(value);
  return typeof text === 'string' ? text : String(text);
};

const normalizePayload = (payload: unknown): Record<string, unknown> => {
  if (payload && typeof payload === 'object' && !Array.isArray(payload)) {
    return payload as Record<string, unknown>;
  }
  if (typeof payload === 'string') {
    try {
      const parsed = JSON.parse(payload);
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return parsed as Record<string, unknown>;
      }
    } catch {
      return { value: payload };
    }
  }
  return payload == null ? {} : { value: payload };
};

const DetailDashboard: FC<DetailDashboardProps> = ({ instanceInfo }) => {
  const formatMessage = useTranslate();
  const theme = useThemeStore((state) => state.theme);
  const [dates, setDates] = useState<[Dayjs, Dayjs]>([dayjs().add(-5, 'minute'), dayjs()]);
  const [fields, setFields] = useState<DashboardField[]>([]);
  const [list, setList] = useState<DashboardPoint[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const rangePresets: TimeRangePickerProps['presets'] = [
    {
      label: <span title={formatMessage('uns.last5minutes')}>{formatMessage('uns.last5minutes')}</span>,
      value: [dayjs().add(-5, 'm'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last30minutes')}>{formatMessage('uns.last30minutes')}</span>,
      value: [dayjs().add(-30, 'm'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last1hour')}>{formatMessage('uns.last1hour')}</span>,
      value: [dayjs().add(-1, 'h'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last6hours')}>{formatMessage('uns.last6hours')}</span>,
      value: [dayjs().add(-6, 'h'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last24hours')}>{formatMessage('uns.last24hours')}</span>,
      value: [dayjs().add(-24, 'h'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last1week')}>{formatMessage('uns.last1week')}</span>,
      value: [dayjs().add(-1, 'w'), dayjs()],
    },
  ];

  const loadDashboard = useCallback(async () => {
    if (!instanceInfo?.id) return;
    setLoading(true);
    setError('');
    try {
      const data = await getUnsDashboardData({
        nodeId: instanceInfo.id,
        timeStart: dates[0].valueOf(),
        timeEnd: dates[1].valueOf(),
        limit: 1000,
      });
      setFields(Array.isArray(data?.fields) ? data.fields : []);
      setList(Array.isArray(data?.list) ? data.list : []);
    } catch (err) {
      console.error(err);
      setFields([]);
      setList([]);
      setError(formatMessage('uns.dashboardLoadFailed'));
    } finally {
      setLoading(false);
    }
  }, [dates, formatMessage, instanceInfo?.id]);

  useEffect(() => {
    loadDashboard();
  }, [loadDashboard]);

  const normalizedList = useMemo(
    () => list.map((point) => ({ ...point, payload: normalizePayload(point.payload) })),
    [list]
  );

  const displayFields = useMemo<DashboardField[]>(() => {
    const byName = new Map<string, DashboardField>(fields.map((field) => [field.name, field] as const));
    normalizedList.forEach((point) => {
      Object.keys(point.payload).forEach((name) => {
        if (!byName.has(name)) {
          byName.set(name, { name });
        }
      });
    });
    return Array.from(byName.values());
  }, [fields, normalizedList]);

  // 与云端保持一致：只有所有展示字段都明确是数值类型时才画折线。
  // State/Action 的 _json 可能包含 schema 外字段，必须切到表格，避免静默隐藏数据。
  const chartCompatible = useMemo(
    () => displayFields.length > 0 && displayFields.every(isNumericField),
    [displayFields]
  );

  const chartData = useMemo(() => {
    if (!chartCompatible) return [];
    return normalizedList.flatMap((point) => {
      const timestamp = Number(point.timestamp);
      return displayFields
        .map((field) => ({
          time: dayjs(timestamp).format('MM-DD HH:mm:ss'),
          timestamp,
          field: field.name,
          value: Number(point.payload[field.name]),
        }))
        .filter((item) => Number.isFinite(item.value));
    });
  }, [chartCompatible, displayFields, normalizedList]);

  const lineConfig = useMemo(() => {
    const isDark = theme === ThemeType.Dark;
    const labelColor = isDark ? '#c6c6c6' : '#525252';
    const axisColor = isDark ? '#525252' : '#d9d9d9';
    const gridColor = isDark ? '#393939' : '#e8e8e8';

    return {
      data: chartData,
      xField: 'time',
      yField: 'value',
      colorField: 'field',
      height: 320,
      autoFit: true,
      smooth: true,
      point: { size: 2 },
      axis: {
        x: { title: false, labelFill: labelColor, lineStroke: axisColor, tickStroke: axisColor },
        y: { title: false, labelFill: labelColor, gridStroke: gridColor },
      },
      legend: {
        color: {
          position: 'bottom',
          layout: { justifyContent: 'center' },
          itemLabelFill: labelColor,
        },
      },
    } as any;
  }, [chartData, theme]);

  const tableColumns = useMemo(
    () => [
      {
        title: formatMessage('common.latestUpdate'),
        dataIndex: '__displayTimestamp',
        key: '__displayTimestamp',
        width: 190,
        fixed: 'left' as const,
        defaultSortOrder: 'descend' as const,
        sorter: (left: Record<string, any>, right: Record<string, any>) =>
          left.__timestampSortValue - right.__timestampSortValue,
      },
      ...displayFields.map((field) => ({
        title: field.name,
        dataIndex: ['payload', field.name],
        key: field.name,
        ellipsis: true,
        render: (value: unknown) => {
          const text = normalizeValue(value);
          return text ? (
            <Typography.Text copyable={{ text }} ellipsis style={{ maxWidth: 220 }}>
              {text}
            </Typography.Text>
          ) : (
            <span style={{ color: 'var(--ui-text-color)', opacity: 0.45 }}>-</span>
          );
        },
      })),
    ],
    [displayFields, formatMessage]
  );

  const tableData = useMemo(
    () =>
      normalizedList.map((point, index) => ({
        key: `${point.timestamp}-${index}`,
        __displayTimestamp: formatTimestamp(point.timestamp) || dayjs(point.timestamp).format('YYYY-MM-DD HH:mm:ss'),
        __timestampSortValue: Number(point.timestamp) || 0,
        payload: point.payload,
      })),
    [normalizedList]
  );

  const refreshRange = () => {
    const duration = dates[1].valueOf() - dates[0].valueOf();
    const end = dayjs();
    setDates([end.add(-duration, 'millisecond'), end]);
  };

  return (
    <Flex vertical gap={12}>
      <Flex gap={10} wrap="wrap">
        <RangePicker
          showTime
          format="YYYY-MM-DD HH:mm:ss"
          value={dates}
          onChange={(nextDates) => {
            if (nextDates?.[0] && nextDates?.[1]) {
              setDates([nextDates[0], nextDates[1]]);
            }
          }}
          presets={rangePresets}
        />
        <Button icon={<Renew {...toolbarIconProps} />} onClick={refreshRange} />
      </Flex>

      {error ? <Alert type="error" showIcon message={error} /> : null}
      <Spin spinning={loading}>
        {normalizedList.length === 0 ? (
          <ComEmptyState variant="inline" description={formatMessage('uns.dashboardNoData')} />
        ) : chartCompatible && chartData.length > 0 ? (
          <Line {...lineConfig} />
        ) : tableData.length > 0 && displayFields.length > 0 ? (
          <Table
            bordered
            size="middle"
            pagination={{ pageSize: 10, showSizeChanger: false }}
            columns={tableColumns}
            dataSource={tableData}
            scroll={{ x: Math.max(760, displayFields.length * 180 + 190) }}
          />
        ) : (
          <ComEmptyState variant="inline" description={formatMessage('uns.dashboardNoNumericFields')} />
        )}
      </Spin>
    </Flex>
  );
};

export default DetailDashboard;
