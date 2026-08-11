import { type CSSProperties, type FC, type ReactNode } from 'react';
import styles from './TabLayout.module.scss';

interface TabLayoutProps {
  toolbarActions?: ReactNode;
  children: ReactNode;
  style?: CSSProperties;
}

/**
 * Tab 内容布局组件
 * 用于统一处理 Tab 内容和工具栏按钮的布局
 *
 * @param toolbarActions - 工具栏按钮（可选），会显示在 tabs 右侧
 * @param children - Tab 内容
 */
const TabLayout: FC<TabLayoutProps> = ({ toolbarActions, children, style }) => {
  return (
    <>
      {toolbarActions && <div className={styles.toolbarActions}>{toolbarActions}</div>}
      <div className={styles.tabContent} style={style}>
        {children}
      </div>
    </>
  );
};

export default TabLayout;
