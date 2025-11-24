import { useState } from 'react';
import styles from './McpToolCall.module.scss';

interface ToolCallProps {
  status: 'complete' | 'inProgress' | 'executing';
  name?: string;
  args?: any;
  result?: any;
}

export default function MCPToolCall({ status, name = '', args, result }: ToolCallProps) {
  const [isOpen, setIsOpen] = useState(false);

  // Format content for display
  const format = (content: any): string => {
    if (!content) return '';
    const text = typeof content === 'object' ? JSON.stringify(content, null, 2) : String(content);
    return text.replace(/\\n/g, '\n').replace(/\\t/g, '\t').replace(/\\"/g, '"').replace(/\\\\/g, '\\');
  };

  const getStatusClass = () => {
    if (status === 'complete') return styles.complete;
    if (status === 'inProgress' || status === 'executing') return styles.inProgress;
    return styles.default;
  };

  return (
    <div className={styles.mcpToolCall}>
      <div className={styles.header} onClick={() => setIsOpen(!isOpen)}>
        <span className={styles.title}>MCP: {name || 'MCP Tool Call'}</span>
        <div className={styles.statusIndicator}>
          <div className={`${styles.indicator} ${getStatusClass()}`} />
        </div>
      </div>

      {isOpen && (
        <div className={styles.content}>
          {args && (
            <div className={styles.section}>
              <div className={styles.label}>Parameters:</div>
              <pre className={styles.pre}>{format(args)}</pre>
            </div>
          )}

          {status === 'complete' && result && (
            <div className={styles.section}>
              <div className={styles.label}>Result:</div>
              <pre className={styles.pre}>{format(result)}</pre>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
