import { Flex, Tag } from 'antd';

export type RoleItem = {
  id: string;
  name: string;
  description: string;
  memberCount: number | string;
  canAccess: string[];
};

type FormatMessage = (...args: any[]) => string;
type RolesStyles = Record<string, string>;

export const createRolesColumns = ({
  formatMessage,
  styles,
}: {
  formatMessage: FormatMessage;
  styles: RolesStyles;
}) => {
  return [
    {
      title: () => formatMessage('project.role'),
      dataIndex: 'name',
      key: 'name',
      ellipsis: true,
      width: '20%',
      render: (value: string) => <span className={styles.roleName}>{value}</span>,
    },
    {
      title: () => formatMessage('project.description'),
      dataIndex: 'description',
      key: 'description',
      width: '30%',
      ellipsis: true,
    },
    {
      title: () => formatMessage('project.memberCount'),
      dataIndex: 'memberCount',
      key: 'memberCount',
      width: '15%',
    },
    {
      title: () => formatMessage('project.canAccess'),
      dataIndex: 'canAccess',
      key: 'canAccess',
      width: '30%',
      render: (apps: string[]) => (
        <Flex gap={8} wrap="wrap">
          {apps.map((app) => (
            <Tag key={app} className={styles.accessTag}>
              {app}
            </Tag>
          ))}
        </Flex>
      ),
    },
  ];
};
