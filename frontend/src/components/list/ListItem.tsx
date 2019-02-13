import { css } from '@emotion/core';
import React from 'react';

export const ListItem = ({ game }) => {
  return <li>{game.name}</li>;
};
