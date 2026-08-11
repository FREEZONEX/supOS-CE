import { useState, useEffect, type FC, type CSSProperties, useRef } from 'react';
import { getModelInfo, updateModelSubscribe } from '@/apis/core-api/uns';
import { useClipboard, useTranslate } from '@/hooks';
import { Collapse, App, theme, Typography, Flex, Tag } from 'antd';
import { CaretRight, ChartLine, ClipboardList, Copy, Folder, Route, SendAlt } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { UNS_TOPIC_ICON_COLORS, UnsTopicTypeIconWrap } from '@/pages/uns/components/uns-tree/tree-icons';
import DocumentList from '@/pages/uns/components/DocumentList.tsx';
import UploadButton from '@/pages/uns/components/UploadButton.tsx';
import EditButton from '@/pages/uns/components/EditButton.tsx';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import EditDetailButton from '@/pages/uns/components/EditDetailButton.tsx';
import type { InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import { formatTimestamp } from '@/utils/format';
import ProTable from '@/components/pro-table';
import { useBaseStore } from '@/stores/base';
import Subscribe from '@/pages/uns/components/subscribe';
import CustomParagraph from '@/components/custom-paragraph';

const { Title } = Typography;

const panelStyle: CSSProperties = {
  background: 'val(--ui-bg-color)',
  border: 'none',
};

export interface FolderDetailProps {
  currentNode: UnsTreeNode;
  initTreeData: InitTreeDataFnType;
  readOnly?: boolean;
}

const Module: FC<FolderDetailProps> = (props) => {
  const { message } = App.useApp();
  const {
    currentNode: { id, countChildren, mount },
    readOnly = false,
  } = props;
  const folderReadOnly = readOnly || Boolean(mount);
  const documentListRef = useRef(null);
  const formatMessage = useTranslate();
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const [activeList, setActiveList] = useState<string[]>(['detail', 'definition', 'document']);
  const { token } = theme.useToken();

  const [modelInfo, setModelInfo] = useState<{ [key: string]: any }>({});

  const getModel = (id: string) => {
    getModelInfo({ id })
      .then((data: any) => {
        setModelInfo(data || {});
      })
      .catch(() => {});
  };

  useEffect(() => {
    if (id) {
      getModel(id as string);
    }
  }, [id]);
  const { copy } = useClipboard();

  const folderTypeMap: { [key: number]: string } = {
    1: formatMessage('uns.state'),
    2: formatMessage('uns.action'),
    3: formatMessage('uns.metric'),
  };

  const items = [
    {
      key: 'detail',
      label: <span>{formatMessage('common.detail')}</span>,
      children: (
        <>
          {modelInfo?.subscribeEnable && (
            <div className="detailItem">
              <div className="detailKey">Topic</div>
              <div>
                {modelInfo.topic}
                <span
                  style={{ marginLeft: '5px', verticalAlign: 'sub', cursor: 'pointer' }}
                  onClick={() => copy(modelInfo.topic)}
                  title={formatMessage('common.copy')}
                >
                  <Copy {...toolbarIconProps} />
                </span>
              </div>
            </div>
          )}
          <div className="detailItem">
            <div className="detailKey"> {formatMessage('uns.alias')}</div>
            <div>{modelInfo.alias}</div>
          </div>
          <div className="detailItem">
            <div className="detailKey">{formatMessage('uns.displayName')}</div>
            <div>{modelInfo.displayName}</div>
          </div>
          <div className="detailItem">
            <div className="detailKey"> {formatMessage('uns.description')}</div>
            <div>{modelInfo.description}</div>
          </div>
          {modelInfo?.subscribeEnable && (
            <div className="detailItem">
              <div className="detailKey"> {formatMessage('uns.subscriptionFrequency')}</div>
              {folderReadOnly ? (
                <div>{modelInfo.subscribeFrequency || '-'}</div>
              ) : (
                <CustomParagraph
                  value={modelInfo.subscribeFrequency}
                  onChange={(value) => {
                    updateModelSubscribe({ id, ...value }).then(() => {
                      message.success(formatMessage('uns.editSuccessful'));
                      getModel(id as string);
                    });
                  }}
                />
              )}
            </div>
          )}
          <div className="detailItem">
            <div className="detailKey">{formatMessage('common.creationTime')}</div>
            <div>{formatTimestamp(modelInfo.createTime)}</div>
          </div>
          <div className="detailItem">
            <div className="detailKey">{formatMessage('uns.namespace')}</div>
            <div>{modelInfo.path}</div>
          </div>
          <div className="detailItem">
            <div className="detailKey">{formatMessage('uns.instanceCountStatistics')}</div>
            <div>{countChildren}</div>
          </div>
          <div className="detailItem">
            <div className="detailKey">{formatMessage('uns.originalName')}</div>
            <div>{modelInfo.name}</div>
          </div>
          <div className="detailItem">
            <div className="detailKey">{formatMessage('common.latestUpdate')}</div>
            {(modelInfo.updateTime || modelInfo.createTime) && (
              <div>{formatTimestamp(modelInfo.updateTime || modelInfo.createTime)}</div>
            )}
          </div>
          {modelInfo.extend &&
            Object.keys(modelInfo.extend).map((item: string, index: number) => (
              <div className="detailItem" key={index}>
                <div className="detailKey">{item}</div>
                <div>
                  {modelInfo.extend[item]}
                  <Tag style={{ marginLeft: '8px' }}>{formatMessage('uns.expandedInformation')}</Tag>
                </div>
              </div>
            ))}
          {modelInfo.dataType > 0 && (
            <div className="detailItem">
              <div className="detailKey">{formatMessage('uns.folderType')}</div>
              <div>{folderTypeMap[modelInfo.dataType]}</div>
            </div>
          )}
        </>
      ),
      style: panelStyle,
      extra: folderReadOnly ? null : (
        <EditDetailButton
          auth={ButtonPermission['uns.folderDetail']}
          type="folder"
          modelInfo={modelInfo}
          getModel={() => getModel(id as string)}
        />
      ),
    },
    {
      key: 'definition',
      label: formatMessage('uns.definition'),
      extra: folderReadOnly ? null : (
        <EditButton
          auth={ButtonPermission['uns.folderDetail']}
          modelInfo={modelInfo}
          getModel={() => getModel(id as string)}
          editType="folder"
        />
      ),
      children: (
        <ProTable
          rowKey={'name'}
          bordered={true}
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
              render: (text) => <span style={{ color: 'var(--ui-theme-color)' }}>{text}</span>,
            },
            {
              title: formatMessage('common.length'),
              dataIndex: 'maxLen',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--ui-theme-color)' }}>{text}</span>,
            },
            {
              title: formatMessage('uns.displayName'),
              dataIndex: 'displayName',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--ui-theme-color)' }}>{text}</span>,
            },
            {
              title: formatMessage('uns.remark'),
              dataIndex: 'remark',
              width: '20%',
              render: (text) => <span style={{ color: 'var(--ui-theme-color)' }}>{text}</span>,
            },
          ]}
          dataSource={modelInfo?.fields || []}
          hiddenEmpty
          pagination={false}
          size="middle"
        />
      ),
      style: panelStyle,
    },
    {
      key: 'document',
      label: formatMessage('common.document'),
      children: (
        <DocumentList alias={modelInfo.alias} ownerId={modelInfo.id} readOnly={folderReadOnly} ref={documentListRef} />
      ),
      style: panelStyle,
      extra: folderReadOnly ? null : (
        <UploadButton
          auth={ButtonPermission['uns.folderDetail']}
          alias={modelInfo.alias}
          ownerId={modelInfo.id}
          documentListRef={documentListRef}
          setActiveList={setActiveList}
        />
      ),
    },
  ];

  const handleChangeSubscribe = async (enable: boolean, frequency?: string) => {
    await updateModelSubscribe({ id, enable, frequency });
    getModel(id as string);
    message.success(enable ? formatMessage('uns.subscribeSuccessful') : formatMessage('uns.unsubscribeSuccessful'));
  };

  const getFolderIcon = () => {
    switch (modelInfo.dataType) {
      case 0:
        return <Route size={20} strokeWidth={1.75} aria-hidden />;
      case 1:
        return (
          <ClipboardList size={20} strokeWidth={1.75} aria-hidden style={{ color: UNS_TOPIC_ICON_COLORS.state }} />
        );
      case 2:
        return <SendAlt size={20} strokeWidth={1.75} aria-hidden style={{ color: UNS_TOPIC_ICON_COLORS.action }} />;
      case 3:
        return <ChartLine size={20} strokeWidth={1.75} aria-hidden style={{ color: UNS_TOPIC_ICON_COLORS.metric }} />;
      default:
        return null;
    }
  };

  return (
    <div className="topicDetailWrap">
      <div className="topicDetailContent">
        <Flex className="detailTitle" gap={8} align="center">
          {systemInfo.enableAutoCategorization ? (
            <UnsTopicTypeIconWrap topicType={modelInfo.dataType ?? 0}>{getFolderIcon()}</UnsTopicTypeIconWrap>
          ) : (
            <Folder size={20} />
          )}
          <Title level={2} style={{ margin: 0, width: '100%', insetInlineStart: 0 }} editable={false}>
            {modelInfo.pathName}
          </Title>
          {!folderReadOnly && (
            <Subscribe
              hidden
              showModal
              value={modelInfo.subscribeEnable}
              topic={modelInfo.topic}
              subscribeFrequency={modelInfo.subscribeFrequency}
              fileCount={countChildren}
              onChange={handleChangeSubscribe}
            />
          )}
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
      </div>
    </div>
  );
};
export default Module;
