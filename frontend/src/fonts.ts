import { css } from '@emotion/core';

export const charterFont = css`
  @font-face {
    font-family: 'Charter';
    src: url('/static/fonts/charter/web/charter_regular.woff2') format('woff2'),
      url('/static/fonts/charter/web/charter_regular.woff') format('woff');
    font-weight: normal;
    font-style: normal;
    font-display: auto;
  }

  @font-face {
    font-family: 'Charter';
    src: url('/static/fonts/charter/web/charter_italic.woff2') format('woff2'),
      url('/static/fonts/charter/web/charter_italic.woff') format('woff');
    font-weight: normal;
    font-style: italic;
    font-display: auto;
  }

  @font-face {
    font-family: 'Charter';
    src: url('/static/fonts/charter/web/charter_bold.woff2') format('woff2'),
      url('/static/fonts/charter/web/charter_bold.woff') format('woff');
    font-weight: bold;
    font-style: normal;
    font-display: auto;
  }

  @font-face {
    font-family: 'Charter';
    src: url('/static/fonts/web/charter_bold_italic.woff2') format('woff2'),
      url('/static/fonts/charter/web/charter_bold_italic.woff') format('woff');
    font-weight: bold;
    font-style: italic;
    font-display: auto;
  }
`;
