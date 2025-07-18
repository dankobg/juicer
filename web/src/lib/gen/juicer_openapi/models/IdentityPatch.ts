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
import type { CreateIdentityBody } from './CreateIdentityBody';
import {
    CreateIdentityBodyFromJSON,
    CreateIdentityBodyFromJSONTyped,
    CreateIdentityBodyToJSON,
    CreateIdentityBodyToJSONTyped,
} from './CreateIdentityBody';

/**
 * Payload for patching an identity
 * @export
 * @interface IdentityPatch
 */
export interface IdentityPatch {
    /**
     * 
     * @type {CreateIdentityBody}
     * @memberof IdentityPatch
     */
    create?: CreateIdentityBody;
    /**
     * The ID of this patch.
     * 
     * The patch ID is optional. If specified, the ID will be returned in the
     * response, so consumers of this API can correlate the response with the
     * patch.
     * @type {string}
     * @memberof IdentityPatch
     */
    patchId?: string;
}

/**
 * Check if a given object implements the IdentityPatch interface.
 */
export function instanceOfIdentityPatch(value: object): value is IdentityPatch {
    return true;
}

export function IdentityPatchFromJSON(json: any): IdentityPatch {
    return IdentityPatchFromJSONTyped(json, false);
}

export function IdentityPatchFromJSONTyped(json: any, ignoreDiscriminator: boolean): IdentityPatch {
    if (json == null) {
        return json;
    }
    return {
        
        'create': json['create'] == null ? undefined : CreateIdentityBodyFromJSON(json['create']),
        'patchId': json['patch_id'] == null ? undefined : json['patch_id'],
    };
}

export function IdentityPatchToJSON(json: any): IdentityPatch {
    return IdentityPatchToJSONTyped(json, false);
}

export function IdentityPatchToJSONTyped(value?: IdentityPatch | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'create': CreateIdentityBodyToJSON(value['create']),
        'patch_id': value['patchId'],
    };
}

