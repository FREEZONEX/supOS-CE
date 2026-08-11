import { Empty, type EmptyProps } from 'antd';
import type { FC } from 'react';

export const EMPTY_IMAGE = Empty.PRESENTED_IMAGE_SIMPLE;

const ComEmpty: FC<EmptyProps> = ({ image = EMPTY_IMAGE, ...props }) => {
  return <Empty image={image} {...props} />;
};

export default ComEmpty;
