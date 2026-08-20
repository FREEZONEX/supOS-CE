import { ProTable } from '@/components';
import type { ProTableColumns } from '@/components/pro-table/types';
import { Folder, Star } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import type { Project } from '@/pages/launchpad/type';
import { formatTimestamp } from '@/utils/format';
import classNames from 'classnames';
import { type FC } from 'react';
import HighlightMatch from './HighlightMatch';
import styles from './LaunchpadTable.module.scss';

interface ProjectTableProps {
  projects: Project[];
  onProjectClick: (project: Project) => void;
  isProjectFavorite: (id: string | number) => boolean;
  onToggleFavorite: (project: Project) => void;
  /** 搜索关键词，用于项目名命中高亮 */
  keyword?: string;
}

const ProjectTable: FC<ProjectTableProps> = ({
  projects,
  onProjectClick,
  isProjectFavorite,
  onToggleFavorite,
  keyword,
}) => {
  const formatMessage = useTranslate();

  const columns: ProTableColumns<Project> = [
    {
      title: '',
      dataIndex: 'isFavorite',
      width: 28,
      maxWidth: 28,
      align: 'center',
      render: (_: unknown, record: Project) => {
        const favorite = isProjectFavorite(record.projectId);
        return (
          <button
            type="button"
            className={classNames(styles.starButton, favorite && styles.starActive)}
            aria-label={favorite ? 'Remove favorite' : 'Add favorite'}
            onClick={(event) => {
              event.stopPropagation();
              onToggleFavorite(record);
            }}
          >
            <Star size={16} />
          </button>
        );
      },
    },
    {
      title: formatMessage('common.name'),
      dataIndex: 'projectName',
      render: (_: unknown, record: Project) => {
        const name = record.displayName || record.projectName;
        return (
          <div className={styles.nameCell} title={name}>
            <span className={`${styles.nameIcon} ${styles.nameIconFolder}`}>
              <Folder size={16} className={styles.folderIcon} />
            </span>
            <span className={styles.nameText}>
              <HighlightMatch text={name} keyword={keyword} />
            </span>
          </div>
        );
      },
    },
    {
      title: formatMessage('common.apps'),
      dataIndex: 'apps',
      width: 160,
      render: (_: unknown, record: Project) =>
        formatMessage('Launchpad.appsCount', { count: record.apps?.length || 0 }),
    },
    {
      title: formatMessage('common.latestUpdate'),
      dataIndex: 'updatedAt',
      width: 220,
      render: (updatedAt?: string) => (updatedAt ? formatTimestamp(updatedAt, 'YYYY/MM/DD HH:mm', true) : '-'),
    },
  ];

  return (
    <ProTable
      columns={columns}
      dataSource={projects}
      rowKey="projectId"
      pagination={false}
      onRow={(record: Project) => ({
        onClick: () => onProjectClick(record),
        style: { cursor: 'pointer' },
      })}
    />
  );
};

export default ProjectTable;
