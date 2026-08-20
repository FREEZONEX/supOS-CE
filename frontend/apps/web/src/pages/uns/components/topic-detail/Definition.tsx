import type { FC } from 'react';
import { useMemo } from 'react';
import { App } from 'antd';
import { TrashCan } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useTranslate } from '@/hooks';
import { editModel } from '@/apis/core-api/uns';
import ProTable from '@/components/pro-table';
import { AuthWrapper } from '@/components/auth';
import type { FieldItem } from '@/pages/uns/types';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import DefinitionTypeTag from './DefinitionTypeTag';
import styles from './Definition.module.scss';

interface DefinitionProps {
  instanceInfo: { [key: string]: any };
  modelInfo?: { [key: string]: any };
  getModel?: (savedFields?: FieldItem[]) => void | Promise<unknown>;
  auth?: string;
  editable?: boolean;
}

const Definition: FC<DefinitionProps> = ({ instanceInfo, modelInfo, getModel, auth, editable = false }) => {
  const formatMessage = useTranslate();
  const { message, modal } = App.useApp();
  const { id } = instanceInfo || {};
  const tableFields = useMemo(() => {
    if (Array.isArray(modelInfo?.fields)) {
      return modelInfo.fields;
    }
    return Array.isArray(instanceInfo?.fields) ? instanceInfo.fields : [];
  }, [instanceInfo, modelInfo?.fields]);

  const handleDeleteField = (record: FieldItem) => {
    const sourceInfo = modelInfo || instanceInfo;
    const nextFields = tableFields.filter((field) => field.name !== record.name);

    modal.confirm({
      ...createDeleteConfirmOptions({
        title: formatMessage('common.deleteConfirm'),
        name: record.name,
        formatMessage,
      }),
      onOk: () =>
        editModel({
          ...sourceInfo,
          id: sourceInfo.id || instanceInfo.id,
          alias: sourceInfo.alias,
          fields: nextFields,
          extendFieldUsed: sourceInfo.extendFieldUsed || [],
        })
          .then(() => {
            message.success(formatMessage('uns.editSuccessful'));
            return getModel?.(nextFields);
          })
          .catch(() => {
            message.error(formatMessage('common.failed'));
          }),
    });
  };

  return (
    <ProTable
      key={id}
      bordered
      resizeable={false}
      columnFit={false}
      className={styles.schemaTable}
      columns={[
        {
          title: formatMessage('uns.key'),
          dataIndex: 'name',
          ellipsis: true,
          render: (text: string) => (
            <span className={styles.schemaKeyCell} title={text}>
              {text}
            </span>
          ),
        },
        {
          title: formatMessage('uns.type'),
          dataIndex: 'type',
          width: 128,
          ellipsis: true,
          render: (text: string) => <DefinitionTypeTag type={text} />,
        },
        {
          title: formatMessage('uns.unit'),
          dataIndex: 'unit',
          width: 96,
          ellipsis: true,
          render: (text: string) => (
            <span className={styles.schemaUnitCell} title={text}>
              {text || ''}
            </span>
          ),
        },
        {
          title: formatMessage('common.operation'),
          dataIndex: 'operation',
          width: 104,
          render: (_: unknown, record: FieldItem) =>
            editable ? (
              <AuthWrapper auth={auth}>
                <span className={styles.schemaActionCell} title={formatMessage('common.delete')}>
                  <TrashCan
                    {...toolbarIconProps}
                    className={styles.schemaDeleteBtn}
                    onClick={(event) => {
                      event.stopPropagation();
                      handleDeleteField(record);
                    }}
                  />
                </span>
              </AuthWrapper>
            ) : null,
        },
      ]}
      dataSource={tableFields}
      rowKey="name"
      pagination={false}
      size="middle"
      hiddenEmpty
    />
  );
};

export default Definition;
