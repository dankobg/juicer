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
 * Update Login Flow with OpenID Connect Method
 * @export
 * @interface UpdateLoginFlowWithOidcMethod
 */
export interface UpdateLoginFlowWithOidcMethod {
    /**
     * The CSRF Token
     * @type {string}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    csrfToken?: string;
    /**
     * IDToken is an optional id token provided by an OIDC provider
     * 
     * If submitted, it is verified using the OIDC provider's public key set and the claims are used to populate
     * the OIDC credentials of the identity.
     * If the OIDC provider does not store additional claims (such as name, etc.) in the IDToken itself, you can use
     * the `traits` field to populate the identity's traits. Note, that Apple only includes the users email in the IDToken.
     * 
     * Supported providers are
     * Apple
     * Google
     * @type {string}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    idToken?: string;
    /**
     * IDTokenNonce is the nonce, used when generating the IDToken.
     * If the provider supports nonce validation, the nonce will be validated against this value and required.
     * @type {string}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    idTokenNonce?: string;
    /**
     * Method to use
     * 
     * This field must be set to `oidc` when using the oidc method.
     * @type {string}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    method: string;
    /**
     * The provider to register with
     * @type {string}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    provider: string;
    /**
     * The identity traits. This is a placeholder for the registration flow.
     * @type {object}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    traits?: object;
    /**
     * Transient data to pass along to any webhooks
     * @type {object}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    transientPayload?: object;
    /**
     * UpstreamParameters are the parameters that are passed to the upstream identity provider.
     * 
     * These parameters are optional and depend on what the upstream identity provider supports.
     * Supported parameters are:
     * `login_hint` (string): The `login_hint` parameter suppresses the account chooser and either pre-fills the email box on the sign-in form, or selects the proper session.
     * `hd` (string): The `hd` parameter limits the login/registration process to a Google Organization, e.g. `mycollege.edu`.
     * `prompt` (string): The `prompt` specifies whether the Authorization Server prompts the End-User for reauthentication and consent, e.g. `select_account`.
     * @type {object}
     * @memberof UpdateLoginFlowWithOidcMethod
     */
    upstreamParameters?: object;
}

/**
 * Check if a given object implements the UpdateLoginFlowWithOidcMethod interface.
 */
export function instanceOfUpdateLoginFlowWithOidcMethod(value: object): value is UpdateLoginFlowWithOidcMethod {
    if (!('method' in value) || value['method'] === undefined) return false;
    if (!('provider' in value) || value['provider'] === undefined) return false;
    return true;
}

export function UpdateLoginFlowWithOidcMethodFromJSON(json: any): UpdateLoginFlowWithOidcMethod {
    return UpdateLoginFlowWithOidcMethodFromJSONTyped(json, false);
}

export function UpdateLoginFlowWithOidcMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateLoginFlowWithOidcMethod {
    if (json == null) {
        return json;
    }
    return {
        
        'csrfToken': json['csrf_token'] == null ? undefined : json['csrf_token'],
        'idToken': json['id_token'] == null ? undefined : json['id_token'],
        'idTokenNonce': json['id_token_nonce'] == null ? undefined : json['id_token_nonce'],
        'method': json['method'],
        'provider': json['provider'],
        'traits': json['traits'] == null ? undefined : json['traits'],
        'transientPayload': json['transient_payload'] == null ? undefined : json['transient_payload'],
        'upstreamParameters': json['upstream_parameters'] == null ? undefined : json['upstream_parameters'],
    };
}

export function UpdateLoginFlowWithOidcMethodToJSON(json: any): UpdateLoginFlowWithOidcMethod {
    return UpdateLoginFlowWithOidcMethodToJSONTyped(json, false);
}

export function UpdateLoginFlowWithOidcMethodToJSONTyped(value?: UpdateLoginFlowWithOidcMethod | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'csrf_token': value['csrfToken'],
        'id_token': value['idToken'],
        'id_token_nonce': value['idTokenNonce'],
        'method': value['method'],
        'provider': value['provider'],
        'traits': value['traits'],
        'transient_payload': value['transientPayload'],
        'upstream_parameters': value['upstreamParameters'],
    };
}

