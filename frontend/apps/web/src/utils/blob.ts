export const downloadFn = ({ data, name }: { data: any; name: string }) => {
  const url = window.URL.createObjectURL(new Blob([data]));
  const link = document.createElement('a');
  link.href = url;
  link.download = name;
  document.body.appendChild(link);
  link.click(); // 模拟点击下载
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
};
