import { css } from 'lit';

export default css`
  :host {
    display: block;
    z-index: 99;

    position: absolute;
    width: calc(30rem / 8);
    height: calc(30rem / 8);
  }

  .piece {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    background-size: contain;
    background-position: center;
    background-repeat: no-repeat;
    overflow: hidden;
  }

  .piece[data-selected='true'] {
    transform: scale(1.2);
  }

  .piece[data-piece-fen='r'] {
    background-image: var(--br-piece-set);
  }
  .piece[data-piece-fen='b'] {
    background-image: var(--bb-piece-set);
  }
  .piece[data-piece-fen='n'] {
    background-image: var(--bn-piece-set);
  }
  .piece[data-piece-fen='q'] {
    background-image: var(--bq-piece-set);
  }
  .piece[data-piece-fen='k'] {
    background-image: var(--bk-piece-set);
  }
  .piece[data-piece-fen='p'] {
    background-image: var(--bp-piece-set);
  }
  .piece[data-piece-fen='R'] {
    background-image: var(--wr-piece-set);
  }
  .piece[data-piece-fen='B'] {
    background-image: var(--wb-piece-set);
  }
  .piece[data-piece-fen='N'] {
    background-image: var(--wn-piece-set);
  }
  .piece[data-piece-fen='Q'] {
    background-image: var(--wq-piece-set);
  }
  .piece[data-piece-fen='K'] {
    background-image: var(--wk-piece-set);
  }
  .piece[data-piece-fen='P'] {
    background-image: var(--wp-piece-set);
  }
`;
