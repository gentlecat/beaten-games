import React from 'react';
import { css } from '@emotion/core';

export const Header = () => {
  return (
    <nav>
      <h1
        css={css`
          font-size: 2em;
        `}
      >
        Game Collector <br />
        <span
          css={css`
            font-size: 0.5em;
            color: red;
          `}
        >
          WIP
        </span>
      </h1>
    </nav>
  );
};
