import { css } from 'lit';

export default css`
  :host {
    display: block;
  }

  .square {
    display: flex;
    justify-content: center;
    align-items: center;
    user-select: none;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .square[data-color='b'] {
    background-color: saddlebrown;
  }
  .square[data-color='w'] {
    background-color: lightgoldenrodyellow;
  }

  .square[data-selected='true'] {
    background-color: #138808;
  }
  .square[data-highlight-from='true'] {
    background-color: plum;
  }
  .square[data-highlight-to='true'] {
    background-color: purple;
  }
  .square[data-under-attack='true'] {
    position: relative;
  }
  .square[data-under-attack='true']::after {
    content: '';
    position: absolute;
    inset: 0;
    background-color: skyblue;
    border-radius: 100vmax;
    transform: scale(0.6);
  }
`;
