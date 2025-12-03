import { usePagination, useTranslate } from '@supos_host/hooks';
import { ButtonPermission } from './common-types/button-permission';
import { Alarm, Edit, Search, TrashCan, View } from '@carbon/icons-react';
import { useActivate } from '@supos_host/tabs-lifecycle-context';
import { getAlertList, addRule, deleteRule, editRule, searchUserManageList } from './apis';
import { type FC, useEffect, useState } from 'react';
import { useLocation } from 'react-router';
import { App, Empty, Flex, Form, Pagination } from 'antd';
import NameSpace from './components/NameSpace.tsx';
import Condition from './components/Condition.tsx';
import DeadZone from './components/DeadZone.tsx';
import {
  ComDrawer,
  useInformationModal,
  ComLayout,
  ComContent,
  ComSearch,
  OperationForm,
  ComSegmented,
  ProCardContainer,
  ProTable,
  AuthButton,
  ProCard,
  SecondaryList,
} from '@supos_host/components';
import { useLocalStorageState } from 'ahooks';

import { validInputPattern } from '@supos_host/utils';
import { REMOTE_NAME } from '../variables';

interface AlertProps {
  title: string;
}
const Alert: FC<AlertProps> = ({ title }) => {
  const location = useLocation();
  const businessId = location?.state?.businessId;
  const [searchForm] = Form.useForm();
  const formatMessage = useTranslate(REMOTE_NAME);
  const commonFormatMessage = useTranslate();
  const [isEdit, setIsEdit] = useState(false);
  const [show, setShow] = useState(false);
  const [mode, setMode] = useLocalStorageState<string>('SUPOS_ALERT_MODE', {
    defaultValue: 'card',
  });
  const { modal } = App.useApp();

  const formItemOptions = [
    {
      label: commonFormatMessage('common.name'),
      name: 'name',
      rules: [
        { required: true, message: commonFormatMessage('rule.required') },
        { pattern: validInputPattern, message: commonFormatMessage('rule.illegality') },
      ],
      properties: {
        placeholder: commonFormatMessage('common.name'),
      },
    },
    {
      label: commonFormatMessage('uns.description'),
      name: 'description',
      type: 'TextArea',
      properties: {
        placeholder: commonFormatMessage('uns.description'),
      },
    },
    {
      label: formatMessage('acceptType'),
      name: 'withFlags',
      type: 'Select',
      rules: [{ required: true, message: commonFormatMessage('rule.required') }],
      properties: {
        options: [
          { label: formatMessage('person'), value: 16 },
          // { label: formatMessage('alert.workflow'), value: 32 },
        ],
        placeholder: formatMessage('acceptType'),
        onClear: () => {
          form.setFieldsValue({ accept: undefined });
        },
      },
    },
    {
      label: formatMessage('accept'),
      name: 'accept',
      type: 'Select',
      rules: [{ required: true, message: commonFormatMessage('rule.required') }],
      properties: {
        isRequest: show,
        placeholder: formatMessage('person'),
        mode: 'multiple',
        filterOption: (input: any, option: any) =>
          ((option?.label as string) ?? '').toLowerCase().includes(input.toLowerCase()),
        showSearch: true,
        onChange: (_: any, option: any) => {
          const selected = Array.isArray(option)
            ? option.map((o) => ({ value: o.value, label: o.label }))
            : { value: option.value, label: option.label };

          form.setFieldsValue({ accept: selected });
        },
        api: (key: string) => searchUserManageList({ preferredUsername: key, page: 1, pageSize: 999 }),
      },
    },
    {
      type: 'divider',
    },
    {
      render: () => <NameSpace isEdit={isEdit} />,
    },
    {
      render: Condition,
    },
    {
      render: DeadZone,
    },
    {
      label: formatMessage('overDuration'),
      name: ['protocol', 'overTime'],
      type: 'Number',
      properties: {
        min: 1,
        precision: 0,
        placeholder: formatMessage('regularValue'),
      },
    },
    {
      label: 'id',
      name: 'id',
      hidden: true,
    },
    {
      type: 'divider',
    },
  ];
  const [form] = Form.useForm();
  const { loading, pagination, data, refreshRequest, reload, setSearchParams } = usePagination({
    fetchApi: getAlertList,
    initPageSize: 18,
    defaultParams: {
      type: 5,
    },
  });
  useEffect(() => {
    data?.map((e: any) => {
      if (e.id === businessId) {
        onOpen({
          unsId: e.id,
        });
      }
    });
  }, [businessId, data]);
  useActivate(() => {
    refreshRequest?.();
  });
  const onAddHandle = () => {
    setIsEdit(false);
    form.resetFields();
    if (show) return;
    setShow(true);
  };

  const onClose = () => {
    setShow(false);
    setIsEdit(false);
    form.resetFields();
    // setOptionsData([]);
  };
  const onSave = async () => {
    const values = await form.validateFields();
    console.log(values, 'values');
    let addData: any;
    if (values.withFlags === 16) {
      addData = {
        ...values,
        userList: values.accept.map((e: any) => {
          return {
            id: e.value,
            preferredUsername: e.label,
          };
        }),
        refers: values.refers?.map((e: any) => ({ id: e?.refer?.value, field: e.field })),
      };

      delete addData.accept;
    } else if (values.withFlags === 32) {
      addData = { ...values, extend: values.accept.value };
      delete addData.accept;
    }
    const api = isEdit ? editRule : addRule;
    await api({
      ...addData,
      expression: `a1${values?.protocol?.condition}${values?.protocol?.limitValue}`,
    })
      .then(() => {
        refreshRequest();
        onClose();
      })
      .finally(() => {});
  };
  const onDeleteHandle = (item: any) => {
    return deleteRule({
      id: item.id,
      withFlow: false,
      withDashboard: false,
    }).then(() => {
      reload();
    });
  };

  const { ModalDom, onOpen } = useInformationModal({
    onCallBack: refreshRequest,
  });

  const actions: any = (item: any) => {
    return [
      {
        key: 'edit',
        auth: ButtonPermission['Alert.edit'],
        label: formatMessage('edit'),
        button: {
          type: 'primary',
        },
        extra: (
          <Flex justify="center" align="center">
            <Edit />
          </Flex>
        ),
        onClick: () => {
          setIsEdit(true);
          setShow(true);
          form.setFieldsValue({
            id: item.id,
            refers: [
              {
                refer: { value: item.refUns },
                field: item.field,
              },
            ],
            name: item.name,
            description: item.description,
            protocol: item.alarmRuleDefine,
            withFlags: item.withFlags,
            accept:
              item.withFlags === 32
                ? { value: item.processDefinition.id, label: item.processDefinition.id.processDefinitionName }
                : item.handlerList.map((item: any) => {
                    return {
                      value: item.userId,
                      label: item.username,
                    };
                  }),
          });
        },
      },
      {
        key: 'show',
        auth: ButtonPermission['Alert.show'],
        label: formatMessage('show'),
        extra: (
          <Flex justify="center" align="center">
            <View />
          </Flex>
        ),
        onClick: () => {
          onOpen({
            // 报警自己的topic id
            unsId: item.id,
          });
        },
      },
      {
        key: 'delete',
        auth: ButtonPermission['Alert.delete'],
        label: formatMessage('delete'),
        extra: (
          <Flex justify="center" align="center">
            <TrashCan />
          </Flex>
        ),
        onClick: () => {
          return modal.confirm({
            title: commonFormatMessage('common.deleteConfirm'),
            onOk: () => onDeleteHandle?.(item),
            okButtonProps: {
              title: commonFormatMessage('common.confirm'),
            },
            cancelButtonProps: {
              title: commonFormatMessage('common.cancel'),
            },
          });
        },
      },
    ];
  };
  return (
    <ComLayout loading={loading}>
      <ComContent
        title={
          <Flex gap={8} align="center">
            <Alarm size={20} />
            <span title={title}>{title}</span>
          </Flex>
        }
        extra={
          <Flex gap={8}>
            <ComSearch
              form={searchForm}
              formItemOptions={[
                {
                  name: 'k',
                  properties: {
                    prefix: <Search />,
                    placeholder: formatMessage('search'),
                    style: { width: 300 },
                    allowClear: true,
                  },
                },
              ]}
              formConfig={{
                onFinish: () => {
                  setSearchParams(searchForm.getFieldsValue());
                },
              }}
              onSearch={() => {
                setSearchParams(searchForm.getFieldsValue());
              }}
            />
            <AuthButton auth={ButtonPermission['Alert.add']} type="primary" onClick={onAddHandle}>
              + {formatMessage('createAlert')}
            </AuthButton>
          </Flex>
        }
        hasBack={false}
        style={{
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
        }}
      >
        <ComSegmented value={mode} onChange={setMode} defaultValue="card" />
        <div style={{ flex: 1, padding: '0 16px 16px', overflow: 'auto', alignItems: 'center' }}>
          {mode === 'card' ? (
            data.length > 0 ? (
              <ProCardContainer>
                {data?.map((d: any) => {
                  return (
                    <ProCard
                      key={d.id}
                      header={{
                        title: d.name,
                        customIcon: (
                          <Flex align="center" justify="center">
                            <Alarm color="#DA1E28" size="28" />
                          </Flex>
                        ),
                      }}
                      description={d?.description}
                      item={d}
                      actions={actions}
                      secondaryDescription={
                        <SecondaryList
                          options={[
                            {
                              label: formatMessage('alterCount'),
                              content: d?.alarmCount,
                              span: 24,
                              key: 'alterCount',
                            },
                          ]}
                        />
                      }
                    />
                  );
                })}
              </ProCardContainer>
            ) : (
              <Empty />
            )
          ) : (
            <ProTable
              resizeable
              style={{ height: '100%' }}
              scroll={{ y: 'calc(100vh  - 285px)', x: 'max-content' }}
              dataSource={data as any}
              columns={[
                {
                  dataIndex: 'name',
                  ellipsis: true,
                  title: () => commonFormatMessage('common.name'),
                  width: '30%',
                },
                {
                  dataIndex: 'noReadCount',
                  ellipsis: true,
                  title: () => formatMessage('alterCount'),
                  width: '20%',
                },
                {
                  dataIndex: 'description',
                  ellipsis: true,
                  title: () => commonFormatMessage('uns.description'),
                  width: '30%',
                },
              ]}
              pagination={{
                total: pagination?.total,
                pageSize: pagination?.pageSize || 18,
                current: pagination?.page,
                showQuickJumper: true,
                pageSizeOptions: pagination.pageSizes,
                showSizeChanger: true,
                onChange: pagination.onChange,
                onShowSizeChange: (current: number, size: number) => {
                  pagination.onChange({ page: current, pageSize: size });
                },
              }}
              operationOptions={{
                render: actions,
              }}
            />
          )}
        </div>
        {mode === 'card' && (
          <Pagination
            size="small"
            className="custom-pagination"
            align="center"
            style={{ margin: '20px 0' }}
            showSizeChanger={false}
            total={pagination?.total}
            pageSize={pagination?.pageSize || 20}
            current={pagination?.page}
            onChange={pagination.onChange}
          />
        )}
        {ModalDom}
      </ComContent>
      <ComDrawer title=" " width={680} open={show} onClose={onClose}>
        {show && (
          <OperationForm
            formConfig={{
              labelCol: { span: 7 },
              wrapperCol: { span: 17 },
            }}
            title={isEdit ? formatMessage('editAlert') : formatMessage('createAlert')}
            form={form}
            onCancel={onClose}
            onSave={onSave}
            formItemOptions={formItemOptions}
          />
        )}
      </ComDrawer>
    </ComLayout>
  );
};

export default Alert;
