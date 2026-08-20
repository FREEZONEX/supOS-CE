import { Flex } from 'antd';
import ViewModeSegmented from '@/components/lucide-icon/ViewModeSegmented';
import useTranslate from '@/hooks/useTranslate';
import usePropsValue from '@/hooks/usePropsValue.ts';

const ComSegmented = ({
  value,
  onChange,
  defaultValue,
}: {
  value?: string;
  onChange?: (v: string) => void;
  defaultValue?: string;
}) => {
  const [mode, setMode] = usePropsValue({
    value,
    defaultValue,
    onChange,
  });
  const formatMessage = useTranslate();
  return (
    <Flex justify="flex-end" align="center" style={{ marginBottom: 16, marginTop: 16, paddingRight: 16 }}>
      <ViewModeSegmented
        value={mode}
        onChange={setMode}
        cardTitle={formatMessage('common.cardMode')}
        listTitle={formatMessage('common.listMode')}
      />
    </Flex>
  );
};

export default ComSegmented;
