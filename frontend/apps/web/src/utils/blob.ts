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

export async function blobToJsonUsingTextMethod(blob: Blob) {
  try {
    const jsonString = await blob.text();
    return JSON.stringify(JSON.parse(jsonString), null, 2);
  } catch (error) {
    throw new Error('转换失败，原因：' + error);
  }
}
