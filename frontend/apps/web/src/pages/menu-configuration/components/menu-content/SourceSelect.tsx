import { Flex } from 'antd';
import ComSelect from '@/components/com-select';
import { useTranslate } from '@/hooks';
import usePropsValue from '@/hooks/usePropsValue.ts';
import { getKongRoutesApi } from '@/apis/inter-api';

const SourceSelect = ({ value, onChange }: { value?: any; onChange?: (v: any) => void }) => {
  const formatMessage = useTranslate();

  const [v, setV] = usePropsValue<{ routeSource: number; route?: any }>({
    value,
    onChange,
  });
  return (
    <Flex justify="space-between">
      <ComSelect
        style={{ width: '30%' }}
        value={v?.routeSource}
        onChange={(v) => {
          setV({
            route: null,
            routeSource: v,
          });
        }}
        options={[
          {
            label: formatMessage('MenuConfiguration.manual'),
            value: 1,
          },
          {
            label: formatMessage('MenuConfiguration.routeFetching'),
            value: 2,
          },
        ]}
      />
      {v?.routeSource === 2 && (
        <ComSelect
          allowClear
          showSearch
          value={v?.route?.name}
          filterOption={(input, option) => ((option?.name as string) ?? '').toLowerCase().includes(input.toLowerCase())}
          onChange={(_, option) => {
            setV((pre: any) => ({
              ...pre,
              route: option,
            }));
          }}
          onBlurRequest={false}
          isRequest={v?.routeSource}
          fieldNames={{ value: 'name', label: 'name' }}
          api={getKongRoutesApi}
          style={{ width: 'calc(70% - 8px)' }}
        />
      )}
    </Flex>
  );
};

export default SourceSelect;
