import { LitElement, html } from 'lit';
import { customElement } from 'lit/decorators.js';

import './board/my-board';
import './web/web';

@customElement('my-wat')
export class MyWat extends LitElement {
  render() {
    return html`
      <!-- <my-web></my-web> -->
      <my-board pieceMovement="click"></my-board>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'my-wat': MyWat;
  }
}
