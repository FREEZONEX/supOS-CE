import { Flex, Segmented } from 'antd';
import { Grid, List } from '@carbon/icons-react';
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
      <Segmented
        size="small"
        value={mode}
        onChange={(v) => setMode(v)}
        options={[
          {
            value: 'card',
            icon: (
              <span title={'common.cardMode'}>
                <Grid />
              </span>
            ),
          },
          {
            value: 'list',
            icon: (
              <span title={formatMessage('common.listMode')}>
                <List />
              </span>
            ),
          },
        ]}
      />
    </Flex>
  );
};

export default ComSegmented;
