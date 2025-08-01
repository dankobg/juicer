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
 * Used when an administrator creates a recovery link for an identity.
 * @export
 * @interface RecoveryLinkForIdentity
 */
export interface RecoveryLinkForIdentity {
    /**
     * Recovery Link Expires At
     * 
     * The timestamp when the recovery link expires.
     * @type {Date}
     * @memberof RecoveryLinkForIdentity
     */
    expiresAt?: Date;
    /**
     * Recovery Link
     * 
     * This link can be used to recover the account.
     * @type {string}
     * @memberof RecoveryLinkForIdentity
     */
    recoveryLink: string;
}

/**
 * Check if a given object implements the RecoveryLinkForIdentity interface.
 */
export function instanceOfRecoveryLinkForIdentity(value: object): value is RecoveryLinkForIdentity {
    if (!('recoveryLink' in value) || value['recoveryLink'] === undefined) return false;
    return true;
}

export function RecoveryLinkForIdentityFromJSON(json: any): RecoveryLinkForIdentity {
    return RecoveryLinkForIdentityFromJSONTyped(json, false);
}

export function RecoveryLinkForIdentityFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecoveryLinkForIdentity {
    if (json == null) {
        return json;
    }
    return {
        
        'expiresAt': json['expires_at'] == null ? undefined : (new Date(json['expires_at'])),
        'recoveryLink': json['recovery_link'],
    };
}

export function RecoveryLinkForIdentityToJSON(json: any): RecoveryLinkForIdentity {
    return RecoveryLinkForIdentityToJSONTyped(json, false);
}

export function RecoveryLinkForIdentityToJSONTyped(value?: RecoveryLinkForIdentity | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'expires_at': value['expiresAt'] == null ? undefined : ((value['expiresAt']).toISOString()),
        'recovery_link': value['recoveryLink'],
    };
}

