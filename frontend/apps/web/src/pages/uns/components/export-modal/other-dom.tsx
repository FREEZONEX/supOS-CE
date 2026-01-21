import { useBaseStore } from '@/stores/base';
import { ConfigProvider, Form, type FormInstance } from 'antd';
import ProTreeSelect from '@/components/pro-tree-select/ProTreeSelect.tsx';
import { flowPage } from '@/apis/inter-api/flow.ts';
import { flowPage as EventFlowPage } from '@/apis/inter-api/event-flow.ts';
import useTranslate from '@/hooks/useTranslate.ts';
import { getDashboardList } from '@/apis/inter-api/uns';

const OtherDom = ({ form }: { form: FormInstance }) => {
  const { containerList, dashboardType } = useBaseStore((state) => ({
    containerList: state.containerList,
    dashboardType: state.dashboardType,
  }));
  const hasNodeRed = !!containerList?.aboutUs?.some((s) => s.name === 'nodered') || true;
  const hasEventflow = !!containerList?.aboutUs?.some((s) => s.name === 'eventflow') || true;
  const formatMessage = useTranslate();
  return (
    <div style={{ width: '100%' }}>
      <ConfigProvider
        theme={{
          components: {
            Form: {
              itemMarginBottom: 12,
            },
          },
        }}
      >
        <Form
          layout="vertical"
          name="exportForm"
          form={form}
          colon={false}
          style={{ color: 'var(--supos-text-color)' }}
          initialValues={{
            sourceFlowExportParam: [],
            eventFlowExportParam: [],
            dashboardExportParam: [],
          }}
          // disabled={loading}
        >
          {hasNodeRed && (
            <Form.Item label={formatMessage('home.sourceFlow')} name="sourceFlowExportParam">
              <ProTreeSelect
                lazy
                maxTagCount="responsive"
                treeCheckable
                api={(params) =>
                  flowPage({
                    k: params?.searchValue,
                    pageNo: params?.pageNo,
                    pageSize: params?.pageSize,
                  }).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        value: item.id,
                        title: item.flowName,
                        key: item.id,
                      })),
                    };
                  })
                }
              />
            </Form.Item>
          )}
          {hasEventflow && (
            <Form.Item label={formatMessage('home.eventFlow')} name="eventFlowExportParam">
              <ProTreeSelect
                allowClear
                maxTagCount="responsive"
                treeCheckable
                showSwitcherIcon={false}
                lazy={true}
                api={(params) =>
                  EventFlowPage({
                    k: params?.searchValue,
                    pageNo: params?.pageNo,
                    pageSize: params?.pageSize,
                  }).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        value: item.id,
                        title: item.flowName,
                        key: item.id,
                      })),
                    };
                  })
                }
              />
            </Form.Item>
          )}
          {
            <Form.Item label={formatMessage('home.dashboard')} name="dashboardExportParam">
              <ProTreeSelect
                allowClear
                maxTagCount="responsive"
                treeCheckable
                showSwitcherIcon={false}
                lazy={true}
                api={(params) =>
                  getDashboardList({
                    k: params?.searchValue,
                    pageNo: params?.pageNo,
                    pageSize: params?.pageSize,
                    type: dashboardType?.length >= 2 ? undefined : dashboardType?.includes('fuxa') ? 2 : 1,
                  }).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        value: item.id,
                        title: item.name,
                        key: item.id,
                      })),
                    };
                  })
                }
              />
            </Form.Item>
          }
        </Form>
      </ConfigProvider>
    </div>
  );
};

export default OtherDom;
