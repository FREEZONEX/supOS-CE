import React, { type FC, useState } from 'react';
import { AuthButton } from '../auth';
import type { ButtonProps } from 'antd';

type ComButtonProps = Omit<ButtonProps, 'onClick'> & {
  onClick?: React.MouseEventHandler<HTMLElement> | (() => Promise<any>);
};

const ComButton: FC<ComButtonProps> = (props) => {
  const [loading, setLoading] = useState(false);
  const handleClick: React.MouseEventHandler<HTMLButtonElement> = async (e) => {
    if (!props.onClick) return;

    try {
      setLoading(true);
      const result = props.onClick(e);

      if (result instanceof Promise) {
        await result;
      }
    } catch (error) {
      console.error('Button action failed:', error);
    } finally {
      setLoading(false);
    }
  };
  return <AuthButton {...props} loading={props.loading || loading} onClick={handleClick}></AuthButton>;
};

export default ComButton;
