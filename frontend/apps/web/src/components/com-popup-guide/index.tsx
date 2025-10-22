import { Flex, Progress, Tooltip, type TooltipProps } from 'antd';
import { cloneElement, type FC, isValidElement, type ReactNode, useEffect, useRef, useState } from 'react';
import { CloseOutline } from '@carbon/icons-react';
import classNames from 'classnames';
import './index.scss';
import { I18nEnum, useI18nStore } from '@/stores/i18n-store.ts';

type ComPopupGuideProps = TooltipProps & {
  timer?: number;
  onFinish?: (stepName: string | number, nextStepName: string | number | null | undefined, info: any) => void;
  onBegin?: (stepName: string | number, nextStepName: string | number | null | undefined, info: any) => void;
  children: ReactNode;
  currentStep?: string | number;
  stepName: string | number;
  steps?: any[];
};

const ComPopupGuide: FC<ComPopupGuideProps> = ({
  timer = 2000,
  onFinish,
  stepName,
  children,
  onBegin,
  currentStep,
  steps = [],
  ...restProps
}) => {
  const lang = useI18nStore((state) => state.lang);

  const [percent, setPercent] = useState(100);
  const [open, setOpen] = useState(stepName === currentStep);
  const info = steps?.find((f) => f.stepName === stepName) || {};

  const timerRef = useRef<number | undefined>(undefined);
  useEffect(() => {
    if (stepName === currentStep) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setOpen(true);
      onBegin?.(stepName, info.nextStep, info);
    } else {
      setPercent(100);
    }
    if (!(stepName === currentStep)) return;
    setOpen(true);
    const intervalTime = 100;
    const decrement = 100 / (timer / intervalTime);

    if (timerRef.current !== undefined) {
      clearInterval(timerRef.current);
      timerRef.current = undefined;
    }

    timerRef.current = window.setInterval(() => {
      setPercent((prev) => {
        if (prev <= 0) {
          setOpen(false);
          clearInterval(timerRef.current);
          timerRef.current = undefined;
          onFinish?.(stepName, info.nextStep, info);
          return 100;
        }
        return prev - decrement;
      });
    }, intervalTime);

    return () => {
      if (timerRef.current !== undefined) {
        clearInterval(timerRef.current);
        timerRef.current = undefined;
      }
    }; // 清除副作用
  }, [timer, stepName, currentStep]);

  const handleReset = () => {
    clearInterval(timerRef.current);
    setPercent(0);
    setOpen(false);
    onFinish?.(stepName, info.nextStep, info);
  };
  const title = lang !== I18nEnum.EnUS ? info?.titleEnglish : info?.title;
  const customTitle = typeof title === 'function' ? title?.() : title;
  const Children = isValidElement(children) ? (
    cloneElement(children, {
      ...(children.props || {}),
      className: classNames(children.props?.className, { 'com-popup-guide-wrapper': open }),
    })
  ) : open ? (
    <div className="com-popup-guide-wrapper">{children}</div>
  ) : (
    children
  );

  const Title = (
    <div className="title">
      {customTitle ?? 'no found'}
      <Flex style={{ width: '100%' }} align="center" gap={8}>
        <Progress success={{ strokeColor: '#0f62fe' }} showInfo={false} percent={percent} style={{ flex: 1 }} />
        <CloseOutline size={20} style={{ cursor: 'pointer' }} onClick={handleReset} />
      </Flex>
    </div>
  );
  return (
    <Tooltip
      styles={{
        root: { '--antd-arrow-background-color': '#0f62fe' },
      }}
      classNames={{
        root: 'com-popup-guide',
      }}
      trigger={[]}
      color="#fff"
      title={Title}
      open={open}
      {...restProps}
    >
      {Children}
    </Tooltip>
  );
};

export default ComPopupGuide;
