import type { FC } from 'react';
import CodeSnippet from '@/components/code-snippet';
import classNames from 'classnames';
import { useTranslate } from '@/hooks';

const RawData: FC<{ payload?: string | Record<string, unknown>; className?: string }> = ({ payload, className }) => {
  const formatMessage = useTranslate();
  let parsedPayload;
  try {
    if (typeof payload === 'string') {
      const firstParse = JSON.parse(payload);
      parsedPayload = typeof firstParse === 'string' ? JSON.parse(firstParse) : firstParse;
    } else {
      parsedPayload = payload;
    }
  } catch (error) {
    console.error('Failed to parse payload:', error);
    return null;
  }

  if (!parsedPayload) {
    return (
      <CodeSnippet
        className={classNames('codeViewWrap', className)}
        type="multi"
        minCollapsedNumberOfRows={3}
        maxCollapsedNumberOfRows={24}
      >
        <span style={{ fontSize: 15 }}>{formatMessage('uns.awaitingDataInput')}</span>
      </CodeSnippet>
    );
  }

  const formattedPayload = JSON.stringify(parsedPayload, null, 2);

  return (
    <CodeSnippet
      className={classNames('codeViewWrap', className)}
      type="multi"
      minCollapsedNumberOfRows={6}
      maxCollapsedNumberOfRows={24}
    >
      {formattedPayload}
    </CodeSnippet>
  );
};

export default RawData;
