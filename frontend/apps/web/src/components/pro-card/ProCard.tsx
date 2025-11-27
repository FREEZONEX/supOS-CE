import { Checkbox, Divider, Flex, Spin, Image, Typography, Dropdown, Button } from 'antd';
import { type FC, useState } from 'react';
import type { ProCardProps } from '@/components/pro-card/type.ts';
import cx from 'classnames';
import { ChevronRight, OverflowMenuVertical, Pin, PinFilled } from '@carbon/icons-react';
import { AuthButton } from '@/components/auth';
import defaultUrl from '@/assets/home-icons/default.svg';
import InlineLoading from '@/components/inline-loading';
import { hasPermission } from '@/utils/auth.ts';
import ComButton from '../com-button';
import useTranslate from '@/hooks/useTranslate.ts';
import './index.scss';

const { Paragraph } = Typography;

const ProCard: FC<ProCardProps> = ({
  loading,
  styles,
  classNames,
  value,
  onChange,
  statusHeader,
  header,
  onClick,
  description,
  secondaryDescription,
  allowHover = true,
  actions: _actions,
  iconBg = true,
  item,
  actionConfig,
  border = false,
}) => {
  const formatMessage = useTranslate();
  const [checked, setChecked] = useState(false);
  const [clickLoading, setClickLoading] = useState(false);
  const cardClassName = cx(
    'pro-card',
    classNames?.card,
    checked && 'pro-card-checked',
    border && 'pro-card-border',
    allowHover && 'pro-card-hover'
  );
  const { allowCheck, statusInfo, statusTag, pinOptions } = statusHeader || {};
  const actionNum = actionConfig?.num ?? 2;
  const {
    customIcon,
    defaultIconUrl = defaultUrl,
    iconSrc,
    title,
    titleDescription,
    onClick: headerClick,
  } = header || {};
  const actions = _actions
    ? typeof _actions === 'function'
      ? _actions(item)
          ?.filter((f) => f?.key)
          .filter((item: any) => {
            return item && (!item.auth || hasPermission(item.auth));
          })
      : _actions
          ?.filter((f) => f?.key)
          .filter((item: any) => {
            return item && (!item.auth || hasPermission(item.auth));
          })
    : _actions;
  const handleClick = async (e: any, onClick: any) => {
    if (!onClick) return;

    try {
      setClickLoading(true);
      const result = onClick(e);

      if (result instanceof Promise) {
        await result;
      }
    } catch (error) {
      console.error('Button action failed:', error);
    } finally {
      setClickLoading(false);
    }
  };
  const isPin = pinOptions?.renderPinIcon?.(item) ?? false;
  const _description =
    description !== false ? (typeof description === 'string' ? description : description?.content) : false;
  return (
    <div style={styles?.root} className={classNames?.root}>
      <Spin spinning={loading || clickLoading || false}>
        <Flex
          vertical
          className={cardClassName}
          onClick={() => onClick?.(item)}
          style={{ cursor: onClick ? 'pointer' : 'inherit', ...styles?.card }}
        >
          {/* statusHeader */}
          {statusHeader ? (
            <Flex
              className={cx('pro-card-status-header', classNames?.statusHeader)}
              style={{ ...styles?.statusHeader }}
              justify="space-between"
              align="center"
              gap={4}
            >
              <Flex
                align="center"
                style={{ flex: 1, overflow: 'hidden', ...styles?.statusInfo }}
                className={cx(styles?.statusInfo)}
                justify="space-between"
              >
                <Flex align="center" justify="flex-start" style={{ flex: 1, overflow: 'hidden' }}>
                  {statusInfo && (
                    <Flex
                      style={{ overflow: 'hidden' }}
                      justify="flex-start"
                      align="center"
                      gap={8}
                      title={`${statusInfo.title}: ${statusInfo.label}`}
                    >
                      <div
                        style={{
                          width: 8,
                          height: 8,
                          borderRadius: '50%',
                          background: statusInfo?.color,
                          flexShrink: 0,
                        }}
                      />
                      <Paragraph
                        ellipsis={{
                          rows: 1,
                        }}
                        style={{ margin: 0, wordBreak: 'break-all', color: 'var(--supos-table-first-color)' }}
                      >
                        {statusInfo.label}
                      </Paragraph>
                    </Flex>
                  )}
                  {statusInfo && statusTag && (
                    <Divider
                      type="vertical"
                      style={{
                        backgroundColor: 'rgba(0, 0, 0, 0.06)',
                      }}
                    />
                  )}
                  {statusTag && (
                    <Flex align="center" style={styles?.statusTag} className={classNames?.statusTag}>
                      {statusTag}
                    </Flex>
                  )}
                </Flex>
                {pinOptions && (
                  <ComButton
                    title={isPin ? formatMessage('common.pin') : formatMessage('common.unPin')}
                    auth={pinOptions?.auth}
                    disabled={pinOptions?.disabled}
                    onClick={() => pinOptions?.onClick?.(item)}
                    icon={isPin ? <Pin style={{ color: '#525252' }} /> : <PinFilled />}
                    size="small"
                    type={'text'}
                  />
                )}
              </Flex>
              {allowCheck ? (
                <Checkbox
                  value={value}
                  className="card-title"
                  onChange={(e) => {
                    setChecked(e.target.checked);
                    onChange?.(e);
                  }}
                />
              ) : (
                <span></span>
              )}
            </Flex>
          ) : null}
          {/* Header */}
          <Flex
            className={cx('pro-card-header', classNames?.header)}
            style={{ ...styles?.header }}
            justify="space-between"
            gap={16}
          >
            {(iconSrc || customIcon) && (
              <Flex
                style={
                  iconBg
                    ? {
                        borderRadius: 3,
                        backgroundColor: 'var(--supos-image-card-color)',
                        padding: 6,
                      }
                    : undefined
                }
              >
                {customIcon ? (
                  customIcon
                ) : (
                  <Image preview={false} src={`${iconSrc ?? ''}`} width={28} height={28} fallback={defaultIconUrl} />
                )}
              </Flex>
            )}
            <Flex
              style={{ cursor: headerClick ? 'pointer' : 'inherit', flex: 1, overflow: 'hidden' }}
              onClick={() => headerClick?.(item)}
              align="center"
            >
              <Flex
                style={{ flex: 1, height: '100%', overflow: 'hidden' }}
                vertical
                justify={titleDescription ? 'space-between' : 'center'}
              >
                <Flex align="center" justify="space-between" title={typeof title === 'string' ? title : ''}>
                  <div className={cx('header-title', classNames?.headerTitle)} style={styles?.headerTitle}>
                    {title}
                  </div>
                </Flex>
                {titleDescription && (
                  <div
                    title={typeof titleDescription === 'string' ? titleDescription : ''}
                    className="title-description"
                  >
                    {titleDescription}
                  </div>
                )}
              </Flex>
              {(onClick || headerClick) && <ChevronRight style={{ flexShrink: 0 }} size={26} />}
            </Flex>
          </Flex>
          {/* description */}
          {description !== false && (
            <div
              className="pro-card-description"
              title={typeof description === 'string' ? description : description?.content}
              style={{
                height: typeof description === 'string' ? 40 : description?.rows ? (40 / 2) * description?.rows : 40,
              }}
            >
              {!_description ? (
                <span style={{ color: '#8D8D8D' }}>
                  {formatMessage(
                    typeof description !== 'string' && description?.empty ? description?.empty : 'common.noDescription'
                  )}
                </span>
              ) : (
                <Paragraph
                  ellipsis={{
                    rows: typeof description === 'string' ? 2 : description?.rows || 2,
                  }}
                  style={{ margin: 0, wordBreak: 'break-all', color: 'var(--supos-table-first-color)', fontSize: 12 }}
                >
                  {typeof description === 'string' ? description : description?.content}
                </Paragraph>
              )}
            </div>
          )}
          {/* 二级描述 */}
          {secondaryDescription && (
            <div
              style={styles?.secondaryDescription}
              className={cx('pro-card-secondary-description', classNames?.secondaryDescription)}
              title={typeof secondaryDescription === 'string' ? secondaryDescription : undefined}
            >
              {secondaryDescription}
            </div>
          )}
          {actions && (
            <Divider
              style={{
                margin: '12px 0',
                backgroundColor: 'var(--supos-t-dividr-color)',
              }}
            />
          )}
          {actions && (
            <Flex align="center" justify="space-between">
              <Flex gap={8} align="center" style={{ flex: 1, overflow: 'hidden' }}>
                {Array.isArray(actions) &&
                  actions.slice(0, actionNum).map(({ label, key, title, icon, button, onClick, status, disabled }) =>
                    key === 'loading' ? (
                      <InlineLoading key={key} status={status || 'active'} description={label} />
                    ) : (
                      <AuthButton
                        {...button}
                        style={{
                          ...button?.style,
                          maxWidth: `calc((100% - 16px) / ${actionNum})`,
                        }}
                        icon={icon}
                        size="small"
                        key={key}
                        onClick={(e) => handleClick(e, onClick)}
                        title={title ? title : typeof label === 'string' ? label : ''}
                        disabled={disabled}
                      >
                        {label}
                      </AuthButton>
                    )
                  )}
              </Flex>
              {Array.isArray(actions) && actions.length > actionNum && (
                <Dropdown
                  menu={{
                    items: actions.slice(actionNum).map(({ key, label, icon, title, onClick, disabled, extra }) => ({
                      key,
                      label,
                      icon,
                      title: title ? title : typeof label === 'string' ? label : '',
                      onClick: (e) => handleClick(e, onClick),
                      disabled,
                      extra,
                    })),
                  }}
                >
                  <Button type="text" icon={<OverflowMenuVertical />} size="small" />
                </Dropdown>
              )}
            </Flex>
          )}
        </Flex>
      </Spin>
    </div>
  );
};

export default ProCard;
