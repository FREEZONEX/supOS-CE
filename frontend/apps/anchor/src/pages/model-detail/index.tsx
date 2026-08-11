import { Add, CenterCircle, ChevronLeft, Maximize, Renew, TrashCan } from '@carbon/icons-react';
import { App, Button, Empty, Form, Input, InputNumber, Modal, Slider, Spin, Table, Tooltip, Tree } from 'antd';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { getModel, updateModel, uploadAsset, type ModelInfo } from '../../api/models';
import {
  createInstance,
  deleteInstance,
  listInstances,
  parseBinding,
  type InstanceInfo,
  type SelectedObjectNode,
} from '../../api/instances';
import { useSuspended } from '../../bridge/use-suspended';
import { ResizableSplit } from '../../components/resizable-split';
import { TopicSelect } from '../../components/topic-select';
import { t } from '../../i18n';
import { flattenNodeTree } from '../../viewer/gltf-nodes';
import { SplatViewer } from '../../viewer/splat-viewer';
import { useModelViewer } from '../../viewer/use-model-viewer';
import { getModelFormat } from '../../utils/thumbnail';
import { StatusTag } from './status-tag';

const heightText = (value: number) => (value > 0 ? `${value.toFixed(3)} M` : '-- M');

// 视口右上角按钮（对齐云端 ViewportIconButton 风格）
const viewportButtonStyle: React.CSSProperties = {
  width: 32,
  height: 32,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: 4,
  cursor: 'pointer',
  border: '1px solid var(--ui-header-splitter-color, #e0e0e0)',
  background: 'var(--ui-bg-color)',
  color: 'var(--ui-text-color)',
};

export default function ModelDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
  const [model, setModel] = useState<ModelInfo | null>(null);
  const [instances, setInstances] = useState<InstanceInfo[]>([]);
  const [nodeEditorOpen, setNodeEditorOpen] = useState(false);
  const [checkedPaths, setCheckedPaths] = useState<string[]>([]);
  // 节点树异步加载，defaultExpandAll 不生效——树数据到位后受控全展开
  const [expandedPaths, setExpandedPaths] = useState<string[]>([]);
  const [createOpen, setCreateOpen] = useState(false);
  const [replacing, setReplacing] = useState(false);
  const replaceInputRef = useRef<HTMLInputElement>(null);
  const [createForm] = Form.useForm<{ name: string; height?: number }>();
  const [createTopic, setCreateTopic] = useState<{ unsNodeId: string; topic: string }>();

  const supported = model && getModelFormat(model.originFile || model.name) !== 'splat';
  const { containerRef, viewer } = useModelViewer(supported ? model?.fileUrl : undefined);
  // SPLAT 走独立查看器，高度由其解析包围盒后上报
  const [splatHeight, setSplatHeight] = useState(0);
  const modelHeight = viewer.modelHeight || splatHeight;
  // 页签隐藏/窗口后台 → 暂停渲染循环
  const suspended = useSuspended();
  useEffect(() => {
    viewer.setSuspended(suspended);
  }, [viewer, suspended]);

  const selectedNodes: SelectedObjectNode[] = useMemo(() => {
    try {
      const parsed = JSON.parse(model?.nodesJson || '[]');
      return Array.isArray(parsed) ? parsed : (parsed?.list ?? []);
    } catch {
      return [];
    }
  }, [model?.nodesJson]);

  const refreshModel = useCallback(() => {
    if (!id) return;
    getModel(id)
      .then(setModel)
      .catch((e) => message.error((e as Error).message));
  }, [id, message]);

  const refreshInstances = useCallback(() => {
    if (!id) return;
    listInstances({ modelId: Number(id), page: 1, size: 200 })
      .then((data) => setInstances(data.list || []))
      .catch(() => setInstances([]));
  }, [id]);

  useEffect(() => {
    refreshModel();
    refreshInstances();
  }, [refreshModel, refreshInstances]);

  // 对齐云端：新模型（尚未选过对象节点）首次进入直接展示 Select Object Nodes；选过则进详情
  useEffect(() => {
    if (model && supported && selectedNodes.length === 0) setNodeEditorOpen(true);
  }, [model, supported, selectedNodes.length]);

  useEffect(() => {
    setExpandedPaths(flattenNodeTree(viewer.nodeTree).map((node) => node.path));
  }, [viewer.nodeTree]);

  const saveSelectedNodes = async () => {
    if (!model) return;
    const flat = flattenNodeTree(viewer.nodeTree);
    const byPath = new Map(flat.map((node) => [node.path, node]));
    const list: SelectedObjectNode[] = checkedPaths
      .map((path) => byPath.get(path))
      .filter(Boolean)
      .map((node) => ({ nodeID: node!.tier0Id, name: node!.title, path: node!.path }));
    if (list.length === 0) {
      message.error(t('model.nodes.required'));
      return;
    }
    try {
      await updateModel(model.id, { nodesJson: JSON.stringify(list) });
      setNodeEditorOpen(false);
      viewer.highlight([]);
      refreshModel();
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const handleCreateInstance = async () => {
    if (!model) return;
    const values = await createForm.validateFields();
    if (!createTopic?.topic) {
      message.error(t('instance.topicRequired'));
      return;
    }
    const binding = {
      selectedObjects: selectedNodes,
      motionMappings: [],
      dataTags: [],
    };
    try {
      await createInstance({
        modelId: model.id,
        name: values.name.trim(),
        unsNodeId: createTopic.unsNodeId,
        topic: createTopic.topic,
        bindingJson: JSON.stringify(binding),
        height: values.height || modelHeight || 0,
      });
      setCreateOpen(false);
      createForm.resetFields();
      setCreateTopic(undefined);
      refreshInstances();
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const confirmDeleteInstance = (instance: InstanceInfo) => {
    modal.confirm({
      title: t('instance.delete'),
      content: t('instance.deleteConfirm', { name: instance.name }),
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await deleteInstance(instance.id);
          refreshInstances();
        } catch (e) {
          message.error((e as Error).message);
        }
      },
    });
  };

  const instanceColumns = [
    {
      title: t('instance.name'),
      dataIndex: 'name',
      render: (name: string, record: InstanceInfo) => (
        <a
          className="table-name-link"
          onClick={() => navigate(`/model/${model?.id}/instances/${record.id}`)}
        >
          {name}
        </a>
      ),
    },
    { title: t('instance.topic'), dataIndex: 'topic', ellipsis: true },
    {
      title: t('instance.nodeMapping'),
      render: (_: unknown, record: InstanceInfo) => {
        const binding = parseBinding(record.bindingJson);
        const mapped = new Set(binding.motionMappings.map((item) => item.nodeID)).size;
        return `${mapped}/${binding.selectedObjects.length}`;
      },
    },
    {
      title: t('instance.dataTag'),
      render: (_: unknown, record: InstanceInfo) => parseBinding(record.bindingJson).dataTags.length,
    },
    {
      title: t('instance.status'),
      dataIndex: 'status',
      render: (status: number) => <StatusTag status={status} />,
    },
    {
      title: '',
      width: 50,
      render: (_: unknown, record: InstanceInfo) => (
        <Button
          type="text"
          size="small"
          danger
          icon={<TrashCan size={15} />}
          onClick={() => confirmDeleteInstance(record)}
        />
      ),
    },
  ];

  return (
    <div className="flex h-full flex-col">
      <div className="page-title-bar">
        <div className="flex min-w-0 flex-1 items-center gap-3">
          <Button
            variant="outlined"
            color="default"
            style={{ paddingLeft: '5.5px', gap: '3px' }}
            onClick={() => navigate('/model')}
          >
            <span className="inline-flex items-center gap-2">
              <ChevronLeft size={16} />
              {t('common.back')}
            </span>
          </Button>
          <div className="truncate">{model?.name ?? t('common.loading')}</div>
        </div>
        <input
          ref={replaceInputRef}
          type="file"
          accept=".glb,.gltf,.splat"
          className="hidden"
          onChange={async (e) => {
            const file = e.target.files?.[0];
            if (!file || !model) return;
            if (getModelFormat(file.name) === 'model') {
              message.error(t('model.invalidType'));
              return;
            }
            setReplacing(true);
            try {
              const asset = await uploadAsset(file);
              await updateModel(model.id, {
                fileAssetId: asset.fileId,
                originFile: file.name,
                fileSize: file.size,
              });
              refreshModel();
            } catch (err) {
              message.error((err as Error).message);
            } finally {
              setReplacing(false);
              if (replaceInputRef.current) replaceInputRef.current.value = '';
            }
          }}
        />
        <Button icon={<Renew size={16} />} loading={replacing} onClick={() => replaceInputRef.current?.click()}>
          {t('model.replaceModel')}
        </Button>
      </div>

      <div className="min-h-0 flex-1">
        <ResizableSplit
          defaultLeftWidth={35}
          minLeftWidth={20}
          maxLeftWidth={50}
          left={
            /* 左：3D viewer + 爆炸视图（SPLAT 走独立高斯泼溅查看器），宽度可拖拽 */
            <div className="relative h-full w-full">
              {!supported && model?.fileUrl ? (
                <SplatViewer fileUrl={model.fileUrl} onHeight={setSplatHeight} suspended={suspended} />
              ) : null}
              <div ref={containerRef} className="h-full w-full" style={!supported ? { display: 'none' } : undefined} />
              {supported && viewer.status !== 'ready' ? (
                <div
                  className="absolute inset-0 flex items-center justify-center gap-2 text-sm"
                  style={{ background: 'var(--ui-bg-color)', color: 'var(--ui-description-card-color)' }}
                >
                  {viewer.status === 'error' ? (
                    `${t('model.viewer.error')}: ${viewer.error}`
                  ) : (
                    <>
                      <Spin size="small" /> {t('model.viewer.loading')}
                    </>
                  )}
                </div>
              ) : null}
              {viewer.status === 'ready' ? (
                <div className="absolute right-4 bottom-3 left-4">
                  <div className="mb-1 text-xs font-medium" style={{ color: 'var(--ui-text-color)' }}>
                    {t('model.explodedView')}
                  </div>
                  <Slider min={0} max={100} defaultValue={0} onChange={(value) => viewer.setExploded(value / 100)} />
                </div>
              ) : null}
              {/* 视口右上按钮（对齐云端）：重置视角 / 全屏 */}
              {supported ? (
                <div className="absolute top-3 right-3 z-10 flex gap-1">
                  <Tooltip title={t('viewer.resetView')}>
                    <button type="button" style={viewportButtonStyle} onClick={() => viewer.resetView()}>
                      <CenterCircle size={16} />
                    </button>
                  </Tooltip>
                  <Tooltip title={t('viewer.fullscreen')}>
                    <button
                      type="button"
                      style={viewportButtonStyle}
                      onClick={() => {
                        const el = containerRef.current?.parentElement;
                        if (document.fullscreenElement) void document.exitFullscreen();
                        else void el?.requestFullscreen();
                      }}
                    >
                      <Maximize size={16} />
                    </button>
                  </Tooltip>
                </div>
              ) : null}
            </div>
          }
          right={
            /* 右：高度 / 节点选择 / 实例列表；编辑节点时整个面板原地切换为节点树（对齐云端流程，不用弹窗） */
            <div className="flex h-full min-w-0 flex-col overflow-hidden">
              {nodeEditorOpen ? (
                <div className="flex min-h-0 flex-1 flex-col p-5">
                  <div className="mb-2 text-base font-semibold">{t('model.nodes.editTitle')}</div>
                  <div
                    className="min-h-0 flex-1 overflow-auto rounded border p-2"
                    style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
                  >
                    <Tree
                      checkable
                      checkStrictly
                      expandedKeys={expandedPaths}
                      onExpand={(keys) => setExpandedPaths(keys.map(String))}
                      checkedKeys={{ checked: checkedPaths, halfChecked: [] }}
                      treeData={viewer.nodeTree}
                      fieldNames={{ key: 'path', title: 'title', children: 'children' }}
                      onCheck={(checked) => {
                        const keys = Array.isArray(checked) ? checked : checked.checked;
                        setCheckedPaths(keys.map(String));
                      }}
                      titleRender={(node) => (
                        <span
                          onMouseEnter={() => viewer.highlight([(node as { path: string }).path])}
                          onMouseLeave={() => viewer.highlight(checkedPaths)}
                        >
                          {(node as { title: string }).title}
                        </span>
                      )}
                    />
                  </div>
                  <div className="mt-3 flex shrink-0 justify-end gap-2">
                    {/* 首次选择（还没有已保存节点）不允许取消，必须先确认选择（对齐云端） */}
                    {selectedNodes.length > 0 ? (
                      <Button
                        onClick={() => {
                          setNodeEditorOpen(false);
                          viewer.highlight([]);
                        }}
                      >
                        {t('common.cancel')}
                      </Button>
                    ) : null}
                    <Button type="primary" onClick={saveSelectedNodes}>
                      {t('common.confirm')}
                    </Button>
                  </div>
                </div>
              ) : (
                <>
                  {/* 顶部元信息条：源文件 + 高度，双列铺满 */}
                  <div
                    className="flex items-center gap-4 border-b px-5 py-4"
                    style={{ borderColor: 'var(--ui-line-color, #e0e0e0)' }}
                  >
                    <div className="min-w-0 flex-1">
                      <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                        {t('model.detail.originFile')}
                      </div>
                      <div
                        className="mt-0.5 truncate text-sm font-medium"
                        title={model?.originFile || undefined}
                      >
                        {model?.originFile || '--'}
                      </div>
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                        {t('model.height')}
                      </div>
                      <div className="mt-0.5 flex items-baseline gap-1 text-sm font-medium">
                        <span>{modelHeight > 0 ? modelHeight.toFixed(3) : '--'}</span>
                        <span style={{ color: 'var(--ui-description-card-color)' }}>M</span>
                      </div>
                    </div>
                  </div>

                  <div className="min-h-0 flex-1 overflow-auto p-5">
                  {/* SPLAT 无节点树，不展示对象节点选择（对齐云端） */}
                  {supported ? (
                    <>
                      <div className="mb-2 flex items-center justify-between">
                        <div>
                          <div className="text-base font-semibold">{t('model.nodes.title')}</div>
                          <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                            {t('model.nodes.count', { count: selectedNodes.length })}
                          </div>
                        </div>
                        <Button
                          disabled={viewer.status !== 'ready'}
                          onClick={() => {
                            setCheckedPaths(selectedNodes.map((node) => node.path));
                            setNodeEditorOpen(true);
                          }}
                        >
                          {t('model.nodes.edit')}
                        </Button>
                      </div>
                      {selectedNodes.length === 0 ? (
                        <div
                          className="mb-6 rounded border border-dashed p-4 text-xs"
                          style={{
                            borderColor: 'var(--ui-header-splitter-color, #e0e0e0)',
                            color: 'var(--ui-description-card-color)',
                          }}
                        >
                          {t('model.nodes.empty')}
                        </div>
                      ) : (
                        <div
                          className="mb-6 max-h-64 overflow-auto rounded border"
                          style={{ borderColor: 'var(--ui-line-color, #e0e0e0)' }}
                        >
                          {selectedNodes.map((node, index) => (
                            <div
                              key={node.path}
                              className="selected-object-node-row flex cursor-pointer items-center gap-3 border-b px-4 last:border-b-0"
                              style={{ borderColor: 'var(--ui-line-color, #e0e0e0)' }}
                              onMouseEnter={() => viewer.highlight([node.path])}
                              onMouseLeave={() => viewer.highlight([])}
                            >
                              <span
                                className="w-6 shrink-0 text-sm tabular-nums"
                                style={{ color: 'var(--ui-description-card-color)' }}
                              >
                                {String(index + 1).padStart(2, '0')}
                              </span>
                              <span className="truncate text-sm font-medium" title={node.path}>
                                {node.path}
                              </span>
                            </div>
                          ))}
                        </div>
                      )}
                    </>
                  ) : null}

                  <div className="mb-2 flex items-center justify-between">
                    <div>
                      <div className="text-base font-semibold">{t('instance.list.title')}</div>
                      <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                        {t('instance.list.count', { count: instances.length })}
                      </div>
                    </div>
                    <Button type="primary" icon={<Add size={16} />} onClick={() => setCreateOpen(true)}>
                      {t('instance.create')}
                    </Button>
                  </div>
                  <Table
                    rowKey="id"
                    columns={instanceColumns}
                    dataSource={instances}
                    pagination={false}
                    locale={{ emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} /> }}
                  />
                  </div>
                </>
              )}
            </div>
          }
        />
      </div>

      {/* 创建实例弹窗 */}
      <Modal
        title={t('instance.create')}
        open={createOpen}
        onCancel={() => setCreateOpen(false)}
        onOk={handleCreateInstance}
        destroyOnHidden
      >
        <Form form={createForm} layout="vertical" className="pt-2">
          <Form.Item label={t('instance.topic')} required>
            <TopicSelect value={createTopic} onChange={setCreateTopic} />
          </Form.Item>
          <Form.Item
            name="name"
            label={t('instance.name')}
            rules={[{ required: true, message: t('instance.nameRequired') }]}
          >
            <Input maxLength={128} />
          </Form.Item>
          <Form.Item name="height" label={t('instance.height')} extra={t('instance.heightHint')}>
            <InputNumber
              min={0.001}
              max={1000}
              precision={3}
              style={{ width: '100%' }}
              placeholder={heightText(modelHeight)}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
