import { LitElement, html } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { ifDefined } from 'lit/directives/if-defined.js';
import styles from './my-piece.styles';
import { Coordinate, Piece } from '../shared/types';

@customElement('my-piece')
export class MyPiece extends LitElement {
  static styles = styles;

  @property()
  coord!: Coordinate;

  @property()
  piece!: Piece;

  @property({ type: Boolean })
  dragEnabled: boolean = false;

  @property({ type: Boolean })
  selected: boolean = false;

  private draggedItem: HTMLElement | null = null;
  private dx: number = 0;
  private dy: number = 0;

  private onDrag(event: DragEvent) {
    if (!this.dragEnabled) {
      return;
    }

    // console.log('onDrag: ', event);

    if (this.draggedItem !== null) {
      this.draggedItem.style.left = `${event.clientX - this.dx}`;
      this.draggedItem.style.top = `${event.clientY - this.dy}`;
    }

    this.requestUpdate();
  }

  private onDragStart(event: DragEvent) {
    if (!this.dragEnabled) {
      return;
    }

    // console.log('onDragStart: ', event);

    this.draggedItem = event.target as HTMLElement;
    const rect = this.draggedItem.getBoundingClientRect();

    this.dx = event.clientX - rect.x;
    this.dy = event.clientX - rect.y;
    this.style.position = 'absolute';

    this.requestUpdate();
  }

  private onDragEnd(event: DragEvent) {
    if (!this.dragEnabled) {
      return;
    }

    // console.log('onDragEnd: ', event);

    const retainPosition = true;

    if (this.draggedItem !== null) {
      if (retainPosition) {
        this.draggedItem.style.left = `${event.clientX - this.dx}`;
        this.draggedItem.style.top = `${event.clientY - this.dy}`;
      } else {
        this.draggedItem.style.position = 'static';
      }
    }

    // this.draggedItem = null;
    // this.dx = 0;
    // this.dy = 0;

    this.requestUpdate();
  }

  render() {
    return html`
      <div
        id="${this.coord}-piece"
        part="base"
        class="piece"
        data-piece-fen="${ifDefined(this.piece.toFENPieceSymbol())}"
        data-selected="${ifDefined(this.selected ? 'true' : undefined)}"
        @drag="${this.onDrag}"
        @dragstart="${this.onDragStart}"
        @dragend="${this.onDragEnd}"
        draggable="${this.dragEnabled ? 'true' : 'false'}"
      ></div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'my-piece': MyPiece;
  }
}
