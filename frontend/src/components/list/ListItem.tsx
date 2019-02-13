import { css } from '@emotion/core';
import React from 'react';

export const ListItem = () => {
  return (
    <nav>
      <h1
        css={css`
          font-size: 2em;
        `}
      >
        Beaten games
      </h1>
    </nav>
  );
};
