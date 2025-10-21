import { Col, Row } from 'antd';
import type { CSSProperties } from 'react';
import './secondaryList.scss';

const SecondaryList = ({
  options,
  colon = true,
}: {
  options: {
    label: string;
    content: string;
    span: number;
    labelStyle?: CSSProperties;
    contentStyle?: CSSProperties;
    wrapperStyle?: CSSProperties;
    key: string | number;
  }[];
  colon?: boolean;
}) => {
  return (
    <Row style={{ overflow: 'hidden', marginTop: 4 }} className="secondaryList">
      {options?.map(({ label, labelStyle, contentStyle, span, content, wrapperStyle, key }) => {
        return (
          <Col span={span} style={{ display: 'flex', ...wrapperStyle }} key={key}>
            <div
              style={labelStyle ? labelStyle : { maxWidth: 'calc(50% - 9px)', color: '#A8A8A8' }}
              className="span-ellipsis"
              title={label}
            >
              {label}
            </div>
            {colon && <span style={{ paddingRight: 4, color: '#A8A8A8' }}>:</span>}
            <div
              style={contentStyle ? contentStyle : { flex: 1, minWidth: 0 }}
              className="span-ellipsis"
              title={content}
            >
              {content}
            </div>
          </Col>
        );
      })}
    </Row>
  );
};

export default SecondaryList;
