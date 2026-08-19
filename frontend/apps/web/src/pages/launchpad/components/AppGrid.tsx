import { type FC } from 'react';
import type { App } from '@/pages/launchpad/type';
import AppCard from './AppCard';
import styles from './AppGrid.module.scss';

interface AppGridProps {
  apps: App[];
  onAppClick: (app: App) => void | Promise<void>;
  isAppFavorite: (id: string | number) => boolean;
  onToggleFavorite: (app: App) => void;
  /** 搜索关键词，用于应用名命中高亮 */
  keyword?: string;
}

const AppGrid: FC<AppGridProps> = ({ apps, onAppClick, isAppFavorite, onToggleFavorite, keyword }) => {
  return (
    <div className={styles.appGrid}>
      {apps.map((app) => (
        <AppCard
          key={app.appId}
          className={styles.appCard}
          app={app}
          isFavorite={isAppFavorite(app.appId)}
          onToggleFavorite={onToggleFavorite}
          onAppClick={onAppClick}
          showUpdatedAt
          keyword={keyword}
        />
      ))}
    </div>
  );
};

export default AppGrid;
