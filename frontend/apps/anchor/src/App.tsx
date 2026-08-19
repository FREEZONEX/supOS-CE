import { useEffect } from 'react';
import { Navigate, Route, Routes, useLocation } from 'react-router';
import { anchorClient } from './bridge';
import InstanceDetailPage from './pages/instance-detail';
import ModelDetailPage from './pages/model-detail';
import ModelLibraryPage from './pages/model-library';
import QrViewerPage from './pages/qr-viewer';
import SceneEditorPage from './pages/scene-editor';
import SceneLibraryPage from './pages/scene-library';

// 宿主桥接：挂载后上报 ready（宿主回推初始上下文），内部路由变化只回传路径展示、不改写宿主路由
function HostRouteBridge() {
  const location = useLocation();
  useEffect(() => {
    anchorClient.sendReady();
  }, []);
  useEffect(() => {
    anchorClient.sendRoute({ path: location.pathname });
  }, [location.pathname]);
  return null;
}

export default function App() {
  return (
    <>
      <HostRouteBridge />
      <Routes>
        <Route path="/" element={<Navigate to="/model" replace />} />
        <Route path="/model" element={<ModelLibraryPage />} />
        <Route path="/model/:id" element={<ModelDetailPage />} />
        <Route path="/model/:modelId/instances/:instanceId" element={<InstanceDetailPage />} />
        <Route path="/scene" element={<SceneLibraryPage />} />
        <Route path="/scene/:sceneId" element={<SceneEditorPage />} />
        {/* 扫码分享 viewer：免登录移动端页面 */}
        <Route path="/viewer" element={<QrViewerPage />} />
        <Route path="*" element={<Navigate to="/model" replace />} />
      </Routes>
    </>
  );
}
