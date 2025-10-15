import type { CSSProperties, FC, ReactNode } from 'react';
import { CheckmarkFilled, ErrorFilled } from '@carbon/icons-react';
import './index.scss';

export interface InlineLoadingProps {
  /** 加载状态描述文本 */
  description?: ReactNode;
  /** 加载状态 */
  status?: 'active' | 'finished' | 'error';
  /** 自定义类名 */
  className?: string;
  style?: CSSProperties;
  title?: string;
}

const InlineLoading: FC<InlineLoadingProps> = ({ description, status = 'active', className = '', style, title }) => {
  return (
    <div className={`inline-loading ${className}`} style={style}>
      {status === 'active' && (
        <div className={`inline-loading-animation ${status}`}>
          <div className="inline-loading-spinner">
            <svg className="ods-loading_svg" viewBox="0 0 100 100">
              <circle className="ods-loading_background" cx="50" cy="50" r="42" />
              <circle className="ods-loading_stroke" cx="50" cy="50" r="42" />
            </svg>
          </div>
        </div>
      )}
      {status === 'finished' && <CheckmarkFilled fill={'#24a148'} />}
      {status === 'error' && <ErrorFilled fill={'#da1e28'} />}
      {description && (
        <div title={title} className="inline-loading-text">
          {description}
        </div>
      )}
    </div>
  );
};

export default InlineLoading;
