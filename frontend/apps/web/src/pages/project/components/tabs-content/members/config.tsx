import { formatTimestamp } from '@/utils/format';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import { Edit, TrashCan } from '@/components/lucide-icon/carbon';
import { Flex } from 'antd';
import type { HookAPI } from 'antd/es/modal/useModal';

export type MemberItem = {
  id: string;
  /** 后端返回的用户唯一标识，用于添加成员时过滤已存在成员 */
  userId: string;
  name: string;
  email: string;
  roles: string[];
  updatedAt: string;
};

type FormatMessage = (...args: any[]) => string;
type MembersStyles = Record<string, string>;

export const createMembersColumns = ({
  formatMessage,
  styles,
  onEdit,
  onDelete,
  modal,
}: {
  formatMessage: FormatMessage;
  styles: MembersStyles;
  onEdit: (record: MemberItem) => void;
  onDelete: (record: MemberItem) => void;
  modal: HookAPI;
}) => {
  return [
    {
      title: () => formatMessage('project.member'),
      dataIndex: 'name',
      key: 'name',
      width: '10%',
      render: (value: string) => (
        <div className={styles.memberCell}>
          <span className={styles.memberName}>{value}</span>
        </div>
      ),
    },
    {
      title: () => formatMessage('project.email'),
      dataIndex: 'email',
      key: 'email',
      width: '20%',
    },
    {
      title: () => formatMessage('project.role'),
      dataIndex: 'roleNames',
      key: 'roleNames',
      width: '35%',
      render: (roleNames: string[]) => (
        <Flex gap={6} wrap="wrap">
          {roleNames.map((role) => (
            <div key={role} className={styles.roleTag} title={role}>
              {role}
            </div>
          ))}
        </Flex>
      ),
    },
    {
      title: () => formatMessage('project.lastUpdated'),
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      width: '20%',
      render: (value: number) => formatTimestamp(value, 'YYYY/MM/DD HH:mm', true),
    },
    {
      title: () => formatMessage('project.operation'),
      key: 'operation',
      width: 120,
      fixed: 'right',
      render: (_: any, record: MemberItem) => (
        <span className={styles['operation']}>
          <i title={formatMessage('common.edit')}>
            <Edit
              className="custom-operation"
              style={{ cursor: 'pointer' }}
              onClick={() => {
                onEdit(record);
              }}
            />
          </i>
          <i title={formatMessage('common.delete')}>
            <TrashCan
              className="custom-operation"
              style={{ cursor: 'pointer' }}
              onClick={() => {
                modal.confirm({
                  ...createDeleteConfirmOptions({
                    title: formatMessage('project.confirmDeleteMemberTitle'),
                    content: formatMessage('project.confirmDeleteMemberDesc'),
                    okText: formatMessage('common.delete'),
                    cancelText: formatMessage('common.cancel'),
                  }),
                  onOk: () => onDelete(record),
                });
              }}
            />
          </i>
        </span>
      ),
    },
  ];
};
