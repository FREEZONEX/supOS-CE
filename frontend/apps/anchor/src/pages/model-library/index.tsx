import { Add, TrashCan, Edit } from '@carbon/icons-react';
import { App, Button, Empty, Input, Modal, Pagination, Spin, Tooltip } from 'antd';
import { Box } from 'lucide-react';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router';
import {
  createDemoModel,
  createModel,
  deleteModel,
  listModels,
  updateModel,
  uploadAsset,
  type ModelInfo,
} from '../../api/models';
import { t } from '../../i18n';
import { formatBytes, formatTime } from '../../utils/format';
import { generateModelThumbnailFile, getModelFormat } from '../../utils/thumbnail';

const MAX_FILE_SIZE = 50 * 1024 * 1024;
const PAGE_SIZE = 20;

export default function ModelLibraryPage() {
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [models, setModels] = useState<ModelInfo[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState('');
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [loadingSample, setLoadingSample] = useState(false);
  const [renaming, setRenaming] = useState<ModelInfo | null>(null);
  const [renameValue, setRenameValue] = useState('');

  const refresh = useCallback(
    async (targetPage: number, kw: string) => {
      setLoading(true);
      try {
        const data = await listModels({ page: targetPage, size: PAGE_SIZE, keyword: kw });
        setModels(data.list || []);
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

  const handleUpload = async (file: File) => {
    if (getModelFormat(file.name) === 'model') {
      message.error(t('model.invalidType'));
      return;
    }
    if (file.size > MAX_FILE_SIZE) {
      message.error(t('model.fileTooLarge'));
      return;
    }
    setUploading(true);
    try {
      const thumbnailFile = await generateModelThumbnailFile(file);
      const asset = await uploadAsset(file);
      const thumbnailAsset = await uploadAsset(thumbnailFile);
      await createModel({
        name: file.name.replace(/\.[^/.]+$/, ''),
        originFile: file.name,
        fileAssetId: asset.fileId,
        thumbnailAssetId: thumbnailAsset.fileId,
        fileSize: file.size,
      });
      setPage(1);
      await refresh(1, keyword);
    } catch (e) {
      message.error(`${t('model.uploadFailed')}: ${(e as Error).message}`);
    } finally {
      setUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  const handleRename = async () => {
    if (!renaming) return;
    const name = renameValue.trim();
    if (!name) return;
    try {
      await updateModel(renaming.id, { name });
      setRenaming(null);
      await refresh(page, keyword);
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  // 一键创建示例（服务端流程：UNS 节点 + mock 数据流 + 模型 + 绑定实例，对齐云端 modelCreateDemo）
  const handleLoadSample = async () => {
    setLoadingSample(true);
    try {
      await createDemoModel();
      message.success(t('model.sampleSuccess'));
      setPage(1);
      await refresh(1, keyword);
    } catch (e) {
      message.error((e as Error).message);
    } finally {
      setLoadingSample(false);
    }
  };

  // 与云端一致：已存在 Demo Model 时不再展示示例卡片
  const hasDemoModel = models.some((model) => model.name === 'Demo Model' || model.originFile === 'demo_model.glb');

  const confirmDelete = (model: ModelInfo) => {
    modal.confirm({
      title: t('model.delete'),
      content: t('model.deleteConfirm', { name: model.name }),
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await deleteModel(model.id);
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
          <Box aria-hidden size={20} strokeWidth={1.75} className="shrink-0" />
          <span className="truncate">{t('model.title')}</span>
        </div>
        <span className="text-xs font-normal" style={{ color: 'var(--ui-description-card-color)' }}>
          {t('model.total', { total })}
        </span>
        <Input.Search
          allowClear
          style={{ width: 220 }}
          placeholder={t('model.search')}
          onSearch={(value) => {
            setPage(1);
            setKeyword(value.trim());
          }}
        />
        <Button
          type="primary"
          icon={<Add size={18} />}
          loading={uploading}
          onClick={() => fileInputRef.current?.click()}
        >
          {uploading ? t('model.uploading') : t('model.upload')}
        </Button>
        <input
          ref={fileInputRef}
          type="file"
          accept=".glb,.gltf,.splat"
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) void handleUpload(file);
          }}
        />
      </div>

      <div className="min-h-0 flex-1 overflow-auto p-6">
        {loading ? (
          <div className="flex h-64 items-center justify-center">
            <Spin />
          </div>
        ) : models.length === 0 && (hasDemoModel || keyword) ? (
          <div className="flex h-64 flex-col items-center justify-center">
            <Empty description={t('model.empty')}>
              <span className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                {t('model.uploadHint')}
              </span>
            </Empty>
          </div>
        ) : (
          <div className="model-card-grid">
            {!hasDemoModel && !keyword ? (
              <div
                className="flex flex-col justify-between rounded-md border p-4"
                style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
              >
                <div
                  className="relative flex aspect-[2/1] w-full items-center justify-center overflow-hidden rounded border border-dashed"
                  style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
                >
                  <img
                    src="/anchor/demo/demo_thumbnail.svg"
                    alt={t('model.sampleTitle')}
                    className="h-full w-full object-cover"
                  />
                </div>
                <div className="mt-3">
                  <div className="text-sm font-semibold">{t('model.sampleTitle')}</div>
                  <div className="mt-1 text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                    {t('model.sampleDesc')}
                  </div>
                </div>
                <Button className="mt-3 w-full" loading={loadingSample} onClick={() => void handleLoadSample()}>
                  {loadingSample ? t('model.sampleLoading') : t('model.sampleLoad')}
                </Button>
              </div>
            ) : null}
            {models.map((model) => (
              <div key={model.id} className="model-card group" onClick={() => navigate(`/model/${model.id}`)}>
                <div className="model-card-thumb">
                  {model.thumbnailUrl ? (
                    <img src={model.thumbnailUrl} alt={model.name} className="h-full w-full object-cover" />
                  ) : null}
                  <div className="model-card-actions">
                    <Tooltip title={t('model.rename')}>
                      <Button
                        type="text"
                        size="small"
                        className="model-card-action"
                        icon={<Edit size={16} />}
                        onClick={(e) => {
                          e.stopPropagation();
                          setRenaming(model);
                          setRenameValue(model.name);
                        }}
                      />
                    </Tooltip>
                    <Tooltip title={t('model.delete')}>
                      <Button
                        type="text"
                        size="small"
                        className="model-card-action"
                        danger
                        icon={<TrashCan size={16} />}
                        onClick={(e) => {
                          e.stopPropagation();
                          confirmDelete(model);
                        }}
                      />
                    </Tooltip>
                  </div>
                </div>
                <div className="model-card-info">
                  <div className="model-card-title" title={model.name}>
                    {model.name}
                  </div>
                  <div className="model-card-meta">
                    {formatBytes(model.fileSize)} · {formatTime(model.createdTime)}
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
        title={t('model.rename')}
        open={Boolean(renaming)}
        okButtonProps={{ disabled: !renameValue.trim() }}
        onCancel={() => setRenaming(null)}
        onOk={handleRename}
        destroyOnHidden
      >
        <Input
          autoFocus
          placeholder={t('model.nameRequired')}
          value={renameValue}
          onChange={(e) => setRenameValue(e.target.value)}
          onPressEnter={() => void handleRename()}
        />
      </Modal>
    </div>
  );
}
