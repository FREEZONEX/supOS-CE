import ProModal from '@/components/pro-modal';
import useTranslate from '@/hooks/useTranslate.ts';
import { type Key, useState } from 'react';
import { useImmer } from 'use-immer';
import ProTree from '@/components/pro-tree/index.ts';
import type { TreeProps } from 'antd';
const defaultData: any[] = [
  {
    title: '首页',
    key: 'home',
    children: [
      {
        title: '概述',
        key: 'overview',
        isLeaf: true,
        parentKey: 'home',
      },
      {
        title: '示例',
        key: 'examples',
        isLeaf: true,
        parentKey: 'home',
      },
      {
        title: '资源监控',
        key: 'resource-monitor',
        isLeaf: true,
        parentKey: 'home',
      },
    ],
  },
  {
    title: '数据管理',
    key: 'data-management',
    children: [
      {
        title: '数据连接',
        key: 'data-connection',
        isLeaf: true,
        parentKey: 'data-management',
      },
      {
        title: '数据建模',
        key: 'data-modeling',
        isLeaf: true,
        parentKey: 'data-management',
      },
      {
        title: '事件流程',
        key: 'event-process',
        isLeaf: true,
        parentKey: 'data-management',
      },
      {
        title: '采集器网关管理',
        key: 'collector-gateway',
        isLeaf: true,
        parentKey: 'data-management',
      },
    ],
  },
  {
    title: '工具集',
    key: 'toolset',
    children: [
      {
        title: 'CICD',
        key: 'cicd',
        isLeaf: true,
        parentKey: 'toolset',
      },
      {
        title: '数据源链接',
        key: 'data-source-link',
        isLeaf: true,
        parentKey: 'toolset',
      },
      {
        title: '报警管理',
        key: 'alarm-management',
        isLeaf: true,
        parentKey: 'toolset',
      },
      {
        title: 'SQL编辑器',
        key: 'sql-editor',
        isLeaf: true,
        parentKey: 'toolset',
      },
    ],
  },
  {
    title: '应用集',
    key: 'application-set',
    children: [
      {
        title: 'PRIDE智能监控',
        key: 'pride-monitor',
        isLeaf: true,
        parentKey: 'application-set',
      },
    ],
  },
  {
    title: '系统配置',
    key: 'system-config',
    children: [
      {
        title: '用户管理',
        key: 'user-management',
        isLeaf: true,
        parentKey: 'system-config',
      },
      {
        title: 'APP管理',
        key: 'app-management',
        isLeaf: true,
        parentKey: 'system-config',
      },
    ],
  },
];
const HOME_KEY = 'home';
const SYSTEM_CONFIG_KEY = 'system-config';
const disableDND = [HOME_KEY, SYSTEM_CONFIG_KEY];

const loop = (data: any[], key: Key, callback: (node: any, i: number, data: any[]) => void) => {
  for (let i = 0; i < data.length; i++) {
    if (data[i].key === key) {
      return callback(data[i], i, data);
    }
    if (data[i].children) {
      loop(data[i].children!, key, callback);
    }
  }
};

const useMenuSetting = () => {
  const [open, setOpen] = useState(false);
  // const [gData, setData] = useState(defaultData);
  const [gData, setData] = useImmer(defaultData);
  const formatMessage = useTranslate();

  const onMenuModalOpen = () => {
    setOpen(true);
  };

  const onDrop: TreeProps['onDrop'] = (info) => {
    const { node, dragNode, dropPosition, dropToGap } = info;
    const dropKey = node.key;
    const dragKey = dragNode.key;
    const dropPos = node.pos.split('-');
    const resolvedDropPosition = dropPosition - Number(dropPos[dropPos.length - 1]);
    setData((draft) => {
      // 1. 查找并从原位置删除拖拽节点
      let dragObj: any;
      loop(draft, dragKey, (item, index, arr) => {
        arr.splice(index, 1);
        dragObj = item;
      });

      // 2. 将拖拽节点插入到新位置
      if (dragObj) {
        if (!dropToGap) {
          // 拖拽到目标节点内部
          loop(draft, dropKey, (item) => {
            item.children = item.children || [];
            item.children.unshift(dragObj);
          });
        } else {
          // 拖拽到目标节点的间隙
          let ar: any[] = [];
          let i: number = 0;
          loop(draft, dropKey, (_item, index, arr) => {
            ar = arr;
            i = index;
          });
          const insertIndex = resolvedDropPosition === -1 ? i : i + 1;
          ar.splice(insertIndex, 0, dragObj);
        }
      }
    });
  };

  const MenuModal = (
    <ProModal title={formatMessage('account.menuSettings')} open={open} onCancel={() => setOpen(false)}>
      <div>
        <ProTree
          multiple
          onDrop={onDrop}
          treeData={gData}
          draggable={{
            icon: false,
            nodeDraggable: (node: any) => {
              return !(disableDND.includes(node.key) || (node.parentKey && disableDND.includes(node.parentKey)));
            },
          }}
          allowDrop={({ dropNode, dropPosition, dragNode }: any) => {
            // 检查父节点是否在禁用列表中
            if (dropNode.parentKey && disableDND.includes(dropNode.parentKey)) return false;
            // 检查特定节点的放置限制
            if (dropNode.key === HOME_KEY) return ![0, -1].includes(dropPosition);
            if (dropNode.key === SYSTEM_CONFIG_KEY) return ![1, 0].includes(dropPosition);
            // 叶子节点不能作为父节点放置
            if (dropNode.isLeaf && dropPosition === 0) return false;
            if (!dragNode?.isLeaf) {
              // 文件夹不能拖到文件夹内
              if (dropNode?.isLeaf) {
                // 文件任何位置不行
                return false;
              } else if (dropPosition === 0) {
                // 文件夹内不行
                return false;
              } else {
                return true;
              }
            }
            return true;
          }}
        />
      </div>
    </ProModal>
  );

  return {
    onMenuModalOpen,
    MenuModal,
  };
};

export default useMenuSetting;
