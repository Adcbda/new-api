import React, { useEffect, useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Avatar,
  Button,
  Card,
  Col,
  Form,
  Row,
  Select,
  Space,
  Spin,
  Tag,
  Typography,
  Banner,
  Descriptions,
} from '@douyinfe/semi-ui';
import {
  IconCalendarClock,
  IconClose,
  IconSave,
  IconUserAdd,
} from '@douyinfe/semi-icons';
import { Clock, RefreshCw, Users } from 'lucide-react';
import { API, showError, showSuccess } from '../../../helpers';
import {
  quotaToDisplayAmount,
  displayAmountToQuota,
} from '../../../helpers/quota';
import { useIsMobile } from '../../../hooks/common/useIsMobile';

const { Text, Title } = Typography;

const durationUnitOptions = [
  { value: 'unlimited', label: '无限期' },
  { value: 'year', label: '年' },
  { value: 'month', label: '月' },
  { value: 'day', label: '日' },
  { value: 'hour', label: '小时' },
  { value: 'custom', label: '自定义(秒)' },
];

const resetPeriodOptions = [
  { value: 'never', label: '不重置' },
  { value: 'daily', label: '每天' },
  { value: 'weekly', label: '每周' },
  { value: 'monthly', label: '每月' },
  { value: 'custom', label: '自定义(秒)' },
];

const SubscriptionsPage = () => {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [binding, setBinding] = useState(false);
  const [groupOptions, setGroupOptions] = useState([]);
  const [groupLoading, setGroupLoading] = useState(false);
  const formApiRef = useRef(null);

  const getInitValues = () => ({
    title: '',
    subtitle: '',
    duration_unit: 'unlimited',
    duration_value: 0,
    custom_seconds: 0,
    quota_reset_period: 'never',
    quota_reset_custom_seconds: 0,
    enabled: true,
    total_amount: 0,
    upgrade_group: '',
  });

  const loadPlan = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/subscription/admin/plan');
      if (res.data?.success) {
        const p = res.data.data?.plan || {};
        const values = {
          title: p.title || '',
          subtitle: p.subtitle || '',
          duration_unit: p.duration_unit || 'unlimited',
          duration_value: Number(p.duration_value || 0),
          custom_seconds: Number(p.custom_seconds || 0),
          quota_reset_period: p.quota_reset_period || 'never',
          quota_reset_custom_seconds: Number(p.quota_reset_custom_seconds || 0),
          enabled: p.enabled !== false,
          total_amount: Number(
            quotaToDisplayAmount(p.total_amount || 0).toFixed(2),
          ),
          upgrade_group: p.upgrade_group || '',
        };
        formApiRef.current?.setValues(values);
      } else {
        showError(res.data?.message || t('加载失败'));
      }
    } catch (e) {
      showError(t('请求失败'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPlan();
    setGroupLoading(true);
    API.get('/api/group')
      .then((res) => {
        if (res.data?.success) {
          setGroupOptions(res.data?.data || []);
        } else {
          setGroupOptions([]);
        }
      })
      .catch(() => setGroupOptions([]))
      .finally(() => setGroupLoading(false));
  }, []);

  const submit = async (values) => {
    if (!values.title || values.title.trim() === '') {
      showError(t('套餐标题不能为空'));
      return;
    }
    setSaving(true);
    try {
      const payload = {
        plan: {
          ...values,
          duration_value: Number(values.duration_value || 0),
          custom_seconds: Number(values.custom_seconds || 0),
          quota_reset_period: values.quota_reset_period || 'never',
          quota_reset_custom_seconds:
            values.quota_reset_period === 'custom'
              ? Number(values.quota_reset_custom_seconds || 0)
              : 0,
          total_amount: displayAmountToQuota(values.total_amount),
          upgrade_group: values.upgrade_group || '',
        },
      };
      const res = await API.put('/api/subscription/admin/plan', payload);
      if (res.data?.success) {
        showSuccess(t('更新成功，所有用户订阅已同步变更'));
      } else {
        showError(res.data?.message || t('更新失败'));
      }
    } catch (e) {
      showError(t('请求失败'));
    } finally {
      setSaving(false);
    }
  };

  const handleBindAll = async () => {
    setBinding(true);
    try {
      const res = await API.post('/api/subscription/admin/bind_all');
      if (res.data?.success) {
        const count = res.data.data?.bound_count || 0;
        showSuccess(t(`已为 ${count} 个用户绑定全局套餐`));
      } else {
        showError(res.data?.message || t('操作失败'));
      }
    } catch (e) {
      showError(t('请求失败'));
    } finally {
      setBinding(false);
    }
  };

  return (
    <Spin spinning={loading}>
      <Card className='!rounded-2xl shadow-sm border-0'>
        <div className='flex items-center justify-between mb-4'>
          <div className='flex items-center gap-2'>
            <Avatar size='small' color='blue' className='mr-1 shadow-md'>
              <IconCalendarClock size={16} />
            </Avatar>
            <Title heading={4} className='m-0'>
              {t('全局套餐管理')}
            </Title>
          </div>
          <Space>
            <Button
              theme='light'
              type='primary'
              icon={<IconUserAdd />}
              onClick={handleBindAll}
              loading={binding}
              size='small'
            >
              {t('为所有用户绑定套餐')}
            </Button>
          </Space>
        </div>

        <Banner
          type='info'
          description={t(
            '全局套餐是所有用户共用的套餐模版。修改套餐后，所有用户的订阅信息将统一变更，已用额度将清零重新计算。新注册用户将自动绑定此套餐。',
          )}
          closeIcon={null}
          className='!rounded-lg mb-4'
        />

        <Form
          initValues={getInitValues()}
          getFormApi={(api) => (formApiRef.current = api)}
          onSubmit={submit}
        >
          {({ values }) => (
            <div className='space-y-4'>
              <Card className='!rounded-2xl shadow-sm border-0'>
                <div className='flex items-center mb-2'>
                  <Avatar
                    size='small'
                    color='blue'
                    className='mr-2 shadow-md'
                  >
                    <IconCalendarClock size={16} />
                  </Avatar>
                  <div>
                    <Text className='text-lg font-medium'>
                      {t('基本信息')}
                    </Text>
                    <div className='text-xs text-gray-600'>
                      {t('套餐的基本信息')}
                    </div>
                  </div>
                </div>

                <Row gutter={12}>
                  <Col span={24}>
                    <Form.Input
                      field='title'
                      label={t('套餐标题')}
                      placeholder={t('例如：基础套餐')}
                      required
                      rules={[
                        { required: true, message: t('请输入套餐标题') },
                      ]}
                      showClear
                    />
                  </Col>

                  <Col span={24}>
                    <Form.Input
                      field='subtitle'
                      label={t('套餐副标题')}
                      placeholder={t('例如：适合轻度使用')}
                      showClear
                    />
                  </Col>

                  <Col span={12}>
                    <Form.InputNumber
                      field='total_amount'
                      label={t('总额度')}
                      required
                      min={0}
                      precision={2}
                      rules={[{ required: true, message: t('请输入总额度') }]}
                      extraText={`${t('0 表示不限')} · ${t('原生额度')}：${displayAmountToQuota(
                        values.total_amount,
                      )}`}
                      style={{ width: '100%' }}
                    />
                  </Col>

                  <Col span={12}>
                    <Form.Select
                      field='upgrade_group'
                      label={t('升级分组')}
                      showClear
                      loading={groupLoading}
                      placeholder={t('不升级')}
                      extraText={t(
                        '用户绑定套餐后会升级到该分组；修改套餐时分组也会同步更新。',
                      )}
                    >
                      <Select.Option value=''>{t('不升级')}</Select.Option>
                      {(groupOptions || []).map((g) => (
                        <Select.Option key={g} value={g}>
                          {g}
                        </Select.Option>
                      ))}
                    </Form.Select>
                  </Col>

                  <Col span={12}>
                    <Form.Switch
                      field='enabled'
                      label={t('启用状态')}
                      size='large'
                    />
                  </Col>
                </Row>
              </Card>

              <Card className='!rounded-2xl shadow-sm border-0'>
                <div className='flex items-center mb-2'>
                  <Avatar
                    size='small'
                    color='green'
                    className='mr-2 shadow-md'
                  >
                    <Clock size={16} />
                  </Avatar>
                  <div>
                    <Text className='text-lg font-medium'>
                      {t('有效期设置')}
                    </Text>
                    <div className='text-xs text-gray-600'>
                      {t('配置套餐的有效时长')}
                    </div>
                  </div>
                </div>

                <Row gutter={12}>
                  <Col span={12}>
                    <Form.Select
                      field='duration_unit'
                      label={t('有效期单位')}
                      required
                      rules={[{ required: true }]}
                    >
                      {durationUnitOptions.map((o) => (
                        <Select.Option key={o.value} value={o.value}>
                          {o.label}
                        </Select.Option>
                      ))}
                    </Form.Select>
                  </Col>

                  <Col span={12}>
                    {values.duration_unit === 'custom' ? (
                      <Form.InputNumber
                        field='custom_seconds'
                        label={t('自定义秒数')}
                        required
                        min={1}
                        precision={0}
                        rules={[
                          { required: true, message: t('请输入秒数') },
                        ]}
                        style={{ width: '100%' }}
                      />
                    ) : values.duration_unit === 'unlimited' ? (
                      <Form.InputNumber
                        field='duration_value'
                        label={t('有效期数值')}
                        disabled
                        style={{ width: '100%' }}
                        extraText={t('无限期，永不过期')}
                      />
                    ) : (
                      <Form.InputNumber
                        field='duration_value'
                        label={t('有效期数值')}
                        required
                        min={1}
                        precision={0}
                        rules={[
                          { required: true, message: t('请输入数值') },
                        ]}
                        style={{ width: '100%' }}
                      />
                    )}
                  </Col>
                </Row>
              </Card>

              <Card className='!rounded-2xl shadow-sm border-0'>
                <div className='flex items-center mb-2'>
                  <Avatar
                    size='small'
                    color='orange'
                    className='mr-2 shadow-md'
                  >
                    <RefreshCw size={16} />
                  </Avatar>
                  <div>
                    <Text className='text-lg font-medium'>
                      {t('额度重置')}
                    </Text>
                    <div className='text-xs text-gray-600'>
                      {t('所有用户统一时间重置额度')}
                    </div>
                  </div>
                </div>

                <Row gutter={12}>
                  <Col span={12}>
                    <Form.Select
                      field='quota_reset_period'
                      label={t('重置周期')}
                    >
                      {resetPeriodOptions.map((o) => (
                        <Select.Option key={o.value} value={o.value}>
                          {o.label}
                        </Select.Option>
                      ))}
                    </Form.Select>
                  </Col>
                  <Col span={12}>
                    {values.quota_reset_period === 'custom' ? (
                      <Form.InputNumber
                        field='quota_reset_custom_seconds'
                        label={t('自定义秒数')}
                        required
                        min={60}
                        precision={0}
                        rules={[
                          { required: true, message: t('请输入秒数') },
                        ]}
                        style={{ width: '100%' }}
                      />
                    ) : (
                      <Form.InputNumber
                        field='quota_reset_custom_seconds'
                        label={t('自定义秒数')}
                        min={0}
                        precision={0}
                        style={{ width: '100%' }}
                        disabled
                      />
                    )}
                  </Col>
                </Row>
              </Card>

              <div className='flex justify-end gap-2'>
                <Button
                  theme='solid'
                  type='primary'
                  onClick={() => formApiRef.current?.submitForm()}
                  icon={<IconSave />}
                  loading={saving}
                  size='large'
                >
                  {t('保存并同步到所有用户')}
                </Button>
              </div>
            </div>
          )}
        </Form>
      </Card>
    </Spin>
  );
};

export default SubscriptionsPage;
