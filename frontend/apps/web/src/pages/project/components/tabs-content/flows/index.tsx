import { Image } from '@/components/lucide-icon/carbon';
import { Button, Flex } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router';

import { getProjectFlows } from '@/apis/core-api';
import ProTable from '@/components/pro-table';
import { useTranslate } from '@/hooks';
import { formatTimestamp } from '@/utils/format';

import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { getSearchParamsString } from '@/utils/url-util';
import TabLayout from '../../TabLayout';

interface FlowsProps {
  projectId?: string;
}

interface FlowItem {
  id?: string;
  appId?: string;
  appName?: string;
  flowId?: string;
  runtimeFlowId?: string;
  flowType?: number | string;
  flowName: string;
  associateApp: string;
  iconUrl?: string;
  template: string;
  updatedAt?: string | number;
}

const Flows: FC<FlowsProps> = ({ projectId }) => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const [flows, setFlows] = useState<FlowItem[]>();
  const [loading, setLoading] = useState(false);

  const getFlows = async (id: string) => {
    setLoading(true);
    try {
      const data = await getProjectFlows(id);
      setFlows(data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!projectId) {
      setFlows([]);
      return;
    }

    getFlows(projectId);
  }, [projectId]);

  useActivate(() => {
    if (!projectId) {
      setFlows([]);
      return;
    }
    getFlows(projectId);
  });

  const handleGoFlowEditor = useCallback(
    (record: FlowItem) => {
      const id = record.id || record.flowId;
      const flowId = record.runtimeFlowId || record.flowId || id;
      if (!id) {
        return;
      }
      const isEventFlow = Number(record.flowType) === 2 || String(record.flowType).toLowerCase() === 'event';
      const editorPath = isEventFlow ? '/event-flow/flow-editor' : '/collection-flow/flow-editor';

      navigate(
        `${editorPath}?${getSearchParamsString({
          id,
          name: record.flowName,
          flowId,
          from: location.pathname,
          copyEnable: 'false',
        })}`
      );
    },
    [navigate]
  );

  const columns: any = useMemo(() => {
    return [
      {
        title: () => formatMessage('common.name'),
        dataIndex: 'flowName',
        key: 'flowName',
        width: '30%',
        render: (value: string, record: FlowItem) => {
          if (!record.flowId) {
            return value;
          }

          return (
            <Button type="link" className="table-link-button" onClick={() => handleGoFlowEditor(record)} title={value}>
              {value}
            </Button>
          );
        },
      },
      {
        title: formatMessage('project.associateApp'),
        dataIndex: 'appName',
        key: 'appName',
        width: '30%',
        render: (value: string, record: FlowItem) => (
          <Flex align="center" gap={8} title={value}>
            <span
              style={{
                width: 18,
                height: 18,
                borderRadius: '50%',
                display: 'inline-flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: '#f4f4f4',
                overflow: 'hidden',
              }}
            >
              {record.iconUrl ? (
                <img style={{ maxWidth: '100%', maxHeight: '100%' }} src={record.iconUrl} />
              ) : (
                <Image size={12} />
              )}
            </span>
            <span
              style={{
                maxWidth: '80%',
                whiteSpace: 'nowrap',
                textOverflow: 'ellipsis',
                overflow: 'hidden',
                display: 'inline-block',
              }}
            >
              {value}
            </span>
          </Flex>
        ),
      },
      {
        title: () => formatMessage('collectionFlow.flowTemplate'),
        dataIndex: 'template',
        key: 'template',
        width: '18%',
        render: (_: unknown, record: FlowItem) =>
          Number(record.flowType) === 2 || String(record.flowType).toLowerCase() === 'event'
            ? 'event-flow'
            : 'source-flow',
      },
      {
        title: () => formatMessage('project.lastUpdated'),
        dataIndex: 'updatedAt',
        key: 'updatedAt',
        width: '18%',
        render: (value: string | number | undefined) =>
          value ? formatTimestamp(value, 'YYYY/MM/DD HH:mm', true) : '-',
      },
    ];
  }, [formatMessage, handleGoFlowEditor]);

  return (
    <TabLayout style={{ padding: '0 24px' }}>
      <ProTable
        loading={loading}
        scroll={{ x: 'max-content', y: 'calc(100vh - 300px)' }}
        columns={columns}
        dataSource={flows}
        rowKey="id"
        resizeable
        pagination={false}
      />
    </TabLayout>
  );
};

export default Flows;
