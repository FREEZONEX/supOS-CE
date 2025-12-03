import type { FC } from 'react';
import type { RenderMessageProps } from '@copilotkit/react-ui';
import classNames from 'classnames';
import { ChatBot } from '@carbon/icons-react';
import ReactMarkdown from 'react-markdown';
import './TextMessage.scss';
import InlineLoading from '@/components/inline-loading';

const TextMessage: FC<RenderMessageProps & { message: any }> = (props) => {
  const { message, inProgress, isCurrentMessage } = props;
  return (
    <div
      className={classNames('text-message', {
        'text-message-assistant': message.role === 'assistant',
        'text-message-user': message.role === 'user',
      })}
    >
      {inProgress && isCurrentMessage ? (
        <InlineLoading status="active" />
      ) : message.role === 'assistant' ? (
        <ReactMarkdown>{message.content}</ReactMarkdown>
      ) : (
        message.content
      )}
      {message.role === 'assistant' && (
        <div className="icon">
          <ChatBot size={16} color="var(--supos-theme-color)" />
        </div>
      )}
    </div>
  );
};

export default TextMessage;
