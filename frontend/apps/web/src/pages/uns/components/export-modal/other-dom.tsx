import { useBaseStore } from '@/stores/base';
import { ConfigProvider, Form, type FormInstance, TreeSelect } from 'antd';
import ProTreeSelect from '@/components/pro-tree-select/ProTreeSelect.tsx';
import { getFlowAndGroupList } from '@/apis/inter-api/flow.ts';
import { getEventFlowAndGroupList } from '@/apis/inter-api/event-flow.ts';
import useTranslate from '@/hooks/useTranslate.ts';
import { getDashboardAndGroupList } from '@/apis/inter-api/dashboard.ts';
import { FlowData, Folder } from '@carbon/icons-react';
const { SHOW_PARENT } = TreeSelect;

const TreeNode = (dataNode: any) => {
  if (dataNode.category === 'group') {
    return <Folder style={{ flexShrink: 0, marginRight: '5px' }} />;
  }
  return <FlowData style={{ flexShrink: 0, marginRight: '5px' }} />;
};

const OtherDom = ({ form }: { form: FormInstance }) => {
  const { containerList } = useBaseStore((state) => ({
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
                showSearch={false}
                loadDataEnable
                lazy
                listHeight={350}
                maxTagCount="responsive"
                showSwitcherIcon
                treeCheckable
                fieldNames={{ label: 'name', value: 'id' }}
                showCheckedStrategy={SHOW_PARENT}
                allowClear
                api={(params, config) =>
                  getFlowAndGroupList(
                    {
                      k: params?.searchValue,
                      groupId: params?.key ? params?.key : undefined,
                      pageNo: params?.pageNo,
                      pageSize: params?.pageSize,
                    },
                    config
                  ).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        isLeaf: !item.hasChildren,
                        value: item.id,
                        title: item.name,
                        key: item.id,
                      })),
                    };
                  })
                }
                treeNodeIcon={TreeNode}
              />
            </Form.Item>
          )}
          {hasEventflow && (
            <Form.Item label={formatMessage('home.eventFlow')} name="eventFlowExportParam">
              <ProTreeSelect
                showSearch={false}
                loadDataEnable
                lazy
                listHeight={350}
                maxTagCount="responsive"
                showSwitcherIcon
                treeCheckable
                fieldNames={{ label: 'name', value: 'id' }}
                showCheckedStrategy={SHOW_PARENT}
                allowClear
                api={(params, config) =>
                  getEventFlowAndGroupList(
                    {
                      k: params?.searchValue,
                      groupId: params?.key ? params?.key : undefined,
                      pageNo: params?.pageNo,
                      pageSize: params?.pageSize,
                    },
                    config
                  ).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        isLeaf: !item.hasChildren,
                        value: item.id,
                        title: item.name,
                        key: item.id,
                      })),
                    };
                  })
                }
                treeNodeIcon={TreeNode}
              />
            </Form.Item>
          )}
          {
            <Form.Item label={formatMessage('home.dashboard')} name="dashboardExportParam">
              <ProTreeSelect
                showSearch={false}
                loadDataEnable
                lazy
                listHeight={350}
                maxTagCount="responsive"
                showSwitcherIcon
                treeCheckable
                fieldNames={{ label: 'name', value: 'id' }}
                showCheckedStrategy={SHOW_PARENT}
                allowClear
                api={(params, config) =>
                  getDashboardAndGroupList(
                    {
                      k: params?.searchValue,
                      groupId: params?.key ? params?.key : undefined,
                      pageNo: params?.pageNo,
                      pageSize: params?.pageSize,
                    },
                    config
                  ).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        isLeaf: !item.hasChildren,
                        value: item.id,
                        title: item.name,
                        key: item.id,
                      })),
                    };
                  })
                }
                treeNodeIcon={TreeNode}
              />
            </Form.Item>
          }
        </Form>
      </ConfigProvider>
    </div>
  );
};

export default OtherDom;
