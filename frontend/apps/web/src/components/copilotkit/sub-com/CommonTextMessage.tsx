import type { FC, ReactNode } from 'react';
import classNames from 'classnames';
import { ChatBot } from '@carbon/icons-react';
import './TextMessage.scss';

const CommonTextMessage: FC<{ children: ReactNode }> = ({ children }) => {
  return (
    <div className={classNames('text-message', 'text-message-assistant')}>
      {children}
      <div className="icon">
        <ChatBot size={16} color="var(--supos-theme-color)" />
      </div>
    </div>
  );
};

export default CommonTextMessage;
