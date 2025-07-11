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
 * Paginated Relationship List
 * @export
 * @interface Relationships
 */
export interface Relationships {
    /**
     * The opaque token to provide in a subsequent request
     * to get the next page. It is the empty string iff this is
     * the last page.
     * @type {string}
     * @memberof Relationships
     */
    nextPageToken?: string;
    /**
     * 
     * @type {Array<Relationship>}
     * @memberof Relationships
     */
    relationTuples?: Array<Relationship>;
}

/**
 * Check if a given object implements the Relationships interface.
 */
export function instanceOfRelationships(value: object): value is Relationships {
    return true;
}

export function RelationshipsFromJSON(json: any): Relationships {
    return RelationshipsFromJSONTyped(json, false);
}

export function RelationshipsFromJSONTyped(json: any, ignoreDiscriminator: boolean): Relationships {
    if (json == null) {
        return json;
    }
    return {
        
        'nextPageToken': json['next_page_token'] == null ? undefined : json['next_page_token'],
        'relationTuples': json['relation_tuples'] == null ? undefined : ((json['relation_tuples'] as Array<any>).map(RelationshipFromJSON)),
    };
}

export function RelationshipsToJSON(json: any): Relationships {
    return RelationshipsToJSONTyped(json, false);
}

export function RelationshipsToJSONTyped(value?: Relationships | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'next_page_token': value['nextPageToken'],
        'relation_tuples': value['relationTuples'] == null ? undefined : ((value['relationTuples'] as Array<any>).map(RelationshipToJSON)),
    };
}

