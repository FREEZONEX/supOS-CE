import './index.scss';
import type { CSSProperties, FC, ReactNode } from 'react';
import UserPopover from './UserPopover';
import { AuthWrapper } from '../auth';

interface OptionProps {
  onClick?: (item: OptionProps) => void;
  title?: ReactNode;
  label?: ReactNode;
  style?: CSSProperties;
  auth?: string | string[];
  noHoverStyle?: boolean;
  key: string;
  hidden?: boolean;
}
interface ComGroupButtonProps {
  options: (OptionProps | undefined)[];
}

const ComGroupButton: FC<ComGroupButtonProps> = ({ options }) => {
  const filterOptions = options?.filter((item) => item !== undefined);
  return (
    <div className="com-group-button">
      {filterOptions?.map((item: any) => {
        if (item.hidden) return null;
        switch (item.key) {
          case 'user':
            return (
              <UserPopover key={item.key}>
                <div style={item.style} className="item">
                  {item.label}
                </div>
              </UserPopover>
            );
          default:
            return (
              <AuthWrapper auth={item.auth} key={item.key}>
                <div
                  style={item.style}
                  onClick={() => item?.onClick?.(item)}
                  className={!item.noHoverStyle ? 'item' : 'no-hover-item'}
                  title={item.title}
                >
                  {item.label}
                </div>
              </AuthWrapper>
            );
        }
      })}
    </div>
  );
};

export default ComGroupButton;
