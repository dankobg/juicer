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
 * Update Login flow using the code method
 * @export
 * @interface UpdateLoginFlowWithCodeMethod
 */
export interface UpdateLoginFlowWithCodeMethod {
    /**
     * Address is the address to send the code to, in case that there are multiple addresses. This field
     * is only used in two-factor flows and is ineffective for passwordless flows.
     * @type {string}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    address?: string;
    /**
     * Code is the 6 digits code sent to the user
     * @type {string}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    code?: string;
    /**
     * CSRFToken is the anti-CSRF token
     * @type {string}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    csrfToken: string;
    /**
     * Identifier is the code identifier
     * The identifier requires that the user has already completed the registration or settings with code flow.
     * @type {string}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    identifier?: string;
    /**
     * Method should be set to "code" when logging in using the code strategy.
     * @type {string}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    method: string;
    /**
     * Resend is set when the user wants to resend the code
     * @type {string}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    resend?: string;
    /**
     * Transient data to pass along to any webhooks
     * @type {object}
     * @memberof UpdateLoginFlowWithCodeMethod
     */
    transientPayload?: object;
}

/**
 * Check if a given object implements the UpdateLoginFlowWithCodeMethod interface.
 */
export function instanceOfUpdateLoginFlowWithCodeMethod(value: object): value is UpdateLoginFlowWithCodeMethod {
    if (!('csrfToken' in value) || value['csrfToken'] === undefined) return false;
    if (!('method' in value) || value['method'] === undefined) return false;
    return true;
}

export function UpdateLoginFlowWithCodeMethodFromJSON(json: any): UpdateLoginFlowWithCodeMethod {
    return UpdateLoginFlowWithCodeMethodFromJSONTyped(json, false);
}

export function UpdateLoginFlowWithCodeMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateLoginFlowWithCodeMethod {
    if (json == null) {
        return json;
    }
    return {
        
        'address': json['address'] == null ? undefined : json['address'],
        'code': json['code'] == null ? undefined : json['code'],
        'csrfToken': json['csrf_token'],
        'identifier': json['identifier'] == null ? undefined : json['identifier'],
        'method': json['method'],
        'resend': json['resend'] == null ? undefined : json['resend'],
        'transientPayload': json['transient_payload'] == null ? undefined : json['transient_payload'],
    };
}

export function UpdateLoginFlowWithCodeMethodToJSON(json: any): UpdateLoginFlowWithCodeMethod {
    return UpdateLoginFlowWithCodeMethodToJSONTyped(json, false);
}

export function UpdateLoginFlowWithCodeMethodToJSONTyped(value?: UpdateLoginFlowWithCodeMethod | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'address': value['address'],
        'code': value['code'],
        'csrf_token': value['csrfToken'],
        'identifier': value['identifier'],
        'method': value['method'],
        'resend': value['resend'],
        'transient_payload': value['transientPayload'],
    };
}

