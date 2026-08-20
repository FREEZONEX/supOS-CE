import {
  cloneNotebook,
  createFolder,
  createNotebook,
  deleteFolder,
  deleteNotebook,
  getFolderTree,
  getNotebookList,
  shutdownNotebook,
  updateFolder,
  updateNotebook,
} from '@/apis/core-api/notebook';
import { ButtonPermission } from '@/common-types/button-permission';
import ComLayout from '@/components/com-layout';
import ComLeft from '@/components/com-layout/ComLeft';
import { AuthButton } from '@/components/auth';
import { useTranslate } from '@/hooks';
import type {
  FolderNode,
  NotebookDetail,
  NotebookItem,
  NotebookTabType,
  NotebookTreeItem,
} from '@/pages/notebook/types';
import { PageTitleIcon } from '@/components/lucide-icon';
import { Run as PlayIcon } from '@/components/lucide-icon/carbon';
import classNames from 'classnames';
import { Monitor } from 'lucide-react';
import { App } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { createSearchParams, useLocation, useNavigate } from 'react-router';

import FolderModal from './components/FolderModal';
import ImportModal from './components/ImportModal';
import NotebookModal from './components/NotebookModal';
import NotebookSidebar from './components/NotebookSidebar';
import NotebookTable from './components/NotebookTable';
import { isNotebookRunning } from './components/NotebookStatusTag';
import styles from './index.module.scss';

const findFolderNode = (nodes: FolderNode[], targetId: number): FolderNode | null => {
  for (const node of nodes) {
    if (node.id === targetId) return node;
    const found = findFolderNode(node.children || [], targetId);
    if (found) return found;
  }
  return null;
};

const collectFolderIds = (node: FolderNode): number[] => {
  const ids = [node.id];
  for (const child of node.children || []) {
    ids.push(...collectFolderIds(child));
  }
  return ids;
};

const NotebookPage: FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [tab, setTab] = useState<NotebookTabType>('all');
  const [search, setSearch] = useState('');
  const [folderSearch, setFolderSearch] = useState('');
  const [selectedFolderId, setSelectedFolderId] = useState(0);
  const [folders, setFolders] = useState<FolderNode[]>([]);
  const [rootNotebooks, setRootNotebooks] = useState<NotebookTreeItem[]>([]);
  const [items, setItems] = useState<NotebookItem[]>([]);
  const [folderModalOpen, setFolderModalOpen] = useState(false);
  const [folderModalLoading, setFolderModalLoading] = useState(false);
  const [editingFolder, setEditingFolder] = useState<FolderNode | null>(null);
  const [notebookModalOpen, setNotebookModalOpen] = useState(false);
  const [notebookModalLoading, setNotebookModalLoading] = useState(false);
  const [editingNotebook, setEditingNotebook] = useState<NotebookDetail | null>(null);
  const [importModalOpen, setImportModalOpen] = useState(false);
  const loadItemsRequestRef = useRef(0);
  const previousPathnameRef = useRef<string | null>(null);

  const loadFolders = useCallback(async () => {
    const response = await getFolderTree();
    setFolders(response?.list || []);
    setRootNotebooks(response?.rootNotebooks || []);
  }, []);

  const loadItems = useCallback(async () => {
    const requestId = ++loadItemsRequestRef.current;
    setLoading(true);
    try {
      const response = await getNotebookList({
        listType: tab,
        folderId: tab === 'running' ? 0 : selectedFolderId,
        search,
      });
      if (requestId !== loadItemsRequestRef.current) {
        return;
      }
      setItems(response?.list || []);
    } finally {
      if (requestId === loadItemsRequestRef.current) {
        setLoading(false);
      }
    }
  }, [search, selectedFolderId, tab]);

  useEffect(() => {
    void loadFolders();
  }, [loadFolders]);

  useEffect(() => {
    void loadItems();
  }, [loadItems]);

  const handleReload = useCallback(async () => {
    await Promise.all([loadFolders(), loadItems()]);
  }, [loadFolders, loadItems]);

  useEffect(() => {
    const previousPathname = previousPathnameRef.current;
    if (location.pathname === '/notebook' && previousPathname && previousPathname !== location.pathname) {
      void handleReload();
    }
    previousPathnameRef.current = location.pathname;
  }, [handleReload, location.pathname]);

  const runningNotebookIds = useMemo(
    () => items.filter((item) => item.type === 'notebook' && isNotebookRunning(item.status)).map((item) => item.id),
    [items]
  );
  const displayedItems = useMemo(() => {
    if (tab !== 'running' || selectedFolderId <= 0) {
      return items;
    }
    const selectedFolder = findFolderNode(folders, selectedFolderId);
    if (!selectedFolder) {
      return items;
    }
    const folderIdSet = new Set(collectFolderIds(selectedFolder));
    return items.filter((item) => item.type === 'notebook' && folderIdSet.has(item.folderId || 0));
  }, [folders, items, selectedFolderId, tab]);

  return (
    <ComLayout className={styles.notebookPage}>
      <div className={styles.pageColumn}>
        <div className={styles.subHeader}>
          <div className={styles.pageTitleWrap}>
            <PageTitleIcon resourceKey="notebook.view" />
            <span className={styles.pageTitle}>{formatMessage('Notebook.title', {}, 'Notebook')}</span>
          </div>
          <AuthButton
            auth={ButtonPermission['Notebook.manage']}
            type="primary"
            className={styles.importBtn}
            onClick={() => setImportModalOpen(true)}
          >
            {formatMessage('common.import', {}, 'Import')}
          </AuthButton>
        </div>

        <div className={styles.globalHeader}>
          <button
            type="button"
            className={classNames(
              styles.globalTab,
              tab === 'all' && styles.globalTabActive,
              tab !== 'all' && styles.globalTabInactive
            )}
            onClick={() => setTab('all')}
          >
            <span className={styles.globalTabInner}>
              <Monitor size={16} strokeWidth={1.75} />
              <span>{formatMessage('Notebook.tabAll', {}, 'All')}</span>
            </span>
            {tab === 'all' ? <span className={styles.globalTabInk} aria-hidden /> : null}
          </button>
          <button
            type="button"
            className={classNames(
              styles.globalTab,
              tab === 'running' && styles.globalTabActive,
              tab !== 'running' && styles.globalTabInactive
            )}
            onClick={() => {
              setTab('running');
              setSelectedFolderId(0);
            }}
          >
            <span className={styles.globalTabInner}>
              <PlayIcon size={16} strokeWidth={1.75} />
              <span>{formatMessage('Notebook.tabRunning', {}, 'Running')}</span>
            </span>
            {tab === 'running' ? <span className={styles.globalTabInk} aria-hidden /> : null}
          </button>
        </div>

        <div className={styles.pageBody}>
          <ComLeft
            className={styles.sidebar}
            resize
            defaultWidth={360}
            style={{ overflow: 'hidden', display: 'flex', flexDirection: 'column', minHeight: 0 }}
          >
            <div className={styles.sidebarInner}>
              <div className={styles.fileBrowserHeader}>
                {formatMessage('Notebook.fileBrowser', {}, 'File Browser')}
              </div>
              <div className={styles.sidebarContent}>
                <NotebookSidebar
                  folders={folders}
                  rootNotebooks={rootNotebooks}
                  runningNotebookIds={runningNotebookIds}
                  selectedFolderId={selectedFolderId}
                  tab={tab}
                  searchValue={folderSearch}
                  onSearchChange={setFolderSearch}
                  onSelect={(folderId) => {
                    setSelectedFolderId(folderId);
                  }}
                  onOpenNotebook={(notebookId) => navigate(`/notebook/editor/${notebookId}`)}
                  onNewFolder={() => {
                    setEditingFolder(null);
                    setFolderModalOpen(true);
                  }}
                  onNewNotebook={() => {
                    setEditingNotebook(null);
                    setNotebookModalOpen(true);
                  }}
                  onRefresh={() => void handleReload()}
                  onEditNotebook={(notebookId) => {
                    const allNotebooks = [
                      ...rootNotebooks,
                      ...folders.flatMap(function collectNotebooks(f: FolderNode): NotebookTreeItem[] {
                        return [...(f.notebooks || []), ...(f.children || []).flatMap(collectNotebooks)];
                      }),
                    ];
                    const nb = allNotebooks.find((n) => n.id === notebookId);
                    if (nb) {
                      setEditingNotebook({
                        id: nb.id,
                        name: nb.name,
                        description: '',
                        folderId: nb.folderId,
                        filePath: nb.path,
                        status: nb.status || 'stopped',
                        fileSize: 0,
                      });
                      setNotebookModalOpen(true);
                    }
                  }}
                  onEditFolder={(folderId) => {
                    const found = findFolderNode(folders, folderId);
                    if (found) {
                      setEditingFolder(found);
                      setFolderModalOpen(true);
                    }
                  }}
                  onDeleteFolder={async (folderId) => {
                    await deleteFolder(folderId);
                    message.success(formatMessage('Notebook.folderDeletedSuccess', {}, 'Folder deleted'));
                    if (selectedFolderId === folderId) {
                      setSelectedFolderId(0);
                    }
                    await handleReload();
                  }}
                  onDeleteNotebook={async (notebookId) => {
                    await deleteNotebook(notebookId);
                    message.success(formatMessage('Notebook.deletedSuccess', {}, 'Notebook deleted'));
                    await handleReload();
                  }}
                />
              </div>
            </div>
          </ComLeft>
          <div className={styles.main}>
            <div className={styles.tableFrame}>
              <div className={styles.tableWrap}>
                <NotebookTable
                  items={displayedItems}
                  loading={loading}
                  onOpenFolder={(item) => {
                    setSelectedFolderId(item.id);
                    setTab('all');
                  }}
                  onOpenNotebook={(item) => {
                    const search =
                      item.path && item.type === 'notebook'
                        ? `?${createSearchParams({ filePath: item.path }).toString()}`
                        : '';
                    navigate(`/notebook/editor/${item.id}${search}`);
                  }}
                  onEditFolder={(item) => {
                    setEditingFolder(findFolderNode(folders, item.id));
                    setFolderModalOpen(true);
                  }}
                  onEditNotebook={(item) => {
                    setEditingNotebook({
                      id: item.id,
                      name: item.name,
                      description: item.description,
                      folderId: item.folderId || 0,
                      filePath: item.path,
                      status: item.status || 'stopped',
                      owner: item.owner,
                      fileSize: item.fileSize || 0,
                    });
                    setNotebookModalOpen(true);
                  }}
                  onDeleteNotebook={async (item) => {
                    await deleteNotebook(item.id);
                    message.success(formatMessage('Notebook.deletedSuccess', {}, 'Notebook deleted'));
                    await handleReload();
                  }}
                  onCloneNotebook={async (item) => {
                    await cloneNotebook(item.id);
                    message.success(formatMessage('Notebook.clonedSuccess', {}, 'Notebook cloned'));
                    await handleReload();
                  }}
                  onShutdownNotebook={async (item) => {
                    await shutdownNotebook(item.id);
                    message.success(formatMessage('Notebook.shutdownSuccess', {}, 'Notebook shutdown'));
                    await loadItems();
                  }}
                  onDeleteFolder={async (item) => {
                    await deleteFolder(item.id);
                    message.success(formatMessage('Notebook.folderDeletedSuccess', {}, 'Folder deleted'));
                    if (selectedFolderId === item.id) {
                      setSelectedFolderId(0);
                    }
                    await handleReload();
                  }}
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <FolderModal
        open={folderModalOpen}
        title={
          editingFolder
            ? formatMessage('Notebook.renameFolder', {}, 'Rename Folder')
            : formatMessage('Notebook.createFolder', {}, 'Create Folder')
        }
        initialName={editingFolder?.name}
        confirmLoading={folderModalLoading}
        onCancel={() => {
          setFolderModalOpen(false);
          setEditingFolder(null);
        }}
        onSubmit={async (name) => {
          setFolderModalLoading(true);
          try {
            if (editingFolder) {
              await updateFolder(editingFolder.id, { name });
              message.success(formatMessage('Notebook.folderUpdatedSuccess', {}, 'Folder updated'));
              setSelectedFolderId(editingFolder.id);
            } else {
              const created = await createFolder({ name, parentId: selectedFolderId });
              setSearch('');
              setSelectedFolderId(created?.id || selectedFolderId);
              message.success(formatMessage('Notebook.folderCreatedSuccess', {}, 'Folder created'));
            }
            setFolderModalOpen(false);
            setEditingFolder(null);
            await handleReload();
          } finally {
            setFolderModalLoading(false);
          }
        }}
      />

      <ImportModal
        open={importModalOpen}
        folders={folders}
        selectedFolderId={selectedFolderId}
        onCancel={() => setImportModalOpen(false)}
        onSuccess={async () => {
          setImportModalOpen(false);
          await handleReload();
        }}
      />

      <NotebookModal
        open={notebookModalOpen}
        title={
          editingNotebook
            ? formatMessage('Notebook.editNotebook', {}, 'Edit Notebook')
            : formatMessage('Notebook.createNotebook', {}, 'Create Notebook')
        }
        folders={folders}
        notebook={editingNotebook}
        selectedFolderId={selectedFolderId}
        confirmLoading={notebookModalLoading}
        onCancel={() => {
          setNotebookModalOpen(false);
          setEditingNotebook(null);
        }}
        onSubmit={async (values) => {
          setNotebookModalLoading(true);
          try {
            if (editingNotebook) {
              await updateNotebook(editingNotebook.id, values);
              message.success(formatMessage('Notebook.updatedSuccess', {}, 'Notebook updated'));
              setSelectedFolderId(values.folderId ?? editingNotebook.folderId ?? 0);
            } else {
              const created = await createNotebook(values);
              setSearch('');
              setSelectedFolderId(created?.folderId ?? values.folderId ?? 0);
              message.success(formatMessage('Notebook.createdSuccess', {}, 'Notebook created'));
            }
            setTab('all');
            setNotebookModalOpen(false);
            setEditingNotebook(null);
            await handleReload();
          } finally {
            setNotebookModalLoading(false);
          }
        }}
      />
    </ComLayout>
  );
};

export default NotebookPage;
