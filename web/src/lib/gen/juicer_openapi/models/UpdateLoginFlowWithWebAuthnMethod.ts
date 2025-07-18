/* tslint:disable */
/* eslint-disable */
/**
 * Juicer schema
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from '../runtime';
/**
 * Update Login Flow with WebAuthn Method
 * @export
 * @interface UpdateLoginFlowWithWebAuthnMethod
 */
export interface UpdateLoginFlowWithWebAuthnMethod {
    /**
     * Sending the anti-csrf token is only required for browser login flows.
     * @type {string}
     * @memberof UpdateLoginFlowWithWebAuthnMethod
     */
    csrfToken?: string;
    /**
     * Identifier is the email or username of the user trying to log in.
     * @type {string}
     * @memberof UpdateLoginFlowWithWebAuthnMethod
     */
    identifier: string;
    /**
     * Method should be set to "webAuthn" when logging in using the WebAuthn strategy.
     * @type {string}
     * @memberof UpdateLoginFlowWithWebAuthnMethod
     */
    method: string;
    /**
     * Transient data to pass along to any webhooks
     * @type {object}
     * @memberof UpdateLoginFlowWithWebAuthnMethod
     */
    transientPayload?: object;
    /**
     * Login a WebAuthn Security Key
     * 
     * This must contain the ID of the WebAuthN connection.
     * @type {string}
     * @memberof UpdateLoginFlowWithWebAuthnMethod
     */
    webauthnLogin?: string;
}

/**
 * Check if a given object implements the UpdateLoginFlowWithWebAuthnMethod interface.
 */
export function instanceOfUpdateLoginFlowWithWebAuthnMethod(value: object): value is UpdateLoginFlowWithWebAuthnMethod {
    if (!('identifier' in value) || value['identifier'] === undefined) return false;
    if (!('method' in value) || value['method'] === undefined) return false;
    return true;
}

export function UpdateLoginFlowWithWebAuthnMethodFromJSON(json: any): UpdateLoginFlowWithWebAuthnMethod {
    return UpdateLoginFlowWithWebAuthnMethodFromJSONTyped(json, false);
}

export function UpdateLoginFlowWithWebAuthnMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateLoginFlowWithWebAuthnMethod {
    if (json == null) {
        return json;
    }
    return {
        
        'csrfToken': json['csrf_token'] == null ? undefined : json['csrf_token'],
        'identifier': json['identifier'],
        'method': json['method'],
        'transientPayload': json['transient_payload'] == null ? undefined : json['transient_payload'],
        'webauthnLogin': json['webauthn_login'] == null ? undefined : json['webauthn_login'],
    };
}

export function UpdateLoginFlowWithWebAuthnMethodToJSON(json: any): UpdateLoginFlowWithWebAuthnMethod {
    return UpdateLoginFlowWithWebAuthnMethodToJSONTyped(json, false);
}

export function UpdateLoginFlowWithWebAuthnMethodToJSONTyped(value?: UpdateLoginFlowWithWebAuthnMethod | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'csrf_token': value['csrfToken'],
        'identifier': value['identifier'],
        'method': value['method'],
        'transient_payload': value['transientPayload'],
        'webauthn_login': value['webauthnLogin'],
    };
}

