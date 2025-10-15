import { Button, Flex, ConfigProvider, Card } from 'antd';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';
import { Tag, WatsonHealth3DMprToggle, Folders, DocumentMultiple_02 } from '@carbon/icons-react';
import FileTable from './FileTable';
import FolderTable from './FolderTable.tsx';
import TagTable from './TagTable.tsx';
import TemplateTable from './TemplateTable.tsx';
import { ProModal } from '@/components';
import { useMemo, useState } from 'react';
import { useTreeStore } from '@/pages/uns/store/treeStore.tsx';

const DataSubscript = () => {
  const formatMessage = useTranslate();
  const [modalKey, setModalKey] = useState<string | null>(null);
  const {
    lazyTree,
    setTreeType,
    setTreeMap,
    scrollTreeNode,
    loadData,
    setSelectedNode,
    setCurrentTreeMapType,
    setSearchValue,
  } = useTreeStore((state) => ({
    setTreeMap: state.setTreeMap,
    setTreeType: state.setTreeType,
    scrollTreeNode: state.scrollTreeNode,
    loadData: state.loadData,
    setSelectedNode: state.setSelectedNode,
    setCurrentTreeMapType: state.setCurrentTreeMapType,
    setSearchValue: state.setSearchValue,
    lazyTree: state.lazyTree,
  }));
  const onMoreClick = (item: any) => {
    setModalKey(item.key);
    setOpen(true);
  };
  const onNameClick = (item: any, type: any, scrollTreeNodeFn = scrollTreeNode) => {
    setOpen(false);
    setTreeMap(false);
    setTreeType(type);
    if (type === 'uns') {
      setSearchValue(item.name);
    } else {
      setSearchValue('');
    }
    loadData({ reset: true }, (data) => {
      const findInfo = data?.find((f) => f.id === item.id);
      if (findInfo || !lazyTree) {
        setSelectedNode(data?.find((f) => f.id === item.id) || item);
      } else {
        setSelectedNode(item, true);
      }
      scrollTreeNodeFn?.(item.id);
      setCurrentTreeMapType('all');
    });
  };
  const [open, setOpen] = useState(false);
  const list = [
    {
      label: 'uns.model',
      icon: <Folders size={18} />,
      onClick: onMoreClick,
      content: FolderTable,
      key: 'uns.model',
    },
    {
      label: 'uns.instance',
      icon: <DocumentMultiple_02 size={18} />,
      onClick: onMoreClick,
      content: FileTable,
      key: 'uns.instance',
    },
    {
      label: 'common.template',
      icon: <WatsonHealth3DMprToggle size={18} />,
      onClick: onMoreClick,
      content: TemplateTable,
      key: 'common.template',
    },
    {
      label: 'common.label',
      icon: <Tag size={18} />,
      onClick: onMoreClick,
      content: TagTable,
      key: 'common.label',
    },
  ];

  const ModalContent = useMemo(() => {
    if (modalKey) {
      return list.find((f) => f.key === modalKey)?.content;
    }
    return null;
  }, [modalKey]);
  return (
    <ConfigProvider
      theme={{
        components: {
          Table: {
            headerBg: 'var(--supos-switchwrap-bg-color)',
            borderColor: '#f0f0f0',
          },
          Card: {
            bodyPadding: 16,
            headerPadding: 16,
            headerBg: 'var(--supos-bg-color)',
            colorBgContainer: 'var(--supos-bg-color)',
          },
        },
      }}
    >
      <div style={{ padding: 20 }}>
        <Flex wrap gap={8} style={{ background: 'var(--supos-switchwrap-bg-color)', padding: 8 }}>
          {list?.map((item) => {
            const Content = item.content;
            return (
              <Card
                title={
                  <Flex align="center" style={{ fontWeight: 500, fontSize: 18 }} gap={8}>
                    {item.icon} {formatMessage(item.label)}
                  </Flex>
                }
                extra={<Button onClick={() => item.onClick(item)}>{formatMessage('common.more')}</Button>}
                key={item.key}
                variant="borderless"
                className={styles['subscript-table'] + ' ' + styles['custom-table']}
              >
                <Content
                  key={scrollTreeNode + ''}
                  onNameClick={(item: any, type: string) => {
                    onNameClick(item, type, scrollTreeNode);
                  }}
                />
              </Card>
            );
          })}
        </Flex>
      </div>
      <ProModal
        title={formatMessage(modalKey || 'uns.dataSubscriptions')}
        open={open}
        onCancel={() => {
          setOpen(false);
          setModalKey(null);
        }}
        maskClosable={false}
        destroyOnHidden
        forceRender={true}
      >
        {ModalContent &&
          ((isFullscreen) => <ModalContent isFullscreen={isFullscreen} isSimple={false} onNameClick={onNameClick} />)}
      </ProModal>
    </ConfigProvider>
  );
};

export default DataSubscript;
