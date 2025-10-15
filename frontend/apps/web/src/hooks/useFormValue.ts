import { Form } from 'antd';
const useFormValue = (name: string | (string | number)[], form: any) => {
  return Form.useWatch(name, form) || form.getFieldValue(name);
};
export default useFormValue;
