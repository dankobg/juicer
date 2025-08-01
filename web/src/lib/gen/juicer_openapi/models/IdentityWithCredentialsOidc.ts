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
import type { IdentityWithCredentialsOidcConfig } from './IdentityWithCredentialsOidcConfig';
import {
    IdentityWithCredentialsOidcConfigFromJSON,
    IdentityWithCredentialsOidcConfigFromJSONTyped,
    IdentityWithCredentialsOidcConfigToJSON,
    IdentityWithCredentialsOidcConfigToJSONTyped,
} from './IdentityWithCredentialsOidcConfig';

/**
 * Create Identity and Import Social Sign In Credentials
 * @export
 * @interface IdentityWithCredentialsOidc
 */
export interface IdentityWithCredentialsOidc {
    /**
     * 
     * @type {IdentityWithCredentialsOidcConfig}
     * @memberof IdentityWithCredentialsOidc
     */
    config?: IdentityWithCredentialsOidcConfig;
}

/**
 * Check if a given object implements the IdentityWithCredentialsOidc interface.
 */
export function instanceOfIdentityWithCredentialsOidc(value: object): value is IdentityWithCredentialsOidc {
    return true;
}

export function IdentityWithCredentialsOidcFromJSON(json: any): IdentityWithCredentialsOidc {
    return IdentityWithCredentialsOidcFromJSONTyped(json, false);
}

export function IdentityWithCredentialsOidcFromJSONTyped(json: any, ignoreDiscriminator: boolean): IdentityWithCredentialsOidc {
    if (json == null) {
        return json;
    }
    return {
        
        'config': json['config'] == null ? undefined : IdentityWithCredentialsOidcConfigFromJSON(json['config']),
    };
}

export function IdentityWithCredentialsOidcToJSON(json: any): IdentityWithCredentialsOidc {
    return IdentityWithCredentialsOidcToJSONTyped(json, false);
}

export function IdentityWithCredentialsOidcToJSONTyped(value?: IdentityWithCredentialsOidc | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'config': IdentityWithCredentialsOidcConfigToJSON(value['config']),
    };
}

