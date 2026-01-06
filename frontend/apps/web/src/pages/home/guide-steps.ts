import type { Step, StepOptions } from 'shepherd.js';
import homeFlow from '@/assets/guide/home-flow.gif';
import homeFlowChartreuse from '@/assets/guide/home-flow-chartreuse.gif';
import './guide-steps.scss';
import { getIntl } from '@/stores/i18n-store.ts';

// 新手导航步骤
export const guideSteps: (navigate: any, opt: any, theme: any) => Array<StepOptions> | Array<Step> = (
  navigate: any,
  opt: any,
  theme: any
) => {
  return [
    {
      id: 'home_start',
      cancelIcon: {
        enabled: false,
      },
      classes: 'guide-home-classes',
      text: `
        <img src=${theme.includes('chartreuse') ? homeFlowChartreuse : homeFlow} class="guide-home-logo"/>
        <div class="guide-home-text">
          <div class="guide-home-title">${getIntl('common.welcome', opt)}</div>
          <div class="guide-home-info">${getIntl('home.guideText')}</div>
        </div>

    `,
      attachTo: undefined,
      buttons: [
        {
          action() {
            return this.next();
          },
          text: getIntl('home.guide1Next'),
          classes: 'home-guide-next',
        },
        {
          action() {
            return this.complete();
          },
          text: getIntl('home.guide1Exit'),
          classes: 'home-guide-exit prev-class',
        },
      ],
    },
    {
      id: 'home_uns',
      classes: 'guide-home-uns-classes',
      cancelIcon: {
        enabled: false,
      },
      title: getIntl('home.guideUnsTitle'),
      attachTo: {
        element: '#home_route_uns',
        on: 'right',
      },
      buttons: [
        {
          action() {
            return this.complete();
          },
          text: getIntl('home.guide1Exit'),
          classes: 'prev-class',
        },
        {
          action() {
            navigate('/uns');
            return this.complete();
          },
          text: getIntl('common.next'),
        },
      ],
    },
    // {
    //   id: 'homepage',
    //   title: getIntl('home.guide1Title'),
    //   text: `
    // <ul class="user-guide-list">
    //   <li>${getIntl('home.guide1Text1')}</li>
    //   <li>${getIntl('home.guide1Text2')}</li>
    // </ul>
    // `,
    //   attachTo: undefined,
    // },
    // {
    //   id: 'home_step1',
    //   title: getIntl('home.guide2Title'),
    //   text: `
    // ${getIntl('home.guide2Text1', opt)}
    // `,
    //   attachTo: {
    //     element: '#custom_menu_left',
    //     on: 'bottom',
    //   },
    // },
    // {
    //   id: 'home_step2',
    //   title: getIntl('home.guide3Title'),
    //   classes: 'guide-home-topBar-classes',
    //   text: `
    // ${getIntl('home.guide3Text1')}
    // <ul class="user-guide-list">
    //   <li>${getIntl('home.guide3Text6')}</li>
    //   <li>${getIntl('home.guide3Text2')}</li>
    //   <li>${getIntl('home.guide3Text7')}</li>
    //   <li>${getIntl('home.guide3Text3')}</li>
    //   <li>${getIntl('home.guide3Text4')}</li>
    //   <li>${getIntl('home.guide3Text5')}</li>
    // </ul>
    // `,
    //   attachTo: {
    //     element: '#custom_menu_right',
    //     on: 'bottom-end',
    //   },
    // },
    // {
    //   id: 'home_step3',
    //   title: getIntl('home.guide4Title'),
    //   text: `
    // ${getIntl('home.guide4Text1', { ...opt, strong: (chunks: string) => `<strong>${chunks}</strong>` })}
    // <ul class="user-guide-list">
    //   <li>${getIntl('home.guide4Text2')}</li>
    //   <li>${getIntl('home.guide4Text3')}</li>
    // </ul>
    // `,
    //   attachTo: {
    //     element: '#home_section_step1',
    //     on: 'auto',
    //   },
    //   buttons: [
    //     {
    //       action() {
    //         return this.back();
    //       },
    //       text: getIntl('common.prev'),
    //       classes: 'prev-class',
    //     },
    //     {
    //       action() {
    //         const userGuideRoute = storageOpt.get(SUPOS_USER_GUIDE_ROUTES);
    //         const currentRoute = find(userGuideRoute, (route) => route?.menu?.url === '/uns' && route?.menu?.picked);
    //         if (currentRoute && currentRoute?.isVisited === false) {
    //           this.addStep({
    //             id: 'home_step4',
    //             title: getIntl('home.guide5Title'),
    //             text: `
    //           ${getIntl('home.guide5Text1')}
    //           `,
    //             attachTo: undefined,
    //             buttons: [
    //               {
    //                 action() {
    //                   return this.complete();
    //                 },
    //                 text: getIntl('common.cancel'),
    //                 classes: 'prev-class',
    //               },
    //               {
    //                 action() {
    //                   navigate('/uns');
    //                   return this.complete();
    //                 },
    //                 text: getIntl('global.tipGo'),
    //               },
    //             ],
    //           });
    //           return this.next();
    //         } else {
    //           return this.complete();
    //         }
    //       },
    //       text: getIntl('global.tipDone'),
    //     },
    //   ],
    // },
  ];
};
