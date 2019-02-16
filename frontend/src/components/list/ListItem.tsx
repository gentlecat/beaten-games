import { css } from '@emotion/core';
import React from 'react';

export const ListItem = ({ game }) => {
  return (
    <div
      css={css`
        flex: 1;
        display: flex;
        margin-top: 20px;
      `}
    >
      <div css={css`
        flex: 1;
        font-weight: bold;
      `}>{game['name']}</div>
      <div css={css`
        flex: 1;
        color: grey;
        text-align: right;
      `}>{game['beaten_on']}</div>
    </div>
  );
};
