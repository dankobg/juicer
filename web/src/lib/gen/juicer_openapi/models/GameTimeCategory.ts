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
 * 
 * @export
 * @interface GameTimeCategory
 */
export interface GameTimeCategory {
    /**
     * Game time category id
     * @type {string}
     * @memberof GameTimeCategory
     */
    id: string;
    /**
     * Game time category name
     * @type {string}
     * @memberof GameTimeCategory
     */
    name: string;
    /**
     * Game time category upper time limit
     * @type {number}
     * @memberof GameTimeCategory
     */
    upperTimeLimitSecs?: number;
}

/**
 * Check if a given object implements the GameTimeCategory interface.
 */
export function instanceOfGameTimeCategory(value: object): value is GameTimeCategory {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('name' in value) || value['name'] === undefined) return false;
    return true;
}

export function GameTimeCategoryFromJSON(json: any): GameTimeCategory {
    return GameTimeCategoryFromJSONTyped(json, false);
}

export function GameTimeCategoryFromJSONTyped(json: any, ignoreDiscriminator: boolean): GameTimeCategory {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'upperTimeLimitSecs': json['upper_time_limit_secs'] == null ? undefined : json['upper_time_limit_secs'],
    };
}

export function GameTimeCategoryToJSON(json: any): GameTimeCategory {
    return GameTimeCategoryToJSONTyped(json, false);
}

export function GameTimeCategoryToJSONTyped(value?: GameTimeCategory | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'id': value['id'],
        'name': value['name'],
        'upper_time_limit_secs': value['upperTimeLimitSecs'],
    };
}

