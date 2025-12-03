import { Flex } from 'antd';
import classNames from 'classnames';
import styles from './Html.module.scss';

interface PropsTypes {
  className?: string;
  value: string;
  label: string;
}

const HtmlField = (props: PropsTypes) => {
  const { value, label, className } = props;
  return (
    <Flex align="baseline" className={classNames('com-copy-content', styles.container, className)}>
      <div className="label">{label}</div>
      <Flex className={'content'} justify="space-between">
        <div className={'text'} dangerouslySetInnerHTML={{ __html: value }} />
      </Flex>
    </Flex>
  );
};

export default HtmlField;
