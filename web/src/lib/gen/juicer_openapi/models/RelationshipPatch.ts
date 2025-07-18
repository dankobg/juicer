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
import type { Relationship } from './Relationship';
import {
    RelationshipFromJSON,
    RelationshipFromJSONTyped,
    RelationshipToJSON,
    RelationshipToJSONTyped,
} from './Relationship';

/**
 * Payload for patching a relationship
 * @export
 * @interface RelationshipPatch
 */
export interface RelationshipPatch {
    /**
     * 
     * @type {string}
     * @memberof RelationshipPatch
     */
    action?: RelationshipPatchActionEnum;
    /**
     * 
     * @type {Relationship}
     * @memberof RelationshipPatch
     */
    relationTuple?: Relationship;
}


/**
 * @export
 */
export const RelationshipPatchActionEnum = {
    Insert: 'insert',
    Delete: 'delete'
} as const;
export type RelationshipPatchActionEnum = typeof RelationshipPatchActionEnum[keyof typeof RelationshipPatchActionEnum];


/**
 * Check if a given object implements the RelationshipPatch interface.
 */
export function instanceOfRelationshipPatch(value: object): value is RelationshipPatch {
    return true;
}

export function RelationshipPatchFromJSON(json: any): RelationshipPatch {
    return RelationshipPatchFromJSONTyped(json, false);
}

export function RelationshipPatchFromJSONTyped(json: any, ignoreDiscriminator: boolean): RelationshipPatch {
    if (json == null) {
        return json;
    }
    return {
        
        'action': json['action'] == null ? undefined : json['action'],
        'relationTuple': json['relation_tuple'] == null ? undefined : RelationshipFromJSON(json['relation_tuple']),
    };
}

export function RelationshipPatchToJSON(json: any): RelationshipPatch {
    return RelationshipPatchToJSONTyped(json, false);
}

export function RelationshipPatchToJSONTyped(value?: RelationshipPatch | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'action': value['action'],
        'relation_tuple': RelationshipToJSON(value['relationTuple']),
    };
}

