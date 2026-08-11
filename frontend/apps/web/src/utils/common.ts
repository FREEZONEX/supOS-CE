export const isJsonString = (value: any): boolean => {
  if (!value) return false;
  try {
    const jsonVal = JSON.parse(value);
    if (['[object Object]', '[object Array]'].includes(Object.prototype.toString.call(jsonVal))) {
      return true;
    } else {
      return false;
    }
    // eslint-disable-next-line
  } catch (err) {
    return false;
  }
};

type CopyCallback = (success: boolean) => void;

const fallbackCopyToClipboard = (text: string): boolean => {
  const textarea = document.createElement('textarea');
  const activeElement = document.activeElement as HTMLElement | null;
  const selection = document.getSelection();
  const originalRange = selection && selection.rangeCount > 0 ? selection.getRangeAt(0) : null;

  textarea.value = text;
  textarea.setAttribute('readonly', 'readonly');
  textarea.style.position = 'fixed';
  textarea.style.top = '0';
  textarea.style.left = '-9999px';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);

  try {
    textarea.focus();
    textarea.select();
    textarea.setSelectionRange(0, text.length);
    return document.execCommand('copy');
  } catch (err) {
    console.log(err);
    return false;
  } finally {
    document.body.removeChild(textarea);
    if (selection) {
      selection.removeAllRanges();
      if (originalRange) {
        selection.addRange(originalRange);
      }
    }
    activeElement?.focus?.();
  }
};

// 使用 const 定义一个函数，并同时指定参数和返回值的类型
export const copyToClipboard: (text: string, callback?: CopyCallback) => void = (
  text: string,
  callback?: CopyCallback
): void => {
  const handleResult = (success: boolean) => {
    if (callback) {
      callback(success);
    }
  };

  if (navigator.clipboard && typeof navigator.clipboard.writeText === 'function' && window.isSecureContext) {
    // 使用 Clipboard API
    const writeResult = navigator.clipboard.writeText(text);
    if (writeResult && typeof writeResult.then === 'function') {
      writeResult
        .then(() => {
          handleResult(true); // 成功时处理结果
        })
        .catch(() => {
          handleResult(fallbackCopyToClipboard(text)); // 失败时回退到旧方式
        });
    } else {
      handleResult(fallbackCopyToClipboard(text));
    }
  } else {
    handleResult(fallbackCopyToClipboard(text));
  }
};

export function canModifyParentHref() {
  try {
    if (window.parent === window || !window.parent) {
      return false;
    }

    const parentOrigin = window.parent.location.origin;
    const currentOrigin = window.location.origin;

    if (parentOrigin === currentOrigin) {
      // 同域
      return true;
    }

    return false;
  } catch (e) {
    console.log(e);
    // 跨域
    return -1;
  }
}

export interface FetchAllPagesParams {
  pageNo: number;
  pageSize: number;
}

export interface FetchAllPagesResponse<T> {
  data: T[];
  total: number;
}

export interface FetchAllPagesOptions {
  pageSize: number;
  maxConcurrency: number;
  startPageNo?: number;
}

export const fetchAllPagesWithConcurrency = async <T>(
  fetchPage: (params: FetchAllPagesParams) => Promise<FetchAllPagesResponse<T>>,
  options: FetchAllPagesOptions
) => {
  const startPageNo = options.startPageNo ?? 1;
  const pageSize = options.pageSize;
  const safeMaxConcurrency = Math.max(1, options.maxConcurrency);

  const pageDataMap: Record<number, T[]> = {};

  const firstPage = await fetchPage({
    pageNo: startPageNo,
    pageSize,
  });
  const firstPageData = Array.isArray(firstPage?.data) ? firstPage.data : [];

  let total = Number(firstPage?.total ?? firstPageData.length);
  let loadedCount = firstPageData.length;
  let nextPageNo = startPageNo + 1;

  pageDataMap[startPageNo] = firstPageData;

  while (loadedCount < total) {
    const remainingPages = Math.ceil((total - loadedCount) / pageSize);
    const batchSize = Math.min(safeMaxConcurrency, remainingPages);
    const batchPageNos = Array.from({ length: batchSize }, (_, index) => nextPageNo + index);

    const batchResults = await Promise.all(
      batchPageNos.map((pageNo) => {
        return fetchPage({ pageNo, pageSize });
      })
    );

    let hasAnyData = false;
    batchResults.forEach((result, index) => {
      const pageNo = batchPageNos[index];
      const pageData = Array.isArray(result?.data) ? result.data : [];

      total = Number(result?.total ?? total);
      pageDataMap[pageNo] = pageData;
      loadedCount += pageData.length;

      if (pageData.length > 0) {
        hasAnyData = true;
      }
    });

    if (!hasAnyData) {
      break;
    }

    nextPageNo += batchSize;
  }

  const maxFetchedPageNo = Math.max(...Object.keys(pageDataMap).map((key) => Number(key)));
  const allData: T[] = [];
  for (let pageNo = startPageNo; pageNo <= maxFetchedPageNo; pageNo += 1) {
    if (Array.isArray(pageDataMap[pageNo])) {
      allData.push(...pageDataMap[pageNo]);
    }
  }

  const normalizedData = total > 0 ? allData.slice(0, total) : allData;

  return {
    data: normalizedData,
    total,
  };
};
