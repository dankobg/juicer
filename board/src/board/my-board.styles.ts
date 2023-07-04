import { css } from 'lit';

export default css`
  :host {
    display: block;
    --board-size: 30rem;
  }

  .outer {
    width: 100%;
    display: grid;
    grid-template-areas:
      'ranks-label board'
      '. files-label';
    place-content: center;
    margin-top: 2rem;
  }

  .ranks-label,
  .files-label {
    display: grid;
    gap: 0;
    justify-content: center;
    justify-items: center;
    background-color: #ddd;
    user-select: none;
  }

  .ranks-label {
    grid-area: ranks-label;
    grid-template-columns: 1fr;
    grid-template-rows: repeat(8, 1fr);
    grid-auto-flow: row;
    width: 1.5rem;
    height: var(--board-size);
  }

  .files-label {
    grid-area: files-label;
    grid-template-columns: repeat(8, 1fr);
    width: var(--board-size);
    height: 1.5rem;
  }

  .ranks-label > div,
  .files-label > div {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
  }

  .board {
    grid-area: board;
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    grid-template-rows: repeat(8, 1fr);
    grid-auto-flow: row;
    justify-content: center;
    gap: 0;
    width: var(--board-size);
    height: var(--board-size);
    border: 1px solid #333;
    position: relative;
  }

  .board[data-orientation='b'] {
    transform: rotate(0.5turn);
  }

  .board[data-orientation='b'] > my-square {
    transform: rotate(0.5turn);
  }
`;
