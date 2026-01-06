import type { Location } from 'react-router';

export interface PageProps {
  location?: Partial<Location>;
  // 路由title
  title?: string;
}
