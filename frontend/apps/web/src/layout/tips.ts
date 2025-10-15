import { getIntl } from '@/stores/i18n-store.ts';
import './tips.scss';

export const tips = (opt?: any) => [
  {
    id: 'tip1',
    classes: 'tips-guide-classes',
    title: getIntl('common.welcome', opt),
    text: `
      <div class="tips-info">${getIntl('global.tipInfo')}</div>
      <div class="tips-card">
        <span class="tips-card-icon-wrap"><i class="tips-card-icon card-icon-1"></i></span>
        <p class="tips-card-content">${getIntl('global.tip1Text1')}</p>
      </div>
      `,
    attachTo: undefined,
  },
  {
    id: 'tip2',
    classes: 'tips-guide-classes',
    title: getIntl('common.welcome', opt),
    text: `
      <div class="tips-info">${getIntl('global.tipInfo')}</div>
      <div class="tips-card">
        <span class="tips-card-icon-wrap"><i class="tips-card-icon card-icon-2"></i></span>
        <p class="tips-card-content">${getIntl('global.tip2Text1')}</p>
      </div>
      `,
    attachTo: undefined,
  },
  {
    id: 'tip3',
    classes: 'tips-guide-classes',
    title: getIntl('common.welcome', opt),
    text: `
      <div class="tips-info">${getIntl('global.tipInfo')}</div>
      <div class="tips-card">
        <span class="tips-card-icon-wrap"><i class="tips-card-icon card-icon-3"></i></span>
        <p class="tips-card-content">${getIntl('global.tip3Text1')}</p>
      </div>
      `,
    attachTo: undefined,
  },
  {
    id: 'tip4',
    classes: 'tips-guide-classes',
    title: getIntl('common.welcome', opt),
    text: `
      <div class="tips-info">${getIntl('global.tipInfo')}</div>
      <div class="tips-card">
        <span class="tips-card-icon-wrap"><i class="tips-card-icon card-icon-4"></i></span>
        <p class="tips-card-content">${getIntl('global.tip4Text1')}</p>
      </div>
      `,
    attachTo: undefined,
  },
  {
    id: 'tip5',
    classes: 'tips-guide-classes',
    title: getIntl('common.welcome', opt),
    text: `
      <div class="tips-info">${getIntl('global.tipInfo')}</div>
      <div class="tips-card">
        <span class="tips-card-icon-wrap"><i class="tips-card-icon card-icon-5"></i></span>
        <p class="tips-card-content">${getIntl('global.tip5Text1')}</p>
      </div>
      `,
    attachTo: undefined,
  },
  {
    id: 'tip6',
    classes: 'tips-guide-classes',
    title: getIntl('common.welcome', opt),
    text: `
      <div class="tips-info">${getIntl('global.tipInfo')}</div>
      <div class="tips-card">
        <span class="tips-card-icon-wrap"><i class="tips-card-icon card-icon-6"></i></span>
        <p class="tips-card-content">${getIntl('global.tip6Text1')}</p>
      </div>
      `,
    attachTo: undefined,
  },
];
