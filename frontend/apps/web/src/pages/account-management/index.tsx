import { Flex, message, App, Tag, Tooltip } from 'antd';
import { useTranslate, usePagination, useMediaSize } from '@/hooks';
import { Edit, FlashOff, Password, PlayOutline, TrashCan } from '@/components/lucide-icon/carbon';
import { deleteUser, getUserManageList, updateUser } from '@/apis/core-api/user-manage';
import useResetPassword from '@/pages/account-management/components/useResetPassword';
import useAddUser from '@/pages/account-management/components/useAddUser';
import { cloneElement, useEffect, useRef, useState, type FC, type ReactElement } from 'react';
import type { PageProps } from '@/common-types';
import { ButtonPermission } from '@/common-types/button-permission';
import type { PaginationProps } from 'antd';
import { AuthButton, AuthWrapper } from '@/components/auth';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ProTable from '@/components/pro-table';
import { useBaseStore } from '@/stores/base';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import styles from './index.module.scss';

const apiObj: any = {
  updateUser,
};

const fixedColumnWidth = {
  status: 110,
  operation: 160,
};
const tableScrollbarReserve = 12;

const flexibleColumnConfig = {
  account: { min: 110, max: 220, padding: 46 },
  displayName: { min: 130, max: 220, padding: 46 },
  phone: { min: 130, max: 180, padding: 46 },
  email: { min: 220, max: 320, padding: 46 },
};
const flexibleColumnKeys = ['account', 'displayName', 'phone', 'email'] as const;
type FlexibleColumnKey = (typeof flexibleColumnKeys)[number];
const minFlexibleWidth = flexibleColumnKeys.reduce((sum, key) => sum + flexibleColumnConfig[key].min, 0);
const absoluteFlexibleMin = 72;

const clampWidth = (width: number, min: number, max: number) => Math.min(max, Math.max(min, width));

const estimateTextWidth = (value: string, config: { min: number; max: number; padding: number }) =>
  clampWidth(config.padding + value.length * 7, config.min, config.max);

const getElementContentWidth = (element: HTMLElement) => {
  const style = window.getComputedStyle(element);
  const horizontalPadding = (Number.parseFloat(style.paddingLeft) || 0) + (Number.parseFloat(style.paddingRight) || 0);
  return Math.max(0, element.getBoundingClientRect().width - horizontalPadding);
};

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
  if (restWidth > 0) {
    widths.email += restWidth;
  }
};

const resolveFlexibleColumnWidths = (records: any[], containerWidth: number) => {
  const widths: Record<FlexibleColumnKey, number> = {
    account: flexibleColumnConfig.account.min,
    displayName: flexibleColumnConfig.displayName.min,
    phone: flexibleColumnConfig.phone.min,
    email: flexibleColumnConfig.email.min,
  };

  records.forEach((record) => {
    widths.account = Math.max(
      widths.account,
      estimateTextWidth(record?.preferredUsername || '', flexibleColumnConfig.account)
    );
    widths.displayName = Math.max(
      widths.displayName,
      estimateTextWidth(record?.firstName || '', flexibleColumnConfig.displayName)
    );
    widths.phone = Math.max(widths.phone, estimateTextWidth(record?.phone || '', flexibleColumnConfig.phone));
    widths.email = Math.max(widths.email, estimateTextWidth(record?.email || '', flexibleColumnConfig.email));
  });

  const fixedWidth = fixedColumnWidth.status + fixedColumnWidth.operation;
  const availableWidth = containerWidth > 0 ? Math.max(0, containerWidth - fixedWidth) : minFlexibleWidth;
  const flexibleWidth = widths.account + widths.displayName + widths.phone + widths.email;

  if (flexibleWidth > availableWidth) {
    const shrinkableWidth =
      widths.account -
      flexibleColumnConfig.account.min +
      widths.displayName -
      flexibleColumnConfig.displayName.min +
      widths.phone -
      flexibleColumnConfig.phone.min +
      widths.email -
      flexibleColumnConfig.email.min;
    const overflow = flexibleWidth - availableWidth;
    if (shrinkableWidth > 0) {
      const shrinkRatio = Math.min(1, overflow / shrinkableWidth);
      widths.account -= (widths.account - flexibleColumnConfig.account.min) * shrinkRatio;
      widths.displayName -= (widths.displayName - flexibleColumnConfig.displayName.min) * shrinkRatio;
      widths.phone -= (widths.phone - flexibleColumnConfig.phone.min) * shrinkRatio;
      widths.email -= (widths.email - flexibleColumnConfig.email.min) * shrinkRatio;
    } else if (availableWidth > 0) {
      const scale = availableWidth / flexibleWidth;
      flexibleColumnKeys.forEach((key) => {
        widths[key] = Math.max(absoluteFlexibleMin, Math.round(widths[key] * scale));
      });
      const currentSum = flexibleColumnKeys.reduce((sum, key) => sum + widths[key], 0);
      widths.email = Math.max(absoluteFlexibleMin, widths.email + (availableWidth - currentSum));
    }
  } else if (containerWidth && flexibleWidth < availableWidth) {
    distributeExtraWidth(widths, availableWidth - flexibleWidth);
  }

  return {
    account: Math.round(widths.account),
    displayName: Math.round(widths.displayName),
    phone: Math.round(widths.phone),
    email: Math.round(widths.email),
  };
};

const getTableScrollX = (widths: ReturnType<typeof resolveFlexibleColumnWidths>) =>
  widths.account +
  widths.displayName +
  widths.phone +
  widths.email +
  fixedColumnWidth.status +
  fixedColumnWidth.operation;

const AccountManagement: FC<PageProps> = ({ title }) => {
  const [tableShellWidth, setTableShellWidth] = useState(0);
  const tableShellRef = useRef<HTMLDivElement>(null);
  const formatMessage = useTranslate();
  const ldapEnable = useBaseStore((state) => state?.systemInfo?.ldapEnable);
  const { modal } = App.useApp();
  const { isH5 } = useMediaSize();
  const { data, pagination, setLoading, loading, refreshRequest } = usePagination({
    initPageSize: 100,
    fetchApi: getUserManageList,
  });
  const records = Array.isArray(data) ? data : [];
  const tableAvailableWidth = tableShellWidth ? Math.max(0, tableShellWidth - tableScrollbarReserve) : 0;
  const flexibleColumnWidth = resolveFlexibleColumnWidths(records, tableAvailableWidth);
  const tableScrollX = getTableScrollX(flexibleColumnWidth);

  useEffect(() => {
    const element = tableShellRef.current;
    if (!element) return undefined;
    const updateWidth = () => setTableShellWidth(Math.floor(getElementContentWidth(element)));
    updateWidth();
    const resizeObserver = new ResizeObserver(updateWidth);
    resizeObserver.observe(element);
    return () => resizeObserver.disconnect();
  }, []);

  const showTotal: PaginationProps['showTotal'] = (total) =>
    isH5 ? null : `${formatMessage('common.total')}  ${total}  ${formatMessage('common.items')}`;
  const handle = (params: any, apiKey: string) => {
    setLoading(true);
    apiObj?.[apiKey]?.(params)
      .then(() => {
        message.success(formatMessage('common.optsuccess'));
      })
      .finally(() => {
        refreshRequest();
        setLoading(false);
      });
  };

  const { ModalDom, onOpen } = useResetPassword({
    onSaveBack: refreshRequest,
  });
  const { ModalAddDom, onAddOpen } = useAddUser({
    onSaveBack: refreshRequest,
  });
  const onAddHandle = () => {
    onAddOpen();
  };
  const renderStatusTag = (enabled: boolean) =>
    enabled ? (
      <Tag bordered={false} className={`${styles.statusPill} ${styles.statusAvailable}`}>
        {formatMessage('account.available')}
      </Tag>
    ) : (
      <Tag bordered={false} className={`${styles.statusPill} ${styles.statusUnavailable}`}>
        {formatMessage('account.unavailable')}
      </Tag>
    );

  const renderOperationIcon = ({
    auth,
    disabled,
    label,
    icon,
    onClick,
  }: {
    auth: string;
    disabled?: boolean;
    label: string;
    icon: ReactElement<any>;
    onClick: () => void;
  }) => (
    <AuthWrapper auth={auth}>
      <Tooltip title={label}>
        <span className={styles.operationIcon}>
          {cloneElement(icon, {
            className: `${icon.props.className ? `${icon.props.className} ` : ''}custom-operation ${
              disabled ? styles.operationDisabled : ''
            }`,
            style: { cursor: disabled ? 'not-allowed' : 'pointer', ...(icon.props.style || {}) },
            onClick: disabled ? undefined : onClick,
            'aria-label': label,
            role: 'button',
          })}
        </span>
      </Tooltip>
    </AuthWrapper>
  );

  const renderStatusToggleIcon = (record: any) => {
    if (record?.preferredUsername === 'tier0') return null;
    const statusDisabled = ldapEnable || record?.source === 'external';
    if (record?.enabled) {
      return renderOperationIcon({
        auth: ButtonPermission['UserManagement.disable'],
        disabled: statusDisabled,
        label: formatMessage('account.disable'),
        icon: <FlashOff size={16} />,
        onClick: () => {
          handle(
            {
              userId: record.id,
              enabled: false,
              roleList: record.roleList,
            },
            'updateUser'
          );
        },
      });
    }
    return renderOperationIcon({
      auth: ButtonPermission['UserManagement.enable'],
      disabled: statusDisabled,
      label: formatMessage('account.enable'),
      icon: <PlayOutline size={16} />,
      onClick: () => {
        handle(
          {
            userId: record.id,
            enabled: true,
            roleList: record.roleList,
          },
          'updateUser'
        );
      },
    });
  };

  const columns: any = [
    {
      dataIndex: 'preferredUsername',
      ellipsis: true,
      titleIntlId: 'account.account',
      width: flexibleColumnWidth.account,
      minWidth: flexibleColumnConfig.account.min,
      maxWidth: flexibleColumnConfig.account.max,
    },
    {
      dataIndex: 'firstName',
      ellipsis: true,
      titleIntlId: 'account.displayName',
      width: flexibleColumnWidth.displayName,
      minWidth: flexibleColumnConfig.displayName.min,
      maxWidth: flexibleColumnConfig.displayName.max,
    },
    {
      dataIndex: 'phone',
      ellipsis: true,
      titleIntlId: 'account.phone',
      width: flexibleColumnWidth.phone,
      minWidth: flexibleColumnConfig.phone.min,
      maxWidth: flexibleColumnConfig.phone.max,
    },
    {
      dataIndex: 'email',
      ellipsis: true,
      titleIntlId: 'account.email',
      width: flexibleColumnWidth.email,
      minWidth: flexibleColumnConfig.email.min,
      maxWidth: flexibleColumnConfig.email.max,
    },
    {
      dataIndex: 'enabled',
      titleIntlId: 'common.status',
      fixed: 'right',
      render: (enabled: boolean) => renderStatusTag(enabled),
      width: fixedColumnWidth.status,
      minWidth: fixedColumnWidth.status,
    },
    {
      dataIndex: 'operation',
      titleIntlId: 'common.operation',
      fixed: 'right',
      width: fixedColumnWidth.operation,
      minWidth: fixedColumnWidth.operation,
      render: (_: any, record: any) => (
        <span className={styles.operation}>
          {renderOperationIcon({
            auth: ButtonPermission['UserManagement.edit'],
            disabled: record?.source === 'external',
            label: formatMessage('common.edit'),
            icon: <Edit size={16} />,
            onClick: () => onAddOpen?.(record),
          })}
          {renderOperationIcon({
            auth: ButtonPermission['UserManagement.resetPassword'],
            disabled: ldapEnable && record?.preferredUsername !== 'tier0',
            label: formatMessage('account.resetpassword'),
            icon: <Password size={16} />,
            onClick: () => onOpen?.(record),
          })}
          {renderStatusToggleIcon(record)}
          {record?.preferredUsername !== 'tier0'
            ? renderOperationIcon({
                auth: ButtonPermission['UserManagement.delete'],
                disabled: (ldapEnable && record?.preferredUsername !== 'tier0') || record?.source === 'external',
                label: formatMessage('common.delete'),
                icon: <TrashCan size={16} />,
                onClick: () => {
                  modal.confirm({
                    ...createDeleteConfirmOptions({
                      title: formatMessage('common.deleteConfirm'),
                      name: record?.preferredUsername,
                      formatMessage,
                    }),
                    onOk: () => {
                      setLoading(true);
                      deleteUser(record.id)
                        .then(() => {
                          message.success(formatMessage('common.optsuccess'));
                          refreshRequest();
                        })
                        .finally(() => {
                          setLoading(false);
                        });
                    },
                  });
                },
              })
            : null}
        </span>
      ),
    },
  ];
  return (
    <ComLayout loading={loading}>
      <ComContent
        className={styles.accountContent}
        hasBack={false}
        title={title}
        extra={
          <Flex align="center" gap={8} className={styles.titleActions}>
            <AuthButton
              auth={ButtonPermission['UserManagement.add']}
              style={{ height: 28 }}
              onClick={onAddHandle}
              type="primary"
              disabled={ldapEnable}
            >
              {formatMessage('account.newUsers')}
            </AuthButton>
          </Flex>
        }
      >
        <div ref={tableShellRef} style={{ width: '100%', minWidth: 0 }}>
          <ProTable
            key={tableAvailableWidth || 'auto'}
            resizeable
            columnFit={false}
            className={styles.accountTable}
            style={{ height: '100%' }}
            scroll={{ y: 'calc(100vh  - 240px)', x: tableAvailableWidth || tableScrollX || '100%' }}
            dataSource={data as any}
            columns={columns}
            pagination={{
              total: pagination?.total,
              showTotal: showTotal,
              style: { display: 'flex', justifyContent: 'flex-end', padding: '10px 0' },
              pageSize: pagination?.pageSize || 20,
              current: pagination?.page,
              showQuickJumper: true,
              pageSizeOptions: pagination?.pageSizes,
              onChange: pagination.onChange,
              onShowSizeChange: (current, size) => {
                pagination.onChange({ page: current, pageSize: size });
              },
            }}
          />
        </div>
        {ModalDom}
        {ModalAddDom}
      </ComContent>
    </ComLayout>
  );
};

export default AccountManagement;
