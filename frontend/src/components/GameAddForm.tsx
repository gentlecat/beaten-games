import React, { Component } from 'react';
import Select from 'react-select';
import AsyncSelect from 'react-select/lib/Async';
import axios from 'axios';
import { css } from '@emotion/core';

interface SelectGameItem {
  value: string; // ID
  label: string;
}

enum AddState {
  NowPlaying = 'Now playing',
  Completed = 'Completed',
}

const loadOptions = (inputValue: string, callback: any) => {
  if (!inputValue) callback([]);

  // TODO: Set timeouts
  return axios
    .get('/api/suggest/games', {
      params: {
        q: inputValue,
      },
    })
    .then(response => {
      console.log(response.data);
      callback(
        // TODO: See if it's possible to easily improve type for the item
        response.data.map((item: any) => {
          return { value: item.name, label: item.name };
        })
      );
    })
    .catch(error => {
      console.error(error);
      callback([]);
    });
};

const submitGame = async (name: string, callback: any) => {
  return axios
    .post('/api/games/add', {
      name,
    })
    .then(response => {
      console.log(response);
    })
    .catch(error => {
      console.error(error);
      callback([]);
    });
};

export class GameAddForm extends Component<any, State> {
  public state = {
    inputValue: '',
    submitting: false,
    selectedGame: undefined,
  };

  private isValidInput = () => !!this.state.selectedGame;

  private handleInputChange = (newValue: string) => {
    this.setState({ inputValue: newValue });
    return newValue;
  };

  private handleSelectionChange = (newSelection: any) => {
    this.setState({
      selectedGame: newSelection.value,
      //inputValue: newSelection.value,
    });
  };

  private handleSubmit = async () => {
    if (!this.isValidInput()) {
      console.error('Unable to submit invalid input');
      return;
    }

    this.setState({ submitting: true });
    await submitGame(this.state.inputValue, () =>
      this.setState({ inputValue: '', submitting: false })
    );
  };

  public render() {
    return (
      <div>
        <span>Quick add:</span>

        <div
          css={css`
            display: flex;
            align-items: center;
          `}
        >
          <Select
            options={[
              { value: AddState.NowPlaying, label: AddState.NowPlaying },
              { value: AddState.Completed, label: AddState.Completed },
            ]}
            defaultValue={AddState.NowPlaying}
            isClearable={false}
            isSearchable={false}
            css={css`
              width: 100px;
            `}
          />

          <AsyncSelect
            loadOptions={loadOptions}
            cacheOptions
            onChange={this.handleSelectionChange}
            onInputChange={this.handleInputChange}
            inputValue={this.state.inputValue}
            components={{
              DropdownIndicator: undefined,
            }}
            noOptionsMessage={() => 'Enter name of the game'}
            css={css`
              flex: 1;
            `}
          />

          <button
            onClick={this.handleSubmit}
            disabled={this.state.submitting || !this.isValidInput()}
          >
            +
          </button>
        </div>
      </div>
    );
  }
}
