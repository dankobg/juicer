import { LitElement, html, nothing } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { ifDefined } from 'lit/directives/if-defined.js';
import styles from './my-square.styles';
import { Square } from '../shared/types';

import '../piece/my-piece';

@customElement('my-square')
export class MySquare extends LitElement {
  static styles = styles;

  @property()
  square!: Square;

  @property({ type: Boolean })
  selected: boolean = false;

  @property({ type: Boolean })
  underAttack: boolean = false;

  @property({ type: Boolean })
  highlightFrom: boolean = false;

  @property({ type: Boolean })
  highlighTo: boolean = false;

  @property({ type: Boolean })
  dragEnabled: boolean = false;

  private onSquareClick(): void {
    if (this.dragEnabled) {
      return;
    }

    const event = new CustomEvent('square-clicked', {
      detail: { square: this.square },
    });

    this.dispatchEvent(event);
  }

  private onDragOver(event: DragEvent): void {
    if (!this.dragEnabled) {
      return;
    }

    event.preventDefault();
    // console.log('onDragOver: ', event.target);
  }

  private onDragEnter(event: DragEvent): void {
    if (!this.dragEnabled) {
      return;
    }

    event.preventDefault();
    console.log('onDragEnter: ', event.target);
  }

  private onDrop(event: DragEvent): void {
    if (!this.dragEnabled) {
      return;
    }

    // console.log('onDrop: ', event.target);
  }

  private onDragLeave(event: DragEvent): void {
    if (!this.dragEnabled) {
      return;
    }

    // console.log('onDragLeave: ', event.target);
  }

  render() {
    return html`
      <div
        part="base"
        class="square"
        data-color="${this.square.color}"
        data-coord="${this.square.coordinate}"
        data-selected="${ifDefined(this.selected ? 'true' : undefined)}"
        data-under-attack="${ifDefined(this.underAttack ? 'true' : undefined)}"
        data-highlight-from="${ifDefined(this.highlightFrom ? 'true' : undefined)}"
        data-highlight-to="${ifDefined(this.highlighTo ? 'true' : undefined)}"
        @click="${this.onSquareClick}"
        @__dragover="${this.onDragOver}"
        @dragenter="${this.onDragEnter}"
        @dragleave="${this.onDragLeave}"
        @drop="${this.onDrop}"
      >
        ${this.square.hasPiece()
          ? html`
              <my-piece
                coord="${this.square.coordinate}"
                ?dragEnabled="${this.dragEnabled}"
                .piece="${this.square.piece}"
                ?selected="${this.selected}"
              ></my-piece>
            `
          : nothing}
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'my-square': MySquare;
  }
}
