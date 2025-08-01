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
 * Used when an administrator creates a recovery code for an identity.
 * @export
 * @interface RecoveryCodeForIdentity
 */
export interface RecoveryCodeForIdentity {
    /**
     * Expires At is the timestamp of when the recovery flow expires
     * 
     * The timestamp when the recovery code expires.
     * @type {Date}
     * @memberof RecoveryCodeForIdentity
     */
    expiresAt?: Date;
    /**
     * RecoveryCode is the code that can be used to recover the account
     * @type {string}
     * @memberof RecoveryCodeForIdentity
     */
    recoveryCode: string;
    /**
     * RecoveryLink with flow
     * 
     * This link opens the recovery UI with an empty `code` field.
     * @type {string}
     * @memberof RecoveryCodeForIdentity
     */
    recoveryLink: string;
}

/**
 * Check if a given object implements the RecoveryCodeForIdentity interface.
 */
export function instanceOfRecoveryCodeForIdentity(value: object): value is RecoveryCodeForIdentity {
    if (!('recoveryCode' in value) || value['recoveryCode'] === undefined) return false;
    if (!('recoveryLink' in value) || value['recoveryLink'] === undefined) return false;
    return true;
}

export function RecoveryCodeForIdentityFromJSON(json: any): RecoveryCodeForIdentity {
    return RecoveryCodeForIdentityFromJSONTyped(json, false);
}

export function RecoveryCodeForIdentityFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecoveryCodeForIdentity {
    if (json == null) {
        return json;
    }
    return {
        
        'expiresAt': json['expires_at'] == null ? undefined : (new Date(json['expires_at'])),
        'recoveryCode': json['recovery_code'],
        'recoveryLink': json['recovery_link'],
    };
}

export function RecoveryCodeForIdentityToJSON(json: any): RecoveryCodeForIdentity {
    return RecoveryCodeForIdentityToJSONTyped(json, false);
}

export function RecoveryCodeForIdentityToJSONTyped(value?: RecoveryCodeForIdentity | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'expires_at': value['expiresAt'] == null ? undefined : ((value['expiresAt']).toISOString()),
        'recovery_code': value['recoveryCode'],
        'recovery_link': value['recoveryLink'],
    };
}

