import { useSize } from 'ahooks';

interface UseMediaSizeType {
  media?: number;
}

const useMediaSize = ({ media = 640 }: UseMediaSizeType = {}) => {
  const size = useSize(document.querySelector('body'));
  return {
    width: size?.width,
    height: size?.height,
    isH5: size?.width && size?.width < media,
  };
};

export default useMediaSize;
