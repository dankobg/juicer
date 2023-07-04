import type {
  LoginFlow,
  RecoveryFlow,
  RegistrationFlow,
  SettingsFlow,
  VerificationFlow,
  UiNode,
  UiNodeAnchorAttributes,
  UiNodeAttributes,
  UiNodeImageAttributes,
  UiNodeInputAttributes,
  UiNodeScriptAttributes,
  UiNodeTextAttributes,
} from '@ory/client';
import type { AxiosError } from '../../app';

export const ORY_KRATOS_UI_TEST_ID_PREFIX = 'ory-kratos-ui';

export function createTestId(id: string): string {
  return `${ORY_KRATOS_UI_TEST_ID_PREFIX}-${id}`;
}

export function isUiNodeAnchorAttributes(attrs: UiNodeAttributes): attrs is UiNodeAnchorAttributes {
  return attrs.node_type === 'a';
}

export function isUiNodeImageAttributes(attrs: UiNodeAttributes): attrs is UiNodeImageAttributes {
  return attrs.node_type === 'img';
}

export function isUiNodeInputAttributes(attrs: UiNodeAttributes): attrs is UiNodeInputAttributes {
  return attrs.node_type === 'input';
}

export function isUiNodeTextAttributes(attrs: UiNodeAttributes): attrs is UiNodeTextAttributes {
  return attrs.node_type === 'text';
}

export function isUiNodeScriptAttributes(attrs: UiNodeAttributes): attrs is UiNodeScriptAttributes {
  return attrs.node_type === 'script';
}

export function getNodeId({ attributes }: UiNode): string {
  if (isUiNodeInputAttributes(attributes)) {
    return attributes.name;
  } else {
    return attributes.id;
  }
}

export function getNodeLabel(node: UiNode): string {
  const attributes = node.attributes;
  if (isUiNodeAnchorAttributes(attributes)) {
    return attributes.title.text;
  }

  if (isUiNodeImageAttributes(attributes)) {
    return node.meta.label?.text || '';
  }

  if (isUiNodeInputAttributes(attributes)) {
    if (attributes.label?.text) {
      return attributes.label.text;
    }
  }

  return node.meta.label?.text || '';
}

export function filterNodesByGroups(
  nodes: UiNode[],
  groups?: string[] | string,
  withoutDefaultGroup?: boolean
): UiNode[] {
  if (!groups || groups.length === 0) {
    return nodes;
  }

  const search = typeof groups === 'string' ? groups.split(',') : groups;
  if (!withoutDefaultGroup) {
    search.push('default');
  }

  return nodes.filter(({ group }) => search.includes(group));
}

type AuthFlow = LoginFlow | RegistrationFlow | RecoveryFlow | VerificationFlow | SettingsFlow;

export function extractCSRFToken(flow: AuthFlow | null): string {
  if (!flow) {
    return '';
  }

  const csrfAttributes = flow?.ui?.nodes?.find(node => {
    return (
      isUiNodeInputAttributes(node.attributes) && node.group === 'default' && node.attributes.name === 'csrf_token'
    );
  })?.attributes;

  return csrfAttributes && isUiNodeInputAttributes(csrfAttributes) ? csrfAttributes?.value : '';
}

export function isAxiosError(err: unknown): err is AxiosError {
  return typeof err === 'object' && err !== null && 'response' in err;
}
