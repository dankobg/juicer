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
 * @interface GameStat
 */
export interface GameStat {
    /**
     * 
     * @type {number}
     * @memberof GameStat
     */
    win: number;
    /**
     * 
     * @type {number}
     * @memberof GameStat
     */
    loss: number;
    /**
     * 
     * @type {number}
     * @memberof GameStat
     */
    draw: number;
    /**
     * 
     * @type {number}
     * @memberof GameStat
     */
    interrupted: number;
    /**
     * 
     * @type {number}
     * @memberof GameStat
     */
    total?: number;
}

/**
 * Check if a given object implements the GameStat interface.
 */
export function instanceOfGameStat(value: object): value is GameStat {
    if (!('win' in value) || value['win'] === undefined) return false;
    if (!('loss' in value) || value['loss'] === undefined) return false;
    if (!('draw' in value) || value['draw'] === undefined) return false;
    if (!('interrupted' in value) || value['interrupted'] === undefined) return false;
    return true;
}

export function GameStatFromJSON(json: any): GameStat {
    return GameStatFromJSONTyped(json, false);
}

export function GameStatFromJSONTyped(json: any, ignoreDiscriminator: boolean): GameStat {
    if (json == null) {
        return json;
    }
    return {
        
        'win': json['win'],
        'loss': json['loss'],
        'draw': json['draw'],
        'interrupted': json['interrupted'],
        'total': json['total'] == null ? undefined : json['total'],
    };
}

export function GameStatToJSON(json: any): GameStat {
    return GameStatToJSONTyped(json, false);
}

export function GameStatToJSONTyped(value?: GameStat | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'win': value['win'],
        'loss': value['loss'],
        'draw': value['draw'],
        'interrupted': value['interrupted'],
        'total': value['total'],
    };
}

