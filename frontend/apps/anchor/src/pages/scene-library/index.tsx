import { Add, Edit, TrashCan } from '@carbon/icons-react';
import { App, Button, Empty, Form, Input, Modal, Pagination, Spin, Tooltip } from 'antd';
import { Clapperboard } from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { createScene, deleteScene, listScenes, updateScene, type SceneInfo } from '../../api/scenes';
import { t } from '../../i18n';
import { formatTime } from '../../utils/format';

const PAGE_SIZE = 20;

export default function SceneLibraryPage() {
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
  const [scenes, setScenes] = useState<SceneInfo[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState('');
  const [loading, setLoading] = useState(true);
  const [createOpen, setCreateOpen] = useState(false);
  const [renaming, setRenaming] = useState<SceneInfo | null>(null);
  const [renameValue, setRenameValue] = useState('');
  const [createForm] = Form.useForm<{ name: string; description?: string }>();

  const refresh = useCallback(
    async (targetPage: number, kw: string) => {
      setLoading(true);
      try {
        const data = await listScenes({ page: targetPage, size: PAGE_SIZE, keyword: kw });
        setScenes(data.list || []);
        setTotal(data.total || 0);
      } catch (e) {
        message.error((e as Error).message);
      } finally {
        setLoading(false);
      }
    },
    [message]
  );

  useEffect(() => {
    void refresh(page, keyword);
  }, [page, keyword, refresh]);

  const handleCreate = async () => {
    const values = await createForm.validateFields();
    try {
      const scene = await createScene({ name: values.name.trim(), description: values.description?.trim() });
      setCreateOpen(false);
      createForm.resetFields();
      navigate(`/scene/${scene.id}`);
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const handleRename = async () => {
    if (!renaming || !renameValue.trim()) return;
    try {
      await updateScene(renaming.id, { name: renameValue.trim() });
      setRenaming(null);
      await refresh(page, keyword);
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const confirmDelete = (scene: SceneInfo) => {
    modal.confirm({
      title: t('scene.delete'),
      content: t('scene.deleteConfirm', { name: scene.name }),
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await deleteScene(scene.id);
          await refresh(page, keyword);
        } catch (e) {
          message.error((e as Error).message);
        }
      },
    });
  };

  return (
    <div className="flex h-full flex-col">
      <div className="page-title-bar">
        <div className="flex min-w-0 flex-1 items-center gap-2">
          <Clapperboard aria-hidden size={20} strokeWidth={1.75} className="shrink-0" />
          <span className="truncate">{t('scene.title')}</span>
        </div>
        <span className="text-xs font-normal" style={{ color: 'var(--ui-description-card-color)' }}>
          {t('scene.total', { total })}
        </span>
        <Input.Search
          allowClear
          style={{ width: 220 }}
          placeholder={t('scene.search')}
          onSearch={(value) => {
            setPage(1);
            setKeyword(value.trim());
          }}
        />
        <Button type="primary" icon={<Add size={18} />} onClick={() => setCreateOpen(true)}>
          {t('scene.create')}
        </Button>
      </div>

      <div className="min-h-0 flex-1 overflow-auto p-6">
        {loading ? (
          <div className="flex h-64 items-center justify-center">
            <Spin />
          </div>
        ) : scenes.length === 0 ? (
          <div className="flex h-64 items-center justify-center">
            <Empty description={t('scene.empty')} />
          </div>
        ) : (
          <div className="model-card-grid">
            {scenes.map((scene) => (
              <div key={scene.id} className="model-card group" onClick={() => navigate(`/scene/${scene.id}`)}>
                <div className="model-card-thumb">
                  {scene.thumbnailUrl ? (
                    <img src={scene.thumbnailUrl} alt={scene.name} className="h-full w-full object-cover" />
                  ) : null}
                  <div className="model-card-actions">
                    <Tooltip title={t('scene.rename')}>
                      <Button
                        type="text"
                        size="small"
                        className="model-card-action"
                        icon={<Edit size={16} />}
                        onClick={(e) => {
                          e.stopPropagation();
                          setRenaming(scene);
                          setRenameValue(scene.name);
                        }}
                      />
                    </Tooltip>
                    <Tooltip title={t('scene.delete')}>
                      <Button
                        type="text"
                        size="small"
                        className="model-card-action"
                        danger
                        icon={<TrashCan size={16} />}
                        onClick={(e) => {
                          e.stopPropagation();
                          confirmDelete(scene);
                        }}
                      />
                    </Tooltip>
                  </div>
                </div>
                <div className="model-card-info">
                  <div className="model-card-title" title={scene.name}>
                    {scene.name}
                  </div>
                  <div className="model-card-meta">
                    {t('scene.items', { count: scene.itemCount })} · {formatTime(scene.createdTime)}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {total > PAGE_SIZE ? (
          <div className="mt-5 flex justify-center">
            <Pagination current={page} pageSize={PAGE_SIZE} total={total} showSizeChanger={false} onChange={setPage} />
          </div>
        ) : null}
      </div>

      <Modal
        title={t('scene.create')}
        open={createOpen}
        onCancel={() => setCreateOpen(false)}
        onOk={handleCreate}
        destroyOnHidden
      >
        <Form form={createForm} layout="vertical" className="pt-2">
          <Form.Item name="name" label={t('scene.name')} rules={[{ required: true, message: t('scene.nameRequired') }]}>
            <Input maxLength={128} />
          </Form.Item>
          <Form.Item name="description" label={t('scene.description')}>
            <Input.TextArea rows={2} maxLength={512} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={t('scene.rename')}
        open={Boolean(renaming)}
        okButtonProps={{ disabled: !renameValue.trim() }}
        onCancel={() => setRenaming(null)}
        onOk={handleRename}
        destroyOnHidden
      >
        <Input
          autoFocus
          value={renameValue}
          onChange={(e) => setRenameValue(e.target.value)}
          onPressEnter={() => void handleRename()}
        />
      </Modal>
    </div>
  );
}
