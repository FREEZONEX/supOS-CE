import type { Project } from '@/pages/launchpad/type';
import { useTranslate } from '@/hooks';
import { Folder, Box, Star } from '@/components/lucide-icon/carbon';
import { Card } from 'antd';
import classNames from 'classnames';
import { type FC } from 'react';
import HighlightMatch from './HighlightMatch';
import styles from './ProjectGrid.module.scss';

interface ProjectGridProps {
  projects: Project[];
  onProjectClick: (project: Project) => void;
  isProjectFavorite: (id: string | number) => boolean;
  onToggleFavorite: (project: Project) => void;
  /** 搜索关键词，用于项目名命中高亮 */
  keyword?: string;
}

const ProjectGrid: FC<ProjectGridProps> = ({
  projects,
  onProjectClick,
  isProjectFavorite,
  onToggleFavorite,
  keyword,
}) => {
  const formatMessage = useTranslate('Launchpad');

  return (
    <div className={styles.projectGrid}>
      {projects.map((project) => {
        const displayName = project.displayName || project.projectName;
        const favorite = isProjectFavorite(project.projectId);

        return (
          <Card key={project.projectId} className={styles.projectCard} onClick={() => onProjectClick(project)}>
            <div className={styles.cardHeader}>
              <div className={styles.projectInfo}>
                <div className={styles.iconWrapper}>
                  <Folder size={26} className={styles.folderIcon} />
                </div>
                <div className={styles.textInfo}>
                  <div className={styles.projectName} title={displayName}>
                    <HighlightMatch text={displayName} keyword={keyword} />
                  </div>
                  <div className={styles.appsCount}>
                    {formatMessage('appsCount', { count: project.apps?.length || 0 })}
                  </div>
                </div>
              </div>
              <button
                type="button"
                className={classNames(styles.starButton, favorite && styles.starActive)}
                aria-label={formatMessage(
                  favorite ? 'rmFav' : 'addFav',
                  undefined,
                  favorite ? 'Remove favorite' : 'Add favorite'
                )}
                onClick={(event) => {
                  event.stopPropagation();
                  onToggleFavorite(project);
                }}
              >
                <Star size={12} />
              </button>
            </div>
            <div className={styles.appIcons}>
              {Array.from({ length: 4 }, (_, index) => {
                const app = project.apps?.[index];
                const isEmpty = !app?.appId;
                return (
                  <div
                    key={app?.appId ?? `placeholder-${project.projectId}-${index}`}
                    className={classNames(styles.appIcon, isEmpty && styles.appIconEmpty)}
                    title={isEmpty ? undefined : app.displayName || app.appName}
                  >
                    {isEmpty ? (
                      <Box size={17} strokeWidth={1.5} className={styles.emptyIcon} />
                    ) : (
                      <>
                        <div className={styles.iconContainer}>
                          {app.iconUrl ? (
                            <img
                              src={app.iconUrl}
                              alt={app.displayName || app.appName}
                              className={styles.appIconImage}
                            />
                          ) : (
                            <Box size={17} strokeWidth={1.75} />
                          )}
                        </div>
                        {app.coverUrl ? (
                          <img
                            src={app.coverUrl}
                            alt={app.displayName || app.appName}
                            className={styles.coverImage}
                            onError={(event) => {
                              event.currentTarget.hidden = true;
                            }}
                          />
                        ) : null}
                        <span className={styles.appNameOverlay}>
                          <span className={styles.appNameText}>{app.displayName || app.appName}</span>
                        </span>
                      </>
                    )}
                  </div>
                );
              })}
            </div>
          </Card>
        );
      })}
    </div>
  );
};

export default ProjectGrid;
