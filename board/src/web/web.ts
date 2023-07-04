import { LitElement, html } from 'lit';
import { customElement } from 'lit/decorators.js';

import '../piece/my-piece';

@customElement('my-web')
export class MyWeb extends LitElement {
  ws: WebSocket | null = null;
  interval: NodeJS.Timer | null = null;

  connectedCallback() {
    super.connectedCallback();

    this.ws = new WebSocket('ws://localhost:1337');

    this.ws.addEventListener('open', () => {
      console.log('CONNECTED');
    });

    this.interval = setInterval(() => {
      if (this.ws) {
        this.ws.send('hello');
      }
    }, 2000);

    this.ws.addEventListener('close', () => {
      this.ws = null;
    });

    this.ws.addEventListener('message', event => {
      console.log(event.data);
    });
  }

  disconnectedCallback(): void {
    this.interval = null;
  }

  render() {
    return html`<div>app</div>`;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'my-web': MyWeb;
  }
}
