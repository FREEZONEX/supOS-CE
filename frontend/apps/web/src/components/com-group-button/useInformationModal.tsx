import { useRef, useState } from 'react';
import ComSelect from '../com-select';
import ProModal from '../pro-modal';
import { AuthButton } from '../auth';
import { App, Empty, Flex, Tag, Pagination, type PaginationProps } from 'antd';
import { useTranslate, usePagination } from '@/hooks';
import styles from './InformationModal.module.scss';
import { confirmAlarm, getAlarmList } from '@/apis/inter-api/alarm.ts';
import { getAlertForSelect } from '@/apis/inter-api/uns.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import InfoList from '../com-group-button/InfoList.tsx';
import Loading from '../loading';
import { useNavigate } from 'react-router';
import classNames from 'classnames';
import ComRadio from '../com-radio/index.tsx';
import { formatTimestamp } from '@/utils/format.ts';

const useInformationModal = ({ onCallBack }: any) => {
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
  const [open, setOpen] = useState(false);
  const payload = useRef<any>(null);
  const [status, setStatus] = useState<string | undefined>(undefined);
  const formatMessage = useTranslate();
  // topic 的id
  const [topicValue, setTopicValue] = useState();
  const {
    data: originData,
    pagination,
    setSearchParams,
    reload,
    loading,
  } = usePagination({
    initPageSize: 10,
    fetchApi: getAlarmList,
    firstNotGetData: true,
  });

  const data = originData?.map?.((item: any) => {
    const { field, path } = JSON.parse(item.refers || '[]')?.[0] || {};
    return {
      ...item,
      field,
      path,
    };
  });

  const formulaObj: any = {
    '>': formatMessage('rule.greaterThanThreshold'),
    '<': formatMessage('rule.lessThanThreshold'),
    '<=': formatMessage('rule.lessEqualThreshold'),
    '>=': formatMessage('rule.greaterEqualThreshold'),
    '=': formatMessage('rule.equalThreshold'),
    '!=': formatMessage('rule.noEqualThreshold'),
  };

  const onOpen = (data: any) => {
    payload.current = data;
    setSearchParams({
      unsId: data.unsId,
      readStatus: undefined,
    });
    setTopicValue(data.unsId);
    setOpen(true);
    navigate('', { replace: true, state: {} });
  };

  const onClose = () => {
    setOpen(false);
    setStatus(undefined);
  };
  const showTotal: PaginationProps['showTotal'] = (total) =>
    `${formatMessage('common.total')}  ${total}  ${formatMessage('common.items')}`;
  const Dom = (
    <ProModal
      open={open}
      onCancel={onClose}
      size="xs"
      title={formatMessage('Alert.alert')}
      className={classNames(styles['information-modal'])}
    >
      <Loading spinning={loading}>
        <Flex vertical style={{ height: '100%', overflow: 'hidden' }}>
          <Flex align="center" gap={10} wrap justify={'flex-start'} style={{ padding: '10px 1rem' }}>
            <ComSelect
              key={topicValue}
              isRequest={open}
              variant="filled"
              style={{ width: 150 }}
              value={topicValue}
              allowClear
              onChange={(v) => {
                setTopicValue(v);
                setSearchParams({
                  readStatus: status ? status !== 'pending' : undefined,
                  unsId: v,
                });
              }}
              api={() => getAlertForSelect({ page: 1, pageSize: 10000, type: 5 })}
            />
            <Flex flex={1}>
              <ComRadio
                options={[
                  { label: formatMessage('common.unconfirmed'), value: 'pending' },
                  { label: formatMessage('common.confirmed'), value: 'processed' },
                ]}
                value={status}
                onClick={(e) => {
                  const val = e.target.value;
                  setStatus((prevState) => {
                    if (prevState === val) {
                      setSearchParams({
                        readStatus: undefined,
                        unsId: topicValue,
                      });
                      return undefined;
                    }
                    setSearchParams({
                      readStatus: val !== 'pending',
                      unsId: topicValue,
                    });
                    return val;
                  });
                }}
              />
            </Flex>

            <AuthButton
              auth={ButtonPermission['Alert.confirm']}
              size="small"
              type="primary"
              disabled={!data.some((item: any) => item.canHandler) || data.every((item: any) => item.readStatus)}
              onClick={() => {
                modal.confirm({
                  zIndex: 9999,
                  title: formatMessage('common.confirmOpt'),
                  onOk: () => {
                    confirmAlarm({
                      confirmType: 2,
                      unsId: topicValue,
                    }).then(() => {
                      onCallBack?.();
                      reload();
                      message.success(formatMessage('common.optsuccess'));
                    });
                  },
                  afterClose: () => {},
                  okButtonProps: {
                    title: formatMessage('common.confirm'),
                  },
                  cancelButtonProps: {
                    title: formatMessage('common.cancel'),
                  },
                });
              }}
            >
              {formatMessage('common.confirmAll')}
            </AuthButton>
          </Flex>
          {data?.length > 0 ? (
            <>
              <InfoList
                items={data?.map((item: any) => ({
                  key: item.id,
                  label: (
                    <span>
                      {!item.readStatus && (
                        <Tag
                          color="#DA1E28"
                          style={{
                            padding: 2,
                            lineHeight: 1,
                            borderRadius: 10,
                            display: 'inline-flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                          }}
                        >
                          {formatMessage('common.new')}
                        </Tag>
                      )}

                      {item.ruleName}
                    </span>
                  ),
                  extra: <span>{formatMessage('Alert.alert')}</span>,
                  children: (
                    <div>
                      <div
                        style={{ opacity: item.readStatus ? 0.8 : 1 }}
                      >{`【${item.unsPath}】.【${item.field}】${formatMessage('rule.in')} ${formatTimestamp(item.createAt)} ${item.isAlarm ? formulaObj?.[item?.condition || '>'] + '【' + item.limitValue + '】' : formatMessage('rule.alertCancel')}，${formatMessage('rule.currentValue')}【${item.currentValue}】${item.isAlarm ? '，' + formatMessage('rule.deal') : ''}`}</div>
                      <Flex justify="flex-end" style={{ marginTop: 4 }}>
                        {!item.readStatus || (item.canHandler && !item.readStatus) ? (
                          <AuthButton
                            auth={ButtonPermission['Alert.confirm']}
                            size="small"
                            type="primary"
                            disabled={!item.canHandler}
                            onClick={() => {
                              modal.confirm({
                                zIndex: 9999,
                                title: formatMessage('common.confirmOpt'),
                                onOk: () => {
                                  confirmAlarm({
                                    confirmType: 1,
                                    unsId: topicValue,
                                    ids: [item.id],
                                  }).then(() => {
                                    onCallBack?.();
                                    reload();
                                    message.success(formatMessage('common.optsuccess'));
                                  });
                                },
                                afterClose: () => {},
                                onCancel: () => {},
                                okButtonProps: {
                                  title: formatMessage('common.confirm'),
                                },
                                cancelButtonProps: {
                                  title: formatMessage('common.cancel'),
                                },
                              });
                            }}
                          >
                            {formatMessage('common.confirm')}
                          </AuthButton>
                        ) : (
                          <AuthButton size="small" disabled style={{ cursor: 'inherit' }}>
                            {formatMessage('common.confirmed')}
                          </AuthButton>
                        )}
                      </Flex>
                    </div>
                  ),
                }))}
              />
            </>
          ) : (
            <Empty style={{ padding: '20px 1rem' }} />
          )}
          <div style={{ padding: '0 1rem' }}>
            <Pagination
              size="small"
              align="end"
              pageSize={pagination?.pageSize || 20}
              current={pagination?.page}
              className="custom-pagination-info"
              total={pagination?.total}
              onChange={pagination.onChange}
              showTotal={showTotal}
              onShowSizeChange={(current, size) => {
                pagination.onChange({ page: current, pageSize: size });
              }}
            />
          </div>
        </Flex>
      </Loading>
    </ProModal>
  );
  return {
    ModalDom: Dom,
    onOpen,
  };
};

export default useInformationModal;
