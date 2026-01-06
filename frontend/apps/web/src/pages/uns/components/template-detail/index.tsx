import { type CSSProperties, type FC, useEffect, useState } from 'react';
import { CaretRight, Copy, WatsonHealth3DMprToggle } from '@carbon/icons-react';
import { App, Button, Collapse, Flex, theme, Typography } from 'antd';
import { useTranslate } from '@/hooks';
import Icon from '@ant-design/icons';
import EditButton from '@/pages/uns/components/EditButton.tsx';
import FileList from './FileList';
import { editTemplateName, getTemplateDetail, updateTemplateSubscribe } from '@/apis/inter-api/uns.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';
const { Paragraph, Title } = Typography;

import type { InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import { AuthWrapper } from '@/components/auth';
import ComDetailList from '@/components/com-detail-list';
import ProTable from '@/components/pro-table';
import FileEdit from '@/components/svg-components/FileEdit';
import { hasPermission } from '@/utils/auth';
import { formatTimestamp } from '@/utils/format';
import Subscribe from '@/pages/uns/components/subscribe';
import CustomParagraph from '@/components/custom-paragraph';
import styles from './index.module.scss';
import useClipboard from '@/hooks/useClipboard.ts';

interface TemplateDetailProps {
  // id
  currentNode: UnsTreeNode;
  handleDelete?: (node: UnsTreeNode) => void;
  initTreeData?: InitTreeDataFnType;
}

const panelStyle: CSSProperties = {
  background: 'val(--supos-bg-color)',
  border: 'none',
};

const TemplateDetail: FC<TemplateDetailProps> = ({ currentNode: { id }, handleDelete, initTreeData }) => {
  const [activeList, setActiveList] = useState<string[]>(['detail', 'definition', 'fileList']);
  const { token } = theme.useToken();
  const { message } = App.useApp();
  const [info, setInfo] = useState<{ [key: string]: any }>({});
  const formatMessage = useTranslate();

  const onDeleteHandle = () => {
    if (id) {
      handleDelete?.({ id: id as string, key: '', pathType: 1 });
    }
  };

  const getModel = (id: string) => {
    if (!id) return;
    getTemplateDetail({ id }).then((data) => {
      setInfo(data);
    });
  };

  useEffect(() => {
    if (id) {
      getModel(id as string);
    }
  }, [id]);
  const { copy } = useClipboard();

  const items = [
    {
      key: 'detail',
      label: <span>{formatMessage('common.detail')}</span>,
      children: (
        <ComDetailList
          list={[
            {
              label: 'Topic',
              key: 'topic',
              hide: !info?.subscribeEnable,
              render: () => {
                return (
                  <div>
                    {info.topic}
                    <span
                      style={{ marginLeft: '5px', verticalAlign: 'sub', cursor: 'pointer' }}
                      onClick={() => copy(info.topic)}
                      title={formatMessage('common.copy')}
                    >
                      <Copy />
                    </span>
                  </div>
                );
              },
            },
            {
              label: formatMessage('uns.alias'),
              key: 'alias',
            },
            {
              label: formatMessage('uns.subscriptionFrequency'),
              key: 'subscribeFrequency',
              hide: !info?.subscribeEnable,
              render: () => {
                return (
                  <CustomParagraph
                    className={styles.paragraph}
                    value={info.subscribeFrequency}
                    onChange={(value) => {
                      updateTemplateSubscribe({ id, ...value }).then(() => {
                        message.success(formatMessage('uns.editSuccessful'));
                        getModel(id as string);
                      });
                    }}
                  />
                );
              },
            },
            {
              label: formatMessage('uns.description'),
              key: 'description',
              render: (item) => (
                <Paragraph
                  style={{ margin: 0, width: '100%' }}
                  editable={
                    hasPermission(ButtonPermission['uns.templateDetail'])
                      ? {
                          icon: (
                            <Icon
                              data-button-auth={ButtonPermission['uns.templateDetail']}
                              component={FileEdit}
                              style={{
                                fontSize: 17,
                                color: 'var(--supos-text-color)',
                              }}
                            />
                          ),
                          onChange: (val) => {
                            if (item === val || (!item && !val)) return;
                            if (val.length > 255) {
                              message.warning(
                                formatMessage('uns.labelMaxLength', {
                                  label: formatMessage('uns.description'),
                                  length: 255,
                                })
                              );
                              return;
                            }
                            editTemplateName({ id, description: val }).then(() => {
                              message.success(formatMessage('uns.editSuccessful'));
                              getModel(id as string);
                            });
                          },
                          // maxLength: 255,
                        }
                      : false
                  }
                >
                  {item}
                </Paragraph>
              ),
            },
            {
              label: formatMessage('common.creationTime'),
              key: 'createTime',
              render: (item) => formatTimestamp(item),
            },
          ]}
          data={info}
        />
      ),
      style: panelStyle,
    },
    {
      key: 'definition',
      label: <span>{formatMessage('uns.definition')}</span>,
      extra: (
        <EditButton
          auth={ButtonPermission['uns.templateDetail']}
          modelInfo={info}
          getModel={() => getModel(id as string)}
          editType="template"
        />
      ),
      children: (
        <ProTable
          rowHoverable={false}
          columns={[
            {
              title: formatMessage('common.name'),
              dataIndex: 'name',
              width: '20%',
            },
            {
              title: formatMessage('uns.type'),
              dataIndex: 'type',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
            },
            {
              title: formatMessage('common.length'),
              dataIndex: 'maxLen',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
            },
            {
              title: formatMessage('uns.displayName'),
              dataIndex: 'displayName',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
            },
            {
              title: formatMessage('uns.remark'),
              dataIndex: 'remark',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
            },
          ]}
          dataSource={info?.fields || []}
          rowKey="name"
          pagination={false}
          size="middle"
          hiddenEmpty
          bordered
        />
      ),
      style: panelStyle,
    },
    {
      key: 'fileList',
      label: <span>{formatMessage('common.fileList')}</span>,
      children: <FileList templateId={id as string} />,
      style: panelStyle,
    },
  ];

  const handleChangeSubscribe = async (enable: boolean, frequency?: string) => {
    await updateTemplateSubscribe({ id, enable, frequency });
    getModel(id as string);
    message.success(enable ? formatMessage('uns.subscribeSuccessful') : formatMessage('uns.unsubscribeSuccessful'));
  };

  return (
    <div className="topicDetailWrap">
      <div className="topicDetailContent">
        <Flex className="detailTitle" gap={8} justify="flex-start" align="center">
          <WatsonHealth3DMprToggle size={30} />
          <Title
            level={2}
            style={{ margin: 0, width: '100%', insetInlineStart: 0 }}
            editable={
              hasPermission(ButtonPermission['uns.templateDetail'])
                ? {
                    icon: (
                      <Icon
                        data-button-auth={ButtonPermission['uns.templateDetail']}
                        component={FileEdit}
                        style={{
                          fontSize: 25,
                          color: '#5A5A5A',
                          marginLeft: 5,
                        }}
                      />
                    ),
                    onChange: (val) => {
                      if (info?.name === val || !val || val.trim() === '') return;
                      if (val.length > 63) {
                        return message.warning(
                          formatMessage('uns.labelMaxLength', { label: formatMessage('common.name'), length: 63 })
                        );
                      }
                      const reg = /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/;
                      if (!reg.test(val)) {
                        return message.error(formatMessage('uns.nameFormat'));
                      }
                      editTemplateName({ id, name: val }).then(() => {
                        message.success(formatMessage('uns.editSuccessful'));
                        getModel(id as string);
                        initTreeData?.({ queryType: 'editTemplateName', newNodeKey: id });
                      });
                    },
                  }
                : false
            }
          >
            {info?.name}
          </Title>
          <Subscribe
            hidden
            showModal
            value={info.subscribeEnable}
            subscribeFrequency={info.subscribeFrequency}
            topic={info?.topic}
            onChange={handleChangeSubscribe}
          />
        </Flex>
        <div className="tableWrap">
          <Collapse
            bordered={false}
            collapsible="header"
            activeKey={activeList}
            onChange={(even) => setActiveList(even)}
            expandIcon={({ isActive }) => (
              <CaretRight
                size={20}
                style={{
                  rotate: isActive ? '90deg' : '0deg',
                  transition: '200ms',
                }}
              />
            )}
            items={items}
            style={{ background: token.colorBgContainer }}
          />
        </div>
        <AuthWrapper auth={ButtonPermission['uns.templateDetail']}>
          <div className="deleteBtnWrap" style={{ marginTop: 0 }}>
            <Button
              type="primary"
              style={{
                width: '100px',
                fontWeight: 'bold',
              }}
              onClick={onDeleteHandle}
            >
              {formatMessage('common.delete')}
            </Button>
          </div>
        </AuthWrapper>
      </div>
    </div>
  );
};

export default TemplateDetail;
