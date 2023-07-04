import { LitElement, html } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { map } from 'lit/directives/map.js';
import {
  Board,
  Coordinate,
  Piece,
  PieceSetStyle,
  PieceSymbol,
  Square,
  WHITE,
  convertCoordinateToRowAndColumn,
  swapColor,
} from '../shared/types';
import { Color } from '../shared/types';
import styles from './my-board.styles';

import '../square/my-square';

function drawBoardASCII(squares: Square[][]): string {
  let s = '   +------------------------+\n';

  for (let r = 0; r < squares.length; r++) {
    for (let c = 0; c < squares[r].length; c++) {
      if (c % 8 == 0) {
        s += ' ' + r + ' |';
      }

      if (squares[r][c].hasPiece()) {
        s += ' ' + squares[r][c].piece?.toFENPieceSymbol() + ' ';
      } else {
        s += ' - ';
      }

      if ((c + 1) % 8 == 0) {
        s += '| \n';
      }
    }
  }

  s += '   +------------------------+\n';
  s += '     a  b  c  d  e  f  g  h';

  return s;
}

@customElement('my-board')
export class MyBoard extends LitElement {
  static styles = styles;

  board: Board = new Board();

  @state()
  clickedSquares: Square[] = [];

  @property()
  pieceSet: PieceSetStyle = 'cburnett';

  @property()
  orientation: Color = WHITE;

  @property()
  pieceMovement: 'click' | 'drag' = 'drag';

  private toggleOrientation() {
    this.orientation = swapColor(this.orientation);
  }

  private addPiece(symbol: PieceSymbol, color: Color, coord: Coordinate) {
    const piece = new Piece(symbol, color);
    this.setPiece(piece, coord);
  }

  private setPiece(piece: Piece, coord: Coordinate) {
    const [row, col] = convertCoordinateToRowAndColumn(coord);
    this.board.squares[row][col].piece = piece;
  }

  private isSelected(square: Square): boolean {
    const sq = this.clickedSquares?.[0];
    if (!sq) {
      return false;
    }
    return square.coordinate === sq.coordinate;
  }

  private isUnderAttack(square: Square): boolean {
    const selectedSquare = this.clickedSquares?.[0];
    const piece = selectedSquare?.piece;
    if (!piece) {
      return false;
    }

    const legalMoves: any[] = [];

    const filtered = legalMoves.filter(move => {
      return (
        move.fromSquare.piece?.symbol === piece.symbol && move.fromSquare.coordinate === selectedSquare?.coordinate
      );
    });

    return filtered.findIndex(move => move.toSquare.coordinate === square.coordinate) !== -1;
  }

  private isHistoryFromMove(square: Square): boolean {
    return false;
    // return this.chess.history.at(-1)?.move.fromSquare.coordinate === square.coordinate;
  }

  private isHistoryToMove(square: Square): boolean {
    // return this.chess.history.at(-1)?.move.toSquare.coordinate === square.coordinate;
    return false;
  }

  private resetSelected(): void {
    this.clickedSquares = [];
  }

  private onSquareClick(e: CustomEvent) {
    if (this.pieceMovement !== 'click') {
      return;
    }

    const square = e.detail.square as Square;

    if (this.clickedSquares.length === 2) {
      this.resetSelected();
    }
    if (this.clickedSquares.length === 0) {
      if (square.isEmpty()) {
        return;
      }
      this.clickedSquares = [...this.clickedSquares, square];
      return;
    }
    if (this.clickedSquares.length === 1) {
      const from = this.clickedSquares[0];
      const to = square;
      if (from.equals(to)) {
        this.resetSelected();
        return;
      }

      if (from.piece?.color === square?.piece?.color) {
        this.resetSelected();
        this.clickedSquares = [...this.clickedSquares, square];
      } else {
        this.clickedSquares = [...this.clickedSquares, to];
        this.move(from, to);
        this.resetSelected();
      }
    }
  }

  private renderRanksLabel() {
    if (this.orientation === WHITE) {
      return html`
        <div>8</div>
        <div>7</div>
        <div>6</div>
        <div>5</div>
        <div>4</div>
        <div>3</div>
        <div>2</div>
        <div>1</div>
      `;
    }

    return html`
      <div>1</div>
      <div>2</div>
      <div>3</div>
      <div>4</div>
      <div>5</div>
      <div>6</div>
      <div>7</div>
      <div>8</div>
    `;
  }

  private renderFilesLabel() {
    if (this.orientation === WHITE) {
      return html`
        <div>a</div>
        <div>b</div>
        <div>c</div>
        <div>d</div>
        <div>e</div>
        <div>f</div>
        <div>g</div>
        <div>h</div>
      `;
    }

    return html`
      <div>h</div>
      <div>g</div>
      <div>f</div>
      <div>e</div>
      <div>d</div>
      <div>c</div>
      <div>b</div>
      <div>a</div>
    `;
  }

  private move(from: Square, to: Square) {
    console.log(`from: ${from.coordinate} to: ${to.coordinate}`);

    const copy = this.board.clone();

    const tmp = copy.squares[from.row][from.column].piece;
    copy.squares[to.row][to.column].piece = tmp;
    copy.squares[from.row][from.column].piece = null;
    this.board = copy;
  }

  private flipBoard() {
    this.toggleOrientation();
  }

  render() {
    return html`
      <pre>${this.board.ascii()}</pre>
      <button @click="${this.move}">MOVE</button>
      <button @click="${this.flipBoard}">FLIP</button>

      <pre>MOVEMENTE: ${JSON.stringify(this.pieceMovement)}</pre>

      <div part="base" class="outer">
        <div class="ranks-label">${this.renderRanksLabel()}</div>

        <div class="board" data-piece-set="${this.pieceSet}" data-orientation="${this.orientation}">
          ${map(this.board.squares, row =>
            map(
              row,
              square => html`
                <my-square
                  .square="${square}"
                  ?selected="${this.isSelected(square)}"
                  ?underAttack="${this.isUnderAttack(square)}"
                  ?highlightFrom="${this.isHistoryFromMove(square)}"
                  ?highlightTo="${this.isHistoryToMove(square)}"
                  @square-clicked="${this.onSquareClick}"
                  ?dragEnabled="${this.pieceMovement === 'drag'}"
                ></my-square>
              `
            )
          )}
        </div>

        <div class="files-label">${this.renderFilesLabel()}</div>
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'my-board': MyBoard;
  }
}
